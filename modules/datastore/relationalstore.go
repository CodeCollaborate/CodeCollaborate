package datastore

import (
	"time"

	"strings"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/utils"
)

var relationalStoreFactoryMap = map[string]func(cfg *config.ConnCfg) RelationalStore{}

// RegisterRelationalStore is the registration point for any relational datastore modules
func RegisterRelationalStore(name string, initFunc func(cfg *config.ConnCfg) RelationalStore) {
	relationalStoreFactoryMap[strings.ToLower(name)] = initFunc
}

// FileMetadata represents the relational file information that doesn't often change
type FileMetadata struct {
	ProjectID    int64
	FileID       int64
	Filename     string
	RelativePath string
	Creator      string
	CreationDate time.Time
}

// ProjectMetadata represents the relational project information stored in the database
type ProjectMetadata struct {
	ProjectID          int64
	Name               string
	ProjectPermissions []*ProjectPermission
}

// ProjectPermission represents the permissions the users have for the projects
type ProjectPermission struct {
	Username        string
	PermissionLevel int
	GrantedBy       string
	GrantedDate     time.Time
}

// UserData represents the data stored for each user
type UserData struct {
	Username  string
	Email     string
	Password  string // MUST be hashed before storage
	FirstName string
	LastName  string
}

// RelationalStore defines the interface for all relational storage class datastores (Google Cloud SQL, Amazon RDS, etc)
type RelationalStore interface {
	// Connect starts this relationalStore's connection to the server
	// If connect fails, it will throw a fatal error
	Connect()

	// Shutdown terminates this relationalStore's connection to the server
	Shutdown()

	// AddUser stores the given UserData using the internal username, without overwriting
	AddUser(data *UserData) error

	// SetUser stores the given UserData using the internal username, overwriting if necessary
	SetUser(data *UserData) error

	// GetUser retrieves the UserData for the user with given username
	GetUser(username string) (*UserData, error)

	// DeleteUser deletes the UserData for the user with the given username
	DeleteUser(username string) error

	// AddProject stores the given ProjectMetadata using the internal projectID, without overwriting
	AddProject(data *ProjectMetadata) error

	// SetProject stores the given ProjectMetadata using the internal projectID, overwriting if necessary
	SetProject(data *ProjectMetadata) error

	// GetProject retrieves the ProjectMetadata for the project with the given projectID
	GetProject(projectID int64) (*ProjectMetadata, error)

	// DeleteProject deletes the ProjectMetadata for the project with the given projectID
	DeleteProject(projectID int64) error

	// TODO(wongb): Create permission map by API version: -1 = blocked; 0 = <delete>; 1 = read; 2 = review; 3=write; 4=admin; 5=owner
	// SetProjectPermission stores (overwriting, if necessary) the new permission entry for the project with given projectID
	// If permissionLevel of the permission entry is 0, any existing entry for that user is deleted
	SetProjectPermission(projectID int64, permission *ProjectPermission) error

	// AddFile stores the given FileMetadata using the internal fileID, without overwriting
	AddFile(data *FileMetadata) error

	// SetFile stores the given FileMetadata using the internal fileID, overwriting if necessary
	SetFile(data *FileMetadata) error

	// GetFile retrieves the FileMetadata for the file with the given fileID
	GetFile(fileID int64) (*FileMetadata, error)

	// DeleteFile deletes the FileMetadata for the file with the given fileID
	DeleteFile(fileID int64) error
}

// InitRelationalStore Initializes the RelationalStore, or throws a fatal error if unsuccessful.
func InitRelationalStore(name string, cfg *config.ConnCfg) RelationalStore {
	name = strings.ToLower(name)

	if relationalStoreFactoryMap[name] == nil {
		utils.LogFatal("Configuration specified unknown RelationalStore", ErrFatalConfigurationErr, utils.LogFields{
			"RelationalStoreName": name,
		})
	}

	store := relationalStoreFactoryMap[name](cfg)
	store.Connect()

	return store
}
