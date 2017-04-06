package bucketstore

import (
	"os"
	"testing"

	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datastore"
	"github.com/stretchr/testify/require"
)

var fsBucketStore = NewFilesystemBucketStore(&config.ConnCfg{Schema: "testFiles"})
var fileID = int64(1)
var fileContents = []byte("testFileContents")
var fileContents2 = []byte("newModifiedData")
var filePath = filepath.Join(fsBucketStore.rootFileDirectory, strconv.FormatInt(fileID, 10))
var swapFilePath = filepath.Join(fsBucketStore.rootFileDirectory, strconv.FormatInt(-1*fileID, 10))

func TestFileSystem_Connect(t *testing.T) {
	fsBucketStore.Connect()
	defer os.RemoveAll(fsBucketStore.rootFileDirectory)

	_, err := os.Stat(fsBucketStore.rootFileDirectory)
	require.Nil(t, err, "Could not find the root file directory")
}

func TestFileSystem_Shutdown(t *testing.T) {
	// Nothing to test
}

func TestFileSystem_AddFileBytes(t *testing.T) {
	fsBucketStore.Connect()
	defer os.RemoveAll(fsBucketStore.rootFileDirectory)

	// Create new file
	err := fsBucketStore.AddFile(fileID, fileContents)
	require.Nil(t, err, "Failed to add a new file to the storage")

	// Check to make sure it's on disk
	_, err = os.Stat(filePath)
	require.Nil(t, err, "Could not find the written file")

	// Check that the contents are correct
	fileBytes, err := ioutil.ReadFile(filePath)
	require.Nil(t, err, "Encountered error reading file")
	require.Equal(t, fileContents, fileBytes, "File contents was incorrect")

	// Attempt to create duplicate file (should fail)
	err = fsBucketStore.AddFile(fileID, fileContents)
	require.Equal(t, datastore.ErrFileAlreadyExists, err, "No error thrown, or incorrect error thrown for duplicate file insertion")
}

func TestFileSystem_SetFileBytes(t *testing.T) {
	fsBucketStore.Connect()
	defer os.RemoveAll(fsBucketStore.rootFileDirectory)

	// Create new file
	err := fsBucketStore.SetFile(fileID, fileContents)
	require.Nil(t, err, "Failed to add a new file to the storage")

	// Check to make sure it's on disk
	_, err = os.Stat(filePath)
	require.Nil(t, err, "Could not find the written file")

	// Check that the contents are correct
	fileBytes, err := ioutil.ReadFile(filePath)
	require.Nil(t, err, "Encountered error reading file")
	require.Equal(t, fileContents, fileBytes, "File contents was incorrect")

	// Attempt to overwrite file (should not fail)
	err = fsBucketStore.SetFile(fileID, fileContents2)
	require.Nil(t, err, "Failed to overwrite file in storage")

	// Check to make sure it is still on disk
	_, err = os.Stat(filePath)
	require.Nil(t, err, "Could not find the written file")

	// Check that the contents were updated
	fileBytes, err = ioutil.ReadFile(filePath)
	require.Nil(t, err, "Encountered error reading file")
	require.Equal(t, fileContents2, fileBytes, "File contents was incorrect")
}

func TestFileSystem_GetFileBytes(t *testing.T) {
	fsBucketStore.Connect()
	defer os.RemoveAll(fsBucketStore.rootFileDirectory)

	// Check to make sure error is thrown if file does not exist
	fileBytes, err := fsBucketStore.GetFile(fileID)
	require.Equal(t, datastore.ErrFileDoesNotExist, err, "GetFile did not throw an error, or threw an incorrect error for nonexistent file")

	// Create new file
	err = fsBucketStore.AddFile(fileID, fileContents)
	require.Nil(t, err, "Failed to add a new file to the storage")

	// Check GetFile returns the correct file
	fileBytes, err = fsBucketStore.GetFile(fileID)
	require.Nil(t, err, "Encountered error reading file")
	require.Equal(t, fileContents, fileBytes, "File contents was incorrect")

	// Remove it via OS calls
	err = os.Remove(filePath)
	require.Nil(t, err, "Failed to delete file")

	// Check to make sure error thrown for the removed file
	fileBytes, err = fsBucketStore.GetFile(fileID)
	require.Equal(t, datastore.ErrFileDoesNotExist, err, "GetFile did not throw an error, or threw an incorrect error for nonexistent file")
}

