package dbfs

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/CodeCollaborate/Server/modules/config"
)

// DatabaseMock is a mock used for testing.
//
// fields are exported in case you're a masochist and wan't to initialize this by hand
type DatabaseMock struct {
	Users    map[string](UserMeta)
	Projects map[string]([]ProjectMeta)
	Files    map[int64]([]FileMeta)

	FileVersion map[int64]int64
	FileChanges map[int64][]string

	ProjectIDCounter int64
	FileIDCounter    int64

	File *[]byte

	// FunctionCallCount is the tracker of how many db functions are called
	FunctionCallCount int
}

// constructor

// NewDBMock is the constructor of the db mock object. It allows us to initialize the maps it holds.
func NewDBMock() *DatabaseMock {
	return &DatabaseMock{
		Users:       make(map[string](UserMeta)),
		Projects:    make(map[string]([]ProjectMeta)),
		Files:       make(map[int64]([]FileMeta)),
		FileVersion: make(map[int64]int64),
		FileChanges: make(map[int64][]string),
	}
}

// couchbase

// CloseCouchbase is a mock of the real implementation
func (dm *DatabaseMock) CloseCouchbase() error {
	dm.FunctionCallCount++
	return nil
}

// CBInsertNewFile is a mock of the real implementation
func (dm *DatabaseMock) CBInsertNewFile(fileID int64, version int64, changes []string) error {
	dm.FileVersion[fileID] = version
	dm.FileChanges[fileID] = changes
	dm.FunctionCallCount++
	return nil
}

// CBDeleteFile is a mock of the real implementation
func (dm *DatabaseMock) CBDeleteFile(fileID int64) error {
	dm.FunctionCallCount++
	return nil
}

// CBGetFileVersion is a mock of the real implementation
func (dm *DatabaseMock) CBGetFileVersion(fileID int64) (int64, error) {
	dm.FunctionCallCount++
	return dm.FileVersion[fileID], nil
}

// CBGetFileChanges is a mock of the real implementation
func (dm *DatabaseMock) CBGetFileChanges(fileID int64) ([]string, error) {
	dm.FunctionCallCount++
	return dm.FileChanges[fileID], nil
}

// CBAppendFileChange is a mock of the real implementation
func (dm *DatabaseMock) CBAppendFileChange(fileID int64, baseVersion int64, changes []string) (int64, error) {
	dm.FunctionCallCount++
	if dm.FileVersion[fileID] > baseVersion {
		return -1, ErrVersionOutOfDate
	}
	dm.FileVersion[fileID]++
	for _, change := range changes {
		dm.FileChanges[fileID] = append(dm.FileChanges[fileID], change)
	}
	return dm.FileVersion[fileID], nil
}

// mysql

// CloseMySQL is a mock of the real implementation
func (dm *DatabaseMock) CloseMySQL() error {
	dm.FunctionCallCount++
	return nil
}

// MySQLUserRegister is a mock of the real implementation
func (dm *DatabaseMock) MySQLUserRegister(user UserMeta) error {
	if _, ok := dm.Users[user.Username]; ok {
		return ErrNoDbChange
	}
	dm.Users[user.Username] = user
	dm.FunctionCallCount++
	return nil
}

// MySQLUserGetPass is a mock of the real implementation
func (dm *DatabaseMock) MySQLUserGetPass(username string) (string, error) {
	dm.FunctionCallCount++
	return dm.Users[username].Password, nil
}

// MySQLUserDelete is a mock of the real implementation
func (dm *DatabaseMock) MySQLUserDelete(username string, pass string) error {
	dm.FunctionCallCount++
	if _, ok := dm.Users[username]; ok {
		delete(dm.Users, username)
		return nil
	}
	return ErrNoDbChange
}

// MySQLUserLookup is a mock of the real implementation
func (dm *DatabaseMock) MySQLUserLookup(username string) (user UserMeta, err error) {
	dm.FunctionCallCount++
	if user, ok := dm.Users[username]; ok {
		return user, nil
	}
	return user, err
}

// MySQLUserProjects is a mock of the real implementation
func (dm *DatabaseMock) MySQLUserProjects(username string) ([]ProjectMeta, error) {
	dm.FunctionCallCount++
	return dm.Projects[username], nil
}

