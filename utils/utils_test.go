package utils

import (
	"errors"
	"testing"
)

func TestFailOnError(t *testing.T) {
	err := errors.New("I'm an error")

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Where's the panic?")
		}
	}()
	FailOnError(err, "Fail me")
}

func TestLogOnError(t *testing.T) {
	err := errors.New("I'm also an error")
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Why did you panic?")
		}
	}()
	LogOnError(err, "Fail me also")
}
