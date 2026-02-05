package internal

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"strconv"
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

	switch layer.LAYERPROTOCOL {
	case "TEXT":

	case "PNG":
		rawImg, err := readBlock(file, layerAddr)
		if err != nil {
			return nil, err
		}

		img, _, err := image.Decode(strings.NewReader(rawImg))
		if err != nil {
			return nil, err
		}

		draw.Draw(layer.Data, bounds, img, image.Point{}, draw.Over)

	case "RATTA_RLE":
		bitmapAddr, err := strconv.ParseInt(layer.LAYERBITMAP, 0, 64)
		if err != nil {
			return nil, err
		}

		encodedBytes, err := readBlockAsBytes(file, bitmapAddr)
		if err != nil {
			return nil, err
		}

		decoded := decodeRLE(encodedBytes, notebook, bounds)
		draw.Draw(layer.Data, bounds, decoded, image.Point{}, draw.Over)

	default:
		fmt.Printf("Unknown layer protocol: %v\n", layer.LAYERPROTOCOL)
	}

	return layer, nil
}
