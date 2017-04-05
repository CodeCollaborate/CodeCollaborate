package datastore

import (
	"github.com/CodeCollaborate/Server/modules/config"
)

// DataStore is the struct containing all the different types of DataStores used by this server
type DataStore struct {
	BucketStore     BucketStore
	DocumentStore   DocumentStore
	RelationalStore RelationalStore
}

// InitDataStore Initializes the DataStore, or throws a fatal error if unsuccessful.
func InitDataStore(cfg *config.DataStoreCfg) (*DataStore, error) {
	bucketStore := InitBucketStore(cfg.BucketStoreName, cfg.BucketStoreCfg)

	documentStore := InitDocumentStore(cfg.DocumentStoreName, cfg.DocumentStoreCfg)

	relationalStore := InitRelationalStore(cfg.RelationalStoreName, cfg.RelationalStoreCfg)

	return &DataStore{
		BucketStore:     bucketStore,
		DocumentStore:   documentStore,
		RelationalStore: relationalStore,
	}, nil
}