func TestFileSystem_DeleteFileBytes(t *testing.T) {
	fsBucketStore.Connect()
	defer os.RemoveAll(fsBucketStore.rootFileDirectory)

	// Attempt to delete a nonexistent file
	err := fsBucketStore.DeleteFile(fileID)
	require.Equal(t, datastore.ErrFileDoesNotExist, err, "DeleteFile did not throw an error, or threw an incorrect error for nonexistent file")

	// Create new file
	err = fsBucketStore.AddFile(fileID, fileContents)
	require.Nil(t, err, "Failed to add a new file to the storage")

	// Check to make sure it's on disk
	_, err = os.Stat(filePath)
	require.Nil(t, err, "Could not find the written file")

	// Make a swap file
	err = fsBucketStore.MakeSwapFile(fileID)
	require.Nil(t, err, "Failed to create the swap file")

	// Delete the file
	err = fsBucketStore.DeleteFile(fileID)
	require.Nil(t, err, "Failed to delete the file")

	// Check to make sure it's no longer on disk
	_, err = os.Stat(filePath)
	require.True(t, os.IsNotExist(err), "File was still present on disk")

	// Also check to make sure swap file was deleted
	_, err = os.Stat(swapFilePath)
	require.True(t, os.IsNotExist(err), "Swap file was still present on disk")

	// Attempt to delete the file again; should give a ErrFileDoesNotExist error
	err = fsBucketStore.DeleteFile(fileID)
	require.Equal(t, datastore.ErrFileDoesNotExist, err, "DeleteFile did not throw an error, or threw an incorrect error for nonexistent file")
}

func TestFileSystem_MakeSwapFile(t *testing.T) {
	fsBucketStore.Connect()
	defer os.RemoveAll(fsBucketStore.rootFileDirectory)

	// Attempt to create a swap file from a nonexistent file
	err := fsBucketStore.MakeSwapFile(fileID)
	require.Equal(t, datastore.ErrFileDoesNotExist, err, "MakeSwapFile did not throw an error, or threw an incorrect error for nonexistent file")

	// Write a temporary swap file
	err = ioutil.WriteFile(swapFilePath, []byte{}, 0744)
	require.Nil(t, err, "Failed to create a temporary swap file")

	// Create new file
	err = fsBucketStore.AddFile(fileID, fileContents)
	require.Nil(t, err, "Failed to add a new file to the storage")

	// Create the swap file
	err = fsBucketStore.MakeSwapFile(fileID)
	require.Nil(t, err, "Failed to make the swap file")

	// Check to make sure it's on disk
	_, err = os.Stat(swapFilePath)
	require.Nil(t, err, "Could not find the swap file")

	// Check that the contents are correct, and overwrote our temporary swap file
	fileBytes, err := ioutil.ReadFile(swapFilePath)
	require.Nil(t, err, "Encountered error reading swap file")
	require.Equal(t, fileContents, fileBytes, "Swap file contents was incorrect")

	// Change original contents
	err = fsBucketStore.SetFile(fileID, fileContents2)
	require.Nil(t, err, "Failed to modify file contents")

	// Check that original file contents changed
	fileBytes, err = ioutil.ReadFile(filePath)
	require.Nil(t, err, "Encountered error reading file")
	require.Equal(t, fileContents2, fileBytes, "File contents were not changed by SetFile")

	// Check that swap file contents did not change
	fileBytes, err = ioutil.ReadFile(swapFilePath)
	require.Nil(t, err, "Encountered error reading swap file")
	require.Equal(t, fileContents, fileBytes, "Swap file contents were changed by SetFile")
}

func TestFileSystem_SetSwapFile(t *testing.T) {
	fsBucketStore.Connect()
	defer os.RemoveAll(fsBucketStore.rootFileDirectory)

	// Write to swap file
	err := fsBucketStore.SetSwapFile(fileID, fileContents)
	require.Nil(t, err, "Failed to add a new file to the storage")

	// Check to make sure it's on disk
	_, err = os.Stat(swapFilePath)
	require.Nil(t, err, "Could not find the swap file")

	// Create new file
	err = fsBucketStore.AddFile(fileID, fileContents)
	require.Nil(t, err, "Failed to add a new file to the storage")

	// Create a swap file, overwriting the one we had previously
	err = fsBucketStore.MakeSwapFile(fileID)
	require.Nil(t, err, "Failed to make the swap file")

	// Check to make sure it's on disk
	_, err = os.Stat(swapFilePath)
	require.Nil(t, err, "Could not find the swap file")

	// Check that the contents are correct
	fileBytes, err := ioutil.ReadFile(swapFilePath)
	require.Nil(t, err, "Encountered error reading swap file")
	require.Equal(t, fileContents, fileBytes, "Swap file contents was incorrect")

	// Change swap file contents
	err = fsBucketStore.SetSwapFile(fileID, fileContents2)
	require.Nil(t, err, "Failed to modify swap file contents")

	// Check that original file contents did not change
	fileBytes, err = ioutil.ReadFile(filePath)
	require.Nil(t, err, "Encountered error reading file")
	require.Equal(t, fileContents, fileBytes, "File contents were changed by SetSwapFile")

	// Check that swap file contents changed
	fileBytes, err = ioutil.ReadFile(swapFilePath)
	require.Nil(t, err, "Encountered error reading swap file")
	require.Equal(t, fileContents2, fileBytes, "Swap file contents were not changed by SetSwapFile")
}

