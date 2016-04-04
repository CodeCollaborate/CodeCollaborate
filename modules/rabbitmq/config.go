package rabbitmq

import (
	"crypto/tls"
	"fmt"
	"github.com/CodeCollaborate/Server/utils"
	"os"
)

// Gets the hostname of this machine, for use in QueueName()
var hostname, _ = os.Hostname()

// ConnectionConfig represents the settings needed to create a new connection, and initialize the required exchanges.
type ConnectionConfig struct {
	Host      string
	Port      int
	User      string
	Pass      string
	Exchanges []ExchangeConfig
	TLSConfig *tls.Config
	Control   *utils.Control
}

// ExchangeConfig represents the basic variables of any exchange
type ExchangeConfig struct {
	ExchangeName string
	Durable      bool
}

// ConnectionString returns the connection string, using amqps:// if TLSConfig has been set, amqp:// otherwise.
func (cfg ConnectionConfig) ConnectionString() string {
	if cfg.TLSConfig != nil {
		return fmt.Sprintf("amqps://%s:%s@%s:%d/", cfg.User, cfg.Pass, cfg.Host, cfg.Port)
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Pass, cfg.Host, cfg.Port)
}

// SubscriberConfig represents the settings needed to create a new subscriber, including the queues and key bindings
type SubscriberConfig struct {
	ExchangeName      string
	QueueID           uint64
	Keys              []string
	IsWorkQueue       bool
	HandleMessageFunc func(AMQPMessage) error
	Control           *utils.Control
}

// QueueName generates the Queue
func (cfg SubscriberConfig) QueueName() string {
	return fmt.Sprintf("%s-%d", hostname, cfg.QueueID)
}

// PublisherConfig represents the settings needed to create a new publisher
type PublisherConfig struct {
	ExchangeName string
	Messages     chan AMQPMessage
	Control      *utils.Control
}

// AMQPMessage represents the information required to send a new message
type AMQPMessage struct {
	Headers     map[string]interface{}
	RoutingKey  string
	ContentType string
	Persistent  bool
	Message     []byte
}
