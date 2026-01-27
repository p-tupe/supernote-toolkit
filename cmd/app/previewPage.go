package app

import (
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func GetPreviewPage(appData *AppData) *fyne.Container {
	filteredList := make([]fyne.URI, 0)
	l, err := appData.inputDir.List()

	if err != nil {
		dialog.NewError(err, appData.window).Show()
		return nil
	}

	for _, n := range l {
		if n.Extension() == ".note" {
			filteredList = append(filteredList, n)
		}
	}
	filteredList = slices.Clip(filteredList)

	notesList := widget.NewList(
		func() int {
			return len(filteredList)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("List Item")
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(filteredList[lii].Name())
		},
	)

	pdfFolder := widget.NewLabel("Output at: " + appData.outputDir.Path())

	pdfFolder.Importance = widget.MediumImportance
	pdfFolder.TextStyle.Bold = true

	convertBtn := widget.NewButton("Convert now!", func() {
		// todo
	})
	convertBtn.Importance = widget.HighImportance

	listLabel := widget.NewLabel("Selected Files: ")
	listLabel.TextStyle.Bold = true

	bottomBar := container.NewHBox(pdfFolder, widget.NewToolbarSpacer().ToolbarObject(), convertBtn)

	return container.NewBorder(
		listLabel,
		bottomBar,
		nil, nil,
		notesList,
	)
}
