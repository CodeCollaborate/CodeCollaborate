package documentstore

import (
	"testing"
	"time"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datastore"
	"github.com/CodeCollaborate/Server/modules/patching"
	"github.com/stretchr/testify/require"
)

var fData = datastore.FileData{
	FileID:              -1,
	Changes:             []string{"Change1", "Change2", "Change3"},
	Version:             1,
	LastModifiedDate:    time.Now().UnixNano(),
	ScrunchedPatchCount: 0,
	PullSwp:             false,
	RemainingChanges:    []string{"Remaining1", "remaining2"},
	TempChanges:         []string{"temp1"},
	UseTemp:             false,
}

func TestCouchbaseDocumentStore_Connect(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewCouchbaseDocumentStore(config.GetConfig().DataStoreConfig.DocumentStoreCfg)

	defer func() {
		require.Nil(t, recover(), "Couchbase connect threw a fatal error")
	}()
	cb.Connect()
}

func TestCouchbaseDocumentStore_Add(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewCouchbaseDocumentStore(config.GetConfig().DataStoreConfig.DocumentStoreCfg)

	defer func() {
		require.Nil(t, recover(), "Couchbase connect threw a fatal error")
	}()
	cb.Connect()

	defer cb.DeleteFileData(fData.FileID)

	err := cb.AddFileData(&fData)
	require.Nil(t, err, "Failed to correctly insert file")

	err = cb.AddFileData(&fData)
	require.Equal(t, err, datastore.ErrFileAlreadyExists, "Failed to throw error on insert of duplicate key")

	val, err := cb.GetFileData(fData.FileID)
	require.Nil(t, err, "Failed to correctly retrieve file")

	require.Equal(t, fData.FileID, val.FileID, "Returned struct fileID did not match original struct's")
	require.Equal(t, fData.Version, val.Version, "Returned struct version did not match original struct's")
	require.Equal(t, fData.Changes, val.Changes, "Returned struct changes array did not match original struct's")
	require.Equal(t, fData.LastModifiedDate, val.LastModifiedDate, "Returned struct lastModifiedDate did not match original struct's")
	require.Equal(t, fData.ScrunchedPatchCount, val.ScrunchedPatchCount, "Returned struct scrunchedPatchCount did not match original struct's")
	require.Equal(t, fData.TempChanges, val.TempChanges, "Returned struct TempChanges array did not match original struct's")
	require.Equal(t, fData.RemainingChanges, val.RemainingChanges, "Returned struct RemainingChanges array did not match original struct's")
	require.Equal(t, fData.PullSwp, val.PullSwp, "Returned struct PullSwp value did not match original struct's")
	require.Equal(t, fData.UseTemp, val.UseTemp, "Returned struct UseTemp value did not match original struct's")
}

func TestCouchbaseDocumentStore_Set(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewCouchbaseDocumentStore(config.GetConfig().DataStoreConfig.DocumentStoreCfg)

	defer func() {
		require.Nil(t, recover(), "Couchbase connect threw a fatal error")
	}()
	cb.Connect()

	defer cb.DeleteFileData(fData.FileID)

	err := cb.SetFileData(&fData)
	require.Nil(t, err, "Failed to correctly insert file")

	val, err := cb.GetFileData(fData.FileID)
	require.Nil(t, err, "Failed to correctly retrieve file")

	fData.Version = 2

	err = cb.SetFileData(&fData)
	require.Nil(t, err, "Couchbase threw error on setting new file data")

	val, err = cb.GetFileData(fData.FileID)
	require.Nil(t, err, "Failed to correctly retrieve file")

	require.Equal(t, fData.FileID, val.FileID, "Returned struct fileID did not match original struct's")
	require.Equal(t, fData.Version, val.Version, "Returned struct version did not match original struct's")
	require.Equal(t, fData.Changes, val.Changes, "Returned struct changes array did not match original struct's")
	require.Equal(t, fData.LastModifiedDate, val.LastModifiedDate, "Returned struct lastModifiedDate did not match original struct's")
	require.Equal(t, fData.ScrunchedPatchCount, val.ScrunchedPatchCount, "Returned struct scrunchedPatchCount did not match original struct's")
	require.Equal(t, fData.TempChanges, val.TempChanges, "Returned struct TempChanges array did not match original struct's")
	require.Equal(t, fData.RemainingChanges, val.RemainingChanges, "Returned struct RemainingChanges array did not match original struct's")
	require.Equal(t, fData.PullSwp, val.PullSwp, "Returned struct PullSwp value did not match original struct's")
	require.Equal(t, fData.UseTemp, val.UseTemp, "Returned struct UseTemp value did not match original struct's")
}

