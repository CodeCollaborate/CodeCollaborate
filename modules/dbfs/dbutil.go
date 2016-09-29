package dbfs

import (
	"errors"
	"time"
)

// ErrNoDbChange : No rows or values in the DB were changed, which was an unexpected result
var ErrNoDbChange = errors.New("No entries were correctly altered")

// ErrNoData : No rows or values were found for this value in the database
var ErrNoData = errors.New("No entries were found")

// ErrVersionOutOfDate : The request attempted to mutate an out of date resource
var ErrVersionOutOfDate = errors.New("The request attempted to modify an out of date resource")

// ErrResourceNotFound : The request attempted to mutate an out of date resource
var ErrResourceNotFound = errors.New("No such resource was found")

// ErrDbNotInitialized : Active db connection does not exist
var ErrDbNotInitialized = errors.New("The database was not propperly initialized before execution")

// ErrMaliciousRequest : The request attempted to directly tamper with our filesystemp / database
var ErrMaliciousRequest = errors.New("The request attempted to directly tamper with our filesystemp / database")

// ProjectPermission is the type which represents the permission relationship on projects
type ProjectPermission struct {
	Username        string
	PermissionLevel int8
	GrantedBy       string
	GrantedDate     time.Time
}

// ProjectMeta is the type which represents a row in the MySQL `Project` table
type ProjectMeta struct {
	ProjectID       int64
	Name            string
	PermissionLevel int8
}

// FileMeta is the type that contains all the metadata about a file
type FileMeta struct {
	FileID       int64
	Creator      string
	CreationDate time.Time
	RelativePath string
	ProjectID    int64
	Filename     string
}

// UserMeta is the type that contains all the metadata about a user
type UserMeta struct {
	Username  string
	Password  string
	Email     string
	FirstName string
	LastName  string
}
