package datahandling

type RabbitMQConfig struct {
	ExchangeName string
	QueueName    string
	RoutingKeys  []string
	Messages     chan string
}

const ConnectionString = "amqp://guest:guest@localhost:5672/"
const ExchangeName = "CodeCollaborate"
