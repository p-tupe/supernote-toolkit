package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binPath string

var testFixtures = filepath.Join("..", "..", "test-files")

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "cli-test-bin")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)

	binPath = filepath.Join(tmp, "supernote-toolkit")
	out, err := exec.Command("go", "build", "-o", binPath, ".").CombinedOutput()
	if err != nil {
		panic("build failed: " + string(out))
	}

	os.Exit(m.Run())
}

func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func runCLI(t *testing.T, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	cmd := exec.CommandContext(t.Context(), binPath, args...)
	var outb, errb strings.Builder
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return outb.String(), errb.String(), exitErr.ExitCode()
	} else if err != nil {
		t.Fatal(err)
	}
	return outb.String(), errb.String(), 0
}

func assertPDF(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected PDF at %s: %v", path, err)
	}
	if !strings.HasPrefix(string(data), "%PDF-") {
		t.Fatalf("invalid PDF header: %s", path)
	}
}

func assertPNG(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected PNG at %s: %v", path, err)
	}
	if len(data) < 4 || data[0] != 0x89 || data[1] != 'P' || data[2] != 'N' || data[3] != 'G' {
		t.Fatalf("invalid PNG header: %s", path)
	}
}

func assertNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Errorf("expected %s to not exist", path)
	}
}

// setupA5X2Dir creates:
//
//	tmp/
//	  Standard.note
//	  sub/
//	    Artifacts.note
func setupA5X2Dir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	copyFile(t, filepath.Join(testFixtures, "A5X2", "Standard.note"), filepath.Join(dir, "Standard.note"))

	sub := filepath.Join(dir, "sub")
	os.MkdirAll(sub, 0o755)
	copyFile(t, filepath.Join(testFixtures, "A5X2", "Artifacts.note"), filepath.Join(sub, "Artifacts.note"))

	return dir
}

func setupA6X2File(t *testing.T, name string) string {
	t.Helper()
	dir := t.TempDir()
	copyFile(t, filepath.Join(testFixtures, "A6X2", name), filepath.Join(dir, name))
	return dir
}

// --- Validation ---

func TestNoInput(t *testing.T) {
	t.Parallel()
	stdout, _, code := runCLI(t)
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if !strings.Contains(stdout, "Example:") {
		t.Error("expected usage output")
	}
}

func TestInvalidInput(t *testing.T) {
	t.Parallel()
	_, _, code := runCLI(t, "-input", "/nonexistent/path")
	if code == 0 {
		t.Error("expected non-zero exit code")
	}
}

func TestInvalidDevice(t *testing.T) {
	t.Parallel()
	input := setupA5X2Dir(t)
	output := t.TempDir()
	_, _, code := runCLI(t, "-input", input, "-output", output, "-device", "INVALID")
	if code == 0 {
		t.Error("expected non-zero exit code for invalid device")
	}
}

// --- A5X2 ---

func TestA5X2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		args         []string
		outputSuffix string
		wantPDFs     []string
		wantPNGs     []string
		notExist     []string
	}{
		{
			name:     "pdf default",
			wantPDFs: []string{"Standard.pdf", "sub/Artifacts.pdf"},
			notExist: []string{"Standard", "sub/Artifacts"},
		},
		{
			name:     "png only",
			args:     []string{"-pdf=false", "-png"},
			wantPNGs: []string{"Standard/PAGE0.png", "sub/Artifacts/PAGE0.png"},
			notExist: []string{"Standard.pdf"},
		},
		{
			name:     "both formats",
			args:     []string{"-png"},
			wantPDFs: []string{"Standard.pdf"},
			wantPNGs: []string{"Standard/PAGE0.png"},
		},
		{
			name:     "neither format",
			args:     []string{"-pdf=false"},
			notExist: []string{"Standard.pdf", "Standard"},
		},
		{
			name:     "recurse false",
			args:     []string{"-recurse=false"},
			wantPDFs: []string{"Standard.pdf"},
			notExist: []string{"sub/Artifacts.pdf"},
		},
		{
			name:         "custom nested output",
			outputSuffix: "nested/out",
			wantPDFs:     []string{"Standard.pdf"},
		},
		{
			name:     "explicit device A5X2",
			args:     []string{"-device", "A5X2"},
			wantPDFs: []string{"Standard.pdf"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input := setupA5X2Dir(t)
			output := t.TempDir()
			if tt.outputSuffix != "" {
				output = filepath.Join(output, tt.outputSuffix)
			}

			args := append([]string{"-input", input, "-output", output}, tt.args...)
			_, _, code := runCLI(t, args...)
			if code != 0 {
				t.Fatalf("exit code %d", code)
			}

			for _, p := range tt.wantPDFs {
				assertPDF(t, filepath.Join(output, p))
			}
			for _, p := range tt.wantPNGs {
				assertPNG(t, filepath.Join(output, p))
			}
			for _, p := range tt.notExist {
				assertNotExists(t, filepath.Join(output, p))
			}
		})
	}
}

