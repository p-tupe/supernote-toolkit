package internal

import (
	"os"
)

type Header struct {
	APPLY_EQUIPMENT  string
	FILE_RECOGN_TYPE string
	HORIZONTAL_CHECK string
}

func NewHeader(file *os.File, notebook *Notebook) error {
	headerStr, err := readBlock(file, notebook.Footer.FILE_FEATURE)
	if err != nil {
		return err
	}

	metadata := parseMetadata(headerStr)

	notebook.Header = &Header{
		APPLY_EQUIPMENT:  metadata["APPLY_EQUIPMENT"],
		FILE_RECOGN_TYPE: metadata["FILE_RECOGN_TYPE"],
		HORIZONTAL_CHECK: metadata["HORIZONTAL_CHECK"],
	}

	return nil
}
