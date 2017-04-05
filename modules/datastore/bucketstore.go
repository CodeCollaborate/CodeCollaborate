package datastore

import (
	"github.com/CodeCollaborate/Server/modules/config"
)

var bucketStoreFactoryMap = map[string]func(cfg *config.ConnCfg) BucketStore{}

func registerBucketStore(name string, initFunc func(cfg *config.ConnCfg) BucketStore) {
	bucketStoreFactoryMap[name] = initFunc
}

// BucketStore defines the interface for all bucket storage class datastores (Google Cloud Storage, AWS S3, etc)
type BucketStore interface {
	// Connect starts this bucketStore's connection to the server
	// If connect fails, it will throw a fatal error
	Connect()

	// Shutdown terminates this bucketStore's connection to the server
	Shutdown()

	// AddFileBytes saves the file to the bucket, keyed on the given fileID, without overwriting
	AddFileBytes(fileID int64, fileBytes []byte) error

	// SetFileBytes saves the file to the bucket, keyed on the given fileID, overwriting if necessary
	SetFileBytes(fileID int64, fileBytes []byte) error

	// GetFileBytes retrieves the file from the bucket for a given fileID key
	GetFileBytes(fileID int64) ([]byte, error)

	// DeleteFileBytes deletes the file from the bucket for a given fileID key
	DeleteFileBytes(fileID int64) error
}

// InitBucketStore Initializes the BucketStore, or throws a fatal error if unsuccessful.
func InitBucketStore(name string, cfg *config.ConnCfg) BucketStore {
	store := bucketStoreFactoryMap[name](cfg)
	store.Connect()

	return store
}
