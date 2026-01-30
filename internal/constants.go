package internal

import "image/color"

var X2_CODE_TO_COLOR = map[byte]color.Color{
	0x61: color.Black,
	0x62: color.Transparent,
	0x63: color.RGBA{0x9d, 0x9d, 0x9d, 255},
	0x64: color.RGBA{0xc9, 0xc9, 0xc9, 255},
	0x65: color.White,
}

const BLANK_LINE_LENGTH = 0x4000

var A5X2 = &Device{
	Name:       "Supernote Manta",
	Model:      "A5X2",
	PageWidth:  1920,
	PageHeight: 2560,
	CodeMap:    X2_CODE_TO_COLOR,
}

// var A5X1 = &Device{
// 	Name:       "Supernote Manta",
// 	Model:      "A5X1",
// 	PageWidth:  1920,
// 	PageHeight: 2560,
// 	CodeMap:    X1_CODES,
// }

// var A6X2 = &Device{
// 	Name:       "Supernote Nomad",
// 	Model:      "A6X2",
// 	PageWidth:  1920,
// 	PageHeight: 2560,
// 	CodeMap:    X2_CODES,
// }