func TestCouchbaseDocumentStore_Get(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewCouchbaseDocumentStore(config.GetConfig().DataStoreConfig.DocumentStoreCfg)

	defer func() {
		require.Nil(t, recover(), "Couchbase connect threw a fatal error")
	}()
	cb.Connect()

	defer cb.DeleteFileData(fData.FileID)

	val, err := cb.GetFileData(fData.FileID)
	require.Equal(t, err, datastore.ErrFileDoesNotExist, "Failed to throw error for missing file")

	err = cb.SetFileData(&fData)
	require.Nil(t, err, "Failed to correctly insert file")

	val, err = cb.GetFileData(fData.FileID)
	require.Nil(t, err, "Failed to correctly retrieve file")

	require.Equal(t, fData.FileID, val.FileID, "Returned struct fileID did not match original struct's")
	require.Equal(t, fData.Version, val.Version, "Returned struct version did not match original struct's")
	require.Equal(t, fData.Changes, val.Changes, "Returned struct changes array did not match original struct's")
	require.Equal(t, fData.LastModifiedDate, val.LastModifiedDate, "Returned struct lastModifiedDate did not match original struct's")
	require.Equal(t, fData.ScrunchedPatchCount, val.ScrunchedPatchCount, "Returned struct scrunchedPatchCount did not match original struct's")
	require.Equal(t, fData.TempChanges, val.TempChanges, "Returned struct TempChanges array did not match original struct's")
	require.Equal(t, fData.RemainingChanges, val.RemainingChanges, "Returned struct RemainingChanges array did not match original struct's")
	require.Equal(t, fData.PullSwp, val.PullSwp, "Returned struct PullSwp value did not match original struct's")
	require.Equal(t, fData.UseTemp, val.UseTemp, "Returned struct UseTemp value did not match original struct's")
}

func TestCouchbaseDocumentStore_Delete(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewCouchbaseDocumentStore(config.GetConfig().DataStoreConfig.DocumentStoreCfg)

	defer func() {
		require.Nil(t, recover(), "Couchbase connect threw a fatal error")
	}()
	cb.Connect()

	defer cb.DeleteFileData(fData.FileID)

	err := cb.SetFileData(&fData)
	require.Nil(t, err, "Failed to correctly insert file")

	cb.DeleteFileData(fData.FileID)
	require.Nil(t, err, "Failed to correctly delete file")

	_, err = cb.GetFileData(fData.FileID)
	require.Equal(t, err, datastore.ErrFileDoesNotExist, "File was still found after deletion")
}

