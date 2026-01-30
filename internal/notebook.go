package internal

import (
	"errors"
	"os"
)

type Notebook struct {
	Device *Device
	Footer *Footer
	Header *Header
	Pages  []*Page
}

var ErrUnsupported = errors.New("Unsupported file format")

func NewNotebook(file *os.File) (*Notebook, error) {
	if ok, err := isNote(file); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrUnsupported
	}

	notebook := &Notebook{Pages: make([]*Page, 1)}

	if err := NewFooter(file, notebook); err != nil {
		return nil, err
	}

	if err := NewHeader(file, notebook); err != nil {
		return nil, err
	}

	NewDevice(notebook)

	for _, addr := range notebook.Footer.PAGES {
		if err := NewPage(file, notebook, addr); err != nil {
			return nil, err
		}
	}

	return notebook, nil
}
