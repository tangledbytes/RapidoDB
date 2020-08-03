package eventbus

import "sync"

// EventBus has subscribers which can subscribe to several topics
type EventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}

// DataChannelSlice is the slice of DataChannel
type DataChannelSlice []DataChannel

// DataChannel is the go channel for transferring data
type DataChannel chan DataEvent

// DataEvent encapsualtes the event published by the publishers
// throught the event bus
//
// The data is the payload and the Topic is the event name
type DataEvent struct {
	Data  interface{}
	Topic string
}

// Subscribe method can be used to subscribe particular topics
func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append(DataChannelSlice{}, ch)
	}
	eb.rm.Unlock()
}

// Publish method can be used to publish an event
func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.rm.RLock()
	if chans, found := eb.subscribers[topic]; found {
		// A copy of the slice is created here because different slices
		// can share same arrays even though they are passed value.
		// This will ensure that a new slice is created with different
		// base array and hence mutex locking will work properly
		channels := append(DataChannelSlice{}, chans...)

		go func(data DataEvent, dataChannelSlice DataChannelSlice) {
			for _, ch := range dataChannelSlice {
				ch <- data
			}
		}(DataEvent{Data: data, Topic: topic}, channels)
	}
	eb.rm.RUnlock()
}

// Instance is an event bus singleton
var Instance = &EventBus{
	subscribers: make(map[string]DataChannelSlice),
}
