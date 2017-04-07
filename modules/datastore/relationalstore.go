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

	// UserRegister registers a new user
	UserRegister(data UserData) error

	// UserLookup returns the user's information (except password)
	UserLookup(username string) (*UserData, error)

	// UserGetProjects returns the Project Metadata for all projects the user has permissions for
	UserGetProjects([]ProjectMetadata, error)

	// UserGetProjectPermissions returns the permissions that a user has for a given project
	UserGetProjectPermissions(int, error)

	// UserGetPassword retrieves the hash of the user's password
	UserGetPassword(username string) (string, error)

	// UserDelete deletes the given user
	UserDelete(username string) error

	// ProjectCreate creates a new project, assigning it a projectID
	ProjectCreate(username string, projectName string) (int64, error)

	// ProjectLookup returns the project metadata, including permissions
	ProjectLookup(projectID int64) (*ProjectMetadata, error)

	// ProjectGetFiles returns the list of files that the project contains
	ProjectGetFiles(projectID int64) ([]FileMetadata, error)

	// ProjectGrantPermissions grants the given permissions to the user with provided username
	ProjectGrantPermissions(projectID int64, username string, grantUsername string, permissionLevel int) error // TODO(wongb): CHANGE TO ALLOW BULK UPDATES

	// ProjectRevokePermissions revokes all permissions from the user with provided username
	ProjectRevokePermissions(projectID int64, username string) error // TODO(wongb): CHANGE TO ALLOW BULK UPDATES

	// ProjectRename renames the project
	ProjectRename(projectID int64, newName string) error

	// ProjectDelete deletes a project from MySQL
	ProjectDelete(projectID int64) error

	// FileCreate creates a new file, assigning it a new fileID
	FileCreate(username string, projectID int64, filename string, relativePath string) (int64, error)

	// FileGet returns the metadata for the given fileID
	FileGet(fileID int64) (*FileMetadata, error)

	// FileMove updates the filepath and filename of for the given fileID
	FileMove(fileID int64, newRelativePath string, newName string) error

	// FileDelete deletes the stored metadata for the given fileID
	FileDelete(fileID int64) error
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
