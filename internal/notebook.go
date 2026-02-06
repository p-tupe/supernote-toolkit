package internal

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type Notebook struct {
	Device *Device
	Footer *Footer
	Header *Header
	Pages  []*Page
	Name   string
}

var ErrUnsupported = errors.New("Unsupported file format")

func NewNotebook(input string) (*Notebook, error) {
	file, err := os.Open(input)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()


	if ok, err := isNote(file); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrUnsupported
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	notebook := &Notebook{Pages: make([]*Page, 0, 1), Name: fileInfo.Name()}

	if err := NewFooter(file, notebook); err != nil {
		return nil, err
	}

	if err := NewHeader(file, notebook); err != nil {
		return nil, err
	}

	NewDevice(notebook)

	var wg sync.WaitGroup
	for _, addr := range notebook.Footer.PAGES {
		wg.Go(func() {
			if err := NewPage(file, notebook, addr); err != nil {
				log.Println(err)
			}
		})
	}

	wg.Wait()

	return notebook, nil
}

func (notebook *Notebook) ToPNG() {
	opDir := filepath.Join("output", notebook.Name)
	err := os.MkdirAll(opDir, 0o755)
	if err != nil {
		log.Fatalln(err)
	}

	for i, p := range notebook.Pages {
		op, err := os.Create(filepath.Join(opDir, "page-"+strconv.Itoa(i)+".png"))
		if err != nil {
			log.Fatalln(err)
		}
		defer op.Close()

		bounds := image.Rect(0, 0, notebook.Device.PageWidth, notebook.Device.PageHeight)
		data := image.NewRGBA(bounds)
		draw.Draw(data, bounds, &image.Uniform{color.White}, image.Point{}, draw.Src)

		for _, l := range p.LAYERSEQ {
			draw.Draw(data, bounds, l.Data, image.Point{}, draw.Over)
		}

		err = png.Encode(op, data)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
