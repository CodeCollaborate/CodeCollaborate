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
	configSetup(t)
	di := new(DatabaseImpl)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)

	projectParentPath := filepath.Clean(config.GetConfig().ServerConfig.ProjectPath)
	filepath1 := filepath.Join(projectParentPath, "10", "myFile1.txt")
	filepath2 := filepath.Join(projectParentPath, "10", "hi", "myFile2.txt")

	fileText := []byte("Hello World!\nWelcome to my file\n")

	// in case the test fails before resolving
	defer os.Remove(projectParentPath)
	defer os.Remove(filepath.Join(projectParentPath, "10"))
	defer os.Remove(filepath.Join(projectParentPath, "10", "hi"))
	defer os.Remove(filepath2)
	defer os.Remove(filepath1)

	loc, err := di.FileWrite(".", "myFile1.txt", 10, fileText)
	if err != nil {
		t.Fatal(err)
	}
	if loc != filepath1 {
		t.Fatalf("wrong file location\nexpected:\n%v\nactual:\n%v", filepath1, loc)
	}
	loc, err = di.FileWrite("./hi/", "myFile2.txt", 10, fileText)
	if err != nil {
		t.Fatal(err)
	}
	if loc != filepath2 {
		t.Fatalf("wrong file location\nexpected:\n%v\nactual:\n%v", filepath2, loc)
	}

	// Test a bad path
	_, err = di.FileWrite("..", "myFile.txt", 10, fileText)
	if err != ErrMaliciousRequest {
		t.Fatal("Expected failure to write to bad location")
	}
	// Test a worse but hidden path
	_, err = di.FileWrite("fake/../../../", "myFile.txt", 10, fileText)
	if err != ErrMaliciousRequest {
		t.Fatal("Expected failure to write to bad location")
	}
	// Test with a bad filename
	//_, err = di.FileWrite(".", "../myFile.txt", 10, fileText)
	//if err != ErrMaliciousRequest {
	//	t.Fatal("Expected failure to write to bad location")
	//}

	// check file exists
	if _, err := os.Stat(filepath1); os.IsNotExist(err) {
		t.Fatal(os.ErrNotExist)
	}
	if _, err := os.Stat(filepath2); os.IsNotExist(err) {
		t.Fatal(os.ErrNotExist)
	}

	writtenBytes, _ := ioutil.ReadFile(filepath1)
	if !bytes.Equal(fileText, writtenBytes) {
		t.Fatal("File was not the same")
	}

}

func TestDatabaseImpl_FileRead(t *testing.T) {
	configSetup(t)
	di := new(DatabaseImpl)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)

	projectParentPath := filepath.Clean(config.GetConfig().ServerConfig.ProjectPath)
	filepath1 := filepath.Join(projectParentPath, "10", "myFile1.txt")

	fileText := []byte("Hello World!\nWelcome to my file\n")

	defer os.Remove(projectParentPath)
	defer os.Remove(filepath.Join(projectParentPath, "10"))
	defer os.Remove(filepath1)

	_, err := di.FileWrite(".", "myFile1.txt", 10, fileText)
	if err != nil {
		t.Fatal(err)
	}

	data, err := di.FileRead(".", "myFile1.txt", 10)

	if !bytes.Equal(fileText, *data) {
		t.Fatalf("File was not writen or read correctly\nExpected:\n%v\nActual:\n%v", fileText, data)
	}

}

func TestDatabaseImpl_FileDelete(t *testing.T) {
	configSetup(t)
	di := new(DatabaseImpl)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)

	projectParentPath := filepath.Clean(config.GetConfig().ServerConfig.ProjectPath)
	filepath1 := filepath.Join(projectParentPath, "10", "myFile1.txt")

	fileText := []byte("Hello World!\nWelcome to my file\n")

	defer os.Remove(projectParentPath)
	defer os.Remove(filepath.Join(projectParentPath, "10"))

	_, err := di.FileWrite(".", "myFile1.txt", 10, fileText)
	if err != nil {
		t.Fatal(err)
	}

	err = di.FileDelete(".", "myFile1.txt", 10)
	if err != nil {
		t.Fatal(err)
	}

	if err = os.Remove(filepath1); !os.IsNotExist(err) {
		t.Fatal("File should have been deleted, but was not")
	}

}

