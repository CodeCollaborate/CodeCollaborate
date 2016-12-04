package utils

import (
	"sync"
	"testing"
	"time"
)

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
