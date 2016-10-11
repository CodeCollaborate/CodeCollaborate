package utils

import (
	"errors"
	"sync"
	"time"

	"fmt"
	"runtime"

	log "github.com/Sirupsen/logrus"
)

/**
 * Utility functions for the CodeCollaborate Server.
 */

// LogFields is the logrus.Fields type, but wrapped for convenience.
type LogFields log.Fields

func addFunc(fields LogFields) LogFields {
	if fields == nil {
		fields = LogFields{}
	}

	// pc[0] = runtime.Callers
	// pc[1] = addFunc
	// pc[2] = caller of addFunc (LogDebug, LogInfo, LogWarn...)
	// pc[3] = caller of logging functions
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fields["Location"] = fmt.Sprintf("%s:%d", file, line)
	return fields
}

func logWithFields(fields LogFields) *log.Entry {
	// type assertion is not needed b/c type conversion to log.Fields works for LogFields
	// see http://stackoverflow.com/questions/19577423/how-to-cast-to-a-type-alias-in-go
	return log.WithFields(log.Fields(fields))
}

// LogDebug logs the message, and fields given at DebugLevel
func LogDebug(msg string, fields LogFields) {
	funcFields := addFunc(fields)
	logWithFields(funcFields).Debug(msg)
}

// LogInfo logs the message, and fields given at InfoLevel
func LogInfo(msg string, fields LogFields) {
	funcFields := addFunc(fields)
	logWithFields(funcFields).Info(msg)
}

// LogWarn logs the message, and fields given at WarnLevel
func LogWarn(msg string, fields LogFields) {
	funcFields := addFunc(fields)
	logWithFields(funcFields).Warn(msg)
}

// LogError logs the message, error and fields given at ErrorLevel if the error != nil
func LogError(msg string, err error, fields LogFields) {
	if err == nil {
		return
	}

	funcFields := addFunc(fields)
	funcFields["error"] = err.Error()
	logWithFields(funcFields).Error(msg)
}

// LogFatal logs the message, error and fields given at FatalLevel if the error != nil
func LogFatal(msg string, err error, fields LogFields) {
	if err == nil {
		return
	}

	funcFields := addFunc(fields)
	funcFields["error"] = err.Error()
	logWithFields(funcFields).Fatal(msg)
}

// WaitTimeout will wait on the WaitGroup for a set amount of time,
// returning an error if the wait timed out.
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return nil // completed normally
	case <-time.After(timeout):
		return errors.New("Wait timed out") // timed out
	}
}
