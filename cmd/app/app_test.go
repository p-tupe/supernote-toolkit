package app

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	i "github.com/p-tupe/supernote-toolkit/internal"
)

func collectObjects(obj fyne.CanvasObject) []fyne.CanvasObject {
	var all []fyne.CanvasObject
	if obj == nil {
		return all
	}
	all = append(all, obj)
	if c, ok := obj.(*fyne.Container); ok {
		for _, child := range c.Objects {
			all = append(all, collectObjects(child)...)
		}
	}
	return all
}

func findCheck(objects []fyne.CanvasObject, text string) *widget.Check {
	for _, obj := range objects {
		if c, ok := obj.(*widget.Check); ok && c.Text == text {
			return c
		}
	}
	return nil
}

func findButton(objects []fyne.CanvasObject, text string) *widget.Button {
	for _, obj := range objects {
		if b, ok := obj.(*widget.Button); ok && b.Text == text {
			return b
		}
	}
	return nil
}

func findRadioGroup(objects []fyne.CanvasObject) *widget.RadioGroup {
	for _, obj := range objects {
		if r, ok := obj.(*widget.RadioGroup); ok {
			return r
		}
	}
	return nil
}

func findCheckGroup(objects []fyne.CanvasObject) *widget.CheckGroup {
	for _, obj := range objects {
		if cg, ok := obj.(*widget.CheckGroup); ok {
			return cg
		}
	}
	return nil
}

func findLabel(objects []fyne.CanvasObject, contains string) *widget.Label {
	for _, obj := range objects {
		if l, ok := obj.(*widget.Label); ok && strings.Contains(l.Text, contains) {
			return l
		}
	}
	return nil
}

func newTestAppData(t *testing.T) *AppData {
	t.Helper()
	a := test.NewApp()
	w := a.NewWindow("Test")
	return &AppData{
		app:        a,
		mainWindow: w,
		force:      false,
		recurse:    true,
		convertTo:  []string{convertToOptions[0]},
		device:     deviceOptions[0],
	}
}

// --- AppData defaults ---

func TestAppDataDefaults(t *testing.T) {
	ad := newTestAppData(t)

	if ad.force {
		t.Error("force should default to false")
	}
	if !ad.recurse {
		t.Error("recurse should default to true")
	}
	if ad.device != i.DeviceAuto {
		t.Errorf("device = %q, want %q", ad.device, i.DeviceAuto)
	}
	if len(ad.convertTo) != 1 || ad.convertTo[0] != i.ConvertPDF {
		t.Errorf("convertTo = %v, want [%s]", ad.convertTo, i.ConvertPDF)
	}
	if ad.input != "" {
		t.Errorf("input = %q, want empty", ad.input)
	}
	if ad.output != nil {
		t.Error("output should be nil")
	}
}

// --- Input Page ---

