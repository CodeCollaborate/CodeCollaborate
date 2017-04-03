package rabbitmq

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/utils"
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
	routingKey := LocalWebsocketName(queueID)
	doneTesting := make(chan bool, 1)
	defer close(doneTesting)

	TestMessage := AMQPMessage{
		Headers: map[string]interface{}{
			"Header1": "Value1",
			"Header2": "Value2",
			"Header3": "Value3",
		},
		RoutingKey:  routingKey,
		ContentType: ContentTypeMsg,
		Persistent:  false,
		Message:     []byte("TestMessage1"),
	}

	pubSubCtrl := utils.NewControl(2)

	pubSubCfg := &AMQPPubSubCfg{
		ExchangeName: testExchange.ExchangeName,
		SubCfg: &AMQPSubCfg{
			QueueName:   LocalWebsocketName(queueID),
			Keys:        []string{},
			IsWorkQueue: false,
			HandleMessageFunc: func(msg AMQPMessage) error {
				pubSubCtrl.Shutdown()

				assert.NotNil(t, msg.Ack, "Message does not have Ack")
				assert.NotNil(t, msg.Nack, "Message does not have Nack")
				assert.NotNil(t, msg.ErrHandler, "Message does not have Nack")

				err := msg.Ack()
				assert.Nil(t, err)

				// the messages aren't EXACT copies of each other b/c Nack, Ack, and ErrHandler are nil in TestMessage
				// so we need to fix this
				msg.Ack = nil
				msg.Nack = nil
				msg.ErrHandler = nil

				assert.EqualValues(t, TestMessage, msg, "Sent message does not equal received message")
				doneTesting <- true
				return nil
			},
		},
		PubCfg: &AMQPPubCfg{
			Messages: make(chan AMQPMessage, 1),
		},
		Control: pubSubCtrl,
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		RunSubscriber(pubSubCfg)
	}()

	go func() {
		defer wg.Done()
		RunPublisher(pubSubCfg)
	}()
	pubSubCfg.Control.Ready.Wait()

	// Send message
	pubSubCfg.PubCfg.Messages <- TestMessage

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
	subChanKey := "gene's project"

	doneTesting := make(chan bool, 1)
	defer close(doneTesting)

	TestMessage := AMQPMessage{
		Headers: map[string]interface{}{
			"Header1": "Value1",
			"Header2": "Value2",
			"Header3": "Value3",
		},
		RoutingKey:  subChanKey,
		ContentType: ContentTypeMsg,
		Persistent:  false,
		Message:     []byte("TestMessage1"),
	}

	msgJSON, err := json.Marshal(&RabbitCommandStruct{
		Command: "Subscribe",
		Tag:     1,
		Data: &RabbitQueueData{
			Key: subChanKey,
		},
	})

	TestSubscription := AMQPMessage{
		Headers: map[string]interface{}{
			"Header1": "Value1",
			"Header2": "Value2",
			"Header3": "Value3",
		},
		RoutingKey:  LocalWebsocketName(queueID),
		ContentType: ContentTypeCmd,
		Persistent:  false,
		Message:     msgJSON,
	}

	pubSubCtrl := utils.NewControl(2)
	subWg := &sync.WaitGroup{}
	subWg.Add(1)

	pubSubCfg := &AMQPPubSubCfg{
		ExchangeName: testExchange.ExchangeName,
		SubCfg: &AMQPSubCfg{
			QueueName:   LocalWebsocketName(queueID),
			Keys:        []string{},
			IsWorkQueue: false,
			HandleMessageFunc: func(msg AMQPMessage) error {
				switch msg.ContentType {
				case ContentTypeCmd:
					defer subWg.Done()
					rch := RabbitCommandHandler{
						ExchangeName: testExchange.ExchangeName,
						WSConn:       nil,
						QueueName:    LocalWebsocketName(queueID),
					}

					err := msg.Ack()
					assert.Nil(t, err)

					return rch.HandleCommand(msg)
				default:
					t.Fatalf("Unexpected message type: %d", msg.ContentType)
					return errors.New("Unexpected message type")
				}
			},
		},
		PubCfg: &AMQPPubCfg{
			Messages: make(chan AMQPMessage, 1),
		},
		Control: pubSubCtrl,
	}

	go func() {
		RunSubscriber(pubSubCfg)
	}()
	go func() {
		RunPublisher(pubSubCfg)
	}()
	pubSubCfg.Control.Ready.Wait() // wait for subscribers and publishers to have started up.

	// Send subscribe command
	pubSubCfg.PubCfg.Messages <- TestSubscription
	if err := utils.WaitTimeout(subWg, 5*time.Second); err != nil {
		t.Fatal("Timed out waiting for subscribe command")
	}

	// Change incoming message handler
	pubSubCfg.SubCfg.HandleMessageFunc = func(msg AMQPMessage) error {
		switch msg.ContentType {
		case ContentTypeMsg:
			pubSubCtrl.Shutdown()

			assert.NotNil(t, msg.Ack, "Message does not have Ack")
			assert.NotNil(t, msg.Nack, "Message does not have Nack")
			assert.NotNil(t, msg.ErrHandler, "Message does not have Nack")

			err := msg.Ack()
			assert.Nil(t, err)

			// the messages aren't EXACT copies of each other b/c Nack, Ack, and ErrHandler are nil in TestMessage
			// so we need to fix this
			msg.Ack = nil
			msg.Nack = nil
			msg.ErrHandler = nil

			if !reflect.DeepEqual(msg, TestMessage) {
				t.Fatal("Sent message does not equal received message")
			}
			doneTesting <- true
			return nil
		default:
			t.Fatalf("Unexpected message type: %d", msg.ContentType)
			return errors.New("Unexpected message type")
		}
	}

	pubSubCfg.PubCfg.Messages <- TestMessage

	select {
	case <-doneTesting:
	// success
	case <-time.After(time.Second * 5):
		t.Fatal("control signal timed out")
	}

}
