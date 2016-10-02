package rabbitmq

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/CodeCollaborate/Server/utils"
	"github.com/streadway/amqp"
)

/**
 * RabbitMq manager for CodeCollaborate Server.
 */
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
		cfg.Control = utils.NewControl()
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

// RunSubscriber creates a new subscriber based on the QueueConfig provided. The RabbitMQ Channel used
// is returned, along with a Go Channel of the pushed messages from the RabbitMQ Exchange. Developers should
// remember to defer the closing of the RabbitMQ Channel.
func RunSubscriber(cfg *AMQPSubCfg) error {
	if cfg.Control == nil {
		cfg.Control = NewControl()
	}
	defer func() {
		close(cfg.Control.Exit)
		close(cfg.Control.SubChan)
	}()

	ch, err := GetChannel()
	if err != nil {
		utils.LogError("Failed to get new channel", err, nil)
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
		false,           // exclusive
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
		case subscription := <-cfg.Control.SubChan:
			if subscription.IsSubscribe {
				err = ch.QueueBind(
					cfg.QueueName(),       // queue name
					subscription.GetKey(), // routing key
					cfg.ExchangeName,      // exchange
					false,                 // no-wait
					nil,                   // arguments
				)
				if err != nil {
					utils.LogError("Error binding to key", err, utils.LogFields{
						"Queue":      cfg.QueueName(),
						"RoutingKey": subscription.GetKey(),
					})
					cfg.Control.Exit <- true
				}
			} else {
				err = ch.QueueUnbind(
					cfg.QueueName(),       // queue name
					subscription.GetKey(), // routing key
					cfg.ExchangeName,      // exchange
					nil,                   // arguments
				)
				if err != nil {
					utils.LogError("Error unbinding from key", err, utils.LogFields{
						"Queue":      cfg.QueueName(),
						"RoutingKey": subscription.GetKey(),
					})
					cfg.Control.Exit <- true
				}
			}
		case msg := <-msgs:
			message := AMQPMessage{
				Headers:     msg.Headers,
				RoutingKey:  msg.RoutingKey,
				ContentType: msg.ContentType,
				Message:     msg.Body,
				Persistent:  (msg.DeliveryMode == 2),
			}
			err := cfg.HandleMessageFunc(message)

			utils.LogError("Message handler failed", err, nil)
		}
	}
}

// RunPublisher creates a new publisher, and continually pushes messages submitted to the Go channel
// to RabbitMQ.
func RunPublisher(cfg *AMQPPubCfg) error {
	if cfg.Control == nil {
		cfg.Control = utils.NewControl()
	}
	defer func() {
		close(cfg.Messages)
		close(cfg.Control.Exit)
	}()

	ch, err := GetChannel()
	if err != nil {
		utils.LogError("Failed to get new channel", err, nil)
		// panic so we shut down the subscriber too
		panic(err) // TODO(shapiro): Think of a better way of having publisher and consumer be able to shut each other down
	}
	defer ch.Close()

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
					Body:         message.Message,
				})

			if err != nil {
				utils.LogError("Failed to publish message", err, utils.LogFields{
					"RoutingKey": message.RoutingKey,
				})
				// TODO (shapiro): decide on action at publish error: retry with count?
			}
		}
	}
}
