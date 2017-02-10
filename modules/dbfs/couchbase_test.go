package dbfs

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseImpl_OpenCouchBase(t *testing.T) {
	testConfigSetup(t)
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
	testConfigSetup(t)
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
	testConfigSetup(t)
	di := new(DatabaseImpl)

	// ensure it doesn't actually exist
	di.CBDeleteFile(1)

	f := cbFile{FileID: 1, Version: 2, Changes: []string{"hey there", "sup"}, UseTemp: false}
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
	testConfigSetup(t)
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
	testConfigSetup(t)
	di := new(DatabaseImpl)

	f := cbFile{FileID: 1, Version: 2, Changes: []string{"hey there", "sup"}, UseTemp: false}
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
	testConfigSetup(t)
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
	testConfigSetup(t)
	di := new(DatabaseImpl)

	file := FileMeta{
		FileID:       1,
		Creator:      "_testuser1",
		CreationDate: time.Now(),
		RelativePath: "/.",
		ProjectID:    0,
		Filename:     "_test_file_123",
	}

	di.CBDeleteFile(file.FileID)
	di.CBInsertNewFile(file.FileID, 2, []string{"hey there", "sup"})

	// NOTE: this was added as a need by us changing to dbfs.PullFile
	di.FileWrite(file.RelativePath, file.Filename, file.ProjectID, []byte{})

	raw, changes, err := di.PullFile(file)
	assert.NoError(t, err, "unexpected error getting changes")

	assert.Empty(t, *raw, "we shouldn't have scrunched")

	assert.Len(t, changes, 2, "resultant changes not the correct length")
	assert.Equal(t, "hey there", changes[0], "first change was not correct")
	assert.Equal(t, "sup", changes[1], "second change was not correct")

	di.CBDeleteFile(1)
	di.FileDelete(file.RelativePath, file.Filename, file.ProjectID)
}

func TestDatabaseImpl_CBAppendFileChange(t *testing.T) {
	var originalFileVersion int64 = 2
	file := FileMeta{
		FileID:       1,
		Creator:      "_testuser1",
		CreationDate: time.Now(),
		RelativePath: "/.",
		ProjectID:    0,
		Filename:     "_test_file_123",
	}

	testConfigSetup(t)
	di := new(DatabaseImpl)

	di.CBDeleteFile(file.FileID)

	patch1 := fmt.Sprintf("v%d:\n1:+6:patch1:\n4", originalFileVersion-1)
	patch2 := fmt.Sprintf("v%d:\n2:+6:patch2:\n10", originalFileVersion-1)
	patch3 := fmt.Sprintf("v%d:\n3:+6:patch3:\n10", originalFileVersion)
	patch4 := fmt.Sprintf("v%d:\n4:+6:patch4:\n10", originalFileVersion)

	// although these are not valid patches, this is purely a test of the logic, not of the patching
	// because of that this might fail in the future
	di.CBInsertNewFile(file.FileID, originalFileVersion, []string{patch1, patch2})
	// NOTE: this was added as a need by us changing to dbfs.PullFile
	di.FileWrite(file.RelativePath, file.Filename, file.ProjectID, []byte{})

	changes, _, pulledVersion, _, err := di.PullChanges(file)
	assert.Equal(t, originalFileVersion, pulledVersion, "failed set up verification")

	transformed, version, missing, lenChanges, err := di.CBAppendFileChange(file, []string{patch3})
	assert.NoError(t, err, "unexpected error appending changes")
	assert.Empty(t, missing, "Unexpected missing patches")

	assert.Equal(t, originalFileVersion+1, version, "version did not update properly")

	raw, changes, err := di.PullFile(file)
	assert.NoError(t, err, "unexpected error getting changes")

	assert.Empty(t, *raw, "we shouldn't have scrunched")
	assert.Len(t, changes, 3, "resultant changes not the correct length")
	assert.Equal(t, len(changes), lenChanges, "did not return correct change count")

	assert.Equal(t, patch1, changes[0], "first change was not correct")
	assert.Equal(t, patch2, changes[1], "second change was not correct")

	assert.Len(t, transformed, 1, "returned unexpected number of transformed new changes")
	assert.EqualValues(t, transformed[0], changes[2], "newly inserted change was not correct")

	// Expect AppendFileChange to transform patch4, since it was based on the version created by patch2
	changes, _, pulledVersion, _, err = di.PullChanges(file)
	assert.Equal(t, pulledVersion, version, "version pulled from the database does not match the one given when appending the change")

	transformed, version, missing, lenChanges, err = di.CBAppendFileChange(file, []string{patch4})
	assert.NoError(t, err, "unexpected error appending changes")

	assert.Len(t, missing, 1, "Unexpected number of missing patches")
	assert.Contains(t, missing, patch3, "Unexpected missing patches")

	assert.Equal(t, originalFileVersion+2, version, "version did not update properly")

	raw, changes, err = di.PullFile(file)
	assert.NoError(t, err, "unexpected error getting changes")

	assert.Empty(t, *raw, "we shouldn't have scrunched")
	assert.Len(t, changes, 4, "resultant changes not the correct length")
	assert.Equal(t, len(changes), lenChanges, "did not return correct change count")

	assert.Equal(t, patch1, changes[0], "first change was not correct")
	assert.Equal(t, patch2, changes[1], "second change was not correct")
	assert.Equal(t, patch3, changes[2], "third change was not correct")

	assert.Len(t, transformed, 1, "returned unexpected number of transformed new changes")
	assert.EqualValues(t, transformed[0], changes[3], "newly inserted change was not correct")

	ver, err := di.CBGetFileVersion(file.FileID)
	assert.EqualValues(t, 4, ver, "wrong file version")

	di.CBDeleteFile(file.FileID)
	di.FileDelete(file.RelativePath, file.Filename, file.ProjectID)
}
