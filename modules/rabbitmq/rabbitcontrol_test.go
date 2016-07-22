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
			case sub := <-control.Subscription:
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

	control.Subscription <- subsci
	go timeo(timeout)

	select {
	case <-received:
	// success
	case <-timeout:
		t.Fatal("control sygnal timed out")
	}

	control.Exit <- true
	go timeo(timeout)
	select {
	case <-doneTesting:
	// success
	case <-timeout:
		t.Fatal("control sygnal timed out")
	}

}

func timeo(timeout chan bool) {
	time.Sleep(1 * time.Second)
	timeout <- true
}
