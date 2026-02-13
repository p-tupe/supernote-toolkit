package internal

import "image/color"

type Device struct {
	Name       string
	Model      string
	PageWidth  int
	PageHeight int

	// Hex to RGBA map for device
	ToRGBA func(byte) color.RGBA

	// Orientation Value for Horizontal Pages
	HorizontalOrientation int
}

func NewDevice(notebook *Notebook) {
	switch notebook.Header.APPLY_EQUIPMENT {
	case "N5":
		notebook.Device = A5X2
	// case "N6":
	// 	notebook.Device = A6X2
	default:
		notebook.Device = A6X2
	}
}

// Tech Specs from https://supernote.com/products/supernote-manta
var A5X2 = &Device{
	Name:                  "Supernote Manta",
	Model:                 "A5X2",
	PageWidth:             1920,
	PageHeight:            2560,
	HorizontalOrientation: 1090,
	// VerticalOrientation:   1000,
	ToRGBA: X2CodeToRGBA,
}

// Tech Specs https://supernote.com/products/supernote-nomad
var A6X2 = &Device{
	Name:                  "Supernote Nomad",
	Model:                 "A6X2",
	PageWidth:             1404,
	PageHeight:            1872,
	HorizontalOrientation: 1270,
	// VerticalOrientation:   ???,
	ToRGBA: X2CodeToRGBA,
}

func X2CodeToRGBA(b byte) color.RGBA {
	switch b {
	case 0x61:
		return ColorBlack
	case 0x62:
		return ColorTransparent
	case 0x63:
		return ColorDarkGray
	case 0x64:
		return ColorLightGray
	case 0x65:
		return ColorWhite
	default:
		return color.RGBA{b, b, b, 255}
	}
}
