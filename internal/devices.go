package internal

import "image/color"

type Device struct {
	Name       string
	Model      string
	PageWidth  int
	PageHeight int
	CodeMap    map[byte]color.RGBA
}

func NewDevice(notebook *Notebook) {
	switch notebook.Header.APPLY_EQUIPMENT {
	case "N5":
		notebook.Device = A5X2
	default:
		notebook.Device = A5X2
	}
}
