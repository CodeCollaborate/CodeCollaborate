package rabbitmq

import (
	"fmt"
	"testing"
	"time"
)

func TestNewControl(t *testing.T) {
	control := NewControl()
	doneTesting := make(chan bool, 1)
	received := make(chan bool, 1)
	defer close(doneTesting)
	defer close(received)

	timeout := make(chan bool, 1)
	defer close(timeout)

	subsci := Subscription{
		Channel:     "12345",
		IsSubscribe: true,
	}

	go func() {
		for {
			select {
			case sub := <-control.SubChan:
				if sub == subsci {
					received <- true
				} else {
					t.Fatal("somehow the message was corrupted")
				}
			case <-control.Exit:
				fmt.Println("received all jobs")
				doneTesting <- true
				return
			}
		}
	}()

	control.SubChan <- subsci

	select {
	case <-received:
	// success
	case <-time.After(time.Second * 5):
		t.Fatal("control sygnal timed out")
	}

	control.Exit <- true
	select {
	case <-doneTesting:
	// success
	case <-time.After(time.Second * 5):
		t.Fatal("control sygnal timed out")
	}

}
