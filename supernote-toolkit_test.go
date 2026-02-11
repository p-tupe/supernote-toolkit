package main_test

import (
	"os/exec"
	"testing"
)

func TestCLI(t *testing.T) {
	out, err := exec.Command("go", "test", "-count=1", "./cmd/cli/").CombinedOutput()
	if err != nil {
		t.Fatalf("cmd/cli tests failed:\n%s", out)
	}
	t.Log(string(out))
}

func TestApp(t *testing.T) {
	out, err := exec.Command("go", "test", "-count=1", "./cmd/app/").CombinedOutput()
	if err != nil {
		t.Fatalf("cmd/app tests failed:\n%s", out)
	}
	t.Log(string(out))
}
