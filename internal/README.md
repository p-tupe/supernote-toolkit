# .note file format

This is my understanding of how the file is structured, and the steps taken to de-construct it.

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

Here's one possible pseudo code for decoding a RATTA_RLE byte stream -

```
for each current pair of [color code, length code] bytes in an RLE stream:

        if length code == 0xff (is long run / blank line):
            blank line length = 0x4000
            process pair as [color code, blank line length]

        else if length code's most significant bit is set (is multi-byte length):
            get next pair of bytes as [next color code, next length code]
            if next color code == color code:
                length = 1 + int(next length code) + parsed(length code)
                process pair as [col, length]
            else
                length = parsed(length code)
                process pair as [col, length]


        else
            length = 1 + int(length code)
            process pair as [col, length]


where processing a pair of [color, length]:
    // just col repeated length times
    creating an array of [col, ...] for given length

where parsed length code:
    // last 7 bits + 1, left shifted 7 times
    (int(length code & 0x7f) + 1) << 7
```

Now, the color code is not set in our typical RGB format - we have to use a `CodeMap` for convert from the color code bytes to specific color we need to represent. What we do is simple, assume we have a code map like so:

```
{ a = Black, c = Light Gray, e = Dark Gray, g = White }
```

Then we just loop each color code in the array we got from above and replace it with corresponding RGB color.

Sometimes, we may encounter a code not representing a specific color - that is the "intensity" in greyscale format specifically for antialiasing handwriting. Just convert it to RGB directly.

Ideally, the decoded array should equal width \* length of our device. Once we have all the decoded layers, we simply iterate over them from "LAYERSEQ" order and overlay them on top of each other.

> If the lengths differ, can pad with transparent bits as a workaround

Once all the layers are overlaid in order, our "Page" is ready! Follow the process for each page and Voila - our notebook is ready to read as a PNG.
