package rabbitmq

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/utils"
)

/**
 * Configuration structures and variables for RabbitMQ.
 */

// Gets the hostname of this machine, for use in QueueName()
var hostname, _ = os.Hostname()

// AMQPConnCfg represents the settings needed to create a new connection, and initialize the required exchanges.
type AMQPConnCfg struct {
	config.ConnCfg
	Exchanges []AMQPExchCfg
	TLSConfig *tls.Config
	Control   *utils.Control
}

// ConnectionString returns the connection string, using amqps:// if TLSConfig has been set, amqp:// otherwise.
func (cfg AMQPConnCfg) ConnectionString() string {
	if cfg.TLSConfig != nil {
		return fmt.Sprintf("amqps://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
}

// AMQPExchCfg represents the basic variables of any exchange
type AMQPExchCfg struct {
	ExchangeName string
	Durable      bool
}

// AMQPPubSubCfg aggregates the publisher and subscriber into a single configuration, allowing them to shut each other
// down in the event of a unhandled error.
type AMQPPubSubCfg struct {
	ExchangeName string
	PubCfg       *AMQPPubCfg
	SubCfg       *AMQPSubCfg
	Control      *utils.Control // Used for shutting down both publisher and subscriber
}

// NewAMQPPubSubCfg creates a new AMQPPubSubCfg struct, and returns the pointer.
func NewAMQPPubSubCfg(exchangeName string, pubCfg *AMQPPubCfg, subCfg *AMQPSubCfg) *AMQPPubSubCfg {
	return &AMQPPubSubCfg{
		ExchangeName: exchangeName,
		PubCfg:       pubCfg,
		SubCfg:       subCfg,
		Control:      utils.NewControl(2),
	}
}

// AMQPSubCfg represents the settings needed to create a new subscriber, including the queues and key bindings
type AMQPSubCfg struct {
	QueueName         string
	Keys              []string
	IsWorkQueue       bool
	PrefetchCount     int
	HandleMessageFunc func(AMQPMessage) error
}

// RabbitUserQueueName returns the name of the Queue a websocket for the given user would have
func RabbitUserQueueName(username string) string {
	return fmt.Sprintf("User-%s", username)
}

// RabbitWebsocketQueueName returns the name of the Queue a websocket with the given ID would have
func RabbitWebsocketQueueName(queueID uint64) string {
	return fmt.Sprintf("WS-%s-%d", hostname, queueID)
}

// RabbitProjectQueueName returns the name of the Queue a project with the given ID would have
func RabbitProjectQueueName(projectID int64) string {
	return fmt.Sprintf("Project-%d", projectID)
}

// AMQPPubCfg represents the settings needed to create a new publisher
type AMQPPubCfg struct {
	PubErrHandler func(AMQPMessage) // Handler for publish errors
	Messages      chan AMQPMessage
}

// NewPubConfig creates a new AMQPPubCfg, initialized
func NewPubConfig(errHandler func(AMQPMessage), messageBuffer int) *AMQPPubCfg {
	return &AMQPPubCfg{
		PubErrHandler: errHandler,
		Messages:      make(chan AMQPMessage, messageBuffer), // Buffer messages to make sure a latency spike doesn't kill us.
	}
}

// AMQPMessage represents the information required to send a new message
type AMQPMessage struct {
	Headers     map[string]interface{}
	RoutingKey  string
	ContentType int
	Persistent  bool
	Message     []byte
	ErrHandler  func()
	Ack         func() error
	Nack        func() error
}

const (
	// ContentTypeMsg is the message content-type for an AMQPMessage
	ContentTypeMsg = iota

	// ContentTypeCmd is the command content-type for an AMQPMessage
	ContentTypeCmd

	// ContentTypeWork is the work content-type for an AMQPMessage
	ContentTypeWork
)
