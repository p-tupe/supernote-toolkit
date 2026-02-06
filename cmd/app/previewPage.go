package app

import (
	"log"
	"slices"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	i "github.com/p-tupe/supernote-toolkit/internal"
)

func GetPreviewPage(appData *AppData) *fyne.Container {
	filteredList := make([]fyne.URI, 0)
	l, err := appData.inputDir.List()

	if err != nil {
		dialog.NewError(err, appData.mainWindow).Show()
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
		fyne.DoAndWait(func() {
			var wg sync.WaitGroup
			for _, input := range filteredList {
				wg.Go(func() {
					notebook, err := i.NewNotebook(input.Path())
					if err != nil {
						dialog.NewError(err, appData.mainWindow).Show()
					} else {
						if slices.Contains(appData.convertTo, "Convert to PNG") {
							notebook.ToPNG(appData.outputDir.Path())
						}

						if slices.Contains(appData.convertTo, "Convert to PDF") {
							log.Println("TODO: Converting ", notebook.Name, " to PDF")
						}
					}
				})
			}
			wg.Wait()
			dialog.NewInformation("Done!", "All .note files have been converted successfully!", appData.mainWindow).Show()
		})
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
