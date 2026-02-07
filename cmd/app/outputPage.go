package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var convertToOptions = []string{"Convert to PNG", "Convert to PDF"}

func GetOutputPage(appData *AppData, cb func()) *fyne.Container {
	formatCheckbox := widget.NewCheckGroup(convertToOptions, func(s []string) {
		appData.convertTo = s
	})
	formatCheckbox.Horizontal = true
	formatCheckbox.Required = true
	formatCheckbox.Selected = convertToOptions

	outputDialog := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
		if err != nil {
			dialog.NewError(err, appData.mainWindow).Show()
			return
		}

		if lu == nil {
			return
		}

		appData.outputDir = lu
		cb()
	}, appData.mainWindow)
	outputDialog.Resize(MIN_SIZE)

	selectPDFFolderBtn := widget.NewButton("Select folder for output files", func() { outputDialog.Show() })
	selectPDFFolderBtn.Importance = widget.HighImportance

	return container.NewBorder(nil, formatCheckbox, nil, nil, container.NewCenter(selectPDFFolderBtn))
}
