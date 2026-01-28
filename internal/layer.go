package internal

import (
	"fmt"
	"os"
)

type Layer struct {
	name          string
	addr          int64
	LAYERPROTOCOL string
	LAYERBITMAP   uint64
}

func NewLayer(file *os.File, layerAddr int64) (*Layer, error) {

	layerStr, err := readBlock(file, layerAddr)
	if err != nil {
		return nil, err
	}

	layer, err := parseLayerStr(layerStr)
	if err != nil {
		return nil, err
	}
	fmt.Println(layer)

	return layer, nil
}

func parseLayerStr(layerStr string) (*Layer, error) {
	layer := &Layer{}

	// TODO:

	return layer, nil
}