// MySQLProjectCreate is a mock of the real implementation
func (dm *DatabaseMock) MySQLProjectCreate(username string, projectName string) (int64, error) {
	dm.FunctionCallCount++

	permlvl, _ := config.GetPermissionLevel("Owner")
	proj := ProjectMeta{
		PermissionLevel: permlvl,
		ProjectID:       dm.ProjectIDCounter,
		Name:            projectName,
	}
	dm.ProjectIDCounter++
	dm.Projects[username] = append(dm.Projects[username], proj)
	return proj.ProjectID, nil
}

// MySQLProjectDelete is a mock of the real implementation
func (dm *DatabaseMock) MySQLProjectDelete(projectID int64, senderID string) error {
	dm.FunctionCallCount++
	// so this is kinda horrible, but it is easy to follow what's going on
	for username, projects := range dm.Projects {
		var index = int64(-1)
		for i, proj := range projects {
			if proj.ProjectID == projectID {
				index = int64(i)
			}
		}
		if int64(len(dm.Projects[username])) > index+1 {
			dm.Projects[username] = append(dm.Projects[username][:index], dm.Projects[username][(index+1):]...)
		} else {
			dm.Projects[username] = dm.Projects[username][:index]
		}
		delete(dm.Files, index)
	}
	return nil
}

// MySQLProjectGetFiles is a mock of the real implementation
func (dm *DatabaseMock) MySQLProjectGetFiles(projectID int64) ([]FileMeta, error) {
	dm.FunctionCallCount++
	return dm.Files[projectID], nil
}

// MySQLProjectGrantPermission is a mock of the real implementation
func (dm *DatabaseMock) MySQLProjectGrantPermission(projectID int64, grantUsername string, permissionLevel int8, grantedByUsername string) error {
	dm.FunctionCallCount++
	found := false

	// check if you're changing permission rather than adding
	for _, proj := range dm.Projects[grantUsername] {
		if proj.ProjectID == projectID {
			proj.PermissionLevel = permissionLevel
			found = true
			break
		}
	}
	if !found {
		// add if not found
		for _, proj := range dm.Projects[grantedByUsername] {
			if proj.ProjectID == projectID {
				dm.Projects[grantUsername] = append(dm.Projects[grantUsername], ProjectMeta{
					PermissionLevel: permissionLevel,
					ProjectID:       projectID,
					Name:            proj.Name,
				})
				found = true
			}
		}
		if !found {
			return ErrNoDbChange
		}
	}
	return nil
}

// MySQLProjectRevokePermission is a mock of the real implementation
func (dm *DatabaseMock) MySQLProjectRevokePermission(projectID int64, revokeUsername string, revokedByUsername string) error {
	dm.FunctionCallCount++
	index := -1
	for i, proj := range dm.Projects[revokeUsername] {
		if proj.ProjectID == projectID {
			index = i
		}
	}
	if index < 0 {
		return errors.New("project not found")
	}
	if len(dm.Projects[revokeUsername]) > index+1 {
		dm.Projects[revokeUsername] = append(
			dm.Projects[revokeUsername][:index],
			dm.Projects[revokeUsername][index+1:]...)
	} else {
		dm.Projects[revokeUsername] = dm.Projects[revokeUsername][:index]
	}
	return nil
}

// MySQLUserProjectPermissionLookup returns the permission level of `username` on the project with the given projectID
func (dm *DatabaseMock) MySQLUserProjectPermissionLookup(projectID int64, username string) (int8, error) {
	for _, proj := range dm.Projects[username] {
		if proj.ProjectID == projectID {
			return proj.PermissionLevel, nil
		}
	}
	return 0, ErrNoData
}

// MySQLProjectRename is a mock of the real implementation
func (dm *DatabaseMock) MySQLProjectRename(projectID int64, newName string) error {
	dm.FunctionCallCount++
	// so inefficient but whatever, it's a mock
	found := false
	for _, projects := range dm.Projects {
		for _, project := range projects {
			if project.ProjectID == projectID {
				project.Name = newName
				found = true
			}
		}
	}
	if !found {
		return ErrNoDbChange
	}
	return nil
}

