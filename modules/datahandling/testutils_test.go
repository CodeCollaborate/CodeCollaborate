package datahandling

import (
	"testing"
	"time"

	"github.com/CodeCollaborate/Server/modules/config"
)

func configSetup(t *testing.T) {
	config.SetConfigDir("../../config")
	err := config.InitConfig()
	if err != nil {
		t.Fatal(err)
	}
}

func testToken(t *testing.T, username string) string {
	return signedTokenOrDie(t, username, time.Now().Unix(), time.Now().Add(1*time.Minute).Unix(), privKey)
}