func TestCouchbaseDocumentStore_AppendPatch(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewCouchbaseDocumentStore(config.GetConfig().DataStoreConfig.DocumentStoreCfg)

	defer func() {
		require.Nil(t, recover(), "Couchbase connect threw a fatal error")
	}()
	cb.Connect()

	defer cb.DeleteFileData(fData.FileID)

	// Reset the testing fileData; this test is picky.
	fData = datastore.FileData{
		FileID:              -1,
		Changes:             []string{},
		Version:             1,
		LastModifiedDate:    time.Now().UnixNano(),
		ScrunchedPatchCount: 0,
		PullSwp:             false,
		RemainingChanges:    []string{},
		TempChanges:         []string{},
		UseTemp:             false,
	}

	// Insert the file
	err := cb.SetFileData(&fData)
	require.Nil(t, err, "Failed to correctly insert file")

	// Check that the data is the same at the start
	val, err := cb.GetFileData(fData.FileID)
	require.Nil(t, err, "Failed to correctly retrieve file")
	require.Equal(t, fData.FileID, val.FileID, "Pre-change: Returned struct fileID did not match original struct's")
	require.Equal(t, fData.Version, val.Version, "Pre-change: Returned struct version did not match original struct's")
	require.Equal(t, fData.Changes, val.Changes, "Pre-change: Returned struct changes array did not match original struct's")
	require.Equal(t, fData.TempChanges, val.TempChanges, "Pre-change: Returned struct TempChanges array did not match original struct's")
	require.Equal(t, fData.RemainingChanges, val.RemainingChanges, "Pre-change: Returned struct RemainingChanges array did not match original struct's")
	require.Equal(t, fData.PullSwp, val.PullSwp, "Pre-change: Returned struct PullSwp value did not match original struct's")
	require.Equal(t, fData.UseTemp, val.UseTemp, "Pre-change: Returned struct UseTemp value did not match original struct's")
	require.Equal(t, fData.LastModifiedDate, val.LastModifiedDate, "Pre-change: Returned struct lastModifiedDate did not match original struct's")
	require.Equal(t, fData.ScrunchedPatchCount, val.ScrunchedPatchCount, "Pre-change: Returned struct scrunchedPatchCount did not match original struct's")

	// Add a patch
	//baseText := "hello"
	patch := patching.GetPatchOrDie(t, "v1:2:\n1:+5:test1:\n5")
	result, missingPatches, err := cb.AppendPatch(fData.FileID, patch)
	require.Nil(t, err, "Failed to append patches")
	require.Empty(t, missingPatches)

	fData.Changes = append(fData.Changes, patch.String())
	fData.Version++
	fData.LastModifiedDate = val.LastModifiedDate
	require.Equal(t, fData.Changes, result.Changes, "Post-change 2: Changes AppendPatch returned struct did not match fData")
	require.Equal(t, fData.TempChanges, result.TempChanges, "Post-change 2: TempChanges AppendPatch returned struct did not match fData")
	require.Equal(t, fData.Version, result.Version, "Post-change 2: TempChanges AppendPatch returned struct did not match fData")
	require.NotEqual(t, fData.LastModifiedDate, result.LastModifiedDate, "Post-change 2: Returned result struct lastModifiedDate was not updated")

	// Check that the data is the same at the start
	val, err = cb.GetFileData(fData.FileID)
	require.Nil(t, err, "Failed to correctly retrieve file")
	require.Equal(t, fData.FileID, val.FileID, "Post-change 1: Returned struct fileID did not match original struct's")
	require.Equal(t, fData.Version, val.Version, "Post-change 1: Returned struct version did not match original struct's")
	require.Equal(t, fData.Changes, val.Changes, "Post-change 1: Returned struct changes array did not match original struct's")
	require.Equal(t, fData.TempChanges, val.TempChanges, "Post-change 1: Returned struct TempChanges array did not match original struct's")
	require.Equal(t, fData.RemainingChanges, val.RemainingChanges, "Post-change 1: Returned struct RemainingChanges array did not match original struct's")
	require.Equal(t, fData.PullSwp, val.PullSwp, "Post-change 1: Returned struct PullSwp value did not match original struct's")
	require.Equal(t, fData.UseTemp, val.UseTemp, "Post-change 1: Returned struct UseTemp value did not match original struct's")
	require.NotEqual(t, fData.LastModifiedDate, val.LastModifiedDate, "Post-change 1: Returned struct lastModifiedDate was not updated")
	require.Equal(t, fData.ScrunchedPatchCount, val.ScrunchedPatchCount, "Post-change 1: Returned struct scrunchedPatchCount did not match original struct's")

	// Add a second patch that should be transformed
	patch = patching.GetPatchOrDie(t, "v1:2:\n2:+5:test2:\n5")
	result, missingPatches, err = cb.AppendPatch(fData.FileID, patch)
	require.Nil(t, err, "Failed to append patches")
	require.Contains(t, missingPatches, "v1:2:\n1:+5:test1:\n5", "Did not contain missing patch that it transformed against")

	fData.Changes = append(fData.Changes, "v2:3:\n7:+5:test2:\n10")
	fData.Version++
	fData.LastModifiedDate = val.LastModifiedDate
	require.Equal(t, fData.Changes, result.Changes, "Post-change 2: Changes AppendPatch returned struct did not match fData")
	require.Equal(t, fData.TempChanges, result.TempChanges, "Post-change 2: TempChanges AppendPatch returned struct did not match fData")
	require.Equal(t, fData.Version, result.Version, "Post-change 2: TempChanges AppendPatch returned struct did not match fData")
	require.NotEqual(t, fData.LastModifiedDate, result.LastModifiedDate, "Post-change 2: Returned result struct lastModifiedDate was not updated")

	// Check that the data was inserted correctly
	val, err = cb.GetFileData(fData.FileID)
	require.Nil(t, err, "Failed to correctly retrieve file")
	require.Equal(t, fData.FileID, val.FileID, "Post-change 2: Returned struct fileID did not match original struct's")
	require.Equal(t, fData.Version, val.Version, "Post-change 2: Returned struct version did not match original struct's")
	require.Equal(t, fData.Changes, val.Changes, "Post-change 2: Returned struct changes array did not match original struct's")
	require.Equal(t, fData.TempChanges, val.TempChanges, "Post-change 2: Returned struct TempChanges array did not match original struct's")
	require.Equal(t, fData.RemainingChanges, val.RemainingChanges, "Post-change 2: Returned struct RemainingChanges array did not match original struct's")
	require.Equal(t, fData.PullSwp, val.PullSwp, "Post-change 2: Returned struct PullSwp value did not match original struct's")
	require.Equal(t, fData.UseTemp, val.UseTemp, "Post-change 2: Returned struct UseTemp value did not match original struct's")
	require.NotEqual(t, fData.LastModifiedDate, val.LastModifiedDate, "Post-change 2: Returned struct lastModifiedDate was not updated")
	require.Equal(t, fData.ScrunchedPatchCount, val.ScrunchedPatchCount, "Post-change 2: Returned struct scrunchedPatchCount did not match original struct's")

	// Add a third patch that should not be transformed
	patch = patching.GetPatchOrDie(t, "v3:4:\n13:+5:test3:\n15")
	result, missingPatches, err = cb.AppendPatch(fData.FileID, patch)
	require.Nil(t, err, "Failed to append patches")
	require.Empty(t, missingPatches)

	fData.Changes = append(fData.Changes, patch.String())
	fData.Version++
	fData.LastModifiedDate = val.LastModifiedDate
	require.Equal(t, fData.Changes, result.Changes, "Post-change 2: Changes AppendPatch returned struct did not match fData")
	require.Equal(t, fData.TempChanges, result.TempChanges, "Post-change 2: TempChanges AppendPatch returned struct did not match fData")
	require.Equal(t, fData.Version, result.Version, "Post-change 2: TempChanges AppendPatch returned struct did not match fData")
	require.NotEqual(t, fData.LastModifiedDate, result.LastModifiedDate, "Post-change 2: Returned result struct lastModifiedDate was not updated")

	// Check that the data was inserted correctly
	val, err = cb.GetFileData(fData.FileID)
	require.Nil(t, err, "Failed to correctly retrieve file")
	require.Equal(t, fData.FileID, val.FileID, "Post-change 3: Returned struct fileID did not match original struct's")
	require.Equal(t, fData.Version, val.Version, "Post-change 3: Returned struct version did not match original struct's")
	require.Equal(t, fData.Changes, val.Changes, "Post-change 3: Returned struct changes array did not match original struct's")
	require.Equal(t, fData.TempChanges, val.TempChanges, "Post-change 3: Returned struct TempChanges array did not match original struct's")
	require.Equal(t, fData.RemainingChanges, val.RemainingChanges, "Post-change 3: Returned struct RemainingChanges array did not match original struct's")
	require.Equal(t, fData.PullSwp, val.PullSwp, "Post-change 3: Returned struct PullSwp value did not match original struct's")
	require.Equal(t, fData.UseTemp, val.UseTemp, "Post-change 3: Returned struct UseTemp value did not match original struct's")
	require.NotEqual(t, fData.LastModifiedDate, val.LastModifiedDate, "Post-change 3: Returned struct lastModifiedDate was not updated")
	require.Equal(t, fData.ScrunchedPatchCount, val.ScrunchedPatchCount, "Post-change 3: Returned struct scrunchedPatchCount did not match original struct's")
}

