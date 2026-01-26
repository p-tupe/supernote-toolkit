package main

import (
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var MIN_SIZE = fyne.NewSize(1024, 768)

func main() {
	a := app.NewWithID("supernote-toolkit-v0.1")
	w := a.NewWindow("Supernote Toolkit v0.1")

	var inputDir, outputDir fyne.ListableURI
	var statusTxt *widget.Label
	var selectNoteFolderBtn, selectPDFFolderBtn, ctaBtn *widget.Button
	var content, inputContainer, outputContainer, previewContainer, bottomContainer *fyne.Container

	// First Page (Select .note Files)
	inputDialog := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
		if err != nil {
			statusTxt.SetText(err.Error())
			return
		}

		if lu == nil {
			return
		}

		inputDir = lu
		content.Remove(inputContainer)
		content.Add(outputContainer)
	}, w)
	inputDialog.Resize(MIN_SIZE)

	selectNoteFolderBtn = widget.NewButton("Select folder of .note files", func() { inputDialog.Show() })
	selectNoteFolderBtn.Importance = widget.HighImportance

	inputContainer = container.NewCenter(selectNoteFolderBtn)

	// Second Page (Select .pdf Folder)
	outputDialog := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
		if err != nil {
			statusTxt.SetText(err.Error())
			return
		}

		if lu == nil {
			return
		}

		outputDir = lu

		// Third Page (Start Conversion)
		notesList := widget.NewLabel(getInputList(inputDir))
		pdfFolder := widget.NewLabel("Output folder:" + outputDir.Name())
		previewContainer = container.NewBorder(widget.NewLabel("Selected .note files:"), pdfFolder, nil, nil, notesList)
		ctaBtn.Enable()

		content.Remove(outputContainer)
		content.Add(previewContainer)
	}, w)
	outputDialog.Resize(MIN_SIZE)

	selectPDFFolderBtn = widget.NewButton("Select folder for .pdf files", func() { outputDialog.Show() })
	selectPDFFolderBtn.Importance = widget.HighImportance

	outputContainer = container.NewCenter(selectPDFFolderBtn)

	spacer := widget.NewToolbarSpacer().ToolbarObject()
	ctaBtn = widget.NewButton("Convert now", func() { log.Println("Conversion started...") })
	ctaBtn.Importance = widget.HighImportance
	ctaBtn.Disable()

	bottomContainer = container.NewPadded(container.NewHBox(spacer, ctaBtn))

	content = container.NewBorder(nil, bottomContainer, nil, nil, inputContainer)
	w.Resize(MIN_SIZE)

	w.SetContent(content)
	w.ShowAndRun()
}

func getInputList(list fyne.ListableURI) string {
	allnotes, err := list.List()
	if err != nil {
		return err.Error()
	}

	var listOfNotes strings.Builder
	for _, note := range allnotes {
		if note.Extension() == ".note" {
			listOfNotes.Write([]byte(note.Name() + "\n"))
		}
	}

	return listOfNotes.String()
}