func TestInputPage(t *testing.T) {
	t.Run("creation", func(t *testing.T) {
		ad := newTestAppData(t)
		page := GetInputPage(ad, func() {})
		if page == nil {
			t.Fatal("page is nil")
		}
	})

	t.Run("recurse checked by default", func(t *testing.T) {
		ad := newTestAppData(t)
		page := GetInputPage(ad, func() {})
		objs := collectObjects(page)
		c := findCheck(objs, "Recurse folders")
		if c == nil {
			t.Fatal("recurse checkbox not found")
		}
		if !c.Checked {
			t.Error("recurse should be checked by default")
		}
	})

	t.Run("device radio default", func(t *testing.T) {
		ad := newTestAppData(t)
		page := GetInputPage(ad, func() {})
		objs := collectObjects(page)
		r := findRadioGroup(objs)
		if r == nil {
			t.Fatal("device radio not found")
		}
		if r.Selected != i.DeviceAuto {
			t.Errorf("selected = %q, want %q", r.Selected, i.DeviceAuto)
		}
		if !slices.Equal(r.Options, deviceOptions) {
			t.Errorf("options = %v, want %v", r.Options, deviceOptions)
		}
	})

	t.Run("recurse toggle updates appData", func(t *testing.T) {
		ad := newTestAppData(t)
		page := GetInputPage(ad, func() {})
		c := findCheck(collectObjects(page), "Recurse folders")
		if c == nil {
			t.Fatal("recurse checkbox not found")
		}

		c.SetChecked(false)
		if ad.recurse {
			t.Error("expected recurse=false after uncheck")
		}
		c.SetChecked(true)
		if !ad.recurse {
			t.Error("expected recurse=true after recheck")
		}
	})

	t.Run("device selection updates appData", func(t *testing.T) {
		ad := newTestAppData(t)
		page := GetInputPage(ad, func() {})
		r := findRadioGroup(collectObjects(page))
		if r == nil {
			t.Fatal("device radio not found")
		}

		for _, opt := range []string{i.DeviceManta, i.DeviceNomad, i.DeviceAuto} {
			r.SetSelected(opt)
			if ad.device != opt {
				t.Errorf("after selecting %q: device = %q", opt, ad.device)
			}
		}
	})

	t.Run("select button", func(t *testing.T) {
		ad := newTestAppData(t)
		page := GetInputPage(ad, func() {})
		b := findButton(collectObjects(page), "Select input folder of .note files")
		if b == nil {
			t.Fatal("select button not found")
		}
		if b.Importance != widget.HighImportance {
			t.Errorf("importance = %v, want HighImportance", b.Importance)
		}
	})

	t.Run("callback not fired on creation", func(t *testing.T) {
		ad := newTestAppData(t)
		called := false
		page := GetInputPage(ad, func() { called = true })
		if page == nil {
			t.Fatal("page is nil")
		}
		if called {
			t.Error("callback should not fire on page creation")
		}
	})
}

// --- Output Page ---

func TestOutputPage(t *testing.T) {
	t.Run("creation", func(t *testing.T) {
		ad := newTestAppData(t)
		page := GetOutputPage(ad, func() {})
		if page == nil {
			t.Fatal("page is nil")
		}
	})

	t.Run("convert defaults", func(t *testing.T) {
		ad := newTestAppData(t)
		page := GetOutputPage(ad, func() {})
		cg := findCheckGroup(collectObjects(page))
		if cg == nil {
			t.Fatal("convert checkgroup not found")
		}
		if !cg.Horizontal {
			t.Error("checkgroup should be horizontal")
		}
		if !cg.Required {
			t.Error("checkgroup should be required")
		}
		if !slices.Equal(cg.Selected, []string{convertToOptions[0]}) {
			t.Errorf("selected = %v, want %v", cg.Selected, convertToOptions[0])
		}
	})

	t.Run("force unchecked by default", func(t *testing.T) {
		ad := newTestAppData(t)
		page := GetOutputPage(ad, func() {})
		c := findCheck(collectObjects(page), "Force convert stale notes")
		if c == nil {
			t.Fatal("force checkbox not found")
		}
		if c.Checked {
			t.Error("force should be unchecked by default")
		}
	})

	t.Run("force toggle updates appData", func(t *testing.T) {
		ad := newTestAppData(t)
		page := GetOutputPage(ad, func() {})
		c := findCheck(collectObjects(page), "Force convert stale notes")
		if c == nil {
			t.Fatal("force checkbox not found")
		}

		c.SetChecked(true)
		if !ad.force {
			t.Error("expected force=true after check")
		}
		c.SetChecked(false)
		if ad.force {
			t.Error("expected force=false after uncheck")
		}
	})

	t.Run("convert toggle updates appData", func(t *testing.T) {
		ad := newTestAppData(t)
		page := GetOutputPage(ad, func() {})
		cg := findCheckGroup(collectObjects(page))
		if cg == nil {
			t.Fatal("convert checkgroup not found")
		}

		cg.SetSelected([]string{i.ConvertPDF})
		if len(ad.convertTo) != 1 || ad.convertTo[0] != i.ConvertPDF {
			t.Errorf("convertTo = %v, want [%s]", ad.convertTo, i.ConvertPDF)
		}

		cg.SetSelected(convertToOptions)
		if !slices.Equal(ad.convertTo, convertToOptions) {
			t.Errorf("convertTo = %v, want %v", ad.convertTo, convertToOptions)
		}
	})

	t.Run("select button", func(t *testing.T) {
		ad := newTestAppData(t)
		page := GetOutputPage(ad, func() {})
		b := findButton(collectObjects(page), "Select folder for output files")
		if b == nil {
			t.Fatal("select button not found")
		}
		if b.Importance != widget.HighImportance {
			t.Errorf("importance = %v, want HighImportance", b.Importance)
		}
	})
}

