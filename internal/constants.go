package internal

import "image/color"

const (
	ConvertPNG string = "Convert To PNG"
	ConvertPDF string = "Convert To PDF"
	ConvertTXT string = "Extract TXT"

	DeviceAuto  string = "Infer by note"
	DeviceManta string = "Supernote Manta A5X2"
	DeviceNomad string = "Supernote Nomad A6X2"
	DeviceA5X   string = "Supernote Manta A5X"
	DeviceA6X   string = "Supernote Nomad A6X"
)

var (
	ColorBlack       = color.RGBA{0x00, 0x00, 0x00, 0xff}
	ColorTransparent = color.RGBA{0x00, 0x00, 0x00, 0x00}
	ColorDarkGray    = color.RGBA{0x9d, 0x9d, 0x9d, 0xff}
	ColorLightGray   = color.RGBA{0xc9, 0xc9, 0xc9, 0xff}
	ColorWhite       = color.RGBA{0xff, 0xff, 0xff, 0xff}
)
