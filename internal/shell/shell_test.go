package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uda/uda/internal/config"
)

func TestInitBashContainsPipInstallRouting(t *testing.T) {
	script := Init("bash", "/tmp/uda-bin")
	if !strings.Contains(script, "uda install") {
		t.Fatalf("expected bash init script to contain uda install routing")
	}
	if !strings.Contains(script, `case "$cmd" in`) {
		t.Fatalf("expected bash init case statement")
	}
}

func TestGenerateActivateScriptRemovesCurrentEnvWhenPresent(t *testing.T) {
	envName := "testenv"
	oldHomeDir := config.HomeDir
	config.HomeDir = filepath.Join(t.TempDir(), ".uda")
	defer func() { config.HomeDir = oldHomeDir }()

	envPath := filepath.Join(config.EnvsPath(), envName)
	if err := os.MkdirAll(config.EnvsPath(), 0755); err != nil {
		t.Fatalf("prepare env dir: %v", err)
	}
	if err := os.MkdirAll(envPath, 0755); err != nil {
		t.Fatalf("prepare env path: %v", err)
	}
	defer os.RemoveAll(config.EnvsPath())
	script, err := GenerateActivateScript(envName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `if [ -n "$VIRTUAL_ENV" ] && command -v _uda_remove_path_entry >/dev/null 2>&1; then`
	if !strings.Contains(script, expected) {
		t.Fatalf("expected pre-cleanup guard, got: %s", script)
	}
	if !strings.Contains(script, `export _UDA_ACTIVE_ENV="`+envName+`"`) {
		t.Fatalf("expected active env export")
	}
	if !strings.Contains(script, `export VIRTUAL_ENV="`+envPath+`"`) {
		t.Fatalf("expected virtual env export")
	}
}