// --- Preview Page ---

func setupPreviewAppData(t *testing.T, ad *AppData, inputDir, outputDir string) {
	t.Helper()
	ad.input = inputDir
	ad.force = true
	uri := storage.NewFileURI(outputDir)
	lu, err := storage.ListerForURI(uri)
	if err != nil {
		t.Fatalf("ListerForURI: %v", err)
	}
	ad.output = lu
}

func TestPreviewPage(t *testing.T) {
	t.Run("empty dir shows no files message", func(t *testing.T) {
		ad := newTestAppData(t)
		dir := t.TempDir()
		setupPreviewAppData(t, ad, dir, dir)

		page := GetPreviewPage(ad)
		if page == nil {
			t.Fatal("page is nil")
		}

		label := findLabel(collectObjects(page), "No .note files")
		if label == nil {
			t.Fatal("no-files label not found")
		}
		if label.Text != "No .note files to convert!" {
			t.Errorf("label = %q", label.Text)
		}
	})

	t.Run("empty dir shows quit button", func(t *testing.T) {
		ad := newTestAppData(t)
		dir := t.TempDir()
		setupPreviewAppData(t, ad, dir, dir)

		page := GetPreviewPage(ad)
		if page == nil {
			t.Fatal("page is nil")
		}

		b := findButton(collectObjects(page), "Quit")
		if b == nil {
			t.Fatal("expected Quit button when no files")
		}
	})

	t.Run("output label contains path", func(t *testing.T) {
		ad := newTestAppData(t)
		dir := t.TempDir()
		setupPreviewAppData(t, ad, dir, dir)

		page := GetPreviewPage(ad)
		if page == nil {
			t.Fatal("page is nil")
		}

		label := findLabel(collectObjects(page), "Output at:")
		if label == nil {
			t.Fatal("output label not found")
		}
		if !strings.Contains(label.Text, dir) {
			t.Errorf("label %q does not contain dir %q", label.Text, dir)
		}
	})

	t.Run("with notes shows file list and convert button", func(t *testing.T) {
		src := filepath.Join("..", "..", "test-files", "A5X2", "Standard.note")
		data, err := os.ReadFile(src)
		if err != nil {
			t.Skip("test fixture not found")
		}

		ad := newTestAppData(t)
		inputDir := t.TempDir()
		outputDir := t.TempDir()
		os.WriteFile(filepath.Join(inputDir, "Standard.note"), data, 0o644)
		setupPreviewAppData(t, ad, inputDir, outputDir)

		page := GetPreviewPage(ad)
		if page == nil {
			t.Fatal("page is nil")
		}

		objs := collectObjects(page)

		label := findLabel(objs, "Selected Files")
		if label == nil {
			t.Fatal("selected-files label not found")
		}

		b := findButton(objs, "Convert now!")
		if b == nil {
			t.Fatal("convert button not found")
		}
		if b.Importance != widget.HighImportance {
			t.Errorf("importance = %v, want HighImportance", b.Importance)
		}
	})
}
