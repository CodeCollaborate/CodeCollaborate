package dbfs

import (
	"testing"

	"path/filepath"

	"github.com/CodeCollaborate/Server/modules/config"
)

func testConfigSetup(t *testing.T) {
	config.SetConfigDir("../../config")
	err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	config.GetConfig().ServerConfig.ProjectPath = filepath.Clean(filepath.Join(config.GetConfig().ServerConfig.ProjectPath, "_testFiles"))
}
