package config

import "time"

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
	Host            string
	Port            uint16
	ProjectPath     string
	DisableAuth     bool
	UseTLS          bool
	LogLevel        string
	TokenValidity   string
	MinBufferLength int
	MaxBufferLength int

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
