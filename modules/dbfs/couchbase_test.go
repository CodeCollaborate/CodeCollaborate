package dbfs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseImpl_OpenCouchBase(t *testing.T) {
	configSetup(t)
	di := new(DatabaseImpl)

	cb, err := di.openCouchBase()
	if err != nil {
		t.Fatal(err)
	}
	defer di.CloseCouchbase()

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

func TestDatabaseImpl_CloseCouchbase(t *testing.T) {
	configSetup(t)
	di := new(DatabaseImpl)

	db, err := di.openCouchBase()
	if err != nil || db == nil {
		t.Fatal(err)
	}
	err = di.CloseCouchbase()
	if err != nil {
		t.Fatal(err)
	}
	err = di.CloseCouchbase()
	if err != ErrDbNotInitialized {
		t.Fatal("Wrong error recieved")
	}
}

func TestDatabaseImpl_CBInsertNewFile(t *testing.T) {
	configSetup(t)
	di := new(DatabaseImpl)

	// ensure it doesn't actually exist
	di.CBDeleteFile(1)

	f := cbFile{FileID: 1, Version: 2, Changes: []string{"hey there", "sup"}, Lock: false}
	err := di.cbInsertNewFile(f)
	if err != nil {
		t.Fatal(err)
	}

	err = di.cbInsertNewFile(f)
	if err == nil {
		t.Fatal("Insert should have failed when inserting into an existing key")
	}

	//cleanup
	di.CBDeleteFile(1)
}

func TestDatabaseImpl_CBInsertNewFileByDetails(t *testing.T) {
	configSetup(t)
	di := new(DatabaseImpl)

	di.CBDeleteFile(1)

	err := di.CBInsertNewFile(1, 2, []string{"hey there", "sup"})
	if err != nil {
		t.Fatal(err)
	}
	err = di.CBInsertNewFile(1, 2, []string{"wow"})
	if err == nil {
		t.Fatal("Insert should have failed when inserting into an existing key")
	}

	// cleanup
	di.CBDeleteFile(1)
}

func TestDatabaseImpl_CBDeleteFile(t *testing.T) {
	configSetup(t)
	di := new(DatabaseImpl)

	f := cbFile{FileID: 1, Version: 2, Changes: []string{"hey there", "sup"}, Lock: false}
	di.cbInsertNewFile(f)

	err := di.CBDeleteFile(1)
	if err != nil {
		t.Fatal(err)
	}

	err = di.CBDeleteFile(1)
	if err == nil {
		t.Fatal("Delete should have failed here")
	}
}

func TestDatabaseImpl_CBGetFileVersion(t *testing.T) {
	configSetup(t)
	di := new(DatabaseImpl)

	di.CBDeleteFile(1)
	di.CBInsertNewFile(1, 2, []string{"hey there", "sup"})

	ver, err := di.CBGetFileVersion(1)
	if err != nil {
		t.Fatal(err)
	}

	if ver != 2 {
		t.Fatal(err)
	}

	di.CBDeleteFile(1)
}

func TestDatabaseImpl_CBGetFileChanges(t *testing.T) {
	// setup
	configSetup(t)
	di := new(DatabaseImpl)

	di.CBDeleteFile(1)
	di.CBInsertNewFile(1, 2, []string{"hey there", "sup"})

	changes, err := di.CBGetFileChanges(1)
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

	di.CBDeleteFile(1)
}

func TestDatabaseImpl_CBAppendFileChange(t *testing.T) {
	var originalFileVersion int64 = 2
	var fileID int64 = 1
	configSetup(t)
	di := new(DatabaseImpl)

	di.CBDeleteFile(fileID)

	// although these are not valid patches, this is purely a test of the logic, not of the patching
	di.CBInsertNewFile(fileID, originalFileVersion, []string{"hey there", "sup"})

	version, err := di.CBAppendFileChange(fileID, originalFileVersion, []string{"yooooo"})
	assert.NoError(t, err, "unexpected error appending changes")

	assert.Equal(t, originalFileVersion+1, version, "version did not update properly")

	changes, err := di.CBGetFileChanges(fileID)
	assert.NoError(t, err, "unexpected error getting changes")

	assert.Len(t, changes, 3, "resultant changes not the correct length")

	assert.Equal(t, "hey there", changes[0], "first change was not correct")
	assert.Equal(t, "sup", changes[1], "second change was not correct")

	assert.Equal(t, "yooooo", changes[2], "newly inserted change was not correct")

	ver, err := di.CBGetFileVersion(fileID)
	assert.EqualValues(t, 3, ver, "wrong file version")

	di.CBDeleteFile(fileID)
}
