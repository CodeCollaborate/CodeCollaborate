package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/CodeCollaborate/Server/utils"
)

/**
 * Configuration for the CodeCollaborate Server.
 */

var config *Config
var configDir = "./config"

// SetConfigDir sets config directory to be read from.
func SetConfigDir(dir string) {
	configDir = dir
}

// LoadConfig gets the configuration from the configDir, defaulting to ./config
// if not explicitly set by SetConfigDir. Will parse from json, and return
// a pointer to a Config struct, or error if it failed.
func LoadConfig() error {
	var err error
	utils.LogInfo("Reading Configuration", utils.LogFields{
		"ConfigDir": configDir,
	})
	config, err = parseConfig(configDir)

	if err == nil {
		utils.LogInfo("Loaded Configuration", utils.LogFields{
		//"ServerConfig": pretty.Sprint(config.ServerConfig),
		// TODO: remove secret fields from config and then print again
		})
		setLogLevel()
	}

	return err
}

func setLogLevel() {
	switch {
	case config.ServerConfig.LogLevel == "Panic":
		log.Info("Logger level set to Panic")
		log.SetLevel(log.PanicLevel)
	case config.ServerConfig.LogLevel == "Fatal":
		log.Info("Logger level set to Fatal")
		log.SetLevel(log.FatalLevel)
	case config.ServerConfig.LogLevel == "Error":
		log.Info("Logger level set to Error")
		log.SetLevel(log.ErrorLevel)
	case config.ServerConfig.LogLevel == "Warn":
		log.Info("Logger level set to Warn")
		log.SetLevel(log.WarnLevel)
	case config.ServerConfig.LogLevel == "Info":
		log.Info("Logger level set to Info")
		log.SetLevel(log.InfoLevel)
	case config.ServerConfig.LogLevel == "Debug":
		log.Info("Logger level set to Debug")
		log.SetLevel(log.DebugLevel)
	default:
		log.Info("Logger level set to Warn")
		log.SetLevel(log.WarnLevel) // Default to Warn
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{DisableColors: true})
}

// EnableLoggingToFile redirects logger output to a logfile in the config's LogDir.
// A new logfile will be created each time this method is called.
func EnableLoggingToFile(logDir string) {
	if logDir != "" {
		os.MkdirAll(logDir, 0755)
		logFile := filepath.Join(logDir, fmt.Sprintf("%d.%02d.%02d.%02d.%02d.log", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute()))

		log.Infof("Logging to %s", logFile)
		log.SetFormatter(&log.JSONFormatter{})
		f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Error("Failed to setup logging to file")
			return
		}
		log.SetOutput(f)
		log.AddHook(utils.MakeConsoleHook())
	} else {
		log.Error("No logging directory specified, logging to console")
	}
}

// GetConfig gets the configuration from the configDir, defaulting to ./config
// if not explicitly set by SetConfigDir. Will parse from json, and return
// a pointer to a Config struct, or error if it failed.
func GetConfig() *Config {
	return config
}
