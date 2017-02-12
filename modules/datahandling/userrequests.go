package datahandling

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/CodeCollaborate/Server/modules/datahandling/messages"
	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
)

var userRequestsSetup = false

// initProjectRequests populates the requestMap from requestmap.go with the appropriate constructors for the project methods
func initUserRequests() {
	if userRequestsSetup {
		return
	}

	unauthenticatedRequestMap["User.Register"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(userRegisterRequest), req)
	}

	unauthenticatedRequestMap["User.Login"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(userLoginRequest), req)
	}

	authenticatedRequestMap["User.Delete"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(userDeleteRequest), req)
	}

	authenticatedRequestMap["User.Lookup"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(userLookupRequest), req)
	}

	authenticatedRequestMap["User.Projects"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(userProjectsRequest), req)
	}

	userRequestsSetup = true
}

// User.Register
type userRegisterRequest struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
	Password  string
	abstractRequest
}

func (f *userRegisterRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f userRegisterRequest) process(db dbfs.DBFS, ack func() error) ([]dhClosure, error) {
	defer ack() // ack regardless of success or failure

	hashed, err := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	newUser := dbfs.UserMeta{
		Username:  f.Username,
		FirstName: f.FirstName,
		LastName:  f.LastName,
		Email:     f.Email,
		Password:  string(hashed),
	}

	// TODO (non-immediate/required): password validation

	err = db.MySQLUserRegister(newUser)

	if err != nil {
		if err == dbfs.ErrNoDbChange {
			return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusNotFound, f.Tag)}}, err
		}
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}
	return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusSuccess, f.Tag)}}, err
}

// User.Login
type userLoginRequest struct {
	Username string
	Password string
	abstractRequest
}

func (f *userLoginRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f userLoginRequest) process(db dbfs.DBFS, ack func() error) ([]dhClosure, error) {
	defer ack() // ack regardless of success or failure

	hashed, err := db.MySQLUserGetPass(f.Username)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	if hashed == "" {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, f.Tag)}}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(f.Password)); err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, f.Tag)}}, err
	}

	signed, err := newAuthToken(f.Username)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	res := messages.Response{
		Status: messages.StatusSuccess,
		Tag:    f.Tag,
		Data: struct {
			Token string
		}{
			Token: signed,
		},
	}.Wrap()

	return []dhClosure{toSenderClosure{msg: res},
		// Subscribe user to their own username channel
		// TODO(wongb): What happens if they re-login? Or login as a different user?
		rabbitCommandClosure{
			Command: "Subscribe",
			Tag:     -1,
			Data: rabbitmq.RabbitQueueData{
				Key: rabbitmq.RabbitUserQueueName(f.Username),
			},
		},
	}, nil
}

// User.Delete
type userDeleteRequest struct {
	abstractRequest
}

func (f *userDeleteRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f userDeleteRequest) process(db dbfs.DBFS, ack func() error) ([]dhClosure, error) {
	defer ack() // ack regardless of success or failure

	deletedIDs, err := db.MySQLUserDelete(f.SenderID)

	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}
	closures := []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusSuccess, f.Tag)}}

	// TODO (shapiro): invalidate token
	for _, projectID := range deletedIDs {
		not := messages.Notification{
			Resource:   "Project",
			Method:     "Delete",
			ResourceID: projectID,
			Data:       struct{}{},
		}.Wrap()

		closures = append(closures, toRabbitChannelClosure{msg: not, key: rabbitmq.RabbitProjectQueueName(projectID)})
	}

	return closures, nil
}

// User.Lookup
type userLookupRequest struct {
	Usernames []string
	abstractRequest
}

func (f *userLookupRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f userLookupRequest) process(db dbfs.DBFS, ack func() error) ([]dhClosure, error) {
	defer ack() // ack regardless of success or failure

	users := make([]dbfs.UserMeta, len(f.Usernames))
	index := 0
	var erro error
	for _, username := range f.Usernames {
		usr, err := db.MySQLUserLookup(username)
		if err != nil {
			erro = err
		} else {
			users[index] = usr
			index++
		}
	}
	// shrink as needed
	users = users[:index]

	if len(users) < 0 {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, erro
	} else if erro != nil {
		// at least 1 value failed
		// return what we can but
		// tell the client whatever they don't get back failed
		res := messages.Response{
			Status: messages.StatusPartialFail,
			Tag:    f.Tag,
			Data: struct {
				Users []dbfs.UserMeta
			}{
				Users: users,
			},
		}.Wrap()
		return []dhClosure{toSenderClosure{msg: res}}, erro
	}

	res := messages.Response{
		Status: messages.StatusSuccess,
		Tag:    f.Tag,
		Data: struct {
			Users []dbfs.UserMeta
		}{
			Users: users,
		},
	}.Wrap()

	return []dhClosure{toSenderClosure{msg: res}}, erro
}

// User.Projects
type userProjectsRequest struct {
	abstractRequest
}

func (f *userProjectsRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f userProjectsRequest) process(db dbfs.DBFS, ack func() error) ([]dhClosure, error) {
	defer ack() // ack regardless of success or failure

	var errOut error
	projects, errOut := db.MySQLUserProjects(f.SenderID)

	resultData := make([]projectLookupResult, len(projects))

	i := 0
	for _, project := range projects {
		lookupResult, err := projectLookup(f.SenderID, project.ProjectID, db)

		if err != nil {
			utils.LogError("Project lookup error", err, utils.LogFields{
				"Resource":  f.Resource,
				"Method":    f.Method,
				"SenderID":  f.SenderID,
				"ProjectID": project.ProjectID,
			})
			errOut = err
		} else {
			resultData[i] = lookupResult
			i++
		}
	}
	resultData = resultData[:i]

	if errOut != nil {
		res := messages.Response{
			Status: messages.StatusPartialFail,
			Tag:    f.Tag,
			Data: struct {
				Projects []projectLookupResult
			}{
				Projects: resultData,
			},
		}.Wrap()
		return []dhClosure{toSenderClosure{msg: res}}, errOut
	}

	res := messages.Response{
		Status: messages.StatusSuccess,
		Tag:    f.Tag,
		Data: struct {
			Projects []projectLookupResult
		}{
			Projects: resultData,
		},
	}.Wrap()

	return []dhClosure{toSenderClosure{msg: res}}, nil
}
