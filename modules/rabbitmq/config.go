package rabbitmq

import (
	"crypto/tls"
	"fmt"
	"os"
)

/**
 * Configuration constants for the CodeCollaborate server.
 * @author: Austin Fahsl and Benedict Wong
 */

// Gets the hostname of this machine, for use in QueueName()
var hostname, _ = os.Hostname()

// ConnectionConfig represents the settings needed to create a new connection, and initialize the required exchanges.
type ConnectionConfig struct {
	Host          string
	Port          int
	User          string
	Pass          string
	ExchangeNames []string
	TLSConfig     *tls.Config
}

// ConnectionString returns the connection string, using amqps:// if TLSConfig has been set, amqp:// otherwise.
func (cfg ConnectionConfig) ConnectionString() string {
	if cfg.TLSConfig != nil {
		return fmt.Sprintf("amqps://%s:%s@%s:%d/", cfg.User, cfg.Pass, cfg.Host, cfg.Port)
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Pass, cfg.Host, cfg.Port)
}

// QueueConfig represents the settings needed to create a new queue, and bind it to rhe relevant keys.
type QueueConfig struct {
	ExchangeName string
	QueueId      uint64
	Keys         []string
	IsWorkQueue  bool
}

// QueueName generates the Queue
func (cfg QueueConfig) QueueName() string {
	return fmt.Sprintf("%s-%d", hostname, cfg.QueueId)
}
