package internal

type Device struct {
	name       string
	code       string
	pageWidth  int
	pageHeight int
}

// Creates a new [Device] from a [Header]
func NewDevice(header *Header) *Device {
	if header.APPLY_EQUIPMENT == "N5" {
		return A5X2()
	}

	return A5X2()
}

func A5X2() *Device {
	return &Device{
		name:       "Supernote Manta",
		code:       "A5X2",
		pageWidth:  1920,
		pageHeight: 2560,
	}
}
