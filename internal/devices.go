package internal

type Device struct {
	Name       string
	Code       string
	PageWidth  int
	PageHeight int
}

func NewDevice(notebook *Notebook) {
	if notebook.Header.APPLY_EQUIPMENT == "N5" {
		notebook.Device = A5X2()
	}

	notebook.Device = A5X2()
}

func A5X2() *Device {
	return &Device{
		Name:       "Supernote Manta",
		Code:       "A5X2",
		PageWidth:  1920,
		PageHeight: 2560,
	}
}
