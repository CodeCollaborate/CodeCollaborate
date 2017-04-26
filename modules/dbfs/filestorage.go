package dbfs

import (
	"os"
	"strconv"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datastore/bucketstore"
)

var filePathSeparator = strconv.QuoteRune(os.PathSeparator)[1:2]

func (di *DatabaseImpl) initFileSystemIfNeeded() {
	if di.bucketStore == nil {
		di.bucketStore = bucketstore.NewFilesystemBucketStore(&config.ConnCfg{
			Schema: config.GetConfig().ServerConfig.ProjectPath,
		})
	}

	di.bucketStore.Connect()
}

// FileWrite writes the file with the given bytes to a calculated path, and
// returns that path so it can be put in MySQL
func (di *DatabaseImpl) FileWrite(fileID int64, raw []byte) error {
	di.initFileSystemIfNeeded()

	return di.bucketStore.SetFile(fileID, raw)
}

// FileDelete deletes the file with the given metadata from the file system
// Couple this with dbfs.MySQLFileDelete and dbfs.CBDeleteFile
func (di *DatabaseImpl) FileDelete(fileID int64) error {
	di.initFileSystemIfNeeded()

	return di.bucketStore.DeleteFile(fileID)
}

// FileRead returns the project file from the calculated location on the disk
func (di *DatabaseImpl) FileRead(fileID int64) ([]byte, error) {
	di.initFileSystemIfNeeded()

	return di.bucketStore.GetFile(fileID)
}

// returns the swap file contents and any error
func (di *DatabaseImpl) makeSwp(fileID int64) ([]byte, error) {
	di.initFileSystemIfNeeded()

	if err := di.bucketStore.MakeSwapFile(fileID); err != nil {
		return nil, err
	}

	return di.swapRead(fileID)
}

// swapRead returns the swap file from the calculated location on the disk
func (di *DatabaseImpl) swapRead(fileID int64) ([]byte, error) {
	di.initFileSystemIfNeeded()

	return di.bucketStore.GetSwapFile(fileID)
}

// FileWriteToSwap writes the swapfile for the file with the given info
func (di *DatabaseImpl) FileWriteToSwap(fileID int64, raw []byte) error {
	di.initFileSystemIfNeeded()

	return di.bucketStore.SetSwapFile(fileID, raw)
}

// returns any error
func (di *DatabaseImpl) deleteSwp(fileID int64) error {
	di.initFileSystemIfNeeded()

	return di.bucketStore.DeleteSwapFile(fileID)
}

// swaps the swapfile to the location of the real file
func (di *DatabaseImpl) swapSwp(fileID int64) error {
	di.initFileSystemIfNeeded()

	return di.bucketStore.RestoreSwapFile(fileID)
}
