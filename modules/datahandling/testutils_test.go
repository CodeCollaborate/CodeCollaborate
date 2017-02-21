package datahandling

import (
	"testing"
	"time"

	"github.com/CodeCollaborate/Server/modules/config"
)

func configSetup(t *testing.T) {
	config.SetConfigDir("../../config")
	err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
}

func testToken(t *testing.T, username string) string {
	rsa, _ := config.GenRSA(1024) // make it small so it's faster for tests
	// fun story: I originally had this set to 4096 by default and accidentally
	// made the tests take 10 minutes
	return signedTokenOrDie(t, username, time.Now().Unix(), time.Now().Add(1*time.Minute).Unix(), rsa)
}
