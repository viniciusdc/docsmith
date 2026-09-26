package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

// binDir holds the hello binary built once by TestMain.
var binDir string

// TestMain builds the binary once into a temporary directory so every script
// runs the same executable a user would install.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "hello-script-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "mktemp: %v\n", err)
		os.Exit(1)
	}
	exe := filepath.Join(dir, "hello")
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}
	if out, err := exec.Command("go", "build", "-o", exe, ".").CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build hello: %v\n%s", err, out)
		_ = os.RemoveAll(dir)
		os.Exit(1)
	}
	binDir = dir

	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

// TestScript runs every testdata/examples/*.txtar as a testscript, with the
// freshly built binary first on PATH. These are the same files genexamples
// renders into the docs, so every example shown there is exercised here.
func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: filepath.Join("..", "..", "testdata", "examples"),
		Setup: func(env *testscript.Env) error {
			env.Setenv("PATH", binDir+string(filepath.ListSeparator)+env.Getenv("PATH"))
			return nil
		},
	})
}
