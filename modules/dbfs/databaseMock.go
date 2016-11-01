package dbfs

import (
	"errors"
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
	Swp  *[]byte

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

// GetForScrunching gets all but the remainder entries for a file and creates a temp swp file.
// Returns the changes for scrunching, location of the swap file, and any errors
func (dm *DatabaseMock) GetForScrunching(fileID int64, remainder int) ([]string, []byte, error) {
	dm.FunctionCallCount++
	changes := dm.FileChanges[fileID]
	dm.Swp = new([]byte)
	return changes[0 : len(changes)-remainder], *dm.Swp, nil
}

// DeleteForScrunching deletes `num` elements from the front of `changes` for file with `fileID` and deletes the
// swp file
func (dm *DatabaseMock) DeleteForScrunching(fileID int64, num int) error {
	dm.FunctionCallCount++
	dm.File = dm.Swp
	dm.Swp = nil
	dm.FileChanges[fileID] = dm.FileChanges[fileID][num:]
	return nil
}

// PullFile pulls the changes and the file bytes from the databases
func (dm *DatabaseMock) PullFile(meta FileMeta) (*[]byte, []string, error) {
	dm.FunctionCallCount++
	changes := dm.FileChanges[meta.FileID]
	if dm.File == nil {
		return new([]byte), []string{}, ErrNoData
	}
	return dm.File, changes, nil
}

// CBAppendFileChange is a mock of the real implementation
func (dm *DatabaseMock) CBAppendFileChange(fileID int64, baseVersion int64, changes []string) (int64, error) {
	dm.FunctionCallCount++
	if dm.FileVersion[fileID] > baseVersion {
		return -1, ErrVersionOutOfDate
	}
	dm.FileVersion[fileID]++

	newChanges := append(dm.FileChanges[fileID], changes...)
	dm.FileChanges[fileID] = newChanges

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
func (dm *DatabaseMock) MySQLUserDelete(username string) ([]int64, error) {
	dm.FunctionCallCount += 2

	var deletedIDs []int64
	ownerPerm, err := config.PermissionByLabel("owner")
	if err != nil {
		return deletedIDs, err
	}

	for _, project := range dm.Projects[username] {
		if project.PermissionLevel == ownerPerm.Level {
			deletedIDs = append(deletedIDs, project.ProjectID)
		}
	}

	if _, ok := dm.Users[username]; ok {
		delete(dm.Users, username)
	} else {
		return []int64{}, ErrNoDbChange
	}

	// go through everyone's projects and delete ones which were owned by `username`
	// this is pretty gross... I can't wait for the refactor of this :/
	for username, projects := range dm.Projects {
		var indices []int
		for i, project := range projects {
			for _, deletedID := range deletedIDs {
				if project.ProjectID == deletedID {
					indices = append(indices, i)
				}
				break // continue to next i
			}
		}
		// remove the found projects
		initialLen := len(dm.Projects[username])
		for offset, i := range indices {
			projectsSlice := dm.Projects[username]
			if i != initialLen-1 {
				dm.Projects[username] = append(projectsSlice[0:i-offset], projectsSlice[i+1-offset:]...)
			} else {
				dm.Projects[username] = projectsSlice[0 : i-offset]
			}
		}
	}

	return deletedIDs, nil
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

	perm, _ := config.PermissionByLabel("owner")
	proj := ProjectMeta{
		PermissionLevel: perm.Level,
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
	dm.FunctionCallCount++
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
	return "./this_path_shouldnt_be_used_anywhere", nil
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
	dm.FunctionCallCount++
	// we only keep track of one file anyway
	return nil
}

// FileWriteToSwap writes the swapfile for the file with the given info
func (dm *DatabaseMock) FileWriteToSwap(relpath string, filename string, projectID int64, raw []byte) error {
	dm.FunctionCallCount++
	dm.Swp = &raw
	return nil
}
