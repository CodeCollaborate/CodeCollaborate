package bucketstore

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datastore"
	"github.com/CodeCollaborate/Server/utils"
)

// FilesystemBucketStore is a concrete implementation of the BucketStore, using the filesystem as the backing bucket store of choice
type FilesystemBucketStore struct {
	cfg               *config.ConnCfg
	rootFileDirectory string
}

// NewFilesystemBucketStore creates a new instance of the FilesystemBucketStore, setting the configuration
func NewFilesystemBucketStore(cfg *config.ConnCfg) *FilesystemBucketStore {
	return &FilesystemBucketStore{
		cfg:               cfg,
		rootFileDirectory: filepath.Clean(cfg.Schema),
	}
}

// Connect starts this bucketStore's connection to the server
// If connect fails, it will throw a fatal error
func (store *FilesystemBucketStore) Connect() {
	err := os.MkdirAll(store.rootFileDirectory, 0744)
	if err != nil {
		utils.LogFatal("Could not initialize filesystem bucketStore directory", err, utils.LogFields{
			"Directory": store.rootFileDirectory,
		})
	}
}

// Shutdown terminates this bucketStore's connection to the server
func (store *FilesystemBucketStore) Shutdown() {
	// Do nothing
}

// AddFile saves the file to the bucket, keyed on the given fileID, without overwriting
func (store *FilesystemBucketStore) AddFile(fileID int64, fileBytes []byte) error {
	if fileID == 0 {
		return datastore.ErrInvalidFileID
	}

	relFilePath := filepath.Join(store.rootFileDirectory, strconv.FormatInt(fileID, 10))

	if _, err := ioutil.ReadFile(relFilePath); err == nil {
		return datastore.ErrFileAlreadyExists
	}

	err := ioutil.WriteFile(relFilePath, fileBytes, 0744)
	if err != nil {
		return err
	}

	return err
}

// SetFile saves the file to the bucket, keyed on the given fileID, overwriting if necessary
func (store *FilesystemBucketStore) SetFile(fileID int64, fileBytes []byte) error {
	if fileID == 0 {
		return datastore.ErrInvalidFileID
	}

	relFilePath := filepath.Join(store.rootFileDirectory, strconv.FormatInt(fileID, 10))

	err := ioutil.WriteFile(relFilePath, fileBytes, 0744)
	if err != nil {
		return err
	}

	return err
}

// GetFile retrieves the file from the bucket for a given fileID key
func (store *FilesystemBucketStore) GetFile(fileID int64) ([]byte, error) {
	if fileID == 0 {
		return nil, datastore.ErrInvalidFileID
	}

	relFilePath := filepath.Join(store.rootFileDirectory, strconv.FormatInt(fileID, 10))

	bytes, err := ioutil.ReadFile(relFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, datastore.ErrFileDoesNotExist
		}
		return nil, err
	}

	return bytes, nil
}

// DeleteFile deletes the file from the bucket for a given fileID key
func (store *FilesystemBucketStore) DeleteFile(fileID int64) error {
	if fileID == 0 {
		return datastore.ErrInvalidFileID
	}

	relFilePath := filepath.Join(store.rootFileDirectory, strconv.FormatInt(fileID, 10))
	relSwpFilePath := filepath.Join(store.rootFileDirectory, strconv.FormatInt(-1*fileID, 10))

	err := os.Remove(relFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return datastore.ErrFileDoesNotExist
		}
		return err
	}

	// Make sure that we don't delete the original file if we delete the swapFile
	if _, err := os.Stat(relSwpFilePath); fileID > 0 && !os.IsNotExist(err) {
		err := os.Remove(relSwpFilePath)
		if err != nil {
			utils.LogError("FileSystemBucketStore encountered an error in removing a swap file", err, utils.LogFields{
				"FileID": fileID,
			})
			// Do not return this error to the user; their actual operation succeeded
		}
	}

	return nil
}

// MakeSwapFile makes a copy of the current file, naming it according to -1 * fileID
func (store *FilesystemBucketStore) MakeSwapFile(fileID int64) error {
	if fileID == 0 {
		return datastore.ErrInvalidFileID
	}

	relFilePath := filepath.Join(store.rootFileDirectory, strconv.FormatInt(fileID, 10))
	relSwpFilePath := filepath.Join(store.rootFileDirectory, strconv.FormatInt(-1*fileID, 10))

	err := fileCopy(relFilePath, relSwpFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return datastore.ErrFileDoesNotExist
		}
		return err
	}

	return err
}

// SetSwapFile reads from the swap file for the given fileID
func (store *FilesystemBucketStore) SetSwapFile(fileID int64, fileBytes []byte) error {
	if fileID == 0 {
		return datastore.ErrInvalidFileID
	}

	return store.SetFile(-1*fileID, fileBytes)
}

// GetSwapFile reads from the swap file for the given fileID
func (store *FilesystemBucketStore) GetSwapFile(fileID int64) ([]byte, error) {
	if fileID == 0 {
		return nil, datastore.ErrInvalidFileID
	}

	return store.GetFile(-1 * fileID)
}

// DeleteSwapFile deletes the swap file for the given fileID
func (store *FilesystemBucketStore) DeleteSwapFile(fileID int64) error {
	if fileID == 0 {
		return datastore.ErrInvalidFileID
	}

	return store.DeleteFile(-1 * fileID)
}

// RestoreSwapFile copies the swap file over to the actual file
func (store *FilesystemBucketStore) RestoreSwapFile(fileID int64) error {
	if fileID == 0 {
		return datastore.ErrInvalidFileID
	}

	relFilePath := filepath.Join(store.rootFileDirectory, strconv.FormatInt(fileID, 10))
	relSwpFilePath := filepath.Join(store.rootFileDirectory, strconv.FormatInt(-1*fileID, 10))

	err := fileCopy(relSwpFilePath, relFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return datastore.ErrFileDoesNotExist
		}
		return err
	}

	return err
}

func fileCopy(src string, dst string) error {
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
