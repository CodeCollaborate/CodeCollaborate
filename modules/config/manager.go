package config

var config *Config
var configDir = "./config"

// SetConfigDir sets config directory to be read from.
func SetConfigDir(dir string) {
	configDir = dir
}

// InitConfig gets the configuration from the configDir, defaulting to ./config
// if not explicitly set by SetConfigDir. Will parse from json, and return
// a pointer to a Config struct, or error if it failed.
func InitConfig() error {
	var err error
	config, err = parseConfig(configDir)
	return err
}

// GetConfig gets the configuration from the configDir, defaulting to ./config
// if not explicitly set by SetConfigDir. Will parse from json, and return
// a pointer to a Config struct, or error if it failed.
func GetConfig() *Config {
	return config
}
