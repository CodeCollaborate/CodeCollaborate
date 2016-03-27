package datahandling

import "fmt"

type DataHandler struct{

}

func (dh DataHandler) Handle(wsId uint64, json []byte) error {
	fmt.Printf("Handling JSON: %s\n", json)
	return nil
}