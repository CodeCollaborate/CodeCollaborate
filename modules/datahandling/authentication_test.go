package datahandling

import (
	"crypto/rsa"
	"math/rand"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"

	"github.com/CodeCollaborate/Server/modules/config"
)

func TestAuthenticateRandomUsernames(t *testing.T) {
	for i := 0; i < 100; i++ {
		username := randomString(20)

		token := jwt.NewWithClaims(jwt.SigningMethodRS512, tokenPayload{
			Username:     username,
			CreationTime: time.Now().Unix(),
			Validity:     time.Now().Add(1 * time.Hour).Unix(),
		})

		key, err := config.GenRSA(1024) // make it small so it's faster
		assert.Nil(t, err, "error generating rsa")
		assert.NoError(t, key.Validate(), "Unable to validate RSA key")

		signed, err := token.SignedString(key)
		if err != nil {
			t.Fatal(err)
		}

		req := abstractRequest{
			SenderID:    username,
			SenderToken: signed,
		}

		assert.Nil(t, authenticate(req))
	}
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		desc     string
		senderID string
		token    string
		err      string
	}{
		{
			desc:     "Valid token",
			senderID: "TestUser1",
			token: signedTokenOrDie(t,
				"TestUser1",
				time.Now().Unix(),
				time.Now().Add(1*time.Minute).Unix(),
				getRsaKey(),
			),
		},
		{
			desc:     "Token username case different",
			senderID: "testUser1",
			token: signedTokenOrDie(t,
				"TestUser1",
				time.Now().Unix(),
				time.Now().Add(1*time.Minute).Unix(),
				getRsaKey(),
			),
		},
		{
			desc:     "Token username does not match senderID",
			senderID: "user1",
			token: signedTokenOrDie(t,
				"TestUser1",
				time.Now().Unix(),
				time.Now().Add(1*time.Minute).Unix(),
				getRsaKey(),
			),
			err: "authenticate - senderID did not match token username",
		},
		{
			desc:     "Token username does not match senderID reversed",
			senderID: "TestUser1",
			token: signedTokenOrDie(t,
				"user1",
				time.Now().Unix(),
				time.Now().Add(1*time.Minute).Unix(),
				getRsaKey(),
			),
			err: "authenticate - senderID did not match token username",
		},
		{
			desc:     "Expired token",
			senderID: "TestUser1",
			token: signedTokenOrDie(t,
				"TestUser1",
				time.Now().Unix(),
				time.Now().Add(-1*time.Minute).Unix(),
				getRsaKey(),
			),
			err: "authenticate - expired token",
		},
		{
			desc:     "Expired token",
			senderID: "TestUser1",
			token: signedTokenOrDie(t,
				"TestUser1",
				time.Now().Add(1*time.Minute).Unix(),
				time.Now().Add(1*time.Minute).Unix(),
				getRsaKey(),
			),
			err: "authenticate - token not valid yet",
		},
		{
			desc:     "Invalid token",
			senderID: "TestUser1",
			token:    "Invalid token",
			err:      "authenticate - failed to parse token: token contains an invalid number of segments",
		},
		{
			desc:     "Invalid token",
			senderID: "TestUser1",
			token:    "Invalid. .token",
			err:      "authenticate - failed to parse token: unexpected end of JSON input",
		},
		{
			desc:     "Wrong key",
			senderID: "TestUser1",
			token:    "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6IlRlc3RVc2VyMSIsIkNyZWF0aW9uVGltZSI6MTQ3NDQzOTQ2NywiVmFsaWRpdHkiOjE0NzQ0Mzk0NjZ9.6HK6VyBbXqIwJnRD2fCWIWTM6q466o56QhftJgcywawoi43kN-gEiwdx7K2EaGrDzxz9yd5jJHib_3n-_P9rxA",
			err:      "authenticate - failed to parse token: crypto/rsa: verification error",
		},
		{
			desc:     "Checksum mismatch",
			senderID: "TestUser1",
			token:    "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6IlRlc3RVc2VyMSIsIkNyZWF0aW9uVGltZSI6MTQ3NDQzOTQ2NywiVmFsaWRpdHkiOjE0NzQ0Mzk0NjZ9.6HK6VyBbXqIwJnRD2fCWIWTM6q466o56QhftJgcywawoi43kN-gEiwdx7K2EaGrDzxz9yd5jJw3b_3n-_P9rxA",
			err:      "authenticate - failed to parse token: crypto/rsa: verification error",
		},
	}

	for _, test := range tests {
		req := abstractRequest{
			SenderID:    test.senderID,
			SenderToken: test.token,
		}

		err := authenticate(req)
		if test.err != "" {
			if err == nil {
				t.Errorf("TestAuthenticate[%s]: Expected error: %q", test.desc, test.err)
				continue
			}
			if want, got := test.err, err.Error(); want != got {
				t.Error(pretty.Sprintf("TestAuthenticate[%s]: Expected %q, got %q. Diffs: %v", test.desc, want, got, pretty.Diff(want, got)))
				continue
			}
		} else if err != nil {
			t.Errorf("TestAuthenticate[%s]: Unexpected error: %q", test.desc, err)
			continue
		}
	}
}

func signedTokenOrDie(t *testing.T, username string, creationDate, validity int64, key *rsa.PrivateKey) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, tokenPayload{
		Username:     username,
		CreationTime: creationDate,
		Validity:     validity,
	})

	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}

	return signed
}

func randomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_- "
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
