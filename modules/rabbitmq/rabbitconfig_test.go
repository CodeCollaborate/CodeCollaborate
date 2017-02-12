package rabbitmq

import (
	"crypto/tls"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/CodeCollaborate/Server/modules/config"
)

func TestConnectionString(t *testing.T) {
	connCfg := AMQPConnCfg{
		ConnCfg: config.ConnCfg{
			Host:     "host",
			Port:     80,
			Username: "username",
			Password: "password",
		},
	}

	if !strings.HasPrefix(connCfg.ConnectionString(), "amqp://") {
		t.Fatal("Connection protocol incorrect")
	} else if !strings.HasSuffix(connCfg.ConnectionString(), "username:password@host:80/") {
		t.Fatal("Connection string incorrectly generated")
	}

	connCfgWithTLS := AMQPConnCfg{
		ConnCfg: config.ConnCfg{
			Host:     "host",
			Port:     80,
			Username: "username",
			Password: "password",
		},
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}

	if !strings.HasPrefix(connCfgWithTLS.ConnectionString(), "amqps://") {
		t.Fatal("Connection protocol incorrect")
	} else if !strings.HasSuffix(connCfgWithTLS.ConnectionString(), "username:password@host:80/") {
		t.Fatal("Connection string incorrectly generated")
	}
}

func TestQueueName(t *testing.T) {

	hostname, err := os.Hostname()
	if err != nil {
		t.Fatal("Could not get hostname")
	}

	for i := uint64(0); i < 20; i++ {
		queueID := i

		queueCfg := AMQPSubCfg{
			QueueName:   RabbitWebsocketQueueName(queueID),
			Keys:        []string{"Key1", "Key2"},
			IsWorkQueue: false,
		}

		expected := "WS-" + hostname + "-" + strconv.FormatUint(queueID, 10)
		if queueCfg.QueueName != expected {
			t.Fatalf("QueueName incorrect; expected [%s], got [%s]", expected, queueCfg.QueueName)
		}
	}
}
