package datahandling

import (
	"time"

	"github.com/CodeCollaborate/Server/modules/rabbitmq"
)

type dhClosure interface {
	call(dh DataHandler) error
}

type toSenderClosure struct {
	msg *serverMessageWrapper
}

func (cont toSenderClosure) call(dh DataHandler) error {
	return dh.sendToSender(cont.msg)
}

type toChannelClosure struct {
	msg *serverMessageWrapper
}

func (cont toChannelClosure) call(dh DataHandler) error {
	return dh.sendToChannel(cont.msg)
}

type chanSubscribeClosure struct {
	key string
	tag int64
}

func (cont chanSubscribeClosure) call(dh DataHandler) error {
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	// I (joel) don't believe we actually have a way to know here if this subscribe throws an error
	dh.SubscriptionChan <- rabbitmq.Subscription{
		Channel:     cont.key,
		IsSubscribe: true,
	}

	//if err != nil {
	//	res.ServerMessage = response{
	//		Status: fail,
	//		Tag:    p.Tag,
	//		Data:   struct{}{}}
	//} else {
	res.ServerMessage = response{
		Status: success,
		Tag:    cont.tag,
		Data:   struct{}{}}
	//}
	err := toSenderClosure{msg: res}.call(dh) // go ahead and send from here
	return err

}

type chanUnsubscribeClosure struct {
	key string
	tag int64
}

func (cont chanUnsubscribeClosure) call(dh DataHandler) error {
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	// I (joel) don't believe we actually have a way to know here if this subscribe throws an error
	dh.SubscriptionChan <- rabbitmq.Subscription{
		Channel:     cont.key,
		IsSubscribe: false,
	}

	//if err != nil {
	//	res.ServerMessage = response{
	//		Status: fail,
	//		Tag:    p.Tag,
	//		Data:   struct{}{}}
	//} else {
	res.ServerMessage = response{
		Status: success,
		Tag:    cont.tag,
		Data:   struct{}{}}
	//}
	err := toSenderClosure{msg: res}.call(dh) // go ahead and send from here
	return err
}

/*
 *
 *util
 *
 */

// simple helper to clean up some of the syntax when creating multiple closures
func accumulate(calls ...(dhClosure)) [](dhClosure) {
	return calls
}
