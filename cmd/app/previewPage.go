package app

import (
	"errors"
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
	allNotes, err := i.GetNotesFromDir(appData.input, appData.recurse, "")
	if err != nil {
		dialog.NewError(err, appData.mainWindow).Show()
		return nil
	}

	if !appData.force {
		allNotes, err = i.FilterFreshNotes(allNotes, appData.output.Path(), appData.convertTo)
		if err != nil {
			dialog.NewError(err, appData.mainWindow).Show()
			return nil
		}
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

	pdfFolder := widget.NewLabel("Output at: " + appData.output.Path())
	pdfFolder.Importance = widget.MediumImportance
	pdfFolder.TextStyle.Bold = true

	var convertBtn *widget.Button
	convertBtn = widget.NewButton("Convert now!", func() {
		convertBtn.Disable()
		convertBtn.SetText("Converting...")

		var device *i.Device
		switch appData.device {
		case i.DeviceManta:
			device = i.A5X2
		case i.DeviceNomad:
			device = i.A6X2
		case i.DeviceA5X:
			device = i.A5X
		case i.DeviceA6X:
			device = i.A6X
		}

		var wg sync.WaitGroup
		for _, input := range allNotes {
			wg.Add(1)
			go func() {
				defer wg.Done()
				notebook, err := i.NewNotebook(input, device)
				if err != nil {
					dialog.NewError(err, appData.mainWindow).Show()
				} else {
					op := filepath.Join(appData.output.Path(), input.Parents)

					if slices.Contains(appData.convertTo, i.ConvertPNG) {
						notebook.ToPNG(op)
					}

					if slices.Contains(appData.convertTo, i.ConvertPDF) {
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

	if len(allNotes) == 0 {
		listLabel.SetText("No .note files to convert!")
		dialog.NewError(errors.New("No .note files to convert!"), appData.mainWindow).Show()
		convertBtn.SetText("Quit")
		convertBtn.OnTapped = func() { appData.app.Quit() }
	}

	return container.NewBorder(
		listLabel,
		bottomBar,
		nil, nil,
		notesList,
	)
}
