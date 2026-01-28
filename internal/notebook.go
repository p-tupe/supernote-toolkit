package internal

import (
	"errors"
	"os"
)

type Notebook struct {
	device *Device
	footer *Footer
	header *Header
	pages  []*Page
}

var ErrUnsupported = errors.New("Unsupported file format")

func NewNotebook(file *os.File) (*Notebook, error) {
	ok, err := isNote(file)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrUnsupported
	}

	notebook := &Notebook{}

	notebook.footer, err = NewFooter(file)
	if err != nil {
		return nil, err
	}

	notebook.header, err = NewHeader(file, notebook.footer)
	if err != nil {
		return nil, err
	}

	notebook.device = NewDevice(notebook.header)

	for _, addr := range notebook.footer.PageAddr {
		page, err := NewPage(file, notebook.device, addr)
		if err != nil {
			return nil, err
		}
		notebook.pages = append(notebook.pages, page)
	}

	return notebook, nil
}
