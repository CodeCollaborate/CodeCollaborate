package datastore

import (
	"strings"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/utils"
)

var bucketStoreFactoryMap = map[string]func(cfg *config.ConnCfg) BucketStore{}

// RegisterBucketStore is the registration point for any bucket datastore modules
func RegisterBucketStore(name string, initFunc func(cfg *config.ConnCfg) BucketStore) {
	bucketStoreFactoryMap[strings.ToLower(name)] = initFunc
}

// BucketStore defines the interface for all bucket storage class datastores (Google Cloud Storage, AWS S3, etc)
type BucketStore interface {
	// Connect starts this bucketStore's connection to the server
	// If connect fails, it will throw a fatal error
	Connect()

	// Shutdown terminates this bucketStore's connection to the server
	Shutdown()

	// AddFile saves the file to the bucket, keyed on the given fileID, without overwriting
	AddFile(fileID int64, fileBytes []byte) error

	// SetFile saves the file to the bucket, keyed on the given fileID, overwriting if necessary
	SetFile(fileID int64, fileBytes []byte) error

	// GetFile retrieves the file from the bucket for a given fileID key
	GetFile(fileID int64) ([]byte, error)

	// DeleteFile deletes the file from the bucket for a given fileID key
	DeleteFile(fileID int64) error

	// MakeSwapFile makes a copy of the current file, naming it according to -1 * fileID
	MakeSwapFile(fileID int64) error

	// SetSwapFile reads from the swap file for the given fileID
	SetSwapFile(fileID int64, fileBytes []byte) error

	// GetSwapFile reads from the swap file for the given fileID
	GetSwapFile(fileID int64) ([]byte, error)

	// DeleteSwapFile deletes the swap file for the given fileID
	DeleteSwapFile(fileID int64) error

	// RestoreSwapFile copies the swap file over to the actual file
	RestoreSwapFile(fileID int64) error
}

// InitBucketStore Initializes the BucketStore, or throws a fatal error if unsuccessful.
func InitBucketStore(name string, cfg *config.ConnCfg) BucketStore {
	name = strings.ToLower(name)

	if bucketStoreFactoryMap[name] == nil {
		utils.LogFatal("Configuration specified unknown BucketStore", ErrFatalConfigurationErr, utils.LogFields{
			"BucketStoreName": name,
		})
	}

	store := bucketStoreFactoryMap[name](cfg)
	store.Connect()

	return store
}
