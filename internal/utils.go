package internal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"regexp"
	"slices"
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

// TODO: Partial
// decodeRLE converts a stream of bytes compressed using
// Ratta's RLE algorithm into corresponding color-codes
// and maps them onto a canvas using the device code-map.
func decodeRLE(data []byte, notebook *Notebook, canvas *image.RGBA) {
	expectedLen := notebook.Device.PageWidth * notebook.Device.PageHeight
	decompressed := make([]byte, expectedLen)

	var prevColCode, currColCode, prevLenCode, currLenCode byte
	for i := 0; i < len(data)-1; {
		currColCode, currLenCode = data[0], data[1]
		currLen := int(((currLenCode & 0x7f) + 1) << 7)

		if prevColCode != 0 {
			prevLen := int(((prevLenCode & 0x7f) + 1) << 7)
			if prevColCode == data[0] {
				currLen = currLen + prevLen
			} else {
				decompressed = append(decompressed, processPair(prevColCode, prevLen)...)
			}
		} else if isLongRun(currLenCode) {
			currLen = BLANK_LINE_LENGTH
		} else if isMultiByte(currLenCode) {
			prevColCode = currColCode
			prevLenCode = currLenCode
			continue
		} else {
			currLen = int(data[1]&0x7f) + 1
		}

		decompressed = append(decompressed, processPair(currColCode, currLen)...)
		i += 2
	}

	decompressed = slices.Clip(decompressed)

	if len(decompressed) != expectedLen {
		log.Println("Length did not match")
	}

	converted := make([]color.Color, len(decompressed))
	for _, b := range decompressed {
		converted = append(converted, notebook.Device.ByteToColor(b))
	}

	fmt.Println(converted)
}

func processPair(colorCode byte, length int) []byte {
	d := make([]byte, length)
	for range length {
		d = append(d, colorCode)
	}
	return d
}

func isLongRun(l byte) bool { return l == 0xff }

func isMultiByte(l byte) bool { return l&0x8f != 0 }
