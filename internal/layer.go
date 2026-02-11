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
	LAYERADDR     int64
	LAYERPROTOCOL string
	LAYERNAME     string
	LAYERBITMAP   string
	// LAYERPATH        string
	// LAYERTYPE        string
	// LAYERVECTORGRAPH string
	// LAYERRECOGN      string

	Data *image.RGBA
}

func NewLayer(file *os.File, notebook *Notebook, layerAddr int64, isHorizontal bool) (*Layer, error) {
	layerStr, err := readBlock(file, layerAddr)
	if err != nil {
		return nil, err
	}

	metadata := parseMetadata(layerStr)

	layer := &Layer{
		LAYERADDR:     layerAddr,
		LAYERPROTOCOL: metadata["LAYERPROTOCOL"],
		LAYERNAME:     metadata["LAYERNAME"],
		LAYERBITMAP:   metadata["LAYERBITMAP"],
		// LAYERTYPE:        metadata["LAYERTYPE"],
		// LAYERPATH:        metadata["LAYERPATH"],
		// LAYERVECTORGRAPH: metadata["LAYERVECTORGRAPH"],
		// LAYERRECOGN:      metadata["LAYERRECOGN"],
	}

	var bounds image.Rectangle
	if isHorizontal {
		bounds = image.Rect(0, 0, notebook.Device.PageHeight, notebook.Device.PageWidth)
	} else {
		bounds = image.Rect(0, 0, notebook.Device.PageWidth, notebook.Device.PageHeight)
	}

	layer.Data = image.NewRGBA(bounds)

	switch layer.LAYERPROTOCOL {
	case "TEXT":
	// TODO: Real-time text comes here

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

		decodeRLE(encodedBytes, notebook, layer.Data)

	default:
		fmt.Printf("Unknown layer protocol: %v\n", layer.LAYERPROTOCOL)
	}

	return layer, nil
}
