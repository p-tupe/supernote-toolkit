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

The first 4 bytes of a note are always the constant string "note". I use this to filter out non-.note files, in the `isNote` function.

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

At this point, we have parsed all the info we need to decrypt and convert the .note file into something else.

### Decrypter
