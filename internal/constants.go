package internal

import "image/color"

var X2_CODE_TO_COLOR = map[byte]color.RGBA{
	0x61: {0, 0, 0, 255},          // Black
	0x62: {0, 0, 0, 0},            // Transparent
	0x63: {0x9d, 0x9d, 0x9d, 255}, // Dark Gray
	0x64: {0xc9, 0xc9, 0xc9, 255}, // Light Gray
	0x65: {255, 255, 255, 255},    // White
}

const BLANK_LINE_LENGTH = 0x4000 // 16384

var A5X2 = &Device{
	Name:       "Supernote Manta",
	Model:      "A5X2",
	PageWidth:  1920,
	PageHeight: 2560,
	CodeMap:    X2_CODE_TO_COLOR,
}
