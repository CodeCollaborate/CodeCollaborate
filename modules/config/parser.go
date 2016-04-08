package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func parseConfig(filepath string) (*Config, error) {

	serverConfig, err := parseServerConfig(filepath)
	if err != nil {
		return nil, err
	}
	connectionConfig, err := parseConnectionConfig(filepath)
	if err != nil {
		return nil, err
	}

	config := &Config{
		ServerConfig:     *serverConfig,
		ConnectionConfig: *connectionConfig,
	}

	return config, nil
}

func parseServerConfig(configDir string) (*ServerCfg, error) {
	file, err := os.Open(filepath.Join(configDir, "server.cfg"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &ServerCfg{}

	jsonParser := json.NewDecoder(file)
	if err := jsonParser.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

func parseConnectionConfig(configDir string) (*ConnCfgMap, error) {
	file, err := os.Open(filepath.Join(configDir, "conn.cfg"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &ConnCfgMap{}

	jsonParser := json.NewDecoder(file)
	if err := jsonParser.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
