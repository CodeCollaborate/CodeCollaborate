package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthConfigSetup(t *testing.T) {
	SetConfigDir("../../config")
	err := LoadConfig()
	assert.NoError(t, err, "error initializing config needed for password")

	key, err := rsaConfigSetup("../../config/id_rsa", config.ServerConfig.RSAPrivateKeyPassword)
	assert.Nil(t, err, "error loading rsa key")
	assert.NotNil(t, key, "key was nil")
	assert.NoError(t, key.Validate(), "key could not be validated")
}
