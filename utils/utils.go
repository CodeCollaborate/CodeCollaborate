package utils

import (
	"errors"
	"log"
	"sync"
	"time"
)

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

// FailOnError will throw a panic if err is not nil, printing msg and err to log
// CAUTION: Will cause program to exit.
func FailOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		panic(err)
	}
}

// LogOnError will print msg and err to log
// CAUTION: Will log, and do nothing else.
func LogOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
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
