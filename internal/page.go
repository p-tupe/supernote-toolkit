package internal

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Page struct {
	LAYERSEQ []*Layer
}

func NewPage(file *os.File, notebook *Notebook, pageAddr int64) (*Page, error) {
	pageStr, err := readBlock(file, pageAddr)
	if err != nil {
		return nil, err
	}

	return parsePageStr(file, notebook, pageStr)
}

func parsePageStr(file *os.File, notebook *Notebook, pageStr string) (*Page, error) {
	page := &Page{}

	matches := parseMetadata(pageStr)

	layerSeq := []string{}
	layerAddr := map[string]int64{}
	for k, v := range matches {
		switch k {
		case "LAYERSEQ":
			layerSeq = strings.Split(v, ",")
		case "BGLAYER", "MAINLAYER", "LAYER1", "LAYER2", "LAYER3":
			val, err := strconv.ParseInt(v, 0, 64)
			if err != nil {
				return nil, err
			}
			if val > 0 {
				layerAddr[k] = val
			}
		}
	}

	if len(layerSeq) < 1 {
		return nil, errors.New("Could not find any layers")
	}

	for _, l := range layerSeq {
		newLayer, err := NewLayer(file, notebook, layerAddr[l])
		if err != nil {
			return nil, err
		}
		page.LAYERSEQ = append(page.LAYERSEQ, newLayer)
	}

	return page, nil
}
