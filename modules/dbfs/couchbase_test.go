package dbfs

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

func TestOpenCouchBase(t *testing.T) {
	configSetup()

	cb, err := openCouchBase()
	if err != nil {
		t.Fatal(err)
	}
	defer CloseCouchbase()

	_, err = cb.bucket.Upsert("testingDocumentPleaseIgnore", "mydoc", 0)
	if err != nil {
		t.Fatal(err)
	}

	var value interface{}
	_, err = cb.bucket.Get("testingDocumentPleaseIgnore", &value)
	if err != nil {
		t.Fatal(err)
	}

	if value != "mydoc" {
		t.Fatal("couchbase testing value wrong somehow")
	}

	_, err = cb.bucket.Remove("testingDocumentPleaseIgnore", 0)
	if err != nil {
		t.Fatal(err)
	}

}

func TestCloseCouchbase(t *testing.T) {
	configSetup()
	_, err := openCouchBase()
	if err != nil {
		t.Fatal(err)
	}
	err = CloseCouchbase()
	if err != nil {
		t.Fatal(err)
	}
	err = CloseCouchbase()
	if err != ErrDbNotInitialized {
		t.Fatal("Wrong error recieved")
	}
}

func TestCBInsertNewFile(t *testing.T) {
	configSetup()

	// ensure it doesn't actually exist
	CBDeleteFile(1)

	f := cbFile{FileID: 1, Version: 2, Changes: []string{"hey there", "sup"}}
	err := cbInsertNewFile(f)
	if err != nil {
		t.Fatal(err)
	}

	err = cbInsertNewFile(f)
	if err == nil {
		t.Fatal("Insert should have failed when inserting into an existing key")
	}

	//cleanup
	CBDeleteFile(1)
}

func TestCBInsertNewFileByDetails(t *testing.T) {
	configSetup()
	CBDeleteFile(1)

	err := CBInsertNewFile(1, 2, []string{"hey there", "sup"})
	if err != nil {
		t.Fatal(err)
	}
	err = CBInsertNewFile(1, 2, []string{"wow"})
	if err == nil {
		t.Fatal("Insert should have failed when inserting into an existing key")
	}

	// cleanup
	CBDeleteFile(1)
}

func TestCBDeleteFile(t *testing.T) {
	configSetup()

	f := cbFile{FileID: 1, Version: 2, Changes: []string{"hey there", "sup"}}
	cbInsertNewFile(f)

	err := CBDeleteFile(1)
	if err != nil {
		t.Fatal(err)
	}

	err = CBDeleteFile(1)
	if err == nil {
		t.Fatal("Delete should have failed here")
	}
}

func TestCBGetFileVersion(t *testing.T) {
	// setup
	configSetup()
	CBDeleteFile(1)
	CBInsertNewFile(1, 2, []string{"hey there", "sup"})

	ver, err := CBGetFileVersion(1)
	if err != nil {
		t.Fatal(err)
	}

	if ver != 2 {
		t.Fatal(err)
	}

	CBDeleteFile(1)
}

func TestCBGetFileChanges(t *testing.T) {
	// setup
	configSetup()
	CBDeleteFile(1)
	CBInsertNewFile(1, 2, []string{"hey there", "sup"})

	changes, err := CBGetFileChanges(1)
	if err != nil {
		t.Fatal(err)
	}

	if changes[0] != "hey there" {
		t.Fatal(err)
	}
	if changes[1] != "sup" {
		t.Fatal(err)
	}
	if len(changes) != 2 {
		t.Fatal("resultant changes are not correct")
	}

	CBDeleteFile(1)
}

func TestCBAppendFileChange(t *testing.T) {
	var originalFileVersion int64 = 2
	var fileID int64 = 1

	configSetup()
	CBDeleteFile(fileID)

	// although these are not valid patches, this is purely a test of the logic, not of the patching
	CBInsertNewFile(fileID, originalFileVersion, []string{"hey there", "sup"})

	version, err := CBAppendFileChange(fileID, originalFileVersion, []string{"yooooo"})
	if err != nil {
		t.Fatal(err)
	}

	// new version
	if version != originalFileVersion+1 {
		t.Fatal("version did not update properly")
	}

	changes, err := CBGetFileChanges(fileID)
	if err != nil {
		t.Fatal(err)
	}

	if changes[0] != "hey there" {
		t.Fatal(err)
	}
	if changes[1] != "sup" {
		t.Fatal(err)
	}
	if changes[2] != "yooooo" {
		t.Fatal(err)
	}
	if len(changes) != 3 {
		t.Fatal("resultant changes are not correct")
	}

	ver, err := CBGetFileVersion(fileID)
	if ver != 3 {
		t.Fatal("wrong file version")
	}

	CBDeleteFile(fileID)
}
