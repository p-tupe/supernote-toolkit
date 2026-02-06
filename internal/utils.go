package internal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"os"
	"path/filepath"
	"regexp"
)

// Ensure that file stream starts with 'note' byte,
// otherwise unsupported
func isNote(file *os.File) (bool, error) {
	if filepath.Ext(file.Name()) != ".note" {
		return false, errors.New("Unsupported file format")
	}

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
	pix := img.Pix
	off := 0

	fillRun := func(colorCode byte, length int) {
		col := notebook.Device.ToRGBA(colorCode)
		r, g, b, a := col.R, col.G, col.B, col.A
		end := min(off+length*4, len(pix))
		for off < end {
			pix[off] = r
			pix[off+1] = g
			pix[off+2] = b
			pix[off+3] = a
			off += 4
		}
	}

	i := 0
	holderActive := false
	var holderColor, holderLen byte

	for i+1 < len(data) {
		colorCode := data[i]
		lengthCode := data[i+1]
		i += 2

		var length int

		if holderActive {
			holderActive = false
			if colorCode == holderColor {
				length = 1 + int(lengthCode) + ((int(holderLen&0x7f) + 1) << 7)
			} else {
				fillRun(holderColor, (int(holderLen&0x7f)+1)<<7)
				length = int(lengthCode) + 1
			}
		} else if lengthCode == 0xff {
			length = 0x4000
		} else if lengthCode&0x80 != 0 {
			holderActive = true
			holderColor = colorCode
			holderLen = lengthCode
			continue
		} else {
			length = int(lengthCode) + 1
		}

		fillRun(colorCode, length)
	}

	if holderActive {
		remaining := (len(pix) - off) / 4
		tailLen := min((int(holderLen&0x7f)+1)<<7, remaining)
		if tailLen > 0 {
			fillRun(holderColor, tailLen)
		}
	}
}
