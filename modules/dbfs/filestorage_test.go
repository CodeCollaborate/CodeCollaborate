package dbfs

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseImpl_FileWrite(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)

	projectParentPath := filepath.Clean(config.GetConfig().ServerConfig.ProjectPath)
	relPath1 := filepath.Join(projectParentPath, "10")
	relPath2 := filepath.Join(projectParentPath, "12")

	defer os.Remove(projectParentPath)

	fileText := []byte("Hello World!\nWelcome to my file\n")

	err := di.FileWrite(10, fileText)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(relPath1)

	err = di.FileWrite(12, fileText)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(relPath2)

	// check file exists
	if _, err := os.Stat(relPath1); os.IsNotExist(err) {
		t.Fatal(os.ErrNotExist)
	}
	if _, err := os.Stat(relPath2); os.IsNotExist(err) {
		t.Fatal(os.ErrNotExist)
	}

	writtenBytes, _ := ioutil.ReadFile(relPath1)
	if !bytes.Equal(fileText, writtenBytes) {
		t.Fatal("File was not the same")
	}

}

func TestDatabaseImpl_FileRead(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)

	projectParentPath := filepath.Clean(config.GetConfig().ServerConfig.ProjectPath)

	fileText := []byte("Hello World!\nWelcome to my file\n")

	defer os.Remove(projectParentPath)
	defer os.Remove(filepath.Join(projectParentPath, "10"))

	err := di.FileWrite(10, fileText)
	if err != nil {
		t.Fatal(err)
	}

	data, err := di.FileRead(10)

	if !bytes.Equal(fileText, data) {
		t.Fatalf("File was not writen or read correctly\nExpected:\n%v\nActual:\n%v", fileText, data)
	}

}

func TestDatabaseImpl_FileDelete(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)

	projectParentPath := filepath.Clean(config.GetConfig().ServerConfig.ProjectPath)
	filepath1 := filepath.Join(projectParentPath, "10", "myFile1.txt")

	fileText := []byte("Hello World!\nWelcome to my file\n")

	defer os.Remove(projectParentPath)
	defer os.Remove(filepath.Join(projectParentPath, "10"))

	err := di.FileWrite(10, fileText)
	if err != nil {
		t.Fatal(err)
	}

	err = di.FileDelete(10)
	if err != nil {
		t.Fatal(err)
	}

	if err = os.Remove(filepath1); !os.IsNotExist(err) {
		t.Fatal("File should have been deleted, but was not")
	}

}

var fileText = []byte("Hello World!\nWelcome to my file\n")
var file = FileMeta{
	FileID:       1,
	RelativePath: "./",
	Filename:     "_test_name",
	ProjectID:    0,
}
var filePath string
var swpFilePath string

func setupFileWithSwap(t *testing.T, di *DatabaseImpl) (string, []byte) {
	filePath = filepath.Join(filepath.Clean(config.GetConfig().ServerConfig.ProjectPath), "1")
	swpFilePath = filepath.Join(filepath.Clean(config.GetConfig().ServerConfig.ProjectPath), "-1")

	err := di.FileWrite(file.FileID, fileText)
	assert.NoError(t, err, "error initially writing file")

	_, err = os.Stat(filePath)
	assert.False(t, os.IsNotExist(err), "original file does not exist")

	// make swap file
	raw, err := di.makeSwp(file.FileID)
	assert.NoError(t, err, "error creating swap file")
	assert.EqualValues(t, fileText, raw, "swap file and file are not equal")

	return filePath, raw
}

func TestDatabaseImpl_FileWriteToSwap(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)

	filePath, raw := setupFileWithSwap(t, di)
	defer os.Remove(filePath)

	// test swap read
	swp, err := di.swapRead(file.FileID)
	assert.NoError(t, err, "error reading swp file")
	assert.Equal(t, raw, swp, "swap incorrectly changed")

	// test swap write
	newRawFile := []byte(string(fileText) + "it's a pretty cool file, not going to lie\n")

	err = di.FileWriteToSwap(file.FileID, newRawFile)
	assert.NoError(t, err, "error writing to swap")

	swp, err = di.swapRead(file.FileID)
	assert.NoError(t, err, "error reading swp file")
	assert.Equal(t, newRawFile, swp, "swap incorrectly changed")
}

func TestDatabaseImpl_FileSwapSwap(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)
	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)

	filePath, _ := setupFileWithSwap(t, di)
	defer os.Remove(filePath)

	raw, err := di.FileRead(file.FileID)
	assert.NoError(t, err, "error reading file")
	assert.EqualValues(t, fileText, raw, "swap file was not swapped")

	// test swap write
	newRawFile := []byte(string(fileText) + "it's a pretty cool file, not going to lie\n")
	err = di.FileWriteToSwap(file.FileID, newRawFile)
	assert.NoError(t, err, "error writing to swap")

	err = di.swapSwp(file.FileID)
	assert.NoError(t, err, "error swapping swap")

	raw, err = di.FileRead(file.FileID)
	assert.NoError(t, err, "error reading file")
	assert.EqualValues(t, newRawFile, raw, "swap file was not swapped")

	_, err = os.Stat(swpFilePath)
	assert.False(t, os.IsNotExist(err), "swap does not exists")

	err = di.deleteSwp(file.FileID)
	assert.NoError(t, err, "error deleting swap file")
	_, err = os.Stat(swpFilePath)
	assert.True(t, os.IsNotExist(err), "swap does still exists")
}
