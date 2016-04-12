package utils

import (
	"testing"
	"time"
)

func TestNewControl(t *testing.T) {
	control := NewControl()

	select {
	case control.Exit <- true:
	default:
		t.Fatal("Control's Exit Go-Channel is not buffered")
	}

	select {
	case control.Exit <- true:
		t.Fatal("Control's Exit Go-Channel buffer size > 1")
	default:
	}

	select {
	case <-control.Exit:
	default:
		t.Fatal("Did not get anything from Exit Go-Channel")
	}

	control.Ready.Done()
	if WaitTimeout(&control.Ready, time.Second) != nil {
		t.Fatal("Wait timed out - should have been ready")
	}
}

func TestControlTooManyDoneCalls(t *testing.T) {
	control := NewControl()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Ready called too many times - should have failed.")
		}
	}()

	control.Ready.Done()
	control.Ready.Done()
}
