package config

import "time"

/**
 * Models for the configuration CodeCollaborate Server.
 */

// Config contains all the different config items
type Config struct {
	ServerConfig    *ServerCfg
	DataStoreConfig *DataStoreCfg
}

// ServerCfg contains various config items that pertain to the server
type ServerCfg struct {
	Name            string
	Host            string
	Port            int
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
	Port       int
	Username   string
	Password   string
	Timeout    int
	NumRetries int
	Schema     string
}

// DataStoreCfg represents the information required to initialize the different datastores
type DataStoreCfg struct {
	BucketStoreName string
	BucketStoreCfg  *ConnCfg

	DocumentStoreName string
	DocumentStoreCfg  *ConnCfg

	RelationalStoreName string
	RelationalStoreCfg  *ConnCfg

	MessageBrokerName string
	MessageBrokerCfg  *ConnCfg
}
