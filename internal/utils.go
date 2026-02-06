package internal

import (
	"bytes"
	"encoding/binary"
	"image"
	c "image/color"
	"os"
	"regexp"
)

// Ensure that file stream starts with 'note' byte,
// otherwise unsupported
func isNote(file *os.File) (bool, error) {
	start := make([]byte, 4)
	_, err := file.Read(start)
	if err != nil {
		return false, err
	}

	return bytes.Equal(start, []byte("note")), nil
}

// readBlock takes a file and startAddr and returns a block of bytes
// from startAddr to the end of block formatted as string.
//
// End of block is the length defined by the first 4 bytes from startAddr;
// It must be converted from bytes to little-endian uint32
func readBlock(file *os.File, startAddr int64) (string, error) {
	rawbytes, err := readBlockAsBytes(file, startAddr)
	if err != nil {
		return "", err
	}

	return string(rawbytes), nil
}

// readBlock takes a file and startAddr and returns a block of bytes
// from startAddr to the end of block as is.
func readBlockAsBytes(file *os.File, startAddr int64) ([]byte, error) {
	block := make([]byte, 4)
	_, err := file.ReadAt(block, startAddr)
	if err != nil {
		return nil, err
	}
	lenAddr := binary.LittleEndian.Uint32(block)

	rawBytes := make([]byte, lenAddr)
	_, err = file.ReadAt(rawBytes, startAddr+4)
	if err != nil {
		return nil, err
	}

	return rawBytes, nil
}

var metadataRegex = regexp.MustCompile(`<(\w+):([\w|\d|,]+)>`)

// parseMetadata uses a regex to read a string of form
// <KEY1:VAL1><KEY2:VAL2>... and returns a map of the same as
// [KEY1:VAL1 KEY2:VAL2].
func parseMetadata(str string) map[string]string {
	metaData := map[string]string{}

	matches := metadataRegex.FindAllStringSubmatch(str, -1)
	for _, m := range matches {
		metaData[m[1]] = m[2]
	}

	return metaData
}

// decodeRLE converts a stream of bytes compressed using
// Ratta's RLE algorithm into corresponding color-codes
// and maps them onto a canvas using the device code-map.
func decodeRLE(data []byte, notebook *Notebook, img *image.RGBA) {
	expectedLen := img.Bounds().Dx() * img.Bounds().Dy()
	decompressed := make([]byte, 0, expectedLen)

	i := 0
	for {
		if i+1 >= len(data)-1 {
			break
		}

		colorCode, lengthCode := data[i], data[i+1]
		i += 2

		var length int

		if lengthCode == 0xff {
			length = 0x4000
		} else if lengthCode&0x80 != 0 {
			if i+3 >= len(data)-1 {
				break
			}
			nextColCode, nextLenCode := data[i+2], data[i+3]

			if colorCode == nextColCode {
				length = 1 + int(nextLenCode) + (((int(lengthCode & 0x7f)) + 1) << 7)
				i += 2
			} else {
				heldLen := (((int(lengthCode & 0x7f)) + 1) << 7)
				decompressed = append(decompressed, bytes.Repeat([]byte{colorCode}, heldLen)...)
			}

		} else {
			length = int(lengthCode) + 1
		}

		decompressed = append(decompressed, bytes.Repeat([]byte{colorCode}, length)...)
	}

	width := img.Bounds().Dx()
	for i, d := range decompressed {
		col := getColorFromCode(d, notebook.Device.CodeMap)
		img.SetRGBA(i%width, i/width, col)
	}
}

func getColorFromCode(b byte, codeMap map[byte]c.RGBA) c.RGBA {
	col, ok := codeMap[b]
	if !ok {
		col = c.RGBA{b, b, b, 255}
	}
	return col
}
