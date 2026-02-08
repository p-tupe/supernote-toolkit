package app

import (
	"path/filepath"
	"slices"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	i "github.com/p-tupe/supernote-toolkit/internal"
)

func GetPreviewPage(appData *AppData) *fyne.Container {
	allNotes, err := i.GetNotesFromDir(appData.inputDir, appData.recurse, "")
	if err != nil {
		dialog.NewError(err, appData.mainWindow).Show()
		return nil
	}

	notesList := widget.NewList(
		func() int {
			return len(allNotes)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("List Item")
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(filepath.Base(allNotes[lii].Path))
		},
	)

	pdfFolder := widget.NewLabel("Output at: " + appData.outputDir.Path())
	pdfFolder.Importance = widget.MediumImportance
	pdfFolder.TextStyle.Bold = true

	var convertBtn *widget.Button
	convertBtn = widget.NewButton("Convert now!", func() {
		convertBtn.Disable()
		convertBtn.SetText("Converting...")
		var wg sync.WaitGroup
		for _, input := range allNotes {
			wg.Add(1)
			go func() {
				defer wg.Done()
				notebook, err := i.NewNotebook(input)
				if err != nil {
					dialog.NewError(err, appData.mainWindow).Show()
				} else {
					op := filepath.Join(appData.outputDir.Path(), input.Parents)

					if slices.Contains(appData.convertTo, "Convert to PNG") {
						notebook.ToPNG(op)
					}

					if slices.Contains(appData.convertTo, "Convert to PDF") {
						notebook.ToPDF(op)
					}
				}
			}()
		}
		wg.Wait()
		dialog.NewInformation("Done!", "All .note files have been converted successfully!", appData.mainWindow).Show()
		convertBtn.Enable()
		convertBtn.SetText("Quit")
		convertBtn.OnTapped = func() { appData.app.Quit() }
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
