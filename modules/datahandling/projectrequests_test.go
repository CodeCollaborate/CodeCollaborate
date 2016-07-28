package datahandling

import (
	"log"
	"testing"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/dbfs"
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

// TODO (testing/required): testing... lots of testing

func TestProjectLookupRequest_Process(t *testing.T) {
	configSetup()
	req := *new(projectLookupRequest)

	req.SenderID = "loganga"
	req.Resource = "Project"
	req.Method = "Lookup"
	req.SenderToken = "supersecure"
	req.ProjectIDs = []int64{12345, 38292}

	db := dbfs.NewDBMock()

	continuations, err := req.process(db)

	if err != nil {
		t.Fatal(err)
	}

	if len(continuations) != 1 {
		t.Fatal("did not properly process")
	}
}
