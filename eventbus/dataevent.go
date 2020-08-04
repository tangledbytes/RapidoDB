package eventbus

import "fmt"

// DataEvent is the structure of the message packets
// that are passed through the event bus
type DataEvent struct {
	key   string
	value interface{}
}

// String returns the string representation of the
// DataEvent Object
func (de DataEvent) String() string {
	return fmt.Sprintf("Key: %v Value: %v", de.key, de.value)
}
