package internal

import (
	"os"
	"regexp"
)

type Header struct {
	APPLY_EQUIPMENT  string
	FILE_RECOGN_TYPE string
	HORIZONTAL_CHECK string
}

// Creates a new [Header] from [*os.File] and a [Footer]
func NewHeader(file *os.File, footer *Footer) (*Header, error) {
	headerStr, err := readBlock(file, footer.HeaderAddr)
	if err != nil {
		return nil, err
	}

	header, err := parseHeaderStr(headerStr)
	if err != nil {
		return nil, err
	}

	return header, nil
}

// parseHeaderStr takes in a string with values of
// form <key1:value1><key2:value2>... and returns
// the them as a Header struct.
func parseHeaderStr(headerStr string) (*Header, error) {
	header := &Header{}

	headerRegex, err := regexp.Compile(`<(\w+):(\w+)>`)
	if err != nil {
		return nil, err
	}
	matches := headerRegex.FindAllStringSubmatch(headerStr, -1)

	for _, m := range matches {
		key, value := m[1], m[2]

		switch key {
		case "APPLY_EQUIPMENT":
			header.APPLY_EQUIPMENT = value
		case "FILE_RECOGN_TYPE":
			header.FILE_RECOGN_TYPE = value
		case "HORIZONTAL_CHECK":
			header.HORIZONTAL_CHECK = value
		}
	}

	return header, nil
}
