package datahandling

import (
	"encoding/json"
	"fmt"
	"testing"
	"github.com/CodeCollaborate/Server/modules/config"
	"log"
)

func configSetup() {
	config.SetConfigDir("../../config")
	err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: this is backup for the tests, and will likely fail on a
	// non-local system unless the DB's have been set up to allow for this
	if val := config.GetConfig().ConnectionConfig["MySQL"].Schema; val == "" {
		tempCon := config.GetConfig().ConnectionConfig["MySQL"]
		tempCon.Schema = "testing"
		config.GetConfig().ConnectionConfig["MySQL"] = tempCon
	}
	if val := config.GetConfig().ConnectionConfig["Couchbase"].Schema; val == "" {
		tempCon := config.GetConfig().ConnectionConfig["Couchbase"]
		tempCon.Schema = "testing"
		config.GetConfig().ConnectionConfig["Couchbase"] = tempCon
	}
}

func TestProjectLookupRequest_Process(t *testing.T) {
	configSetup()
	req := *new(abstractRequest)

	req.SenderID = "loganga"
	req.Resource = "Project"
	req.Method = "Lookup"
	req.SenderToken = "supersecure"

	req.Data = json.RawMessage("{\"ProjectIds\": [12345, 38292]}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	response, notification, err2 := newRequest.process()
	if err2 != nil {
		t.Fatal(err2)
	}

	if response == nil || notification == nil {
		// TODO: check if process worked
		fmt.Println("unimplemented response and notification logic!")
	}
}
