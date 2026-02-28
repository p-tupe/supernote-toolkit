package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	i "github.com/p-tupe/supernote-toolkit/internal"
)

var convertToOptions = []string{i.ConvertPDF, i.ConvertPNG, i.ExtractTXT}

func GetOutputPage(appData *AppData, cb func()) *fyne.Container {
	convertCheckbox := widget.NewCheckGroup(convertToOptions, func(s []string) {
		appData.convertTo = s
	})
	convertCheckbox.Horizontal = true
	convertCheckbox.Required = true
	convertCheckbox.Selected = []string{convertToOptions[0]}

	forceCheckbox := widget.NewCheck("Force convert stale notes", func(b bool) {
		appData.force = b
	})
	forceCheckbox.SetChecked(false)

	bottomOptions := container.NewHBox(convertCheckbox, forceCheckbox)

	outputDialog := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
		if err != nil {
			dialog.NewError(err, appData.mainWindow).Show()
			return
		}

		if lu == nil {
			return
		}

		appData.output = lu
		cb()
	}, appData.mainWindow)
	outputDialog.Resize(MIN_SIZE)

	selectPDFFolderBtn := widget.NewButton("Select folder for output files", func() { outputDialog.Show() })
	selectPDFFolderBtn.Importance = widget.HighImportance

	return container.NewBorder(nil, bottomOptions, nil, nil, container.NewCenter(selectPDFFolderBtn))
}
