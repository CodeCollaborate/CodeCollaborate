package datahandling

/**
 * Configuration constants for the CodeCollaborate server.
 * @author: Austin Fahsl and Benedict Wong
 */

type RabbitMQConfig struct {
	ExchangeName string
	QueueName    string
	RoutingKeys  []string
	Messages     chan string
}

// location, username, and password of RabbitMq
const ConnectionString = "amqp://guest:guest@localhost:5672/"

// name of the exchange for RabbitMq
const ExchangeName = "CodeCollaborate"
