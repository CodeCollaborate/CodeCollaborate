package dbfs

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/CodeCollaborate/Server/modules/config"
)

func TestFileWrite(t *testing.T) {
	configSetup()
	projectParentPath := filepath.Clean(config.GetConfig().ServerConfig.ProjectPath)
	filepath1 := filepath.Join(projectParentPath, "10", "myFile1.txt")
	filepath2 := filepath.Join(projectParentPath, "10", "hi", "myFile2.txt")

	fileText := []byte("hello\nWelcome to my file\n")

	loc, err := FileWrite(".", "myFile1.txt", 10, fileText)
	if err != nil {
		t.Fatal(err)
	}
	if loc != filepath1 {
		t.Fatalf("wrong file location\nexpected:\n%v\nactual:\n%v", filepath1, loc)
	}
	loc, err = FileWrite("./hi/", "myFile2.txt", 10, fileText)
	if err != nil {
		t.Fatal(err)
	}
	if loc != filepath2 {
		t.Fatalf("wrong file location\nexpected:\n%v\nactual:\n%v", filepath2, loc)
	}

	// Test a bad path
	_, err = FileWrite("..", "myFile.txt", 10, fileText)
	if err != ErrMaliciousRequest {
		t.Fatal("Expected failure to write to bad location")
	}
	// Test a worse but hidden path
	_, err = FileWrite("fake/../../../", "myFile.txt", 10, fileText)
	if err != ErrMaliciousRequest {
		t.Fatal("Expected failure to write to bad location")
	}
	// Test with a bad filename
	_, err = FileWrite(".", "../myFile.txt", 10, fileText)
	if err != ErrMaliciousRequest {
		t.Fatal("Expected failure to write to bad location")
	}

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

	os.Remove(filepath2)
	os.Remove(filepath1)
	os.Remove(filepath.Join(projectParentPath, "10", "hi/"))
	os.Remove(filepath.Join(projectParentPath, "10"))
	os.Remove(projectParentPath)
}

func TestFileRead(t *testing.T) {
	configSetup()
	projectParentPath := filepath.Clean(config.GetConfig().ServerConfig.ProjectPath)
	filepath1 := filepath.Join(projectParentPath, "10", "myFile1.txt")

	fileText := []byte("hello\nWelcome to my file\n")

	defer os.Remove(projectParentPath)
	defer os.Remove(filepath.Join(projectParentPath, "10"))
	defer os.Remove(filepath1)

	_, err := FileWrite(".", "myFile1.txt", 10, fileText)
	if err != nil {
		t.Fatal(err)
	}

	data, err := FileRead(".", "myFile1.txt", 10)

	if !bytes.Equal(fileText, *data) {
		t.Fatalf("File was not writen or read correctly\nExpected:\n%v\nActual:\n%v", fileText, data)
	}

}
