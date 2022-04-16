package config

import (
	"os"
	"path/filepath"
	"testing"
)

const configJson = `{
	"panel": {
		"settings": {
			"companyName": "My Company"
		}
	}
}`

// chdir changes to the specified directory and restores the previous working
// directory on test cleanup.
func chdir(t *testing.T, dir string) {
	prevDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Couldn't get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Couldn't change working directory: %v", err)
	}
	t.Cleanup(func() { os.Chdir(prevDir) })
}

func TestLoadConfigWorkingDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(configJson), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	chdir(t, dir)

	if err := LoadConfigFile(""); err != nil {
		t.Fatalf("LoadConfigFile error: %v", err)
	}
	if got := GetString("panel.settings.companyName"); got != "My Company" {
		t.Errorf("GetString(\"panel.settings.companyName\") = %q, want \"My Company\"", got)
	}
}
