# .note file format

This is my understanding of how the file is structured, and the steps taken to de-construct it.

> Special thanks to [jya-dev](https://github.com/jya-dev/supernote-tool) and [RohanGautam](https://github.com/RohanGautam/supernote_pdf) for doing the hard work already!

### Structure

```
-----------------
Signature.Header
-----------------
Pages
-----------------
    Page1 Layers
        Layer1
        Layer2
        Layer3
-----------------
    Page2 Layers
        Layer1
        Layer2
        Layer3
-----------------
    Page3 Layers
        Layer1
        Layer2
        Layer3
-----------------
...
-----------------
Footer
-----------------
```

### Parser

The first 4 bytes of a note are always the constant string "note". I use this to filter out non-note files, in the `isNote` function.

If we skip those 4 bytes, then the next 20 hold the `signature` string. We can use this to calculate the header address, but since we don't know for sure if the length of signature will remain the same, we use the header address from elsewhere.

The last 4 bytes of the file contain the starting address of the "Footer" - set as `footerAddr` in code. If we convert the bytes from `footerAddr` to `EOF`, we get some info like: `<FILE_FEATURE:24><PAGE1:1687><PAGE2:4994>...` Where `FILE_FEATURE` represents the starting address of the "Header", and the `<PAGE1:...>...` is named for starting address of each page in the notebook.

I named such strings as `metadata` and use the function `parseMetadata` in `utils.go` that uses a regex to pull the above into a map, which we then parse into structs as applicable.

The address is stored as bytes, which we need to convert into little-endian uint32 format and cast as uint64 numbers to use for seeking location in the file.

Also, while we didn't use it earlier in case of footer, but if we pull 4 bytes from the starting address what we get is the length of the proceeding block in bytes.

We use this info to pull header data - First, get length of header by reading 4 bytes from `headerAddr` and finding `headerLen`, then we pull all bytes from `headerAddr` + 4 till `headerLen`.

This is a frequent enough occurance that I made a common function `readBlock` in `utils.go` that does just that.

Once we convert the header bytes into a string, we find a format similar to footer data, but with different keys. Specifically, we're looking at `APPLICATION_EQUIPMENT`, which we use to determine the type of device we're working with.

Then we refer to the page address in the footer and use `readBlock` on each to get the page data from the file.

Page data generally consists of layers and their details. Some of the keys we pick up from there are: `LAYERSEQ` - the sequence of layers, and the starting address of each layer. Note that the order in which layers are rendered is important.

We follow the layer sequence and use the `layerAdder` with `readBlock` to read the metadata for each layer. The metadata for each layer consists of a couple of keys we need to properly decrypt the layer, namely `LAYERPROTOCOL` and `LAYERBITMAP`.

At this point, we have parsed all the info we need to decode and convert the .note file into something else.

### Decoder

A Layer can be of either "PNG", "TEXT" or "RATTA_RLE" type. The former two are straight-forward, below is an explanation for the last (for X2 devices):

> Check the link for more info about [Run Length Encoding](https://en.wikipedia.org/wiki/Run-length_encoding) algorithm.

Assuming a stream of bytes like so:

```
[a, b, c, d, e, f, ...]
```

The value at a pixel is represented by a successive pair of byte, which here would be `(a, b)`, `(c, d)`, `(e, f)` etc, with the first being `color` and second being `length` - the color `a` is repeated `b` times, then `c` is repeated `d` times, and so on.

However, color codes `a` & `c` may be same when length is too large to be held in a single byte. In that case, we have to extend the length of `a` to equal `b + d`.

Sometimes, the color code may be redundant when the length marker is used to signal other stuff.

The decoder uses a "holder" state machine to handle multi-byte lengths. When the high bit is set on a length code, we can't resolve the run yet -- we hold the pair and read the next one before deciding:

```
for each pair [color code, length code] in the RLE stream:

    if holder is set:
        get [prev color color, prev length code] from holder

        if color code == prev color code:
            combine into one long run as:
                length = 1 + int(length code) + parsed(prev length code)
            then process [color, length]
        else if different color:
            first process [prev color code, parsed(prev length code)]
            then process [color code, int(length code) + 1]

    else if length code == 0xff (blank line):
        process [color, 0x4000]

    else if length code has MSB set (multi-byte marker):
        set holder as [color code, length code]

    else:
        process [color code, int(length code) + 1]

if unresolved holder:
    process [holder color, min(parsed(holder length code), remaining pixels)]

where process [color code, length]:
    fill `length` pixels with `rgba` wrt color code in code-to-rgba map

where parsed(length code):
    (int(length code & 0x7f) + 1) << 7
```

Ideally, the decoded pixel count should equal width \* height of our device. Once we have all the decoded layers, we iterate over them in "LAYERSEQ" order and overlay them on top of each other.

> If the lengths differ, clamping to the pixel buffer size acts as a workaround.

Once all the layers are overlaid in order, our "Page" is ready! Pages are processed concurrently. Follow the process for each page and Voila - our notebook is ready to read as a PNG.

As for PDF - I kinda just copied over what was already present (in other repos) for the structure and representation. It seems to work so I left it at that. Feel free to explore yourself though.
