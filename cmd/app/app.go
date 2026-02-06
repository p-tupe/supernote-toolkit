package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var MIN_SIZE = fyne.NewSize(1024, 768)

type AppData struct {
	app        fyne.App
	mainWindow fyne.Window

	inputDir  fyne.ListableURI
	outputDir fyne.ListableURI
	convertTo []string
}

func Execute() {
	a := app.NewWithID("supernote-toolkit-v0.1")
	w := a.NewWindow("Supernote Toolkit v0.1")

	var content, inputPage, outputPage, previewPage *fyne.Container

	appData := &AppData{app: a, mainWindow: w}

	inputPage = GetInputPage(appData, func() {
		content.Remove(inputPage)
		content.Add(outputPage)
	})

	outputPage = GetOutputPage(appData, func() {
		previewPage = GetPreviewPage(appData)
		content.Remove(outputPage)
		content.Add(previewPage)
	})

	content = container.NewPadded(inputPage)
	w.Resize(MIN_SIZE)
	w.SetContent(content)
	w.ShowAndRun()
}
