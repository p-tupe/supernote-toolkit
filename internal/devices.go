package internal

import "image/color"

type Device struct {
	Name       string
	Model      string
	PageWidth  int
	PageHeight int
	CodeMap    map[byte]color.Color
}

func NewDevice(notebook *Notebook) {
	switch notebook.Header.APPLY_EQUIPMENT {
	case "N5":
		notebook.Device = A5X2
	default:
		notebook.Device = A5X2
	}
}

func (d *Device) ByteToColor(b byte) color.Color {
	if b == 0 {
		return color.Transparent
	}

	c, ok := d.CodeMap[b]
	if !ok {
		return color.RGBA{b, b, b, 255}
	}

	return c
}