func TestDatabaseImpl_FileMove(t *testing.T) {
	configSetup(t)
	di := new(DatabaseImpl)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)

	projectParentPath := filepath.Clean(config.GetConfig().ServerConfig.ProjectPath)
	filepath1 := filepath.Join(projectParentPath, "10", "myFile1.txt")
	filepath2 := filepath.Join(projectParentPath, "10", filepath.Join("newdir", "myFile2.txt"))
	fileText := []byte("Hello World!\nWelcome to my file\n")

	defer os.Remove(projectParentPath)
	defer os.Remove(filepath.Join(projectParentPath, "10"))
	//defer os.Remove(filepath1)

	err := os.MkdirAll(projectParentPath, 0744)
	if err != nil {
		t.Fatal(err)
	}

	_, err = di.FileWrite(".", "myFile1.txt", 10, fileText)
	if err != nil {
		t.Fatal(err)
	}

	err = di.FileMove(".", "myFile1.txt", "newdir", "myFile2.txt", 10)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(filepath2)
	if err != nil {
		t.Fatal("file was not moved")
	}
	err = os.Remove(filepath1)
	if err == nil {
		t.Fatal("file was not moved")
	}

}

var fileText = []byte("Hello World!\nWelcome to my file\n")
var file = FileMeta{
	FileID:       0,
	RelativePath: "./",
	Filename:     "_test_name",
	ProjectID:    0,
}

func setupFileWithSwap(t *testing.T, di *DatabaseImpl) (string, []byte) {
	filePath, err := di.FileWrite(file.RelativePath, file.Filename, file.ProjectID, fileText)
	assert.NoError(t, err, "error initially writing file")

	_, err = os.Stat(filePath)
	assert.False(t, os.IsNotExist(err), "original file does not exist")

	// make swap file
	raw, err := di.makeSwp(file.RelativePath, file.Filename, file.ProjectID)
	assert.NoError(t, err, "error creating swap file")
	assert.EqualValues(t, fileText, raw, "swap file and file are not equal")

	return filePath, raw
}

func TestDatabaseImpl_FileWriteToSwap(t *testing.T) {
	configSetup(t)
	di := new(DatabaseImpl)

	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)

	filePath, raw := setupFileWithSwap(t, di)
	defer os.Remove(filePath)

	// test getSwpLocation (and allow defer-ed deletion)
	swpLoc := di.getSwpLocation(filePath)
	_, err := os.Stat(swpLoc)
	assert.False(t, os.IsNotExist(err), "swap file does not exist")
	defer os.Remove(swpLoc)

	// test swap read
	swp, err := di.swapRead(file.RelativePath, file.Filename, file.ProjectID)
	assert.NoError(t, err, "error reading swp file")
	assert.Equal(t, raw, *swp, "swap incorrectly changed")

	// test swap write
	newRawFile := []byte(string(fileText) + "it's a pretty cool file, not going to lie\n")

	err = di.FileWriteToSwap(file, newRawFile)
	assert.NoError(t, err, "error writing to swap")

	swp, err = di.swapRead(file.RelativePath, file.Filename, file.ProjectID)
	assert.NoError(t, err, "error reading swp file")
	assert.Equal(t, newRawFile, *swp, "swap incorrectly changed")
}

func TestDatabaseImpl_FileSwapSwap(t *testing.T) {
	configSetup(t)
	di := new(DatabaseImpl)
	defer os.RemoveAll(config.GetConfig().ServerConfig.ProjectPath)

	filePath, _ := setupFileWithSwap(t, di)
	defer os.Remove(filePath)

	swpLoc := di.getSwpLocation(filePath)
	defer os.Remove(swpLoc)

	raw, err := di.FileRead(file.RelativePath, file.Filename, file.ProjectID)
	assert.NoError(t, err, "error reading file")
	assert.EqualValues(t, fileText, *raw, "swap file was not swapped")

	// test swap write
	newRawFile := []byte(string(fileText) + "it's a pretty cool file, not going to lie\n")
	err = di.FileWriteToSwap(file, newRawFile)
	assert.NoError(t, err, "error writing to swap")

	err = di.swapSwp(file.RelativePath, file.Filename, file.ProjectID)
	assert.NoError(t, err, "error swapping swap")

	raw, err = di.FileRead(file.RelativePath, file.Filename, file.ProjectID)
	assert.NoError(t, err, "error reading file")
	assert.EqualValues(t, newRawFile, *raw, "swap file was not swapped")

	_, err = os.Stat(swpLoc)
	assert.False(t, os.IsNotExist(err), "swap does not exists")

	err = di.deleteSwp(file.RelativePath, file.Filename, file.ProjectID)
	assert.NoError(t, err, "error deleting swap file")
	_, err = os.Stat(swpLoc)
	assert.True(t, os.IsNotExist(err), "swap does still exists")
}
