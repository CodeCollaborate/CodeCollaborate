package rabbitmq

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"os"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/utils"
	"github.com/streadway/amqp"
)

var rabbitConfig config.ConnCfg

var testExchange = AMQPExchCfg{
	ExchangeName: "TestExchange",
	Durable:      false,
}

func TestMain(m *testing.M) {
	config.SetConfigDir("../../config")
	config.InitConfig()
	rabbitConfig = config.GetConfig().ConnectionConfig["RabbitMQ"]
	os.Exit(m.Run())
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
			ConnCfg: rabbitConfig,
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
	config = rabbitConfig
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
			ConnCfg: rabbitConfig,
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
}
