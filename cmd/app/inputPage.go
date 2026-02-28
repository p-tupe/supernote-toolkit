package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	i "github.com/p-tupe/supernote-toolkit/internal"
)

var deviceOptions = []string{i.DeviceAuto, i.DeviceManta, i.DeviceNomad, i.DeviceA5X, i.DeviceA6X}

func GetInputPage(appData *AppData, cb func()) *fyne.Container {
	recurseCheckbox := widget.NewCheck("Recurse folders", func(b bool) {
		appData.recurse = b
	})
	recurseCheckbox.SetChecked(true)

	deviceLabel := widget.NewLabel("Select Device:")
	deviceRadio := widget.NewRadioGroup(deviceOptions, func(s string) {
		appData.device = s
	})
	deviceRadio.SetSelected(deviceOptions[0])

	deviceContainer := container.NewVBox(deviceLabel, deviceRadio)

	bottomOptions := container.NewHBox(recurseCheckbox, widget.NewToolbarSpacer().ToolbarObject(), deviceContainer)

	noFilesTxt := widget.NewLabel("You won't actually see the .note files; that's expected.")
	noFilesTxt.Importance = widget.LowImportance

	inputDialog := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
		if err != nil {
			dialog.NewError(err, appData.mainWindow).Show()
			return
		}

		if lu == nil {
			return
		}

		appData.input = lu.Path()
		cb()
	}, appData.mainWindow)
	inputDialog.Resize(MIN_SIZE)

	selectNoteFolderBtn := widget.NewButton("Select input folder of .note files", func() { inputDialog.Show() })
	selectNoteFolderBtn.Importance = widget.HighImportance

	return container.NewBorder(nil, bottomOptions, nil, nil, container.NewCenter(container.NewVBox(selectNoteFolderBtn, noFilesTxt)))
}
