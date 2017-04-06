package datastore

import (
	"strings"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/patching"
	"github.com/CodeCollaborate/Server/utils"
)

var documentStoreFactoryMap = map[string]func(cfg *config.ConnCfg) DocumentStore{}

// RegisterDocumentStore is the registration point for any document datastore modules
func RegisterDocumentStore(name string, initFunc func(cfg *config.ConnCfg) DocumentStore) {
	documentStoreFactoryMap[strings.ToLower(name)] = initFunc
}

// FileData is the struct representing short-term changes in document state, and the versions therein
type FileData struct {
	FileID              int64
	Version             int64
	Changes             []string
	TempChanges         []string
	RemainingChanges    []string
	UseTemp             bool
	PullSwp             bool
	ScrunchedPatchCount int
	LastModifiedDate    int64
}

// AggregatedChanges aggregates the Changes, Temp and RemainingChanges
func (fileData *FileData) AggregatedChanges() []string {
	changes := []string{}

	if fileData.PullSwp {
		changes = append(fileData.RemainingChanges, fileData.TempChanges...)
		changes = append(changes, fileData.Changes...)
	} else if fileData.UseTemp {
		changes = append(fileData.Changes, fileData.TempChanges...)
	} else {
		changes = fileData.Changes
	}

	return changes
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

	// AppendPatch appends the patch to the document with the given fileID, and returns the resultant FileData and missing patches if successful
	AppendPatch(fileID int64, patch *patching.Patch) (*FileData, []string, error)

	// ScrunchChanges takes the set of untouched patches, and scrunches them into the base document.
	ScrunchChanges(fileID int64) error
}

// InitDocumentStore Initializes the DocumentStore, or throws a fatal error if unsuccessful.
func InitDocumentStore(name string, cfg *config.ConnCfg) DocumentStore {
	name = strings.ToLower(name)

	if documentStoreFactoryMap[name] == nil {
		utils.LogFatal("Configuration specified unknown DocumentStore", ErrFatalConfigurationErr, utils.LogFields{
			"DocumentStoreName": name,
		})
	}

	store := documentStoreFactoryMap[name](cfg)
	store.Connect()

	return store
}
