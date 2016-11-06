package dbfs

import (
	"testing"
	"time"

	"os"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/stretchr/testify/assert"
)

var baseChanges = []string{"I swear", "I'm not lying"}
var rawBaseFile = "this is a very important file"

func setupFile(t *testing.T) (*DatabaseImpl, FileMeta) {
	configSetup(t)
	di := new(DatabaseImpl)

	file := FileMeta{
		Creator:      "_testuser1",
		CreationDate: time.Now(),
		RelativePath: "./",
		Filename:     "_test_file_123",
	}

	di.MySQLUserDelete(userOne.Username)

	err := di.MySQLUserRegister(userOne)
	assert.NoError(t, err, "error registering mysql user")

	file.ProjectID, err = di.MySQLProjectCreate(userOne.Username, "_test_project_1")
	assert.NoError(t, err, "error creating mysql project")

	file.FileID, err = di.MySQLFileCreate(userOne.Username, file.Filename, file.RelativePath, file.ProjectID)
	assert.NoError(t, err, "error creating mysql file")

	err = di.CBInsertNewFile(file.FileID, 0, []string{})
	assert.NoError(t, err, "error inserting file to couchbase")

	_, err = di.FileWrite(file.RelativePath, file.Filename, file.ProjectID, []byte(rawBaseFile))
	assert.NoError(t, err, "error writing file to disk")

	// TODO: make these actual patches
	_, err = di.CBAppendFileChange(file.FileID, 0, baseChanges)
	assert.NoError(t, err, "error appending change to file")

	return di, file
}

func TestDatabaseImpl_PullFile(t *testing.T) {
	// check normal pull (no scrunching)
	di, file := setupFile(t)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)
	defer di.CBDeleteFile(file.FileID)
	defer di.FileDelete(file.RelativePath, file.Filename, file.ProjectID)
	defer di.MySQLUserDelete(userOne.Username)

	raw, changes1, err := di.PullFile(file)
	assert.NoError(t, err, "Error while pulling file")
	assert.Len(t, changes1, 2, "incorrect number of changes returned from couchbase")
	assert.Contains(t, changes1, baseChanges[0], "file did on contain expected change")
	assert.Contains(t, changes1, baseChanges[1], "file did on contain expected change")
	assert.EqualValues(t, string(*raw), rawBaseFile, "raw file did not match")
}

func TestDatabaseImpl_GetForScrunching(t *testing.T) {
	di, file := setupFile(t)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)
	defer di.CBDeleteFile(file.FileID)
	defer di.FileDelete(file.RelativePath, file.Filename, file.ProjectID)
	defer di.MySQLUserDelete(userOne.Username)

	changes, swp, err := di.GetForScrunching(file.FileID, 1)
	assert.NoError(t, err, "error getting swp or changes")

	assert.Len(t, changes, 1, "changes size was an unexpected length")
	assert.Contains(t, changes, baseChanges[0], "changes didn't contain correct change")

	assert.EqualValues(t, swp, rawBaseFile, "swp file was not cloned properly")

	err = di.deleteSwp(file.RelativePath, file.Filename, file.ProjectID)
	assert.NoError(t, err, "error deletig swp file")
}

// TODO: check delete

// TODO: check mid-delete pulling
