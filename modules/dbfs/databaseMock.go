package dbfs

import (
	"errors"
	"path/filepath"
	"time"
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

	ProjectIDCounter int
	FileIDCounter    int

	File *[]byte
}

// couchbase

// CloseCouchbase is a mock of the real implementation
func (dm *DatabaseMock) CloseCouchbase() error {
	return nil
}

// CBInsertNewFile is a mock of the real implementation
func (dm *DatabaseMock) CBInsertNewFile(fileID int64, version int64, changes []string) error {
	dm.FileVersion[fileID] = version
	dm.FileChanges[fileID] = changes
	return nil
}

// CBDeleteFile is a mock of the real implementation
func (dm *DatabaseMock) CBDeleteFile(fileID int64) error {
	return nil
}

// CBGetFileVersion is a mock of the real implementation
func (dm *DatabaseMock) CBGetFileVersion(fileID int64) (int64, error) {
	return dm.FileVersion[fileID], nil
}

// CBGetFileChanges is a mock of the real implementation
func (dm *DatabaseMock) CBGetFileChanges(fileID int64) ([]string, error) {
	return dm.FileChanges[fileID], nil
}

// CBAppendFileChange is a mock of the real implementation
func (dm *DatabaseMock) CBAppendFileChange(fileID int64, baseVersion int64, changes []string) (int64, error) {
	dm.FileVersion[fileID]++
	for _, change := range changes {
		dm.FileChanges[fileID] = append(dm.FileChanges[fileID], change)
	}
	return dm.FileVersion, nil
}

// mysql

// CloseMySQL is a mock of the real implementation
func (dm *DatabaseMock) CloseMySQL() error {
	return nil
}

// MySQLUserRegister is a mock of the real implementation
func (dm *DatabaseMock) MySQLUserRegister(user UserMeta) error {
	dm.Users[user.Username] = user
	return nil
}

// MySQLUserGetPass is a mock of the real implementation
func (dm *DatabaseMock) MySQLUserGetPass(username string) (string, error) {
	return dm.Users[username].Password, nil
}

// MySQLUserDelete is a mock of the real implementation
func (dm *DatabaseMock) MySQLUserDelete(username string, pass string) error {
	if dm.Users[username] != nil {
		delete(dm.Users, username)
		return nil
	}
	return ErrNoDbChange
}

// MySQLUserLookup is a mock of the real implementation
func (dm *DatabaseMock) MySQLUserLookup(username string) (user UserMeta, err error) {
	user = dm.Users[username]
	if user != nil {
		return user, nil
	}
	return nil, err
}

// MySQLUserProjects is a mock of the real implementation
func (dm *DatabaseMock) MySQLUserProjects(username string) ([]ProjectMeta, error) {
	return dm.Projects[username], nil
}

// MySQLProjectCreate is a mock of the real implementation
func (dm *DatabaseMock) MySQLProjectCreate(username string, projectName string) (int64, error) {
	proj := ProjectMeta{
		PermissionLevel: 10,
		ProjectID:       dm.ProjectIDCounter,
		ProjectName:     projectName,
	}
	dm.ProjectIDCounter++
	dm.Projects[username] = append(dm.Projects[username], proj)
	return proj.ProjectID, nil
}

// MySQLProjectDelete is a mock of the real implementation
func (dm *DatabaseMock) MySQLProjectDelete(projectID int64, senderID string) error {
	// so this is kinda horrible, but it is easy to follow what's going on
	for username, projects := range dm.Projects {
		index := -1
		for i, proj := range projects {
			if proj.ProjectID == projectID {
				index = i
			}
		}
		if len(dm.Projects[username]) > index+1 {
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
	return dm.Files[projectID], nil
}

// MySQLProjectGrantPermission is a mock of the real implementation
func (dm *DatabaseMock) MySQLProjectGrantPermission(projectID int64, grantUsername string, permissionLevel int, grantedByUsername string) error {
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
					ProjectName:     proj.ProjectName,
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

// MySQLProjectRename is a mock of the real implementation
func (dm *DatabaseMock) MySQLProjectRename(projectID int64, newName string) error {
	// so inefficient but whatever, it's a mock
	found := false
	for _, projects := range dm.Projects {
		for _, project := range projects {
			if project.ProjectID == projectID {
				project.ProjectName = newName
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
	for user, projects := range dm.Projects {
		for _, project := range projects {
			if project.ProjectID == projectID {
				name = project.ProjectName
				// NOTE: we're not tracking who the permission is granted by (because I'm lazy)
				permissions = append(permissions, ProjectPermission{
					PermissionLevel: project.PermissionLevel,
					Username:        user,
				})
			}
		}
	}
	return name, permissions, err
}

// MySQLFileCreate is a mock of the real implementation
func (dm *DatabaseMock) MySQLFileCreate(username string, filename string, relativePath string, projectID int64) (int64, error) {
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
	for _, files := range dm.Files {
		for _, file := range files {
			if file.FileID == fileID {
				file.RelativePath = newPath
			}
		}

	}
	return ErrNoDbChange
}

// MySQLFileRename is a mock of the real implementation
func (dm *DatabaseMock) MySQLFileRename(fileID int64, newName string) error {
	for _, files := range dm.Files {
		for _, file := range files {
			if file.FileID == fileID {
				file.Filename = newName
			}
		}

	}
	return ErrNoDbChange
}

// MySQLFileGetInfo is a mock of the real implementation
func (dm *DatabaseMock) MySQLFileGetInfo(fileID int64) (FileMeta, error) {
	for _, files := range dm.Files {
		for _, file := range files {
			if file.FileID == fileID {
				return file, nil
			}
		}

	}
	return nil, nil
}

// FileWrite is a mock of the real implementation
func (dm *DatabaseMock) FileWrite(relpath string, filename string, projectID int64, raw []byte) (string, error) {
	dm.File = &raw
	return filepath.Join(calculateFilePathPath(relpath, filename, projectID), filename), nil
}

// FileDelete is a mock of the real implementation
func (dm *DatabaseMock) FileDelete(relpath string, filename string, projectID int64) error {
	dm.File = nil
	return nil
}

// FileRead is a mock of the real implementation
func (dm *DatabaseMock) FileRead(relpath string, filename string, projectID int64) (*[]byte, error) {
	return dm.File, nil
}
