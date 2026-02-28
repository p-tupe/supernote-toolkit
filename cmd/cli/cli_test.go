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

func assertTXT(t *testing.T, path string, wantContent bool) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected TXT at %s: %v", path, err)
	}
	if wantContent && len(strings.TrimSpace(string(data))) == 0 {
		t.Fatalf("expected non-empty TXT at %s", path)
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
//	  blank_v1000_1p.note
//	  sub/
//	    blank_v1000_1p_artifacts.note
func setupA5X2Dir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	copyFile(t, filepath.Join(testFixtures, "A5X2", "blank_v1000_1p.note"), filepath.Join(dir, "blank_v1000_1p.note"))

	sub := filepath.Join(dir, "sub")
	os.MkdirAll(sub, 0o755)
	copyFile(t, filepath.Join(testFixtures, "A5X2", "blank_v1000_1p_artifacts.note"), filepath.Join(sub, "blank_v1000_1p_artifacts.note"))

	return dir
}

func setupFile(t *testing.T, deviceDir, name string) string {
	t.Helper()
	dir := t.TempDir()
	copyFile(t, filepath.Join(testFixtures, deviceDir, name), filepath.Join(dir, name))
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

// --- A5X2 (Manta X2, N5) ---

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
			wantPDFs: []string{"blank_v1000_1p.pdf", "sub/blank_v1000_1p_artifacts.pdf"},
			notExist: []string{"blank_v1000_1p", "sub/blank_v1000_1p_artifacts"},
		},
		{
			name:     "png only",
			args:     []string{"-pdf=false", "-png"},
			wantPNGs: []string{"blank_v1000_1p/PAGE0.png", "sub/blank_v1000_1p_artifacts/PAGE0.png"},
			notExist: []string{"blank_v1000_1p.pdf"},
		},
		{
			name:     "both formats",
			args:     []string{"-png"},
			wantPDFs: []string{"blank_v1000_1p.pdf"},
			wantPNGs: []string{"blank_v1000_1p/PAGE0.png"},
		},
		{
			name:     "neither format",
			args:     []string{"-pdf=false"},
			notExist: []string{"blank_v1000_1p.pdf", "blank_v1000_1p"},
		},
		{
			name:     "recurse false",
			args:     []string{"-recurse=false"},
			wantPDFs: []string{"blank_v1000_1p.pdf"},
			notExist: []string{"sub/blank_v1000_1p_artifacts.pdf"},
		},
		{
			name:         "custom nested output",
			outputSuffix: "nested/out",
			wantPDFs:     []string{"blank_v1000_1p.pdf"},
		},
		{
			name:     "explicit device A5X2",
			args:     []string{"-device", "A5X2"},
			wantPDFs: []string{"blank_v1000_1p.pdf"},
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

	dir := setupFile(t, "A5X2", "blank_h1090_1p_rtr.note")
	output := t.TempDir()

	_, _, code := runCLI(t, "-input", dir, "-output", output, "-device", "A5X2")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	assertPDF(t, filepath.Join(output, "blank_h1090_1p_rtr.pdf"))
}

// --- A5X (Manta gen1, auto-detected as A6X2) ---

