package internal

import (
	"encoding/binary"
	"os"
	"strconv"
)

type Footer struct {
	// Address of each page in the note,
	// The index corresponds to page number
	PAGES []int64

	// Address of header
	FILE_FEATURE int64
}

func NewFooter(file *os.File, notebook *Notebook) error {
	footerAddr, err := getFooterAddress(file)
	if err != nil {
		return err
	}

	footerStr, err := readBlock(file, footerAddr)
	if err != nil {
		return err
	}

	metadata := parseMetadata(footerStr)

	headerAddr, err := strconv.ParseInt(metadata["FILE_FEATURE"], 0, 64)
	if err != nil {
		return err
	}

	pageAddr := make([]int64, 0, len(metadata))
	pageNum := 1
	for {
		addr, ok := metadata["PAGE"+strconv.Itoa(pageNum)]
		if !ok {
			break
		}

		pgAddr, err := strconv.ParseInt(addr, 0, 64)
		if err != nil {
			return err
		}
		pageAddr = append(pageAddr, pgAddr)
		pageNum++
	}

	notebook.Footer = &Footer{pageAddr, headerAddr}

	return nil
}

// Gets the address of where "Footer" starts,
// It is represented by the last 4 bytes of the file;
// convert those bytes -> little-endian uint32 -> int64
func getFooterAddress(file *os.File) (int64, error) {
	footerAddr := make([]byte, 4)

	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	fileInfo.Size()

	_, err = file.ReadAt(footerAddr, fileInfo.Size()-4)
	if err != nil {
		return 0, err
	}

	addr := binary.LittleEndian.Uint32(footerAddr)
	return int64(addr), nil
}