func TestFileSystem_GetSwapFile(t *testing.T) {
	fsBucketStore.Connect()
	defer os.RemoveAll(fsBucketStore.rootFileDirectory)

	// Attempt to get a nonexistent swap file
	fileBytes, err := fsBucketStore.GetSwapFile(fileID)
	require.Equal(t, datastore.ErrFileDoesNotExist, err, "GetSwapFile did not throw an error, or threw an incorrect error for nonexistent file")

	// Create a temporary swap file
	err = fsBucketStore.SetSwapFile(fileID, fileContents2)
	require.Nil(t, err, "Failed to add a new file to the storage")

	// Check that the contents are correct
	fileBytes, err = fsBucketStore.GetSwapFile(fileID)
	require.Nil(t, err, "Encountered error reading swap file")
	require.Equal(t, fileContents2, fileBytes, "Swap file contents was incorrect")

	// Create new file
	err = fsBucketStore.AddFile(fileID, fileContents)
	require.Nil(t, err, "Failed to add a new file to the storage")

	// Create a swap file, overwriting the one we had previously
	err = fsBucketStore.MakeSwapFile(fileID)
	require.Nil(t, err, "Failed to make the swap file")

	// Check to make sure it's on disk
	_, err = os.Stat(swapFilePath)
	require.Nil(t, err, "Could not find the swap file")

	// Check that the contents are correct
	fileBytes, err = fsBucketStore.GetSwapFile(fileID)
	require.Nil(t, err, "Encountered error reading swap file")
	require.Equal(t, fileContents, fileBytes, "Swap file contents was incorrect")
}

func TestFileSystem_DeleteSwapFile(t *testing.T) {
	fsBucketStore.Connect()
	defer os.RemoveAll(fsBucketStore.rootFileDirectory)

	// Attempt to delete a nonexistent swap file
	err := fsBucketStore.DeleteSwapFile(fileID)
	require.Equal(t, datastore.ErrFileDoesNotExist, err, "DeleteSwapFile did not throw an error, or threw an incorrect error for nonexistent file")

	// Create a temporary swap file
	err = fsBucketStore.SetSwapFile(fileID, fileContents2)
	require.Nil(t, err, "Failed to add a new file to the storage")

	// Attempt to delete the new swap file
	err = fsBucketStore.DeleteSwapFile(fileID)
	require.Nil(t, err, "Failed to delete swap file")

	// Check to make sure it's not on disk
	_, err = os.Stat(swapFilePath)
	require.True(t, os.IsNotExist(err), "Swap file was still on disk")

	// Create new file
	err = fsBucketStore.AddFile(fileID, fileContents)
	require.Nil(t, err, "Failed to add a new file to the storage")

	// Create a swap file
	err = fsBucketStore.MakeSwapFile(fileID)
	require.Nil(t, err, "Failed to make the swap file")

	// Attempt to delete the new swap file
	err = fsBucketStore.DeleteSwapFile(fileID)
	require.Nil(t, err, "Failed to delete swap file")

	// Check to make sure the original file is still on disk
	_, err = os.Stat(filePath)
	require.Nil(t, err, "Original file was deleted by DeleteSwapFile")

	// Check to make sure the swap file has been deleted
	_, err = os.Stat(swapFilePath)
	require.True(t, os.IsNotExist(err), "Swap file was still on disk")

	// Attempt to delete a non-existent swap file
	err = fsBucketStore.DeleteSwapFile(fileID)
	require.Equal(t, datastore.ErrFileDoesNotExist, err, "DeleteSwapFile did not throw an error, or threw an incorrect error for nonexistent file")
}

func TestFileSystem_RestoreSwapFile(t *testing.T) {
	fsBucketStore.Connect()
	defer os.RemoveAll(fsBucketStore.rootFileDirectory)

	// Attempt to restore a nonexistent swap file
	err := fsBucketStore.RestoreSwapFile(fileID)
	require.Equal(t, datastore.ErrFileDoesNotExist, err, "RestoreSwapFile did not throw an error, or threw an incorrect error for nonexistent file")

	// Create new file
	err = fsBucketStore.AddFile(fileID, fileContents)
	require.Nil(t, err, "Failed to add a new file to the storage")

	// Attempt to restore a nonexistent swap file (swap file still does not exist)
	err = fsBucketStore.RestoreSwapFile(fileID)
	require.Equal(t, datastore.ErrFileDoesNotExist, err, "RestoreSwapFile did not throw an error, or threw an incorrect error for nonexistent file")

	// Make the swap file
	err = fsBucketStore.MakeSwapFile(fileID)
	require.Nil(t, err, "Failed to make a swap file")

	// Create new file
	err = fsBucketStore.SetFile(fileID, fileContents2)
	require.Nil(t, err, "Failed to change file contents")

	// Attempt to restore the swap file
	err = fsBucketStore.RestoreSwapFile(fileID)
	require.Nil(t, err, "Failed to restore the swap file")

	// Check file contents of main file
	fileBytes, err := fsBucketStore.GetFile(fileID)
	require.Nil(t, err, "Encountered error reading file")
	require.Equal(t, fileContents, fileBytes, "File contents was incorrect")

	// Check swap file contents have not changed
	fileBytes, err = fsBucketStore.GetSwapFile(fileID)
	require.Nil(t, err, "Encountered error reading swap file")
	require.Equal(t, fileContents, fileBytes, "Swap file contents was incorrect")
}

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
