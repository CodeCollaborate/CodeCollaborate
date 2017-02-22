package config

import (
	"crypto/rsa"
	"time"

	"github.com/CodeCollaborate/Server/utils"
)

/**
 * Models for the configuration CodeCollaborate Server.
 */

// ConnCfgMap is a map of ConnCfgs, keyed on the connection name.
type ConnCfgMap map[string]ConnCfg

// Config contains all the different config items
type Config struct {
	ServerConfig     ServerCfg
	ConnectionConfig ConnCfgMap
}

// ServerCfg contains various config items that pertain to the server
type ServerCfg struct {
	Name            string
	Port            uint16
	ProjectPath     string
	DisableAuth     bool
	LogLevel        string
	TokenValidity   string
	MinBufferLength int
	MaxBufferLength int

	// RSA key
	RSAPrivateKeyLocation string
	RSAPrivateKeyPassword string
	rsaKey                *rsa.PrivateKey

	// Parsed validity
	tokenValidityDuration time.Duration
}

// TokenValidityDuration parses the given duration, and returns the time.Duration struct, or an error.
func (cfg ServerCfg) TokenValidityDuration() (time.Duration, error) {
	if cfg.tokenValidityDuration != 0 {
		return cfg.tokenValidityDuration, nil
	}

	var err error
	cfg.tokenValidityDuration, err = time.ParseDuration(cfg.TokenValidity)
	return cfg.tokenValidityDuration, err
}

// RSAKey returns the RSA key the server should use for signing tokens
func (cfg *ServerCfg) RSAKey() *rsa.PrivateKey {
	if cfg.rsaKey != nil {
		return cfg.rsaKey
	}

	var err error
	cfg.rsaKey, err = rsaConfigSetup(config.ServerConfig.RSAPrivateKeyLocation, config.ServerConfig.RSAPrivateKeyPassword)
	if err != nil {
		utils.LogFatal("Unable to load/generate RSA key", err, utils.LogFields{})
	}
	return cfg.rsaKey
}

// ConnCfg represents the information required to make a connection
type ConnCfg struct {
	Host       string
	Port       uint16
	Username   string
	Password   string
	Timeout    uint16
	NumRetries uint16
	Schema     string
}
