package event

import "github.com/sirupsen/logrus"

type Handler func(dispatcher string, data interface{})

type Event string

type Subscribable interface {
	Subscribe(event Event, handler Handler)
}

type Dispatcher interface {
	Dispatch(event Event, from string, data interface{})
}

type MessageQueue struct {
	topics map[Event][]Handler
}

var mq *MessageQueue

func CreateOrGet() *MessageQueue {
	if mq != nil {
		return mq
	}
	topics := make(map[Event][]Handler)
	mq = &MessageQueue{topics}
	return mq
}

func (mq MessageQueue) Dispatch(event Event, from string, data interface{}) {
	logrus.WithFields(map[string]interface{}{"event": event, "from": from}).Debug("New Event Dispatched")
	subscriber, ok := mq.topics[event]
	if !ok {
		return
	}
	for _, handler := range subscriber {
		handler(from, data)
	}
}

func (mq MessageQueue) Subscribe(event Event, handler Handler) {
	mq.topics[event] = append(mq.topics[event], handler)
}
