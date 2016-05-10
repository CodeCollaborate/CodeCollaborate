package config

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
	Name        string
	Port        uint16
	ProjectPath string
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
