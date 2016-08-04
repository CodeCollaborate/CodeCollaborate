package datahandling

import (
	"time"

	"github.com/CodeCollaborate/Server/modules/rabbitmq"
)

type dhClosure interface {
	call(dh DataHandler) error
}

type toSenderClos struct {
	msg *serverMessageWrapper
}

func (cont toSenderClos) call(dh DataHandler) error {
	return dh.sendToSender(cont.msg)
}

type toChannelClos struct {
	msg *serverMessageWrapper
}

func (cont toChannelClos) call(dh DataHandler) error {
	return dh.sendToChannel(cont.msg)
}

type chanSubscribeClos struct {
	key string
	tag int64
}

func (cont chanSubscribeClos) call(dh DataHandler) error {
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
	err := toSenderClos{msg: res}.call(dh) // go ahead and send from here
	return err

}

type chanUnsubscribeClos struct {
	key string
	tag int64
}

func (cont chanUnsubscribeClos) call(dh DataHandler) error {
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
	err := toSenderClos{msg: res}.call(dh) // go ahead and send from here
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
