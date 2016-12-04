package utils

import "sync"

/**
 * Control utilities for multithreading.
 */

// Control groups common multi-threading control variables, allowing for waiting until thread is ready,
// and setting exit flag.
type Control struct {
	Ready    sync.WaitGroup
	Exit     chan bool
	shutdown sync.Once
}

// Shutdown signals the Exit channel, and closes it once.
// Subsequent calls to this method do nothing
func (ctrl *Control) Shutdown() {
	ctrl.shutdown.Do(func() {
		close(ctrl.Exit)
	})
}

// NewControl creates a new control group, initialized to the not ready state
// (Ready WaitGroup semaphore to 1). Exit Go Channel is also created with a buffer of 1.
func NewControl(wgCount int) *Control {
	control := Control{
		Ready: sync.WaitGroup{},
		Exit:  make(chan bool, 1),
	}
	control.Ready.Add(wgCount)
	return &control
}
