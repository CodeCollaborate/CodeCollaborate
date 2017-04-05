package datastore

import (
	"github.com/CodeCollaborate/Server/modules/config"
)

var documentStoreFactoryMap = map[string]func(cfg *config.ConnCfg) DocumentStore{}

func registerDocumentStore(name string, initFunc func(cfg *config.ConnCfg) DocumentStore) {
	documentStoreFactoryMap[name] = initFunc
}

// FileData is the struct representing short-term changes in document state, and the versions therein
type FileData struct {
	FileID              int64
	Version             int64
	Changes             []string
	ScrunchedPatchCount int
	LastModifiedDate    int64
}

// DocumentStore defines the interface for all document storage class datastores (Google Cloud Datastore, Amazon DynamoDB, etc)
type DocumentStore interface {
	// Connect starts this documentStore's connection to the server
	// If connect fails, it will throw a fatal error
	Connect()

	// Shutdown terminates this documentStore's connection to the server
	Shutdown()

	// AddFileData stores the given FileData using the internal FileID, without overwriting
	AddFileData(data *FileData) error

	// SetFileData stores the given FileData using the internal FileID, overwriting if necessary
	SetFileData(data *FileData) error

	// GetFileData retrieves the FileData for the given fileID
	GetFileData(fileID int64) (*FileData, error)

	// DeleteFileData deletes the FileData for the given fileID
	DeleteFileData(fileID int64) error

	// AppendPatch appends the patch to the document with the given fileID, and returns the resultant FileData if successful
	AppendPatch(fileID int64, patchStr string) (*FileData, error)

	// ScrunchChanges takes the set of untouched patches, and scrunches them into the base document.
	ScrunchChanges(fileID int64) error
}

// InitDocumentStore Initializes the DocumentStore, or throws a fatal error if unsuccessful.
func InitDocumentStore(name string, cfg *config.ConnCfg) DocumentStore {
	store := documentStoreFactoryMap[name](cfg)
	store.Connect()

	return store
}
