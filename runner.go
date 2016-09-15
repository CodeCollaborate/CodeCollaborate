package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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
var addr = flag.String("addr", "0.0.0.0:8000", "http service address")

func main() {

	flag.Parse()
	log.SetFlags(0)

	// START DB CONNECTION
	//managers.ConnectMGo()
	//defer managers.GetPrimaryMGoSession().Close()

	err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	config := config.GetConfig()

	// Get working directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Fatal error: Could not get Working Directory: %s\n", err)
		log.Fatal(err)
	}
	fmt.Println("Running in directory: " + dir)

	// Creates a NewControl block for multithreading control
	AMQPControl := utils.NewControl()

	// RabbitMQ uses "Exchanges" as containers for Queues, and ours is initialized here.
	rabbitmq.SetupRabbitExchange(
		&rabbitmq.AMQPConnCfg{
			ConnCfg: config.ConnectionConfig["RabbitMQ"],
			Exchanges: []rabbitmq.AMQPExchCfg{
				{
					ExchangeName: config.ServerConfig.Name,
					Durable:      true,
				},
			},
			Control: AMQPControl,
		},
	)

	dbfs.Dbfs = new(dbfs.DatabaseImpl)

	http.HandleFunc("/ws/", handlers.NewWSConn)

	fmt.Println("Binding to address: " + *addr)
	err = http.ListenAndServe(*addr, nil)
	utils.FailOnError(err, "Could not bind to port")

	// Kill the SetupRabbitExchange thread (Multithreading control)
	AMQPControl.Exit <- true
}
