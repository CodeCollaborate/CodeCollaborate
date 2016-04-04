package rabbitmq

import (
	"fmt"
	"github.com/CodeCollaborate/Server/utils"
	"github.com/streadway/amqp"
	"reflect"
	"sync"
	"testing"
)

var testExchange = ExchangeConfig{
	ExchangeName: "TestExchange",
	Durable:      false,
}

func TestGetChannel(t *testing.T) {
	var wg sync.WaitGroup

	_, err := GetChannel()
	if err == nil {
		t.Fatal("Channel should have failed; setup has not been called yet.")
	}

	channelQueue = make(chan *amqp.Channel)
	testChannel := &amqp.Channel{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		retVal, err := GetChannel()
		if err != nil {
			t.Fatal("GetChannel threw error for Channel in queue")
		}
		if testChannel != retVal {
			t.Fatal("GetChannel did not return same RabbitMQ Channel that was put inserted into channelQueue")
		}
	}()

	channelQueue <- testChannel

	wg.Wait()
}

func TestSetupRabbitExchange(t *testing.T) {
	channelQueue = nil

	SetupRabbitExchange(
		&ConnectionConfig{
			Host: "localhost",
			Port: 5672,
			User: "guest",
			Pass: "guest",
			Exchanges: []ExchangeConfig{
				testExchange,
			},
		},
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			_, err := GetChannel()
			if err != nil {
				t.Fatal("Rabbit Exchange could not be setup.")
			}
		}
	}()
	wg.Wait()
}

func TestSendMessage(t *testing.T) {
	channelQueue = nil

	SetupRabbitExchange(
		&ConnectionConfig{
			Host: "localhost",
			Port: 5672,
			User: "guest",
			Pass: "guest",
			Exchanges: []ExchangeConfig{
				testExchange,
			},
		},
	)

	queueID := uint64(0)
	routingKey := fmt.Sprintf("%s-%d", hostname, queueID)

	TestMessage := AMQPMessage{
		Headers: map[string]interface{}{
			"Header1": "Value1",
			"Header2": "Value2",
			"Header3": "Value3",
		},
		RoutingKey:  routingKey,
		ContentType: "ContentType1",
		Persistent:  false,
		Message:     []byte("TestMessage1"),
	}

	publisherMessages := make(chan AMQPMessage, 1)
	publisherMessages <- TestMessage
	subscriberControl := utils.NewControl()
	publisherControl := utils.NewControl()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()

		RunSubscriber(&SubscriberConfig{
			ExchangeName: testExchange.ExchangeName,
			QueueID:      queueID,
			Keys:         []string{},
			IsWorkQueue:  false,
			HandleMessageFunc: func(msg AMQPMessage) error {
				subscriberControl.Exit <- true
				publisherControl.Exit <- true
				if !reflect.DeepEqual(msg, TestMessage) {
					t.Fatal("Sent message does not equal received message")
				}
				return nil
			},
			Control: subscriberControl,
		})
	}()
	subscriberControl.Ready.Wait()

	go func() {
		defer wg.Done()
		RunPublisher(&PublisherConfig{
			ExchangeName: testExchange.ExchangeName,
			Messages:     publisherMessages,
			Control:      publisherControl,
		})
	}()
	wg.Wait()
}
