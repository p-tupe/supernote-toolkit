package internal

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Page struct {
	LAYERSEQ []*Layer
}

func NewPage(file *os.File, device *Device, pageAddr int64) (*Page, error) {
	pageStr, err := readBlock(file, pageAddr)
	if err != nil {
		return nil, err
	}

	page, err := parsePageStr(file, pageStr)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func parsePageStr(file *os.File, pageStr string) (*Page, error) {
	page := &Page{}

	layerRegex, err := regexp.Compile(`<(\w*?LAYER\w*?):(.*?)>`)
	if err != nil {
		return nil, err
	}

	matches := layerRegex.FindAllStringSubmatch(pageStr, -1)

	layerSeq := []string{}
	layerAddr := map[string]int64{}
	for _, m := range matches {
		switch m[1] {
		case "LAYERSEQ":
			layerSeq = strings.Split(m[2], ",")
		case "BGLAYER", "MAINLAYER", "LAYER1", "LAYER2", "LAYER3":
			v, err := strconv.ParseInt(m[2], 10, 64)
			if err != nil {
				return nil, err
			}
			if v > 0 {
				layerAddr[m[1]] = v
			}
		}
	}

	if len(layerSeq) < 1 {
		return nil, errors.New("Could not find any layers")
	}

	for _, l := range layerSeq {
		newLayer, err := NewLayer(file, layerAddr[l])
		if err != nil {
			return nil, err
		}
		page.LAYERSEQ = append(page.LAYERSEQ, newLayer)
	}

	return page, nil
}
