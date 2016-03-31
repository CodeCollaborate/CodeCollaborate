package main

import (
	"flag"
	"fmt"
	"github.com/CodeCollaborate/Server/modules/handlers"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", "0.0.0.0:80", "http service address")

func main() {

	flag.Parse()
	log.SetFlags(0)

	// START DB CONNECTION
	//managers.ConnectMGo()
	//defer managers.GetPrimaryMGoSession().Close()

	// Get working directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Fatal error: Could not get Working Directory: %s\n", err)
		log.Fatal(err)
	}
	fmt.Println("Running in directory: " + dir)

	rabbitmq.SetupRabbitExchange(
		rabbitmq.ConnectionConfig{
			Host: "localhost",
			Port: 5672,
			User: "guest",
			Pass: "guest",
			ExchangeNames: []string{
				"CodeCollaborate",
			},
		},
	)

	http.HandleFunc("/ws/", handlers.NewWSConn)

	fmt.Println("Binding to address: " + *addr)
	err = http.ListenAndServe(*addr, nil)
	utils.FailOnError(err, "Could not bind to port")
}
