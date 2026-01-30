package internal

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"
)

type Layer struct {
	LAYERADDR        int64
	LAYERTYPE        string
	LAYERPROTOCOL    string
	LAYERNAME        string
	LAYERPATH        string
	LAYERBITMAP      string
	LAYERVECTORGRAPH string
	LAYERRECOGN      string

	Data *image.RGBA
}

func NewLayer(file *os.File, notebook *Notebook, layerAddr int64) (*Layer, error) {
	layerStr, err := readBlock(file, layerAddr)
	if err != nil {
		return nil, err
	}

	metadata := parseMetadata(layerStr)

	layer := &Layer{
		LAYERADDR:        layerAddr,
		LAYERTYPE:        metadata["LAYERTYPE"],
		LAYERPROTOCOL:    metadata["LAYERPROTOCOL"],
		LAYERNAME:        metadata["LAYERNAME"],
		LAYERPATH:        metadata["LAYERPATH"],
		LAYERBITMAP:      metadata["LAYERBITMAP"],
		LAYERVECTORGRAPH: metadata["LAYERVECTORGRAPH"],
		LAYERRECOGN:      metadata["LAYERRECOGN"],
	}

	bounds := image.Rect(0, 0, notebook.Device.PageWidth, notebook.Device.PageHeight)
	layer.Data = image.NewRGBA(bounds)
	// TODO: White bg?

	switch layer.LAYERPROTOCOL {
	case "PNG":
		fmt.Println("Decoding PNG...")
		rawImg, err := readBlock(file, layerAddr)
		if err != nil {
			return nil, err
		}
		img, _, err := image.Decode(strings.NewReader(rawImg))
		if err != nil {
			return nil, err
		}
		draw.Draw(layer.Data, bounds, img, image.Pt(0, 0), draw.Over)

		// TODO: Remove test check
		file, err := os.Create("output.png")
		if err != nil {
			log.Fatalf("Error creating file: %v", err)
		}
		defer file.Close()
		if err := png.Encode(file, layer.Data); err != nil {
			log.Fatalf("Error encoding image: %v", err)
		}
	case "RATTA_RLE":
		fmt.Println("Decoding RLE...")
		encodedBytes, err := readBlockAsBytes(file, layerAddr)
		if err != nil {
			return nil, err
		}

		decodeRLE(encodedBytes, notebook, layer.Data)
	}

	return layer, nil
}
