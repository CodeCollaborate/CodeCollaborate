package datahandling

import (
	"fmt"

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"time"

	"errors"
	"strings"

	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
	"github.com/dgrijalva/jwt-go"
)

var privKey *ecdsa.PrivateKey

func init() {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic("Could not generate temporary signing key")
	}

	privKey = key
}

/**
 * Data Handling logic for the CodeCollaborate Server.
 */

// DataHandler handles the json data received from the WebSocket connection.
type DataHandler struct {
	MessageChan      chan<- rabbitmq.AMQPMessage
	SubscriptionChan chan<- rabbitmq.Subscription
	WebsocketID      uint64
	Db               dbfs.DBFS
}

// Handle takes the WebSocket Id, MessageType and message in byte-array form,
// processing the data, and updating DB/FS/RabbitMQ as needed.
func (dh DataHandler) Handle(messageType int, message []byte) error {
	fmt.Printf("Handling Message: %s\n", message)

	req, err := createAbstractRequest(message)
	if err != nil {
		utils.LogOnError(err, "Failed to parse json")
		return err
	}

	// automatically determines if the request is authenticated or not
	fullRequest, err := getFullRequest(req)

	var closures []dhClosure

	if err != nil {
		// TODO(shapiro): create response and notification factory
		if err == ErrAuthenticationFailed {
			utils.LogOnError(err, "User not logged in")
			closures = []dhClosure{toSenderClosure{msg: newEmptyResponse(unauthorized, req.Tag)}}
		} else {
			utils.LogOnError(err, "Failed to construct full request")
			closures = []dhClosure{toSenderClosure{msg: newEmptyResponse(unimplemented, req.Tag)}}
		}
	} else {
		closures, err = fullRequest.process(dh.Db)
		if err != nil {
			utils.LogOnError(err, "Failed to handle process request")
			// TODO: forward error message onto client? (or at least inform that error occurred)
		}
	}

	for _, closure := range closures {
		err := closure.call(dh)
		if err != nil {
			utils.LogOnError(err, "Failed to complete continuation")
		}
	}

	return err
}

type tokenPayload struct {
	Username     string
	CreationTime int64
	Validity     int64
}

func (tokenPayload) Valid() error {
	return nil
}

func authenticate(abs abstractRequest) error {
	token, err := jwt.ParseWithClaims(abs.SenderToken, &tokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("ParseWithClaims - Unexpected signing method: %v", token.Header["alg"])
		}
		return &privKey.PublicKey, nil
	})
	if err != nil {
		return fmt.Errorf("authenticate - failed to parse token: %s", err)
	}

	if claims, ok := token.Claims.(*tokenPayload); ok && token.Valid {
		// Check username is the same, and token is still valid
		if !strings.EqualFold(claims.Username, abs.SenderID) {
			return errors.New("authenticate - senderID did not match token username")
		}
		if time.Unix(claims.CreationTime, 0).After(time.Now()) {
			return errors.New("authenticate - token not valid yet")
		}
		if !time.Unix(claims.Validity, 0).After(time.Now()) {
			return errors.New("authenticate - expired token")
		}
		return nil
	}

	return errors.New("authenticate - claims struct was not of tokenPayload type")
}
