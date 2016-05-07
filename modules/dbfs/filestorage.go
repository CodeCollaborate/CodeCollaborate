package dbfs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/CodeCollaborate/Server/modules/config"
)

var pathsep = strconv.QuoteRune(os.PathSeparator)[1:2]

// FileWrite writes the file with the given bytes to a calculated path, and
// returns that path so it can be put in MySQL
func FileWrite(relpath string, filename string, projectID int64, raw []byte) (string, error) {
	relFilePath, err := calculateAndValidatePath(relpath, filename, projectID)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(relFilePath, 0744)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(relFilePath+pathsep+filename, raw, 0744)
	if err != nil {
		return "", err
	}

	return relFilePath + pathsep + filename, err
}

// FileRead returns the project file from the calculated location on the disk
func FileRead(relpath string, filename string, projectID int64) (*[]byte, error) {
	relFilePath, err := calculateAndValidatePath(relpath, filename, projectID)
	if err != nil {
		return new([]byte), err
	}
	fileBytes, err := ioutil.ReadFile(relFilePath + pathsep + filename)
	return &fileBytes, err
}

func calculateAndValidatePath(relpath string, filename string, projectID int64) (string, error) {
	if strings.Contains(filename, pathsep) {
		return "", ErrMalliciousRequest
	}
	cleanPath := filepath.Clean(relpath)
	if strings.HasPrefix(cleanPath, "..") {
		return "", ErrMalliciousRequest
	}

	projectFolderParentPath := filepath.Clean(config.GetConfig().ServerConfig.ProjectPath)
	projectFolderPath := projectFolderParentPath + pathsep + strconv.FormatInt(projectID, 10)
	return filepath.Clean(projectFolderPath + pathsep + cleanPath), nil
}
