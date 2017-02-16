package utils

import "sync"

/**
 * Control utilities for multithreading.
 */

// Control groups common multi-threading control variables, allowing for waiting until thread is ready,
// and setting exit flag.
type Control struct {
	sync.Mutex
	Ready    sync.WaitGroup
	Exit     chan bool
	shutdown sync.Once
	exited   bool
}

// Shutdown signals the Exit channel, and closes it once.
// Subsequent calls to this method do nothing
func (ctrl *Control) Shutdown() {
	ctrl.shutdown.Do(func() {
		ctrl.Lock()
		defer ctrl.Unlock()
		close(ctrl.Exit)
		ctrl.exited = true
	})
}

// HasExited checks to see if the control is still active
func (ctrl *Control) HasExited() bool {
	ctrl.Lock()
	defer ctrl.Unlock()
	return ctrl.exited
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
