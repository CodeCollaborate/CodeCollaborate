package datahandling

import "github.com/CodeCollaborate/Server/modules/rabbitmq"

func toSenderCont(msg *serverMessageWrapper) func(dh DataHandler) error {
	return func(dh DataHandler) error {
		return dh.sendToSender(msg)
	}
}

func toChanCont(msg *serverMessageWrapper) func(dh DataHandler) error {
	return func(dh DataHandler) error {
		return dh.sendToChannel(msg)
	}
}

func chanSubscribe(key string) func(dh DataHandler) error {
	return func(dh DataHandler) error {
		dh.SubscriptionChan <- rabbitmq.Subscription{
			Channel:     key,
			IsSubscribe: true,
		}
		// I (joel) don't believe we actually have a way to know here if that subscribe throws an error
		return nil
	}
}

func chanUnsubscribe(key string) func(dh DataHandler) error {
	return func(dh DataHandler) error {
		dh.SubscriptionChan <- rabbitmq.Subscription{
			Channel:     key,
			IsSubscribe: false,
		}
		// I (joel) don't believe we actually have a way to know here if that subscribe throws an error
		return nil
	}
}

/*
 *
 *util
 *
 */

// simple helper to clean up some of the syntax when creating multiple closures
func accumulate(calls ...(func(dh DataHandler) error)) [](func(dh DataHandler) error) {
	return calls
}
