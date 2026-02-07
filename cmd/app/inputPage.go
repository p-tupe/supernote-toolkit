package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func GetInputPage(appData *AppData, cb func()) *fyne.Container {
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

		appData.inputDir = lu
		cb()
	}, appData.mainWindow)
	inputDialog.Resize(MIN_SIZE)

	selectNoteFolderBtn := widget.NewButton("Select input folder of .note files", func() { inputDialog.Show() })
	selectNoteFolderBtn.Importance = widget.HighImportance

	return container.NewCenter(container.NewVBox(selectNoteFolderBtn, noFilesTxt))
}
