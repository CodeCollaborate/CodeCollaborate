package utils

import (
	"errors"
	"sync"
	"testing"
	"time"
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

func TestWaitTimeout(t *testing.T) {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	err := WaitTimeout(wg, time.Second)
	if err == nil {
		t.Fatal("Should have failed waitTimeout")
	}

	wg.Done()
	err = WaitTimeout(wg, time.Second)
	if err != nil {
		t.Fatal("Should not have failed waitTimeout")
	}
}