// MySQLProjectLookup is a mock of the real implementation
func (dm *DatabaseMock) MySQLProjectLookup(projectID int64, username string) (name string, permissions map[string]ProjectPermission, err error) {
	dm.FunctionCallCount++
	permissions = make(map[string]ProjectPermission)
	for user, projects := range dm.Projects {
		for _, project := range projects {
			if project.ProjectID == projectID {
				name = project.Name
				// NOTE: we're not tracking who the permission is granted by (because I'm lazy)
				permissions[user] = ProjectPermission{
					PermissionLevel: project.PermissionLevel,
					Username:        user,
				}
			}
		}
	}
	return name, permissions, err
}

// MySQLFileCreate is a mock of the real implementation
func (dm *DatabaseMock) MySQLFileCreate(username string, filename string, relativePath string, projectID int64) (int64, error) {
	dm.FunctionCallCount++
	dm.FileIDCounter++
	dm.Files[projectID] = append(
		dm.Files[projectID],
		FileMeta{
			ProjectID:    projectID,
			CreationDate: time.Now(),
			Creator:      username,
			FileID:       dm.FileIDCounter,
			Filename:     filename,
			RelativePath: relativePath,
		})
	return dm.FileIDCounter, nil
}

// MySQLFileDelete is a mock of the real implementation
func (dm *DatabaseMock) MySQLFileDelete(fileID int64) error {
	dm.FunctionCallCount++
	for projectID, files := range dm.Files {
		for i, file := range files {
			if file.FileID == fileID {
				if len(dm.Files[projectID]) > i+1 {
					dm.Files[projectID] = append(
						dm.Files[projectID][:i],
						dm.Files[projectID][i+1:]...)
				} else {
					dm.Files[projectID] = dm.Files[projectID][:i]
				}
				delete(dm.FileVersion, fileID)
				return nil
			}
		}

	}
	return ErrNoDbChange
}

// MySQLFileMove is a mock of the real implementation
func (dm *DatabaseMock) MySQLFileMove(fileID int64, newPath string) error {
	dm.FunctionCallCount++
	for _, files := range dm.Files {
		for _, file := range files {
			if file.FileID == fileID {
				file.RelativePath = newPath
				return nil
			}
		}

	}
	return ErrNoDbChange
}

// MySQLFileRename is a mock of the real implementation
func (dm *DatabaseMock) MySQLFileRename(fileID int64, newName string) error {
	dm.FunctionCallCount++
	for _, files := range dm.Files {
		for _, file := range files {
			if file.FileID == fileID {
				file.Filename = newName
				return nil
			}
		}

	}
	return ErrNoDbChange
}

// MySQLFileGetInfo is a mock of the real implementation
func (dm *DatabaseMock) MySQLFileGetInfo(fileID int64) (filey FileMeta, err error) {
	dm.FunctionCallCount++
	for _, files := range dm.Files {
		for _, file := range files {
			if file.FileID == fileID {
				return file, err
			}
		}

	}
	return filey, err
}

// FileWrite is a mock of the real implementation
func (dm *DatabaseMock) FileWrite(relpath string, filename string, projectID int64, raw []byte) (string, error) {
	dm.FunctionCallCount++
	dm.File = &raw
	filepathy, err := calculateFilePath(relpath, filename, projectID)
	return filepath.Join(filepathy, filename), err
}

// FileDelete is a mock of the real implementation
func (dm *DatabaseMock) FileDelete(relpath string, filename string, projectID int64) error {
	dm.FunctionCallCount++
	dm.File = nil
	return nil
}

// FileRead is a mock of the real implementation
func (dm *DatabaseMock) FileRead(relpath string, filename string, projectID int64) (*[]byte, error) {
	dm.FunctionCallCount++
	if dm.File == nil {
		dm.File = &[]byte{}
	}
	return dm.File, nil
}

// FileMove moves a file form the starting path to the end path
func (dm *DatabaseMock) FileMove(startRelpath string, startFilename string, endRelpath string, endFilename string, projectID int64) error {
	return nil
}
