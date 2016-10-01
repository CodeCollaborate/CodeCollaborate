package main

import (
	"flag"
	"net/http"
	"os"

	"fmt"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/CodeCollaborate/Server/modules/handlers"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
)

/**
 * Runner.go starts the server. It initializes processes and begins listening for websocket requests.
 */

// changed from "0.0.0.0:80" because you need to be root to bind to that port
var logDir = flag.String("log_dir", "./data/logs/", "log file location")

func main() {

	flag.Parse()

	// START DB CONNECTION
	//managers.ConnectMGo()
	//defer managers.GetPrimaryMGoSession().Close()

	config.EnableLoggingToFile(*logDir)
	err := config.LoadConfig()
	if err != nil {
		utils.LogFatal("Failed to load configuration", err, nil)
	}
	cfg := config.GetConfig()

	// Get working directory
	dir, err := os.Getwd()
	utils.LogFatal("Could not get working directory", err, nil)

	utils.LogInfo("Working directory initalized", utils.LogFields{
		"Working Directory": dir,
	})

	// Creates a NewControl block for multithreading control
	AMQPControl := utils.NewControl()

	// RabbitMQ uses "Exchanges" as containers for Queues, and ours is initialized here.
	rabbitmq.SetupRabbitExchange(
		&rabbitmq.AMQPConnCfg{
			ConnCfg: cfg.ConnectionConfig["RabbitMQ"],
			Exchanges: []rabbitmq.AMQPExchCfg{
				{
					ExchangeName: cfg.ServerConfig.Name,
					Durable:      true,
				},
			},
			Control: AMQPControl,
		},
	)

	dbfs.Dbfs = new(dbfs.DatabaseImpl)

	http.HandleFunc("/ws/", handlers.NewWSConn)

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.ServerConfig.Port)
	utils.LogInfo("Starting server", utils.LogFields{
		"Address": addr,
	})
	err = http.ListenAndServe(addr, nil)
	utils.LogError("Could not bind to port", err, nil)

	// Kill the SetupRabbitExchange thread (Multithreading control)
	defer func() {
		AMQPControl.Exit <- true
	}()
}
