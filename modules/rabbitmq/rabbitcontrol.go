package rabbitmq

import (
	"github.com/CodeCollaborate/Server/utils"
)

// RabbitControl controls the RabbitMQ subscriber and allows us to control rabbit subscriptions
type RabbitControl struct {
	*utils.Control
	SubChan chan Subscription
}

// NewControl creates a new control group, initialized to the not ready state
// (Ready WaitGroup semaphore to 1). Exit Go Channel is also created with a buffer of 1.
func NewControl() *RabbitControl {
	return &RabbitControl{
		SubChan: make(chan Subscription, 1),
		Control: utils.NewControl(),
	}
}

// Subscription is the object which is used to control subscription changes to the rabbitmq service
type Subscription struct {
	Channel     string
	IsSubscribe bool
}

// GetKey returns the key of the affected RabbitMQ subscription
func (s Subscription) GetKey() string {
	return s.Channel
}
