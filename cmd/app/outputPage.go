package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func GetOutputPage(appData *AppData, cb func()) *fyne.Container {
	outputDialog := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
		if err != nil {
			dialog.NewError(err, appData.window).Show()
			return
		}

		if lu == nil {
			return
		}

		appData.outputDir = lu
		cb()
	}, appData.window)
	outputDialog.Resize(MIN_SIZE)

	selectPDFFolderBtn := widget.NewButton("Select folder for .pdf files", func() { outputDialog.Show() })
	selectPDFFolderBtn.Importance = widget.HighImportance

	return container.NewCenter(selectPDFFolderBtn)
}
