package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/handlers"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
)

var addr = flag.String("addr", "0.0.0.0:80", "http service address")

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

	AMQPControl := utils.NewControl()

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

	http.HandleFunc("/ws/", handlers.NewWSConn)

	fmt.Println("Binding to address: " + *addr)
	err = http.ListenAndServe(*addr, nil)
	utils.FailOnError(err, "Could not bind to port")

	AMQPControl.Exit <- true
}
