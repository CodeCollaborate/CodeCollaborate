package datahandling

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"strings"

	"github.com/CodeCollaborate/Server/modules/datahandling/messages"
	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
)

var privKey *ecdsa.PrivateKey

func init() {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.LogFatal("Failed to generate signing key", err, nil)

	privKey = key
}

/**
 * Data Handling logic for the CodeCollaborate Server.
 */

// DataHandler handles the json data received from the WebSocket connection.
type DataHandler struct {
	MessageChan chan<- rabbitmq.AMQPMessage
	Db          dbfs.DBFS
}

// Handle takes the MessageType and message in byte-array form,
// processing the data, and updating DBFS/RabbitMQ as needed.
func (dh DataHandler) Handle(message []byte, originID uint64, ack func() error) error {
	// Ignore any request that has a password JSON field
	if !strings.Contains(strings.ToLower(string(message)), "\"password\":") {
		utils.LogDebug("Received Message", utils.LogFields{
			"Message": string(message),
		})
	}

	req, err := createAbstractRequest(message)
	if err != nil {
		utils.LogError("Failed to parse json", err, nil) // Do not log request since passwords may be sent
		return err
	}

	// automatically determines if the request is authenticated or not
	fullRequest, err := getFullRequest(req)

	var closures []dhClosure

	if err != nil {
		// Ignore requests where there
		if req.Resource == "User" && (req.Method == "Register" || req.Method == "Login") {
			utils.LogError("getFullRequest failed for Register/Login", err, nil)
		} else {
			utils.LogError("getFullRequest failed", err, utils.LogFields{
				"Request": string(message),
			})
		}
		if err == ErrAuthenticationFailed {
			utils.LogDebug("User not logged in", utils.LogFields{
				"Resource": req.Resource,
				"Method":   req.Method,
			})
			closures = []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, req.Tag)}}
		} else {
			utils.LogDebug("No such resource/method", utils.LogFields{
				"Resource": req.Resource,
				"Method":   req.Method,
			})
			closures = []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnimplemented, req.Tag)}}
		}
	} else {
		closures, err = fullRequest.process(dh.Db, ack)
		if err != nil {
			utils.LogError("Failed to process request", err, utils.LogFields{
				"Resource": req.Resource,
				"Method":   req.Method,
			})
		}
	}

	for _, closure := range closures {
		err := closure.call(dh, originID)
		if err != nil {
			utils.LogError("Failed to complete continuation", err, utils.LogFields{
				"Resource": req.Resource,
				"Method":   req.Method,
			})
		}
	}

	return err
}
