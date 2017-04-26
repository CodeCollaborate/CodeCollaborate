package datastore

import (
	"os"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/utils"
)

// FilePathSeparator is the string character for the os's filepath separator
var FilePathSeparator = string(os.PathSeparator)

// DataStore is the struct containing all the different types of DataStores used by this server
type DataStore struct {
	BucketStore     BucketStore
	DocumentStore   DocumentStore
	RelationalStore RelationalStore
}

// InitDataStore Initializes the DataStore, or throws a fatal error if unsuccessful.
func InitDataStore(cfg *config.DataStoreCfg) (*DataStore, error) {
	if cfg.BucketStoreName == "" || cfg.BucketStoreCfg == nil {
		utils.LogFatal("Invalid Configuration: Missing BucketStore Name/Config", ErrFatalConfigurationErr, nil)
	}
	if cfg.DocumentStoreName == "" || cfg.DocumentStoreCfg == nil {
		utils.LogFatal("Invalid Configuration: Missing DocumentStore Name/Config", ErrFatalConfigurationErr, nil)
	}
	if cfg.RelationalStoreName == "" || cfg.RelationalStoreCfg == nil {
		utils.LogFatal("Invalid Configuration: Missing DocumentStore Name/Config", ErrFatalConfigurationErr, nil)
	}

	bucketStore := InitBucketStore(cfg.BucketStoreName, cfg.BucketStoreCfg)

	documentStore := InitDocumentStore(cfg.DocumentStoreName, cfg.DocumentStoreCfg)

	relationalStore := InitRelationalStore(cfg.RelationalStoreName, cfg.RelationalStoreCfg)

	return &DataStore{
		BucketStore:     bucketStore,
		DocumentStore:   documentStore,
		RelationalStore: relationalStore,
	}, nil
}
