package rabbitmq

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/CodeCollaborate/Server/utils"
	"github.com/kr/pretty"
	"github.com/streadway/amqp"
)

/**
 * RabbitMq manager for CodeCollaborate Server.
 //*/
const (
	defaultHeartbeat         = 10 * time.Second
	defaultConnectionTimeout = 30
)

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
func SetupRabbitExchange(cfg *AMQPConnCfg) error {
	if cfg.Control == nil {
		cfg.Control = utils.NewControl(1)
	}

	success := true

	if channelQueue == nil {
		channelQueueCreationMutex.Lock()
		if channelQueue == nil {

			ready := make(chan bool)
			go func() {
				// Loop; if connection drops, we should try to restore connection before creating new channels.
				retries := uint16(0)

			redialLoop:
				for {
					conn, err := amqp.DialConfig(cfg.ConnectionString(), amqp.Config{
						Heartbeat: defaultHeartbeat,
						Dial:      getNewDialer(cfg.Timeout),
					})
					if err != nil {
						utils.LogError("Failed to connect to RabbitMQ", err, utils.LogFields{
							"Host": cfg.Host,
							"Port": cfg.Port,
						})
						if retries >= cfg.NumRetries {
							ready <- false
							if channelQueue == nil {
								for {
									select {
									case channelQueue <- nil:
									default:
										channelQueue = nil
										return
									}
								}
							}
							return
						}
						retries++
						continue redialLoop
					}
					retries = 0

					if channelQueue == nil {
						channelQueue = make(chan *amqp.Channel)
					}

					ch, err := conn.Channel()
					if err != nil {
						utils.LogError("Failed to open a channel", err, nil)
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
							utils.LogError("Failed to declare exchange", err, nil)
							ch.Close()
							conn.Close()
							continue redialLoop
						}
					}
					ch.Close()

					for {
						ch, err = conn.Channel()
						if err != nil {
							utils.LogError("Failed to open a channel", err, nil)
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
			success = <-ready

			// Signal that this connection is ready
			cfg.Control.Ready.Done()
		}
		channelQueueCreationMutex.Unlock()
	}

	if !success {
		return errors.New("Failed to connect - timed out.")
	}

	return nil
}

func getNewDialer(timeout uint16) func(network, addr string) (net.Conn, error) {
	if timeout == 0 {
		timeout = defaultConnectionTimeout
	}

	// returns dialer using timeout if non-zero, or dialer using default timeout otherwise.
	return func(network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, time.Duration(timeout)*time.Second)
		if err != nil {
			return nil, err
		}

		// Heartbeating hasn't started yet, don't stall forever on a dead server.
		if err := conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
			return nil, err
		}

		return conn, nil
	}
}

// BindQueue binds this queue to a key.
func BindQueue(ch *amqp.Channel, queueName, key, exchangeName string) error {
	return ch.QueueBind(
		queueName,    // queue name
		key,          // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
}

// UnbindQueue unbinds this queue from a key.
func UnbindQueue(ch *amqp.Channel, queueName, key, exchangeName string) error {
	return ch.QueueUnbind(
		queueName,    // queue name
		key,          // routing key
		exchangeName, // exchange
		nil,          // arguments
	)
}

// RunSubscriber creates a new subscriber based on the QueueConfig provided. The RabbitMQ Channel used
// is returned, along with a Go Channel of the pushed messages from the RabbitMQ Exchange. Developers should
// remember to defer the closing of the RabbitMQ Channel.
func RunSubscriber(cfg *AMQPPubSubCfg) error {
	defer func() {
		cfg.shutdown.Do(func() {
			close(cfg.Control.Exit) // If subscriber exits, kill publisher as well.
		})
	}()

	ch, err := GetChannel()
	if err != nil {
		utils.LogError("Failed to get new channel", err, nil)
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		cfg.SubCfg.QueueName(),  // name (routing key)
		cfg.SubCfg.IsWorkQueue,  // durable - persist data upon restarts?
		!cfg.SubCfg.IsWorkQueue, // delete when unused - no more clients attached
		!cfg.SubCfg.IsWorkQueue, // exclusive - can only be used by this channel
		false, // no-wait - do not wait for server to confirm that the queue has been created
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	for _, key := range append(cfg.SubCfg.Keys, cfg.SubCfg.QueueName()) {
		err = BindQueue(ch,
			cfg.SubCfg.QueueName(), // queue name
			key,              // routing key
			cfg.ExchangeName, // exchange
		)
		if err != nil {
			return err
		}
	}

	msgs, err := ch.Consume(
		cfg.SubCfg.QueueName(), // queue
		"",    // consumer
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
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
			contentType, err := strconv.Atoi(msg.ContentType)
			if err != nil {
				utils.LogError("ContentType not an int", err, utils.LogFields{
					"AMQPMessage": pretty.Sprint(msg),
				})
			}

			message := AMQPMessage{
				Headers:     msg.Headers,
				RoutingKey:  msg.RoutingKey,
				ContentType: contentType,
				Message:     msg.Body,
				Persistent:  (msg.DeliveryMode == 2),
			}
			err = cfg.SubCfg.HandleMessageFunc(message)

			utils.LogError("Message handler failed", err, nil)
		}
	}
}

// RunPublisher creates a new publisher, and continually pushes messages submitted to the Go channel
// to RabbitMQ.
func RunPublisher(cfg *AMQPPubSubCfg) error {
	defer func() {
		close(cfg.PubCfg.Messages)
		cfg.shutdown.Do(func() { // Make sure this is only ever called once.
			close(cfg.Control.Exit) // If subscriber exits, kill publisher as well.
		})
	}()

	ch, err := GetChannel()
	if err != nil {
		// Shut down subscriber if failed here.
		return fmt.Errorf("RunPublisher: Failed to get new channel: %v", err)
	}
	defer ch.Close()

	// Signal that this Publisher is ready
	cfg.Control.Ready.Done()
	for {
		select {
		case <-cfg.Control.Exit:
			return nil
		case message := <-cfg.PubCfg.Messages:

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
					ContentType:  strconv.Itoa(message.ContentType),
					DeliveryMode: deliveryMode, // 0, 1 for transient, 2 for persistent
					Body:         message.Message,
				})

			if err != nil {
				utils.LogError("Failed to publish AMQPMessage", err, utils.LogFields{
					"RoutingKey": message.RoutingKey,
					"Body":       string(message.Message),
				})
				// TODO (shapiro): decide on action at publish error: retry with count?
			}
		}
	}
}
