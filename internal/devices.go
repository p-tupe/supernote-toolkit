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
	ToRGBA: func(b byte) color.RGBA {
		switch b {
		case 0x61: // Black
			return color.RGBA{0, 0, 0, 255}
		case 0x62: // Transparent
			return color.RGBA{0, 0, 0, 0}
		case 0x63: // Dark Gray
			return color.RGBA{0x9d, 0x9d, 0x9d, 255}
		case 0x64: // Light Gray
			return color.RGBA{0xc9, 0xc9, 0xc9, 255}
		case 0x65: // White
			return color.RGBA{255, 255, 255, 255}
		default: // Intensity
			return color.RGBA{b, b, b, 255}
		}
	},
}

// Tech Specs https://supernote.com/products/supernote-nomad
var A6X2 = &Device{
	Name:                  "Supernote Nomad",
	Model:                 "A6X2",
	PageWidth:             1404,
	PageHeight:            1872,
	HorizontalOrientation: 1270,
	ToRGBA: func(b byte) color.RGBA {
		switch b {
		case 0x61: // Black
			return color.RGBA{0, 0, 0, 255}
		case 0x62: // Transparent
			return color.RGBA{0, 0, 0, 0}
		case 0x63: // Dark Gray
			return color.RGBA{0x9d, 0x9d, 0x9d, 255}
		case 0x64: // Light Gray
			return color.RGBA{0xc9, 0xc9, 0xc9, 255}
		case 0x65: // White
			return color.RGBA{255, 255, 255, 255}
		default: // Intensity
			return color.RGBA{b, b, b, 255}
		}
	},
}
