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
}

func TestOpenCouchBase(t *testing.T) {
	configSetup()

	cb, err := openCouchBase()
	if err != nil {
		t.Fatal(err)
	}
	defer cb.close()

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
	configSetup()
	CBDeleteFile(1)
	CBInsertNewFile(1, 2, []string{"hey there", "sup"})

	err := CBAppendFileChange(1, 3, "yooooo")
	if err != nil {
		t.Fatal(err)
	}

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
	if changes[2] != "yooooo" {
		t.Fatal(err)
	}
	if len(changes) != 3 {
		t.Fatal("resultant changes are not correct")
	}

	ver, err := CBGetFileVersion(1)
	if ver != 3 {
		t.Fatal("wrong file version")
	}

	CBDeleteFile(1)
}
