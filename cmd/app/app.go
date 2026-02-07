package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var (
	Version  = "dev"
	MIN_SIZE = fyne.NewSize(1024, 768)
)

type AppData struct {
	app        fyne.App
	mainWindow fyne.Window

	inputDir  fyne.ListableURI
	outputDir fyne.ListableURI
	convertTo []string
}

func Execute() {
	a := app.NewWithID("supernote-toolkit-" + Version)
	w := a.NewWindow("Supernote Toolkit " + Version)

	var content, inputPage, outputPage, previewPage *fyne.Container

	appData := &AppData{app: a, mainWindow: w, convertTo: convertToOptions}

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
