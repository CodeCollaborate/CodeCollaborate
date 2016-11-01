package dbfs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"errors"
	"io"

	"github.com/CodeCollaborate/Server/modules/config"
)

var filePathSeparator = strconv.QuoteRune(os.PathSeparator)[1:2]

// FileWrite writes the file with the given bytes to a calculated path, and
// returns that path so it can be put in MySQL
func (di *DatabaseImpl) FileWrite(relpath string, filename string, projectID int64, raw []byte) (string, error) {
	relFilePath, err := calculateFilePath(relpath, filename, projectID)
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
	relFilePath, err := calculateFilePath(relpath, filename, projectID)
	if err != nil {
		return err
	}
	fileLocation := filepath.Join(relFilePath, filename)
	return os.Remove(fileLocation)
}

// FileRead returns the project file from the calculated location on the disk
func (di *DatabaseImpl) FileRead(relpath string, filename string, projectID int64) (*[]byte, error) {
	relFilePath, err := calculateFilePath(relpath, filename, projectID)
	if err != nil {
		return new([]byte), err
	}
	fileLocation := filepath.Join(relFilePath, filename)
	fileBytes, err := ioutil.ReadFile(fileLocation)
	return &fileBytes, err
}

// FileMove moves a file form the starting path to the end path
func (di *DatabaseImpl) FileMove(startRelpath string, startFilename string, endRelpath string, endFilename string, projectID int64) error {
	startRelFilePath, err := calculateFilePath(startRelpath, startFilename, projectID)
	if err != nil {
		return err
	}
	endRelFilePath, err := calculateFilePath(endRelpath, endFilename, projectID)
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

// returns the swap file path and any error
func (di *DatabaseImpl) makeSwp(relpath string, filename string, projectID int64) ([]byte, error) {
	relFilePath, err := calculateFilePath(relpath, filename, projectID)
	if err != nil {
		return []byte{}, err
	}
	fileLocation := filepath.Join(relFilePath, filename)
	swapLoc := getSwpLocation(fileLocation)

	err = di.fileCopy(fileLocation, swapLoc)
	if err != nil {
		return []byte{}, err
	}

	fileBytes, err := ioutil.ReadFile(swapLoc)
	return fileBytes, err
}

// FileWriteToSwap writes the swapfile for the file with the given info
func (di *DatabaseImpl) FileWriteToSwap(relpath string, filename string, projectID int64, raw []byte) error {
	relFilePath, err := calculateFilePath(relpath, filename, projectID)
	if err != nil {
		return err
	}
	fileLocation := filepath.Join(relFilePath, filename)
	swapLoc := getSwpLocation(fileLocation)

	return ioutil.WriteFile(swapLoc, raw, 0744)
}

// returns any error
func (di *DatabaseImpl) deleteSwp(relpath string, filename string, projectID int64) error {
	relFilePath, err := calculateFilePath(relpath, filename, projectID)
	if err != nil {
		return err
	}
	fileLocation := filepath.Join(relFilePath, filename)
	swapLoc := getSwpLocation(fileLocation)

	return os.Remove(swapLoc)
}

// swaps the swapfile to the location of the real file
func (di *DatabaseImpl) swapSwp(relpath string, filename string, projectID int64) error {
	relFilePath, err := calculateFilePath(relpath, filename, projectID)
	if err != nil {
		return err
	}
	fileLocation := filepath.Join(relFilePath, filename)
	swapLoc := getSwpLocation(fileLocation)

	err = di.fileCopy(swapLoc, fileLocation)
	if err != nil {
		return err
	}
	return os.Remove(swapLoc)
}

func (di *DatabaseImpl) fileCopy(src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.Mode().IsRegular() {
		return errors.New("non-regular source file cannot be copied")
	}
	_, err = os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			err = os.Remove(dst)
			if err != nil {
				return err
			}
		}
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func calculateFilePath(relpath string, filename string, projectID int64) (string, error) {
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

func getSwpLocation(filepath string) string {
	return filepath + ".swp"
}
