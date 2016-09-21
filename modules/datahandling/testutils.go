package datahandling

import (
	"testing"

	"github.com/CodeCollaborate/Server/modules/config"
)

func configSetupUnauthenticated(t *testing.T) {
	config.SetConfigDir("../../config")
	err := config.InitConfig()
	if err != nil {
		t.Fatal(err)
	}

	config.GetConfig().ServerConfig.DisableAuth = true
}

func configSetup(t *testing.T) {
	config.SetConfigDir("../../config")
	err := config.InitConfig()
	if err != nil {
		t.Fatal(err)
	}
}
