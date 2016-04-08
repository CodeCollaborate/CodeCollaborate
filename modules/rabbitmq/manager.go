package rabbitmq

import (
	"errors"
	"sync"

	"github.com/CodeCollaborate/Server/utils"
	"github.com/streadway/amqp"
)

/**
 * RabbitMq manager for CodeCollaborate Server.
 * @author: Austin Fahsl and Benedict Wong
 */

var channelQueueCreationMutex = sync.Mutex{}
var channelQueue chan *amqp.Channel

// GetChannel gets a new RabbitMQ Channel. This function requires that SetupRabbitExchange has been
// previously called. The function will throw an error if the SetupRabbitExchange has not been called.
func GetChannel() (*amqp.Channel, error) {
	if channelQueue == nil {
		return nil, errors.New("Rabbit Exchange not initialized")
	}
	return <-channelQueue, nil
}

// SetupRabbitExchange sets up the RabbitMq exchange, initializing connections, and starting to push RabbitMQ Channels
// into the ChannelQueue. The generation of channels will be done on a new GoRoutine, avoiding blocking, or having
// to pass the RabbitMQ Connection around. This method will also attempt to auto-reconnect if the critical setup steps
// fail.
func SetupRabbitExchange(cfg *AMQPConnCfg) {
	if cfg.Control == nil {
		cfg.Control = utils.NewControl()
	}

	if channelQueue == nil {
		channelQueueCreationMutex.Lock()
		if channelQueue == nil {
			channelQueue = make(chan *amqp.Channel)

			ready := make(chan bool)
			go func() {
				// Loop; if connection drops, we should try to restore connection before creating new channels.
			redialLoop:
				for {
					conn, err := amqp.Dial(cfg.ConnectionString())
					if err != nil {
						utils.LogOnError(err, "Failed to connect to RabbitMQ")
						continue redialLoop
					}

					ch, err := conn.Channel()
					if err != nil {
						utils.LogOnError(err, "Failed to open a channel")
						conn.Close()
						continue redialLoop
					}

					for _, exchange := range cfg.Exchanges {
						err = ch.ExchangeDeclare(
							exchange.ExchangeName, // name
							"direct",              // type
							exchange.Durable,      // durable
							!exchange.Durable,     // auto-deleted
							false,                 // internal
							false,                 // no-wait
							nil,                   // arguments
						)
						if err != nil {
							utils.LogOnError(err, "Failed to declare an exchange")
							ch.Close()
							conn.Close()
							continue redialLoop
						}
					}
					ch.Close()

					for {
						ch, err = conn.Channel()
						utils.LogOnError(err, "Failed to open a channel")
						if err != nil {
							break
						}

						select {
						case ready <- true:
						default:
						}

						select {
						case <-cfg.Control.Exit:
							break
						case channelQueue <- ch:
						}
					}
					conn.Close()
				}
			}()
			<-ready

			// Signal that this connection is ready
			cfg.Control.Ready.Done()
		}
		channelQueueCreationMutex.Unlock()
	}
}

// RunSubscriber creates a new subscriber based on the QueueConfig provided. The RabbitMQ Channel used
// is returned, along with a Go Channel of the pushed messages from the RabbitMQ Exchange. Developers should
// remember to defer the closing of the RabbitMQ Channel.
func RunSubscriber(cfg *AMQPSubCfg) error {
	if cfg.Control == nil {
		cfg.Control = utils.NewControl()
	}

	ch, err := GetChannel()
	if err != nil {
		utils.LogOnError(err, "Failed to get RabbitMQ Channel")
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		cfg.QueueName(),  // name (routing key)
		cfg.IsWorkQueue,  // durable - persist data upon restarts?
		!cfg.IsWorkQueue, // delete when unused - no more clients attached
		!cfg.IsWorkQueue, // exclusive - can only be used by this channel
		false,            // no-wait - do not wait for server to confirm that the queue has been created
		nil,              // arguments
	)
	if err != nil {
		return err
	}

	for _, key := range append(cfg.Keys, cfg.QueueName()) {
		err = ch.QueueBind(
			cfg.QueueName(),  // queue name
			key,              // routing key
			cfg.ExchangeName, // exchange
			false,            // no-wait
			nil,              // arguments
		)
		if err != nil {
			return err
		}
	}

	msgs, err := ch.Consume(
		cfg.QueueName(), // queue
		"",              // consumer
		true,            // auto ack
		true,            // exclusive
		false,           // no local
		false,           // no wait
		nil,             // args
	)
	if err != nil {
		return err
	}

	// Signal that this Subscriber is ready
	cfg.Control.Ready.Done()
	for {
		select {
		case <-cfg.Control.Exit:
			return nil
		case msg := <-msgs:
			message := AMQPMessage{
				Headers:     msg.Headers,
				RoutingKey:  msg.RoutingKey,
				ContentType: msg.ContentType,
				Message:     msg.Body,
				Persistent:  (msg.DeliveryMode == 2),
			}
			err := cfg.HandleMessageFunc(message)
			utils.LogOnError(err, "Failed to handle message")
		}
	}
}

// RunPublisher creates a new publisher, and continually pushes messages submitted to the Go channel
// to RabbitMQ.
func RunPublisher(cfg *AMQPPubCfg) error {
	if cfg.Control == nil {
		cfg.Control = utils.NewControl()
	}

	ch, err := GetChannel()
	if err != nil {
		utils.LogOnError(err, "Failed to get RabbitMQ Channel")
		return err
	}

	// Signal that this Subscriber is ready
	cfg.Control.Ready.Done()
	for {
		select {
		case <-cfg.Control.Exit:
			return nil
		case message := <-cfg.Messages:

			deliveryMode := uint8(0)
			if message.Persistent {
				deliveryMode = 2
			}

			err = ch.Publish(
				cfg.ExchangeName,   // exchange
				message.RoutingKey, // routing key
				false,              // mandatory - must be placed on at least one queue, otherwise return to sender
				false,              // immediate - must be delivered immediately. If no free workers, return to sender
				amqp.Publishing{
					Headers:      message.Headers,
					ContentType:  message.ContentType,
					DeliveryMode: deliveryMode, // 0, 1 for transient, 2 for persistent
					Body:         []byte(message.Message),
				})
			utils.LogOnError(err, "Failed to publish a message")
		}
	}
}
