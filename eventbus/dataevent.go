package eventbus

import "fmt"

// DataEvent is the structure of the message packets
// that are passed through the event bus
type DataEvent struct {
	event string
	key   string
	value interface{}
}

// Event returns the event of the DataEvent
func (de DataEvent) Event() string {
	return de.event
}

// Key returns the key of the DataEvent
func (de DataEvent) Key() string {
	return de.key
}

// Value returns the value of the DataEvent
func (de DataEvent) Value() interface{} {
	return de.value
}

// String returns the string representation of the
// DataEvent Object
func (de DataEvent) String() string {
	return fmt.Sprintf("Event: %v Key: %v Value: %v", de.event, de.key, de.value)
}
