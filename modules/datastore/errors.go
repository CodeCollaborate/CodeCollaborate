package datastore

import "errors"

// ErrFatalServerErr is the error thrown for the failure of critical tasks, such that the server no longer is functioning correctly
var ErrFatalServerErr = errors.New("Server encountered a fatal error")

// ErrFatalConfigurationErr is the error thrown for the failure to parse the configuration, such that the server is unable to startup
var ErrFatalConfigurationErr = errors.New("Server encountered a fatal configuration error")

// ErrInternalServerErr is the error thrown for the failure to process a request where the failure is not caused by invalid data in the request
var ErrInternalServerErr = errors.New("Server encountered an internal error")

// ErrNotYetImplemented is the error thrown when the called method has not yet been implemented
var ErrNotYetImplemented = errors.New("Server encountered an internal error")

// ErrInvalidUsernamePassword is the error thrown when given an incorrect login
var ErrInvalidUsernamePassword = errors.New("Invalid username or password")

// ErrAuthenticationFailed is the error thrown when the provided token cannot be decrypted
var ErrAuthenticationFailed = errors.New("Failed to authenticate token")

// ErrInsufficientPermissions is the error thrown when a user attempts to perform an action he/she does not have the permissions to do
var ErrInsufficientPermissions = errors.New("Insufficient permissions to execute requested operation")

// ErrDuplicateUsername is the error thrown for duplicate usernames when attempting to register a new user
var ErrDuplicateUsername = errors.New("An account with the given username already exists")

// ErrDuplicateEmail is the error thrown for duplicate emails when attempting to register a new user
var ErrDuplicateEmail = errors.New("An account with the given email address already exists")

// ErrFileAlreadyExists is the error thrown when a file creation is requested, but the file already exists
var ErrFileAlreadyExists = errors.New("A file already exists with the given key")

// ErrFileDoesNotExist is the error thrown when a file access/modify is requested, but the file does not exist
var ErrFileDoesNotExist = errors.New("No such file exists for the given fileID")

// ErrFileBaseVersionTooHigh is the error thrown when a client sends a patch that has a base version higher than the server's version
// This represents a logic flaw, since the server is the decider, and should always be the one versioning the files.
var ErrFileBaseVersionTooHigh = errors.New("Patch base version was higher than current document version")

// ErrFileBaseVersionTooLow is the error thrown when a client sends a patch that is too far behind, and thus cannot be transformed up
// to the current version. This generally occurs when the version it depended on has been scrunched.
var ErrFileBaseVersionTooLow = errors.New("Patch base version has been discarded or scrunched")

// ErrInvalidFileID is the error thrown when an invalid fileID was provided. This could be due to the fileID being 0, overflowing, NaN or infinity
var ErrInvalidFileID = errors.New("Invalid fileID provided")

// ErrInvalidFileName is the error thrown when an invalid filename was provided. This could be due to protected characters in the filename (.., /, etc)
var ErrInvalidFileName = errors.New("Invalid filename provided")

// ErrInvalidFilePath is the error thrown when an invalid filepath was provided. This could be due to attempting to write above its root directory
var ErrInvalidFilePath = errors.New("Invalid filepath provided")
