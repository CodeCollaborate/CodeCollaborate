package dbfs

import (
	"os"
	"strconv"
	"testing"
	"time"

	"bytes"
	"fmt"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/stretchr/testify/assert"
)

var defaultBaseFile = "this is a very important file"
var defaultChanges = []string{"v0:\n1:+5:test1", "v0:\n10:+5:test2"}
var transformedChanges = []string{"v0:\n1:+5:test1", "v1:\n10:+5:test2"}

func setupFile(t *testing.T, baseFile string, baseChanges []string) (*DatabaseImpl, FileMeta) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	file := FileMeta{
		Creator:      "_testuser1",
		CreationDate: time.Now(),
		RelativePath: "./",
		Filename:     "_test_file_123",
		ProjectID:    0,
		FileID:       0,
	}

	err := di.CBInsertNewFile(file.FileID, 0, []string{})
	assert.NoError(t, err, "error inserting file to couchbase")

	_, err = di.FileWrite(file.RelativePath, file.Filename, file.ProjectID, []byte(baseFile))
	assert.NoError(t, err, "error writing file to disk")

	// TODO: make these actual patches
	_, _, err = di.CBAppendFileChange(file.FileID, 0, baseChanges, []string{})
	assert.NoError(t, err, "error appending change to file")

	return di, file
}

func TestDatabaseImpl_PullFile(t *testing.T) {
	// check normal pull (no scrunching)
	di, file := setupFile(t, defaultBaseFile, defaultChanges)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)
	defer di.CBDeleteFile(file.FileID)

	checkPullFile(t, di, file, transformedChanges, defaultBaseFile)
}

func TestDatabaseImpl_ScrunchFile(t *testing.T) {
	MinBufferLength = 5
	MaxBufferLength = 30
	patches := make([]string, 50)
	resultPatches := make([]string, 5)
	expectedOutput := bytes.Buffer{}

	for i := 0; i < 50; i++ {
		if i < 10 {
			patches[i] = fmt.Sprintf("v%d:\n2:+1:%d", 0, i)
		} else {
			patches[i] = fmt.Sprintf("v%d:\n2:+2:%d", 0, i)
		}
	}

	for i := 0; i < 5; i++ {
		resultPatches[i] = fmt.Sprintf("v%d:\n2:+2:%d", i+45, i+45)
	}

	expectedOutput.WriteString("te")
	for i := len(patches) - MinBufferLength - 1; i >= 0; i-- {
		if i < len(patches)-MinBufferLength {
			expectedOutput.WriteString(fmt.Sprintf("%d", i))
		}
	}
	expectedOutput.WriteString("st")

	di, file := setupFile(t, "test", patches)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)
	defer di.CBDeleteFile(file.FileID)

	err := di.ScrunchFile(file)
	assert.NoError(t, err, "error getting swp or changes")

	fileBytes, changes, err := di.PullFile(file)
	assert.NoError(t, err, "error pulling file")

	assert.Len(t, changes, MinBufferLength, "changes size was an unexpected length")
	assert.Equal(t, resultPatches, changes, "changes didn't contain correct changes")

	assert.EqualValues(t, expectedOutput.String(), string(*fileBytes), "Scrunched file differed from expected output")
}

func TestDatabaseImpl_GetForScrunching(t *testing.T) {
	di, file := setupFile(t, defaultBaseFile, defaultChanges)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)
	defer di.CBDeleteFile(file.FileID)

	changes, swp, err := di.getForScrunching(file, 1)
	assert.NoError(t, err, "error getting swp or changes")

	assert.Len(t, changes, 1, "changes size was an unexpected length")
	assert.Contains(t, changes, defaultChanges[0], "changes didn't contain correct change")

	assert.EqualValues(t, string(swp), string(defaultBaseFile), "swp file was not cloned properly")

	err = di.deleteSwp(file.RelativePath, file.Filename, file.ProjectID)
	assert.NoError(t, err, "error deleting swp file")
}

func TestDatabaseImpl_DeleteForScrunching(t *testing.T) {
	di, file := setupFile(t, defaultBaseFile, defaultChanges)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)
	defer di.CBDeleteFile(file.FileID)

	// note that this is totally different from what would normally be made from scrunching
	newRawFile := []byte(string(fileText) + "it's a pretty cool file, not going to lie\n")

	err := di.FileWriteToSwap(file, newRawFile)
	assert.NoError(t, err, "Error while writing to swap file")

	di.deleteForScrunching(file, 1)

	raw, changesNew, err := di.PullFile(file)
	assert.NoError(t, err, "Error while pulling file")
	assert.Len(t, changesNew, 1, "incorrect number of changes returned from couchbase")
	assert.Contains(t, changesNew, transformedChanges[1], "file did on contain expected change")
	assert.EqualValues(t, newRawFile, string(*raw), "raw file did not match")
}

