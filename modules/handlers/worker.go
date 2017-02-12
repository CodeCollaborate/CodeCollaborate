package handlers

import (
	"errors"
	"runtime"
	"strconv"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datahandling"
	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
)

// server only has one worker
var serverWorker worker

const workerName string = "datahandling_worker"
const workerOutboundQueueBufferSize int = 512

// Worker objects listen for work jobs coming from RabbitMQ and then run datahandling processes on them
type worker struct {
	cfg *rabbitmq.AMQPPubSubCfg
}

// StartWorker initializes the worker which talks with RabbitMQ for this server
func StartWorker(dbfsImpl dbfs.DBFS) {
	cfg := config.GetConfig()
	prefetchCount := runtime.NumCPU() // note that this is set to `runtime.GOMAXPROCS`

	pubCfg := rabbitmq.NewPubConfig(func(msg rabbitmq.AMQPMessage) {
		// do nothing (for now?)
		msg.ErrHandler()
	}, workerOutboundQueueBufferSize)

	subCfg := &rabbitmq.AMQPSubCfg{
		QueueName:     workerName,
		Keys:          []string{},
		IsWorkQueue:   true,
		PrefetchCount: prefetchCount,
	}

	pubSubCfg := rabbitmq.NewAMQPPubSubCfg(cfg.ServerConfig.Name, pubCfg, subCfg)

	subCfg.HandleMessageFunc = workerMessageHandler(dbfsImpl, pubSubCfg)

	go func() {
		err := rabbitmq.RunPublisher(pubSubCfg)
		if err != nil {
			utils.LogError("Worker publisher error encountered. Exiting", err, nil)
			pubSubCfg.Control.Shutdown()
		}
	}()
	go func() {
		err := rabbitmq.RunSubscriber(pubSubCfg)
		if err != nil {
			utils.LogError("Worker subscriber error encountered. Exiting", err, nil)
			pubSubCfg.Control.Shutdown()
		}
	}()

	pubSubCfg.Control.Ready.Wait()

	serverWorker = worker{
		cfg: pubSubCfg,
	}
}

// WorkerEnqueue takes messages that need to be processed and sends them to RabbitMQ to be assigned to a worker
func WorkerEnqueue(message []byte, wsID uint64) error {
	// can't naturally get wsID's of 0
	if wsID == 0 {
		return errors.New("invalid websocketID given to worker")
	}

	msg := rabbitmq.AMQPMessage{
		Headers: map[string]interface{}{
			"OriginID": strconv.FormatUint(wsID, 16),
		},
		RoutingKey:  workerName,
		ContentType: rabbitmq.ContentTypeWork,
		Persistent:  false,
		Message:     message,
	}

	select {
	case serverWorker.cfg.PubCfg.Messages <- msg:
	default:
		err := errors.New("Channel buffer full")
		utils.LogError("Worker message queue full, failed to add new message", err, utils.LogFields{
			"AMQP Message": msg,
		})
		return err
	}

	return nil
}

func workerMessageHandler(dbfsImpl dbfs.DBFS, cfg *rabbitmq.AMQPPubSubCfg) func(rabbitmq.AMQPMessage) error {
	dh := datahandling.DataHandler{
		MessageChan: cfg.PubCfg.Messages,
		Db:          dbfsImpl,
	}

	return func(msg rabbitmq.AMQPMessage) error {
		switch msg.ContentType {
		case rabbitmq.ContentTypeWork:
			// If notification with self as origin, early-out; ignore our own notifications.
			wsIDRaw, ok := msg.Headers["OriginID"]
			if !ok {
				err := errors.New("Unnown message origin")
				utils.LogError("Worker encountered ", err, utils.LogFields{
					"Message Headers": msg.Headers, // NOTE: message body could contain passwords
				})
				return err
			}
			wsID, err := strconv.ParseUint(wsIDRaw.(string), 16, 64)
			if err != nil {
				utils.LogError("Error converting worker OriginID", err, utils.LogFields{
					"raw OriginID":    wsIDRaw,
					"Message Headers": msg.Headers,
				})
			}

			go dh.Handle(msg.Message, wsID, msg.Ack)
			return nil
		default:
			err := errors.New("Unnable to process RabbitMQ message type")
			utils.LogError("not-work given to worker", err, utils.LogFields{
				"Message Headers": msg.Headers,
				"Message Body":    string(msg.Message),
			})
			return err
		}
	}
}
