package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

// read is this application's translation to the message format, scanning from
// stdin.
func Read(r io.Reader) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			lines <- scan.Text()
		}
	}()
	return lines
}

// write is this application's subscriber of application messages, printing to
// stdout.
func Write(w io.Writer) chan<- string {
	lines := make(chan string)
	go func() {
		for line := range lines {
			fmt.Fprintln(w, line)
		}
	}()
	return lines
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		panic(err)
	}
}

func LogOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
