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
	case "N6":
		notebook.Device = A6X2
	case "A5X":
		notebook.Device = A5X
	case "A6X":
		notebook.Device = A6X
	default:
		notebook.Device = A6X2
	}
}

// Tech Specs from https://supernote.com/products/supernote-manta
var A5X2 = &Device{
	Name:                  "Supernote Manta X2",
	Model:                 "A5X2",
	PageWidth:             1920,
	PageHeight:            2560,
	HorizontalOrientation: 1090,
	ToRGBA:                CodeToRGBA,
}

// Tech Specs from https://supernote.com/products/supernote-nomad
var A6X2 = &Device{
	Name:                  "Supernote Nomad X2",
	Model:                 "A6X2",
	PageWidth:             1404,
	PageHeight:            1872,
	HorizontalOrientation: 1270,
	ToRGBA:                CodeToRGBA,
}

var A5X = &Device{
	Name:                  "Supernote Manta X",
	Model:                 "A5X",
	PageWidth:             1404,
	PageHeight:            1872,
	HorizontalOrientation: 1270, // Unsure
	ToRGBA:                CodeToRGBA,
}

var A6X = &Device{
	Name:                  "Supernote Nomad X",
	Model:                 "A6X",
	PageWidth:             1404,
	PageHeight:            1872,
	HorizontalOrientation: 1270, // Unsure
	ToRGBA:                CodeToRGBA,
}

func CodeToRGBA(b byte) color.RGBA {
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
