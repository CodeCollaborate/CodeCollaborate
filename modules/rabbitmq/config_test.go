package rabbitmq

import (
	"crypto/tls"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestConnectionString(t *testing.T) {
	connCfg := ConnectionConfig{
		Host: "host",
		Port: 80,
		User: "username",
		Pass: "password",
	}

	if !strings.HasPrefix(connCfg.ConnectionString(), "amqp://") {
		t.Fatal("Connection protocol incorrect")
	} else if !strings.HasSuffix(connCfg.ConnectionString(), "username:password@host:80/") {
		t.Fatal("Connection string incorrectly generated")
	}

	connCfgWithTLS := ConnectionConfig{
		Host:      "host",
		Port:      80,
		User:      "username",
		Pass:      "password",
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

		queueCfg := SubscriberConfig{
			ExchangeName: "Exchange",
			QueueID:      queueID,
			Keys:         []string{"Key1", "Key2"},
			IsWorkQueue:  false,
		}

		expected := hostname + "-" + strconv.FormatUint(queueID, 10)
		if queueCfg.QueueName() != expected {
			t.Fatalf("QueueName incorrect; expected [%s], got [%s]", expected, queueCfg.QueueName())
		}
	}
}
