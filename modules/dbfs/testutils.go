package dbfs

import (
	"testing"

	"github.com/CodeCollaborate/Server/modules/config"
)

func configSetup(t *testing.T) {
	config.SetConfigDir("../../config")
	err := config.InitConfig()
	if err != nil {
		t.Fatal(err)
	}
}
