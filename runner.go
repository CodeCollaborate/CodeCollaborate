package main

import (
	"flag"
	"log"
	"os"
	"net/http"
	"fmt"
	"github.com/CodeCollaborate/Server/modules/handlers"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
)

var addr = flag.String("addr", "0.0.0.0:80", "http service address")


func main(){

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

	conn := rabbitmq.SetupRabbitExchange()
	defer conn.Close()

	http.HandleFunc("/ws/", handlers.NewWSConn)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Printf("Fatal error: Could not get bind port: %s\n", err)
		log.Fatal(err)
	}
}