func TestCouchbaseDocumentStore_AppendPatchWithUseTmpOn(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewCouchbaseDocumentStore(config.GetConfig().DataStoreConfig.DocumentStoreCfg)

	defer func() {
		require.Nil(t, recover(), "Couchbase connect threw a fatal error")
	}()
	cb.Connect()

	defer cb.DeleteFileData(fData.FileID)

	// Reset the testing fileData; this test is picky.
	fData = datastore.FileData{
		FileID:              -1,
		Changes:             []string{},
		Version:             1,
		LastModifiedDate:    time.Now().UnixNano(),
		ScrunchedPatchCount: 0,
		PullSwp:             false,
		RemainingChanges:    []string{},
		TempChanges:         []string{},
		UseTemp:             true,
	}

	// Insert the file
	err := cb.SetFileData(&fData)
	require.Nil(t, err, "Failed to correctly insert file")

	// Check that the data is the same at the start
	val, err := cb.GetFileData(fData.FileID)
	require.Nil(t, err, "Failed to correctly retrieve file")
	require.Equal(t, fData.FileID, val.FileID, "Pre-change: Returned struct fileID did not match original struct's")
	require.Equal(t, fData.Version, val.Version, "Pre-change: Returned struct version did not match original struct's")
	require.Equal(t, fData.Changes, val.Changes, "Pre-change: Returned struct changes array did not match original struct's")
	require.Equal(t, fData.TempChanges, val.TempChanges, "Pre-change: Returned struct TempChanges array did not match original struct's")
	require.Equal(t, fData.RemainingChanges, val.RemainingChanges, "Pre-change: Returned struct RemainingChanges array did not match original struct's")
	require.Equal(t, fData.PullSwp, val.PullSwp, "Pre-change: Returned struct PullSwp value did not match original struct's")
	require.Equal(t, fData.UseTemp, val.UseTemp, "Pre-change: Returned struct UseTemp value did not match original struct's")
	require.Equal(t, fData.LastModifiedDate, val.LastModifiedDate, "Pre-change: Returned struct lastModifiedDate did not match original struct's")
	require.Equal(t, fData.ScrunchedPatchCount, val.ScrunchedPatchCount, "Pre-change: Returned struct scrunchedPatchCount did not match original struct's")

	// Add a patch
	//baseText := "hello"
	patch := patching.GetPatchOrDie(t, "v1:2:\n1:+5:test1:\n5")
	result, missingPatches, err := cb.AppendPatch(fData.FileID, patch)
	require.Nil(t, err, "Failed to append patches")
	require.Empty(t, missingPatches)

	fData.TempChanges = append(fData.TempChanges, patch.String())
	fData.Version++
	fData.LastModifiedDate = val.LastModifiedDate
	require.Equal(t, fData.Changes, result.Changes, "Post-change 2: Changes AppendPatch returned struct did not match fData")
	require.Equal(t, fData.TempChanges, result.TempChanges, "Post-change 2: TempChanges AppendPatch returned struct did not match fData")
	require.Equal(t, fData.Version, result.Version, "Post-change 2: TempChanges AppendPatch returned struct did not match fData")
	require.NotEqual(t, fData.LastModifiedDate, result.LastModifiedDate, "Post-change 2: Returned result struct lastModifiedDate was not updated")

	// Check that the data is the same after patch insertion
	val, err = cb.GetFileData(fData.FileID)
	require.Nil(t, err, "Failed to correctly retrieve file")
	require.Equal(t, fData.FileID, val.FileID, "Post-change 1: Returned struct fileID did not match original struct's")
	require.Equal(t, fData.Version, val.Version, "Post-change 1: Returned struct version did not match original struct's")
	require.Equal(t, fData.Changes, val.Changes, "Post-change 1: Returned struct changes array did not match original struct's")
	require.Equal(t, fData.TempChanges, val.TempChanges, "Post-change 1: Returned struct TempChanges array did not match original struct's")
	require.Equal(t, fData.RemainingChanges, val.RemainingChanges, "Post-change 1: Returned struct RemainingChanges array did not match original struct's")
	require.Equal(t, fData.PullSwp, val.PullSwp, "Post-change 1: Returned struct PullSwp value did not match original struct's")
	require.Equal(t, fData.UseTemp, val.UseTemp, "Post-change 1: Returned struct UseTemp value did not match original struct's")
	require.NotEqual(t, fData.LastModifiedDate, val.LastModifiedDate, "Post-change 1: Returned struct lastModifiedDate was not updated")
	require.Equal(t, fData.ScrunchedPatchCount, val.ScrunchedPatchCount, "Post-change 1: Returned struct scrunchedPatchCount did not match original struct's")

	// Add a second patch that should be transformed
	patch = patching.GetPatchOrDie(t, "v1:2:\n2:+5:test2:\n5")
	result, missingPatches, err = cb.AppendPatch(fData.FileID, patch)
	require.Nil(t, err, "Failed to append patches")
	require.Contains(t, missingPatches, "v1:2:\n1:+5:test1:\n5", "Did not contain missing patch that it transformed against")

	fData.TempChanges = append(fData.TempChanges, "v2:3:\n7:+5:test2:\n10")
	fData.Version++
	fData.LastModifiedDate = val.LastModifiedDate
	require.Equal(t, fData.Changes, result.Changes, "Post-change 2: Changes AppendPatch returned struct did not match fData")
	require.Equal(t, fData.TempChanges, result.TempChanges, "Post-change 2: TempChanges AppendPatch returned struct did not match fData")
	require.Equal(t, fData.Version, result.Version, "Post-change 2: TempChanges AppendPatch returned struct did not match fData")
	require.NotEqual(t, fData.LastModifiedDate, result.LastModifiedDate, "Post-change 2: Returned result struct lastModifiedDate was not updated")

	// Check that the data was inserted correctly
	val, err = cb.GetFileData(fData.FileID)
	require.Nil(t, err, "Failed to correctly retrieve file")
	require.Equal(t, fData.FileID, val.FileID, "Post-change 2: Returned struct fileID did not match original struct's")
	require.Equal(t, fData.Version, val.Version, "Post-change 2: Returned struct version did not match original struct's")
	require.Equal(t, fData.Changes, val.Changes, "Post-change 2: Returned struct changes array did not match original struct's")
	require.Equal(t, fData.TempChanges, val.TempChanges, "Post-change 2: Returned struct TempChanges array did not match original struct's")
	require.Equal(t, fData.RemainingChanges, val.RemainingChanges, "Post-change 2: Returned struct RemainingChanges array did not match original struct's")
	require.Equal(t, fData.PullSwp, val.PullSwp, "Post-change 2: Returned struct PullSwp value did not match original struct's")
	require.Equal(t, fData.UseTemp, val.UseTemp, "Post-change 2: Returned struct UseTemp value did not match original struct's")
	require.NotEqual(t, fData.LastModifiedDate, val.LastModifiedDate, "Post-change 2: Returned struct lastModifiedDate was not updated")
	require.Equal(t, fData.ScrunchedPatchCount, val.ScrunchedPatchCount, "Post-change 2: Returned struct scrunchedPatchCount did not match original struct's")
}
