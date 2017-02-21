package datahandling

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/utils"
)

type tokenPayload struct {
	Username     string
	CreationTime int64
	Validity     int64
}

var rsaKey *rsa.PrivateKey

func getRsaKey() *rsa.PrivateKey {
	if rsaKey != nil {
		return rsaKey
	}

	cfg := config.GetConfig()
	if cfg == nil {
		utils.LogFatal("Failed to load RSA key from config", errors.New("config not initialized"), utils.LogFields{})
	}

	rsaKey = cfg.ServerConfig.RSAKey()
	return rsaKey
}

// Valid is the (unused) method to determine if the token is valid. however, since we need to have a reference
// to the abstract request, we cannot do validation here. Token validation has been shifted to the authenticate
// method. This is here for conformance to the token.Claims interface.
func (tokenPayload) Valid() error {
	return nil
}

func authenticate(abs abstractRequest) error {
	token, err := jwt.ParseWithClaims(abs.SenderToken, &tokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("ParseWithClaims - Unexpected signing method: %v", token.Header["alg"])
		}
		return &getRsaKey().PublicKey, nil
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

func newAuthToken(username string) (string, error) {
	tokenValidityDuration, err := config.GetConfig().ServerConfig.TokenValidityDuration()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, tokenPayload{
		Username:     username,
		CreationTime: time.Now().Unix(),
		Validity:     time.Now().Add(tokenValidityDuration).Unix(),
	})

	return token.SignedString(getRsaKey())
}
