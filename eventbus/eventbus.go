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

// DataEvent is the structure of the message packets
// that are passed through the event bus
type DataEvent struct {
	key   string
	value interface{}
}

// NewDataEvent creates a new Data Event from the passed key value pairs
func NewDataEvent(key string, value interface{}) DataEvent {
	return DataEvent{key, value}
}

// Subscribe method can be used to subscribe particular topics
func (eb *EventBus) Subscribe(topic string, buf uint) DataChannel {
	eb.rm.Lock()

	ch := make(DataChannel, buf)
	eb.subscribers[topic] = append(eb.subscribers[topic], ch)

	eb.rm.Unlock()
	return ch
}

// Publish method can be used to publish an event
func (eb *EventBus) Publish(topic string, data DataEvent) {
	eb.rm.RLock()

	for _, ch := range eb.subscribers[topic] {
		go func(ch DataChannel, data DataEvent) {
			ch <- data
		}(ch, data)
	}

	eb.rm.RUnlock()
}

// New returns a new event bus
//
// Although it should be rarely used as singleton
// for EventBus is already being exported
func New() *EventBus {
	return &EventBus{
		subscribers: make(map[string]DataChannelSlice),
	}
}

// Instance is an event bus singleton
var Instance = New()
