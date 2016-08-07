package dbfs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/CodeCollaborate/Server/modules/config"
)

var filePathSeparator = strconv.QuoteRune(os.PathSeparator)[1:2]

// FileWrite writes the file with the given bytes to a calculated path, and
// returns that path so it can be put in MySQL
func (di *DatabaseImpl) FileWrite(relpath string, filename string, projectID int64, raw []byte) (string, error) {
	relFilePath, err := calculateFilePathPath(relpath, filename, projectID)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(relFilePath, 0744)
	if err != nil {
		return "", err
	}
	fileLocation := filepath.Join(relFilePath, filename)
	err = ioutil.WriteFile(fileLocation, raw, 0744)
	if err != nil {
		return "", err
	}

	return fileLocation, err
}

// FileDelete deletes the file with the given metadata from the file system
// Couple this with dbfs.MySQLFileDelete and dbfs.CBDeleteFile
func (di *DatabaseImpl) FileDelete(relpath string, filename string, projectID int64) error {
	relFilePath, err := calculateFilePathPath(relpath, filename, projectID)
	if err != nil {
		return err
	}
	fileLocation := filepath.Join(relFilePath, filename)
	return os.Remove(fileLocation)
}

// FileRead returns the project file from the calculated location on the disk
func (di *DatabaseImpl) FileRead(relpath string, filename string, projectID int64) (*[]byte, error) {
	relFilePath, err := calculateFilePathPath(relpath, filename, projectID)
	if err != nil {
		return new([]byte), err
	}
	fileLocation := filepath.Join(relFilePath, filename)
	fileBytes, err := ioutil.ReadFile(fileLocation)
	return &fileBytes, err
}

// FileMove moves a file form the starting path to the end path
func (di *DatabaseImpl) FileMove(startRelpath string, startFilename string, endRelpath string, endFilename string, projectID int64) error {
	startRelFilePath, err := calculateFilePathPath(startRelpath, startFilename, projectID)
	if err != nil {
		return err
	}
	endRelFilePath, err := calculateFilePathPath(endRelpath, endFilename, projectID)
	if err != nil {
		return err
	}
	err = os.MkdirAll(endRelFilePath, 0744)
	if err != nil {
		return err
	}

	startFileLocation := filepath.Join(startRelFilePath, startFilename)
	endFileLocation := filepath.Join(endRelFilePath, endFilename)

	err = os.Rename(startFileLocation, endFileLocation)
	return err
}

func calculateFilePathPath(relpath string, filename string, projectID int64) (string, error) {
	if strings.Contains(filename, filePathSeparator) {
		return "", ErrMaliciousRequest
	}
	cleanPath := filepath.Clean(relpath)
	if strings.HasPrefix(cleanPath, "..") {
		return "", ErrMaliciousRequest
	}

	projectFolderParentPath := config.GetConfig().ServerConfig.ProjectPath
	return filepath.Join(projectFolderParentPath, strconv.FormatInt(projectID, 10), cleanPath), nil
}