func TestA5X2Horizontal(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	copyFile(t, filepath.Join(testFixtures, "mixed", "horizontal_1090.note"), filepath.Join(dir, "horizontal_1090.note"))
	output := t.TempDir()

	_, _, code := runCLI(t, "-input", dir, "-output", output, "-device", "A5X2")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	assertPDF(t, filepath.Join(output, "horizontal_1090.pdf"))
}

// --- A6X2 ---

func TestA6X2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		file     string
		args     []string
		wantPDFs []string
		wantPNGs []string
	}{
		{
			name:     "auto-detect nomad pdf",
			file:     "1to10.note",
			wantPDFs: []string{"1to10.pdf"},
		},
		{
			name:     "explicit device A6X2 pdf",
			file:     "1to10.note",
			args:     []string{"-device", "A6X2"},
			wantPDFs: []string{"1to10.pdf"},
		},
		{
			name:     "nomad png",
			file:     "1to10.note",
			args:     []string{"-pdf=false", "-png"},
			wantPNGs: []string{"1to10/PAGE0.png"},
		},
		{
			name:     "vertical orientation",
			file:     "vertical_1180.note",
			wantPDFs: []string{"vertical_1180.pdf"},
		},
		{
			name:     "horizontal orientation",
			file:     "horizontal_1270.note",
			wantPDFs: []string{"horizontal_1270.pdf"},
		},
		{
			name:     "horizontal orientation png",
			file:     "horizontal_1270.note",
			args:     []string{"-pdf=false", "-png"},
			wantPNGs: []string{"horizontal_1270/PAGE0.png"},
		},
		{
			name:     "shapes and RTR",
			file:     "nomad-3.15.27-blank-shapes-and-RTR.note",
			wantPDFs: []string{"nomad-3.15.27-blank-shapes-and-RTR.pdf"},
		},
		{
			name:     "multi-page blank",
			file:     "nomad-3.15.27-blank-2p.note",
			args:     []string{"-pdf=false", "-png"},
			wantPNGs: []string{"nomad-3.15.27-blank-2p/PAGE0.png", "nomad-3.15.27-blank-2p/PAGE1.png"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := setupA6X2File(t, tt.file)
			output := t.TempDir()

			args := append([]string{"-input", dir, "-output", output}, tt.args...)
			_, _, code := runCLI(t, args...)
			if code != 0 {
				t.Fatalf("exit code %d", code)
			}

			for _, p := range tt.wantPDFs {
				assertPDF(t, filepath.Join(output, p))
			}
			for _, p := range tt.wantPNGs {
				assertPNG(t, filepath.Join(output, p))
			}
		})
	}
}

// --- Force reconvert ---

func TestForceReconvert(t *testing.T) {
	t.Parallel()
	input := setupA5X2Dir(t)
	output := t.TempDir()

	_, _, code := runCLI(t, "-input", input, "-output", output)
	if code != 0 {
		t.Fatalf("first run: exit code %d", code)
	}

	pdf := filepath.Join(output, "Standard.pdf")
	info1, err := os.Stat(pdf)
	if err != nil {
		t.Fatal(err)
	}

	// without force - output is fresh, should be skipped
	_, _, code = runCLI(t, "-input", input, "-output", output)
	if code != 0 {
		t.Fatalf("second run: exit code %d", code)
	}
	info2, _ := os.Stat(pdf)
	if !info1.ModTime().Equal(info2.ModTime()) {
		t.Error("file was regenerated without -force")
	}

	// with force - should overwrite
	_, _, code = runCLI(t, "-input", input, "-output", output, "-force")
	if code != 0 {
		t.Fatalf("force run: exit code %d", code)
	}
	info3, _ := os.Stat(pdf)
	if info3.ModTime().Equal(info1.ModTime()) {
		t.Error("file was NOT regenerated with -force")
	}
}
