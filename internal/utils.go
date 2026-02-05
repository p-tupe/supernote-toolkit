package internal

import (
	"bytes"
	"encoding/binary"
	"image"
	c "image/color"
	"log"
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
	_, err := file.Seek(startAddr, 0)
	if err != nil {
		return nil, err
	}

	block := make([]byte, 4)
	_, err = file.Read(block)
	if err != nil {
		return nil, err
	}
	lenAddr := binary.LittleEndian.Uint32(block)

	rawBytes := make([]byte, lenAddr)
	_, err = file.Read(rawBytes)
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
func decodeRLE(data []byte, notebook *Notebook, bounds image.Rectangle) image.Image {
	expectedLen := bounds.Dx() * bounds.Dy()
	decompressed := make([]byte, 0, expectedLen)

	holder := []byte{0, 0}
	for i := 0; i < len(data); {
		if i+1 >= len(data) {
			break
		}

		colorCode, lengthCode := data[i], data[i+1]
		i += 2

		var length int

		if holder[1] != 0 {
			prevColCode, prevLenCode := holder[0], holder[1]
			holder = []byte{0, 0}

			if colorCode == prevColCode {
				length = 1 + int(lengthCode) + (((int(prevLenCode & 0x7f)) + 1) << 7)
			} else {
				heldLen := (((int(prevLenCode & 0x7f)) + 1) << 7)
				decompressed = append(decompressed, bytes.Repeat([]byte{prevColCode}, heldLen)...)
				length = int(lengthCode) + 1
			}
		} else if lengthCode == 0xff {
			length = 0x4000
		} else if lengthCode&0x80 != 0 {
			holder = []byte{colorCode, lengthCode}
			continue
		} else {
			length = int(lengthCode) + 1
		}

		decompressed = append(decompressed, bytes.Repeat([]byte{colorCode}, length)...)
	}

	if len(decompressed) != expectedLen {
		log.Println("Length mismatch, expected vs got:", expectedLen, len(decompressed))
	}

	img := image.NewRGBA(bounds)
	i := 0
	for y := range bounds.Max.Y {
		for x := range bounds.Max.X {
			col := getColorFromCode(decompressed[i], notebook.Device.CodeMap)
			img.SetRGBA(x, y, col)
			i++
		}
	}

	return img
}

func getColorFromCode(b byte, codeMap map[byte]c.RGBA) c.RGBA {
	col, ok := codeMap[b]
	if !ok {
		col = c.RGBA{b, b, b, 255}
	}
	return col
}
