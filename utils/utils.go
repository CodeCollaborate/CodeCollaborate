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

//// read is this application's translation to the message format, scanning from
//// stdin.
//func Read(r io.Reader) <-chan string {
//	lines := make(chan string)
//	go func() {
//		defer close(lines)
//		scan := bufio.NewScanner(r)
//		for scan.Scan() {
//			lines <- scan.Text()
//		}
//	}()
//	return lines
//}

//// write is this application's subscriber of application messages, printing to
//// stdout.
//func Write(w io.Writer) chan<- string {
//	lines := make(chan string)
//	go func() {
//		for line := range lines {
//			fmt.Fprintln(w, line)
//		}
//	}()
//	return lines
//}

//func LogWithFunc(fields log.Fields) *log.Entry {
//	// pc[0] = runtime.Callers
//	// pc[1] = LogWithFunc
//	// pc[2] = caller of LogWithFunc
//	pc := make([]uintptr, 1)
//	runtime.Callers(2, pc)
//	f := runtime.FuncForPC(pc[0])
//	file, line := f.FileLine(pc[0])
//	fields["Location"] = fmt.Sprintf("%s:%d", file, line)
//	return log.WithFields(fields)
//}
func addFunc(fields log.Fields) log.Fields {
	if fields == nil {
		fields = log.Fields{}
	}

	// pc[0] = runtime.Callers
	// pc[1] = LogWithFunc
	// pc[2] = caller of LogWithFunc
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fields["Location"] = fmt.Sprintf("%s:%d", file, line)
	return fields
}

// LogDebug logs the message, and fields given at DebugLevel
func LogDebug(msg string, fields log.Fields) {
	funcFields := addFunc(fields)
	log.WithFields(funcFields).Debug(msg)
}

// LogInfo logs the message, and fields given at InfoLevel
func LogInfo(msg string, fields log.Fields) {
	funcFields := addFunc(fields)
	log.WithFields(funcFields).Info(msg)
}

// LogWarn logs the message, and fields given at WarnLevel
func LogWarn(msg string, fields log.Fields) {
	funcFields := addFunc(fields)
	log.WithFields(funcFields).Warn(msg)
}

// LogError logs the message, error and fields given at ErrorLevel
func LogError(msg string, err error, fields log.Fields) {
	funcFields := addFunc(fields)
	funcFields["error"] = err.Error()
	log.WithFields(funcFields).Error(msg)
}

// LogIfError logs the message, error and fields given at ErrorLevel if the error != nil
func LogIfError(msg string, err error, fields log.Fields) {
	if err == nil {
		return
	}

	LogError(msg, err, fields)
}

// LogFatal logs the message, error and fields given at FatalLevel if the error != nil
func LogFatal(msg string, err error, fields log.Fields) {
	funcFields := addFunc(fields)
	funcFields["error"] = err.Error()
	log.WithFields(funcFields).Fatal(msg)
}

// LogIfFatal logs the message, error and fields given at FatalLevel if the error != nil
func LogIfFatal(msg string, err error, fields log.Fields) {
	if err == nil {
		return
	}

	LogFatal(msg, err, fields)
}

//// FailOnError will throw a panic if err is not nil, printing msg and err to log
//// CAUTION: Will cause program to exit.
//func FailOnError(err error, msg string) {
//	if err != nil {
//		log.Printf("%s: %s", msg, err)
//		panic(err)
//	}
//}
//
//// LogOnError will print msg and err to log
//// CAUTION: Will log, and do nothing else.
//func LogOnError(err error, msg string) {
//	if err != nil {
//		log.Printf("%s: %s", msg, err)
//	}
//}

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
