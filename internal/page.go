package internal

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Page struct {
	LAYERSEQ     []*Layer
	ORIENTATION  int
	RECOGNTEXT   int64
	IsHorizontal bool
	RealTimeText string
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
		case "ORIENTATION":
			var err error
			if page.ORIENTATION, err = strconv.Atoi(v); err != nil {
				return nil, err
			}

		case "RECOGNTEXT":
			var err error
			if page.RECOGNTEXT, err = strconv.ParseInt(v, 10, 64); err != nil {
				return nil, err
			}

		case "LAYERSEQ":
			layerSeq = strings.Split(v, ",")
			if len(layerSeq) < 1 {
				return nil, errors.New("Could not find any layers")
			}

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

	page.IsHorizontal = notebook.Device.HorizontalOrientation == page.ORIENTATION

	for _, l := range layerSeq {
		newLayer, err := NewLayer(file, notebook, layerAddr[l], page.IsHorizontal)
		if err != nil {
			return nil, err
		}

		if page.RECOGNTEXT != 0 {
			encodedTxt, err := readBlockAsBytes(file, (page.RECOGNTEXT))
			if err != nil {
				return nil, err
			}

			if err = decodeTXT(encodedTxt, page); err != nil {
				return nil, err
			}
		}

		page.LAYERSEQ = append(page.LAYERSEQ, newLayer)
	}

	return page, nil
}