func TestA5X(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		file     string
		args     []string
		wantPDFs []string
		wantPNGs []string
	}{
		{
			name:     "auto-detect pdf",
			file:     "ruled_v1000_10p_rtr.note",
			wantPDFs: []string{"ruled_v1000_10p_rtr.pdf"},
		},
		{
			name:     "explicit device A6X2",
			file:     "ruled_v1000_10p_rtr.note",
			args:     []string{"-device", "A6X2"},
			wantPDFs: []string{"ruled_v1000_10p_rtr.pdf"},
		},
		{
			name:     "png",
			file:     "ruled_v1000_10p_rtr.note",
			args:     []string{"-pdf=false", "-png"},
			wantPNGs: []string{"ruled_v1000_10p_rtr/PAGE0.png"},
		},
		{
			name:     "old format",
			file:     "ruled8mm_v1000_2p.note",
			wantPDFs: []string{"ruled8mm_v1000_2p.pdf"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := setupFile(t, "A5X", tt.file)
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

// --- A6X (Nomad gen1, auto-detected as A6X2) ---

func TestA6X(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		file     string
		args     []string
		wantPDFs []string
		wantPNGs []string
	}{
		{
			name:     "shapes and RTR",
			file:     "shapes_v1000_1p_rtr.note",
			wantPDFs: []string{"shapes_v1000_1p_rtr.pdf"},
		},
		{
			name:     "multi-page",
			file:     "blank_v1000_2p.note",
			args:     []string{"-pdf=false", "-png"},
			wantPNGs: []string{"blank_v1000_2p/PAGE0.png", "blank_v1000_2p/PAGE1.png"},
		},
		{
			name:     "RTR turkish",
			file:     "blank_v1000_1p_rtr_tr.note",
			wantPDFs: []string{"blank_v1000_1p_rtr_tr.pdf"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := setupFile(t, "A6X", tt.file)
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

// --- A6X2 (Nomad X2, N6) ---

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
			name:     "vertical v1180",
			file:     "blank_v1180_1p.note",
			wantPDFs: []string{"blank_v1180_1p.pdf"},
		},
		{
			name:     "horizontal h1270 pdf",
			file:     "blank_h1270_1p.note",
			wantPDFs: []string{"blank_h1270_1p.pdf"},
		},
		{
			name:     "horizontal h1270 png",
			file:     "blank_h1270_1p.note",
			args:     []string{"-pdf=false", "-png"},
			wantPNGs: []string{"blank_h1270_1p/PAGE0.png"},
		},
		{
			name:     "RTR german task list",
			file:     "task_v1000_1p_rtr_de.note",
			wantPDFs: []string{"task_v1000_1p_rtr_de.pdf"},
		},
		{
			name:     "explicit device A6X2",
			file:     "blank_h1270_1p.note",
			args:     []string{"-device", "A6X2"},
			wantPDFs: []string{"blank_h1270_1p.pdf"},
		},
		{
			name:     "3-layer compositing",
			file:     "multilayer_v1000_1p.note",
			wantPDFs: []string{"multilayer_v1000_1p.pdf"},
		},
		{
			name:     "marker colors",
			file:     "markers_v1000_1p.note",
			wantPDFs: []string{"markers_v1000_1p.pdf"},
		},
		{
			name:     "pressure variation",
			file:     "pressure_v1000_1p.note",
			wantPDFs: []string{"pressure_v1000_1p.pdf"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := setupFile(t, "A6X2", tt.file)
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

// --- TXT extraction ---

func TestTXT(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		deviceDir   string
		file        string
		wantContent bool // true = expect non-empty text
	}{
		// Files with recognized text
		{"A5X lined RTR", "A5X", "lined_v1000_1p_rtr.note", true},
		{"A6X shapes RTR", "A6X", "shapes_v1000_1p_rtr.note", true},
		{"A6X turkish RTR", "A6X", "blank_v1000_1p_rtr_tr.note", true},
		{"A6X2 german RTR", "A6X2", "task_v1000_1p_rtr_de.note", true},
		{"A5X2 text RTR", "A5X2", "text_v1000_1p_rtr.note", true},
		// RTR enabled but no text written — TXT exists, empty
		{"A5X2 RTR no content", "A5X2", "wip_v1000_1p_rtr.note", false},
		{"A5X ruled RTR no content", "A5X", "ruled_v1000_10p_rtr.note", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := setupFile(t, tt.deviceDir, tt.file)
			output := t.TempDir()

			_, _, code := runCLI(t, "-input", dir, "-output", output, "-pdf=false", "-txt")
			if code != 0 {
				t.Fatalf("exit code %d", code)
			}

			stem := strings.TrimSuffix(tt.file, ".note")
			assertTXT(t, filepath.Join(output, stem+".txt"), tt.wantContent)
			assertNotExists(t, filepath.Join(output, stem+".pdf"))
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

	pdf := filepath.Join(output, "blank_v1000_1p.pdf")
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
