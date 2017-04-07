package relationalstore

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datastore"
	"github.com/CodeCollaborate/Server/utils"
)

func init() {
	datastore.RegisterRelationalStore("mysql", NewMySQLRelationalStore)
}

// MySQLRelationalStore is a concrete implementation of the RelationalStore, using MySQL as the backing relational store of choice
type MySQLRelationalStore struct {
	cfg *config.ConnCfg
	db  *sql.DB
}

// NewMySQLRelationalStore creates a new instance of the MySQLRelationalStore, setting the configuration
func NewMySQLRelationalStore(cfg *config.ConnCfg) datastore.RelationalStore {
	return &datastore.RelationalStore(MySQLRelationalStore{
		cfg: cfg,
	})
}

// Connect starts this relationalStore's connection to the server
// If connect fails, it will throw a fatal error
func (store *MySQLRelationalStore) Connect() {
	if store.db != nil {
		err := store.db.Ping()
		if err != nil {
			store.db = nil
		}
	}

	if store.cfg.Schema == "" {
		utils.LogFatal("No MySQL schema found in config", datastore.ErrFatalServerErr, nil)
	}

	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds&parseTime=true",
		store.cfg.Username,
		store.cfg.Password,
		store.cfg.Host,
		store.cfg.Port,
		store.cfg.Schema,
		store.cfg.Timeout)
	db, err := sql.Open("mysql", connString)
	if err == nil {
		for i := 0; i < store.cfg.NumRetries; i++ {
			if err = db.Ping(); err != nil {
				time.Sleep(time.Duration(store.cfg.Timeout) * time.Second)
			} else {
				store.db = db
				break
			}
		}
	}

	if err != nil {
		store.db = nil
		utils.LogFatal("Unable to connect to MySQL", err, utils.LogFields{
			"Host":   store.cfg.Host,
			"Port":   store.cfg.Port,
			"Schema": store.cfg.Schema,
		})
	}
}

// Shutdown terminates this relationalStore's connection to the server
func (store *MySQLRelationalStore) Shutdown() {
	if store.db != nil {
		err := store.db.Close()
		store.db = nil
		utils.LogError("Failed to close MySQLRelationalStore database connection", err, utils.LogFields{
			"Host":   store.cfg.Host,
			"Port":   store.cfg.Port,
			"Schema": store.cfg.Schema,
		})
	} else {
		utils.LogError("Close called on uninitialized MySQLRelationalStore", datastore.ErrInternalServerErr, nil)
	}
}

// UserRegister registers a new user
func (store *MySQLRelationalStore) UserRegister(data datastore.UserData) error {

}

// UserLookup returns the user's information (except password)
func (store *MySQLRelationalStore) UserLookup(username string) (*datastore.UserData, error) {

}

// UserGetProjects returns the Project Metadata for all projects the user has permissions for
func (store *MySQLRelationalStore) UserGetProjects([]datastore.ProjectMetadata, error) {

}

// UserGetProjectPermissions returns the permissions that a user has for a given project
func (store *MySQLRelationalStore) UserGetProjectPermissions(int, error) {

}

// UserGetPassword retrieves the hash of the user's password
func (store *MySQLRelationalStore) UserGetPassword(username string) (string, error) {

}

// UserDelete deletes the given user
func (store *MySQLRelationalStore) UserDelete(username string) error {

}

// ProjectCreate creates a new project, assigning it a projectID
func (store *MySQLRelationalStore) ProjectCreate(username string, projectName string) (int64, error) {

}

// ProjectLookup returns the project metadata, including permissions
func (store *MySQLRelationalStore) ProjectLookup(projectID int64) (*datastore.ProjectMetadata, error) {

}

// ProjectGetFiles returns the list of files that the project contains
func (store *MySQLRelationalStore) ProjectGetFiles(projectID int64) ([]datastore.FileMetadata, error) {

}

// ProjectGrantPermissions grants the given permissions to the user with provided username
// TODO(wongb): CHANGE TO ALLOW BULK UPDATES
func (store *MySQLRelationalStore) ProjectGrantPermissions(projectID int64, username string, grantUsername string, permissionLevel int) error {

}

// ProjectRevokePermissions revokes all permissions from the user with provided username
// TODO(wongb): CHANGE TO ALLOW BULK UPDATES{
func (store *MySQLRelationalStore) ProjectRevokePermissions(projectID int64, username string) error {

}

// ProjectRename renames the project
func (store *MySQLRelationalStore) ProjectRename(projectID int64, newName string) error {

}

// ProjectDelete deletes a project from MySQL
func (store *MySQLRelationalStore) ProjectDelete(projectID int64) error {

}

// FileCreate creates a new file, assigning it a new fileID
func (store *MySQLRelationalStore) FileCreate(username string, projectID int64, filename string, relativePath string) (int64, error) {

}

// FileGet returns the metadata for the given fileID
func (store *MySQLRelationalStore) FileGet(fileID int64) (*datastore.FileMetadata, error) {

}

// FileMove updates the filepath and filename of for the given fileID
func (store *MySQLRelationalStore) FileMove(fileID int64, newRelativePath string, newName string) error {

}

// FileDelete deletes the stored metadata for the given fileID
func (store *MySQLRelationalStore) FileDelete(fileID int64) error {

}
