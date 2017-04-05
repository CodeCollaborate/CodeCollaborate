package config

import (
	"path/filepath"
	"testing"
)

// SetupTestingConfig specifies the config directory, and adds a temporary folder for test files.
func SetupTestingConfig(t *testing.T, dir string) {
	SetConfigDir(dir)
	err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	GetConfig().ServerConfig.ProjectPath = filepath.Clean(filepath.Join(GetConfig().ServerConfig.ProjectPath, "_testFiles"))
}
