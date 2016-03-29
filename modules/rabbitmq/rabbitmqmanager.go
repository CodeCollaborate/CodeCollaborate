package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/CodeCollaborate/Server/modules/datahandling"
	"github.com/CodeCollaborate/Server/utils"
	"fmt"
	"os"
)

/**
 * RabbitMq manager for CodeCollaborate Server.
 * @author: Austin Fahsl and Benedict Wong
 */

// Make a queue of go channels
var ChannelQueue = make(chan *amqp.Channel)
// get the hostname of this machine
var hostname, _ = os.Hostname()

/**
 * Sets up the RabbitMq exchange for managing individual WebSocket queues.
 */
func SetupRabbitExchange() *amqp.Connection{
	conn, err := amqp.Dial(datahandling.ConnectionString)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		datahandling.ExchangeName, // name
		"direct", // type
		true, // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil, // arguments
	)
	utils.FailOnError(err, "Failed to declare an exchange")

	go func() {
		for {
			ch, err = conn.Channel()
			utils.FailOnError(err, "Failed to open a channel")
			ChannelQueue <- ch
		}
	}()

	return conn

}

/**
 * Creates the subscriber that listens to RabbitMq and returns the RabbitMq channel and the message channel.
 */
func RunSubscriber(wsId uint64) (*amqp.Channel, <-chan amqp.Delivery, error) {
	ch := <-ChannelQueue

	queueName := fmt.Sprintf("%s-%d", hostname, wsId)

	_, err := ch.QueueDeclare(
		queueName, // name (routing key)
		false, // durable - persist data upon restarts?
		true, // delete when unused - no more clients attached
		true, // exclusive - can only be used by this channel
		false, // no-wait - do not wait for server to confirm that the queue has been created
		nil, // arguments
	)
	if err != nil {
		ch.Close()
		return nil, nil, err
	}

	err = ch.QueueBind(
		queueName, // queue name
		queueName, // routing key
		datahandling.ExchangeName, // exchange
		false, // no-wait
		nil, // arguments
	)
	if err != nil {
		ch.Close()
		return nil, nil, err
	}

	messages, err := ch.Consume(
		queueName, // queue
		"", // consumer
		true, // auto ack
		true, // exclusive
		false, // no local
		false, // no wait
		nil, // args
	)
	if err != nil {
		ch.Close()
		return nil, nil, err
	}

	return ch, messages, nil
}