func TestDatabaseImpl_PullFile_MidDelete(t *testing.T) {
	di, file := setupFile(t, defaultBaseFile, defaultChanges)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)
	defer di.CBDeleteFile(file.FileID)

	newChanges := []string{"v2:\n2:+1:2", "v2:\n2:+1:3", "v4:\n2:+1:4", "v5:\n2:+1:4", "v6:\n2:+1:5", "v7:\n2:+1:6", "v8:\n2:+1:7", "v9:\n2:+1:8", "v9:\n2:+1:9", "v9:\n2:+2:10"}
	transformedNewChanges := []string{"v2:\n2:+1:2", "v3:\n2:+1:3", "v4:\n2:+1:4", "v5:\n2:+1:4", "v6:\n2:+1:5", "v7:\n2:+1:6", "v8:\n2:+1:7", "v9:\n2:+1:8", "v10:\n2:+1:9", "v11:\n2:+2:10"}
	newRawFile := []byte(string(defaultBaseFile) + "\nit's a pretty cool file, not going to lie\n")

	checkPullFile(t, di, file, transformedChanges, defaultBaseFile)

	// add more changes so it's more visible
	appendChangeToFile(t, di, newChanges[:2])

	checkPullFile(t, di, file, append(transformedChanges, transformedNewChanges[:2]...), defaultBaseFile)

	// arbitrarily saying we're going to scrunch off 2 patches
	num := 2

	// make sure they're right
	changes1, raw1, err := di.getForScrunching(file, num)
	assert.NoError(t, err, "error getting changes for scrunching")
	assert.EqualValues(t, string(defaultBaseFile), string(raw1), "swap was not made correctly")
	assert.Len(t, changes1, num, "pulled wrong number of changes")
	assert.EqualValues(t, transformedChanges, changes1, "changes given for scrunching were not correct")

	// update swap
	err = di.FileWriteToSwap(file, newRawFile)
	assert.NoError(t, err, "Error while writing to swap file")

	// check pull file (expecting old + new changes w/ old base)
	checkPullFile(t, di, file, append(transformedChanges, transformedNewChanges[:2]...), string(defaultBaseFile))

	// START DELETE
	cb, err := di.openCouchBase()
	nativeErr(t, err)

	key := strconv.FormatInt(file.FileID, 10)

	// turn on writing to TempChanges
	builder := cb.bucket.MutateIn(key, 0, 0)
	builder = builder.Upsert("tempchanges", []string{}, false)
	builder = builder.Upsert("usetemp", true, false)
	_, err = builder.Execute()
	nativeErr(t, err)

	// add change
	appendChangeToFile(t, di, newChanges[2:3])
	checkPullFile(t, di, file, append(transformedChanges, transformedNewChanges[:3]...), string(defaultBaseFile))

	// get changes in normal changes
	frag, err := cb.bucket.LookupIn(key).Get("changes").Execute()
	nativeErr(t, err)

	changes := []string{}
	err = frag.Content("changes", &changes)
	nativeErr(t, err)

	// add change
	appendChangeToFile(t, di, newChanges[3:4])
	checkPullFile(t, di, file, append(transformedChanges, transformedNewChanges[:4]...), string(defaultBaseFile))

	// turn off writing to TempChanges & reset normal changes
	builder = cb.bucket.MutateIn(key, 0, 0)
	builder = builder.Upsert("remaining_changes", changes[num:], false)
	builder = builder.Upsert("changes", []string{}, false)
	builder = builder.Upsert("usetemp", false, false)
	builder = builder.Upsert("pullswp", true, false)
	_, err = builder.Execute()
	nativeErr(t, err)

	// add change
	// check switched to swap
	appendChangeToFile(t, di, newChanges[4:5])
	checkPullFile(t, di, file, transformedNewChanges[:5], string(newRawFile))

	// get changes in TempChanges
	frag, err = cb.bucket.LookupIn(key).Get("tempchanges").Execute()
	nativeErr(t, err)

	tempChanges := []string{}
	err = frag.Content("tempchanges", &tempChanges)
	nativeErr(t, err)

	// add change
	// check switched to swap
	appendChangeToFile(t, di, newChanges[5:6])
	checkPullFile(t, di, file, transformedNewChanges[:6], string(newRawFile))

	err = di.swapSwp(file.RelativePath, file.Filename, file.ProjectID)
	assert.NoError(t, err, "Error swapping swap file, NOTE: the server WOULD normally be able to recover from here")

	// add change
	appendChangeToFile(t, di, newChanges[6:7])
	checkPullFile(t, di, file, transformedNewChanges[:7], string(newRawFile))

	// prepend changes and reset temporarily stored changes
	builder = cb.bucket.MutateIn(key, 0, 0)
	builder = builder.ArrayPrependMulti("changes", append(changes[num:], tempChanges...), false)
	builder = builder.Upsert("remaining_changes", []string{}, false)
	builder = builder.Upsert("tempchanges", []string{}, false)
	builder = builder.Upsert("pullswp", false, false)
	_, err = builder.Execute()
	nativeErr(t, err)

	err = di.deleteSwp(file.RelativePath, file.Filename, file.ProjectID)
	assert.NoError(t, err, "Error deleting swap file, NOTE: the server WOULD normally be able to recover from here")

	// add change
	appendChangeToFile(t, di, newChanges[7:])
	checkPullFile(t, di, file, transformedNewChanges, string(newRawFile))
}

func nativeErr(t *testing.T, err error) {
	assert.NoError(t, err, "error in naitive di.DeleteForScrunching code")
}

func appendChangeToFile(t *testing.T, di *DatabaseImpl, changes []string) {
	ver, err := di.CBGetFileVersion(file.FileID)
	assert.NoError(t, err, "Error getting file version")
	_, _, err = di.CBAppendFileChange(file.FileID, ver, changes, []string{})
	assert.NoError(t, err, "Error while appending more changes")
}

func checkPullFile(t *testing.T, di *DatabaseImpl, testFile FileMeta, expectedChanges []string, expectedRaw string) {
	raw, changesNew, err := di.PullFile(testFile)
	assert.NoError(t, err, "Error while pulling file")
	assert.Len(t, changesNew, len(expectedChanges), "incorrect number of changes returned from couchbase")
	assert.EqualValues(t, expectedChanges, changesNew, "file did on contain expected changes")
	assert.EqualValues(t, expectedRaw, string(*raw), "raw file did not match")
}
