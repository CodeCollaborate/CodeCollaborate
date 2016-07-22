package datahandling

import (
	"log"
	"testing"

	"github.com/CodeCollaborate/Server/modules/config"
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

//TODO (testing/required): testing... lots of testing

// TODO (testing/required): switch dbfs to use a dbfs object which owns the database methods and implements an interface so that we can mock the interface for testing

func TestProjectLookupRequest_Process(t *testing.T) {
	configSetup()
	req := *new(projectLookupRequest)

	req.SenderID = "loganga"
	req.Resource = "Project"
	req.Method = "Lookup"
	req.SenderToken = "supersecure"
	req.ProjectIDs = []int64{12345, 38292}

	//dbfs.MySQLProjectLookup = func() {}

	continuations, err := req.process()

	if err != nil {
		t.Fatal(err)
	}

	if len(continuations) != 1 {
		t.Fatal("did not properly process")
	}
}
