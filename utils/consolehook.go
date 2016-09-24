package utils

import (
	"log"

	"github.com/Sirupsen/logrus"
)

// ConsoleHook to duplicate file logs to console
type ConsoleHook struct{}

// Fire is triggered on new log entries
func (hook *ConsoleHook) Fire(entry *logrus.Entry) error {
	str, err := entry.String()
	if err != nil {
		log.Printf("Unable to read entry: %v", err)
		return err
	}

	log.Print(str)
	return nil
}

// Levels returns all levels this hook should be registered to
func (hook *ConsoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
