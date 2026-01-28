package internal

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Footer struct {
	// Address of each page in the note,
	// The index of address corresponds to page number
	PageAddr []int64

	// Address of header
	HeaderAddr int64
}

// Creates a new [Footer] from [*os.File] of type .note
func NewFooter(file *os.File) (*Footer, error) {
	footerAddr, err := getFooterAddress(file)
	if err != nil {
		return nil, err
	}

	footerStr, err := getFooterStr(file, footerAddr)
	if err != nil {
		return nil, err
	}

	footer, err := parseFooterStr(footerStr)
	if err != nil {
		return nil, err
	}

	return footer, nil
}

// Gets the address of where "Footer" starts,
// It is represented by the last 4 bytes of the file;
// convert those bytes -> little-endian uint32 -> int64
func getFooterAddress(file *os.File) (int64, error) {
	end := make([]byte, 4)
	file.Seek(-4, 2)
	_, err := file.Read(end)
	if err != nil {
		return 0, err
	}

	addr := binary.LittleEndian.Uint32(end)
	return int64(addr), nil
}

// Once we have the footer address,
// simply pull all bytes from then until EOF,
// convert 'em into string and voila!
func getFooterStr(file *os.File, footerAddr int64) (string, error) {
	_, err := file.Seek(footerAddr, 0)
	if err != nil {
		return "", err
	}

	var footerStr strings.Builder
	buff := make([]byte, 256)
	for {
		_, err := file.Read(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return "", err
			}
		}
		footerStr.Write(buff)
	}

	return footerStr.String(), nil

}

// Footer string is usually represented in format of:
// <Key1:value1><key2:value2>... and so on.
//
// What we're concerned about is the <PAGE1:ABC><PAGE2:DEF>...
// That represent address of each page in the note;
// And <FILE_FEATURE:xyz> that represents the Header address.
func parseFooterStr(footerStr string) (*Footer, error) {
	headerRegex, err := regexp.Compile(`<FILE_FEATURE:(\d+)>`)
	if err != nil {
		return nil, err
	}
	match := headerRegex.FindStringSubmatch(footerStr)
	if len(match) < 2 {
		return nil, errors.New("Header address not found")
	}
	headerAddr, err := strconv.ParseInt(match[1], 0, 64)
	if err != nil {
		return nil, err
	}

	pageRegex, err := regexp.Compile(`<PAGE(\d+):(\w+)>`)
	if err != nil {
		return nil, err
	}
	matches := pageRegex.FindAllStringSubmatch(footerStr, -1)
	if len(matches) < 1 {
		return nil, errors.New("Page addresses not found")
	}

	pageAddr := make([]int64, 0, len(matches))
	for _, m := range matches {
		pgAddr, err := strconv.ParseInt(m[2], 0, 64)
		if err != nil {
			return nil, err
		}
		pageAddr = append(pageAddr, pgAddr)
	}

	return &Footer{pageAddr, headerAddr}, nil
}
