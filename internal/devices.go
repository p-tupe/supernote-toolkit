package internal

import "image/color"

type Device struct {
	Name       string
	Model      string
	PageWidth  int
	PageHeight int
	ToRGBA     func(byte) color.RGBA
}

func NewDevice(notebook *Notebook) {
	switch notebook.Header.APPLY_EQUIPMENT {
	case "N5":
		notebook.Device = A5X2
	default:
		notebook.Device = A5X2
	}
}

var A5X2 = &Device{
	Name:       "Supernote Manta",
	Model:      "A5X2",
	PageWidth:  1920,
	PageHeight: 2560,
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
