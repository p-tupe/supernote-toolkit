package internal

import (
	"bytes"
	"encoding/binary"
	"os"
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
	_, err := file.Seek(startAddr, 0)
	if err != nil {
		return "", err
	}

	block := make([]byte, 4)
	_, err = file.Read(block)
	if err != nil {
		return "", err
	}
	lenAddr := binary.LittleEndian.Uint32(block)

	headerBytes := make([]byte, lenAddr)
	_, err = file.Read(headerBytes)
	if err != nil {
		return "", err
	}

	return string(headerBytes), nil
}
