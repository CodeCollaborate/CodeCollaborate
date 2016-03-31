package rabbitmq

import (
	"errors"
	"github.com/CodeCollaborate/Server/utils"
	"github.com/streadway/amqp"
)

/**
 * RabbitMq manager for CodeCollaborate Server.
 * @author: Austin Fahsl and Benedict Wong
 */

var channelQueue = make(chan *amqp.Channel)

//
// GetChannel gets a new RabbitMQ Channel. This function requires that SetupRabbitExchange has been
// previously called. The function will throw an error if the SetupRabbitExchange has not been called.
//
func GetChannel() (*amqp.Channel, error) {
	if len(channelQueue) <= 0 {
		return nil, errors.New("Rabbit Exchange not initialized")
	}
	return <-channelQueue, nil
}

//
// SetupRabbitExchange sets up the RabbitMq exchange, initializing connections, and starting to push RabbitMQ Channels
// into the ChannelQueue. The generation of channels will be done on a new GoRoutine, avoiding blocking, or having
// to pass the RabbitMQ Connection around. This method will also attempt to auto-reconnect if the critical setup steps
// fail.
//
func SetupRabbitExchange(cfg ConnectionConfig) {
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

			for _, exchangeName := range cfg.ExchangeNames {
				err = ch.ExchangeDeclare(
					exchangeName, // name
					"direct",     // type
					true,         // durable
					false,        // auto-deleted
					false,        // internal
					false,        // no-wait
					nil,          // arguments
				)
				if err != nil {
					utils.LogOnError(err, "Failed to declare an exchange")
					ch.Close()
					conn.Close()
					continue redialLoop
				}
			}
			ch.Close()

			ready <- true
			for {
				ch, err = conn.Channel()
				utils.LogOnError(err, "Failed to open a channel")
				if err != nil {
					break
				}
				channelQueue <- ch
			}
			conn.Close()
		}
	}()
	<-ready
}

//
// RunSubscriber creates a new subscriber based on the QueueConfig provided. The RabbitMQ Channel used
// is returned, along with a Go Channel of the pushed messages from the RabbitMQ Exchange. Developers should
// remember to defer the closing of the RabbitMQ Channel.
//
func RunSubscriber(cfg QueueConfig) (*amqp.Channel, <-chan amqp.Delivery, error) {
	ch, err := GetChannel()
	if err != nil {
		utils.LogOnError(err, "Failed to get RabbitMQ Channel")
		return nil, nil, err
	}

	_, err = ch.QueueDeclare(
		cfg.QueueName(),  // name (routing key)
		cfg.IsWorkQueue,  // durable - persist data upon restarts?
		!cfg.IsWorkQueue, // delete when unused - no more clients attached
		!cfg.IsWorkQueue, // exclusive - can only be used by this channel
		false,            // no-wait - do not wait for server to confirm that the queue has been created
		nil,              // arguments
	)
	if err != nil {
		ch.Close()
		return nil, nil, err
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
			ch.Close()
			return nil, nil, err
		}
	}

	messages, err := ch.Consume(
		cfg.QueueName(), // queue
		"",              // consumer
		true,            // auto ack
		true,            // exclusive
		false,           // no local
		false,           // no wait
		nil,             // args
	)
	if err != nil {
		ch.Close()
		return nil, nil, err
	}

	return ch, messages, nil
}
