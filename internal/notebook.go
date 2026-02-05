package internal

import (
	"errors"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

	notebook := &Notebook{Pages: make([]*Page, 0, 1)}

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

func (notebook *Notebook) ToPNG() {
	err := os.MkdirAll("./notebook", 0o755)
	if err != nil {
		log.Fatalln(err)
	}

	for i, p := range notebook.Pages {
		op, err := os.Create(filepath.Join("notebook", "page-"+strconv.Itoa(i)+".png"))
		if err != nil {
			log.Fatalln(err)
		}
		defer op.Close()

		bounds := image.Rect(0, 0, notebook.Device.PageWidth, notebook.Device.PageHeight)
		data := image.NewRGBA(bounds)

		for _, l := range p.LAYERSEQ {
			draw.Draw(data, bounds, l.Data, image.Point{}, draw.Over)
		}

		err = png.Encode(op, data)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
