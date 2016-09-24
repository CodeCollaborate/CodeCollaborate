package rabbitmq

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"time"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/utils"
	"github.com/streadway/amqp"
)

var testExchange = AMQPExchCfg{
	ExchangeName: "TestExchange",
	Durable:      false,
}

func getRabbitMQConfig(t *testing.T) config.ConnCfg {
	config.SetConfigDir("../../config")
	err := config.LoadConfig()
	if err != nil {
		t.Fatal("Could not get connection config")
	}
	return config.GetConfig().ConnectionConfig["RabbitMQ"]
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

	err := SetupRabbitExchange(
		&AMQPConnCfg{
			ConnCfg: getRabbitMQConfig(t),
		},
	)
	if err != nil {
		t.Fatal("Failed to connect to Rabbit Exchange: Timed out")
	}

	for i := 0; i < 10; i++ {
		_, err := GetChannel()
		if err != nil {
			t.Fatal("Rabbit Exchange could not be setup.")
		}
	}
}

func TestSetupRabbitExchangeInvalidAddress(t *testing.T) {
	channelQueue = nil

	err := SetupRabbitExchange(
		&AMQPConnCfg{
			ConnCfg: config.ConnCfg{},
		},
	)
	if err == nil {
		t.Fatal("Should have failed to setup exchange")
	}

	_, err = GetChannel()
	if err == nil {
		t.Fatal("Channel should have failed; setup did not succeed.")
	}
}

func TestSetupRabbitExchangeFailConnection(t *testing.T) {
	channelQueue = nil

	config := config.ConnCfg{}
	config = getRabbitMQConfig(t)
	config.Username = ""
	config.Password = ""
	config.Timeout = 1
	config.NumRetries = 1

	err := SetupRabbitExchange(
		&AMQPConnCfg{
			ConnCfg: config,
		},
	)
	if err == nil {
		t.Fatal("Should have failed to setup exchange")
	}

	_, err = GetChannel()
	if err == nil {
		t.Fatal("Channel should have failed; setup did not succeed.")
	}
}

func TestSendMessage(t *testing.T) {
	channelQueue = nil

	err := SetupRabbitExchange(
		&AMQPConnCfg{
			ConnCfg: getRabbitMQConfig(t),
			Exchanges: []AMQPExchCfg{
				testExchange,
			},
		},
	)
	if err != nil {
		t.Fatal("Failed to connect to Rabbit Exchange: Timed out")
	}

	queueID := uint64(0)
	routingKey := fmt.Sprintf("%s-%d", hostname, queueID)
	doneTesting := make(chan bool, 1)
	defer close(doneTesting)

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
	subscriberControl := NewControl()
	publisherControl := utils.NewControl()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()

		RunSubscriber(&AMQPSubCfg{
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
				doneTesting <- true
				return nil
			},
			Control: subscriberControl,
		})
	}()
	subscriberControl.Ready.Wait()

	go func() {
		defer wg.Done()
		RunPublisher(&AMQPPubCfg{
			ExchangeName: testExchange.ExchangeName,
			Messages:     publisherMessages,
			Control:      publisherControl,
		})
	}()
	wg.Wait()

	select {
	case <-doneTesting:
		// success
	case <-time.After(time.Second * 5):
		t.Fatal("control signal timed out")
	}

}

func TestSubscription(t *testing.T) {
	channelQueue = nil

	err := SetupRabbitExchange(
		&AMQPConnCfg{
			ConnCfg: getRabbitMQConfig(t),
			Exchanges: []AMQPExchCfg{
				testExchange,
			},
		},
	)
	if err != nil {
		t.Fatal("Failed to connect to Rabbit Exchange: Timed out")
	}

	queueID := uint64(0)
	subscriptionChannel := "gene's project"

	doneTesting := make(chan bool, 1)
	defer close(doneTesting)

	TestMessage := AMQPMessage{
		Headers: map[string]interface{}{
			"Header1": "Value1",
			"Header2": "Value2",
			"Header3": "Value3",
		},
		RoutingKey:  subscriptionChannel,
		ContentType: "ContentType1",
		Persistent:  false,
		Message:     []byte("TestMessage1"),
	}

	publisherMessages := make(chan AMQPMessage, 1)
	subscriberControl := NewControl()
	publisherControl := utils.NewControl()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		wg.Done()
		RunSubscriber(&AMQPSubCfg{
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
				doneTesting <- true
				return nil
			},
			Control: subscriberControl,
		})
	}()

	subscriberControl.SubChan <- Subscription{
		Channel:     subscriptionChannel,
		IsSubscribe: true,
	}
	subscriberControl.Ready.Wait()

	go func() {
		wg.Done()
		RunPublisher(&AMQPPubCfg{
			ExchangeName: testExchange.ExchangeName,
			Messages:     publisherMessages,
			Control:      publisherControl,
		})
	}()
	wg.Wait()

	publisherMessages <- TestMessage

	select {
	case <-doneTesting:
	// success
	case <-time.After(time.Second * 5):
		t.Fatal("control signal timed out")
	}

}
