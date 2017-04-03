package utils

import (
	"os"

	"github.com/Sirupsen/logrus"
)

// MakeConsoleHook creates the logrus hook to print to the terminal
func MakeConsoleHook() *ConsoleHook {
	consoleLog := logrus.New()
	// use stdout instead of stderr
	consoleLog.Out = os.Stdout
	// text formatter w/ color if it's supported
	consoleLog.Formatter = &logrus.TextFormatter{}
	// Set debug level of this console logger
	consoleLog.Level = logrus.DebugLevel
	return &ConsoleHook{logger: consoleLog}
}

// ConsoleHook to duplicate file logs to console
type ConsoleHook struct {
	logger *logrus.Logger
}

// Fire is triggered on new log entries
func (hook *ConsoleHook) Fire(entry *logrus.Entry) error {
	if entry.Logger.Level >= entry.Level {
		switch entry.Level {
		case logrus.DebugLevel:
			hook.logger.WithFields(entry.Data).Debug(entry.Message)
		case logrus.InfoLevel:
			hook.logger.WithFields(entry.Data).Info(entry.Message)
		case logrus.WarnLevel:
			hook.logger.WithFields(entry.Data).Warn(entry.Message)
		case logrus.ErrorLevel:
			hook.logger.WithFields(entry.Data).Error(entry.Message)
		case logrus.FatalLevel:
			hook.logger.WithFields(entry.Data).Fatal(entry.Message)
		case logrus.PanicLevel:
			hook.logger.WithFields(entry.Data).Panic(entry.Message)
		}
	}
	return nil
}

// Levels returns all levels this hook should be registered to
func (hook *ConsoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
