package rabbitmq

import (
	"testing"
)

func TestGetChannelAndSetupRabbitExchange(t *testing.T) {

	//channy := &amqp.Channel{}
	//cq := make(chan bool)
	//go func (){
	//	channelQueue <- channy
	//	cq <- true
	//}()

	//i := 0
	//for {
	//	i++
	//	if i > 10000 {
	//		break
	//	}
	//}


	//
	//retval, err := GetChannel()
	//if err != nil {
	//	t.Fatal("Failed to get channel")
	//}
	//if retval != channy {
	//	t.Fatal("Channel returned is not the one passed in")
	//}

	//SetupRabbitExchange(
	//	ConnectionConfig{
	//		Host: "localhost",
	//		Port: 5672,
	//		User: "guest",
	//		Pass: "guest",
	//		ExchangeNames: []string{
	//			"CodeCollaborate",
	//		},
	//	},
	//)
	//_, err := GetChannel()
	//if err != nil {
	//	t.Fatal("Failed to get channel")
	//}

}