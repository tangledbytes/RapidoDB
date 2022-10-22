package manage

import (
	"fmt"
	"strings"
)

// Event are the events to which a database
// user can subscribe to
type Event uint

const (
	// NULL event indicates an
	// event
	NULL Event = iota

	// GET event indicates a GET
	// operation on the database
	GET

	// SET event indicates a SET
	// operation on the database
	SET

	// DEL event indicates a DEL
	// operation on the database
	DEL

	// WIPE event indicates a WIPE
	// operation on the database
	WIPE
)

// ConvertStringToEvent takes an event as a string and returns
// an event object for it. If the passed event is not a valid
// event then an error is generated
func ConvertStringToEvent(event string) (Event, error) {
	switch strings.ToLower(event) {
	case "get":
		return GET, nil
	case "set":
		return SET, nil
	case "del":
		return DEL, nil
	case "wipe":
		return WIPE, nil
	default:
		return NULL, fmt.Errorf("Invalid event")
	}
}

// Events is the slice of event
type Events []Event

// Set adds an event to the Events if it doesn't already
// exits in the Events and returns a new Events object
func (e Events) Set(event Event) Events {
	if _, exists := e.Exists(event); !exists {
		e = append(e, event)
	}

	return e
}

// Unset removes an event to the Events if it  already
// exits in the Events and returns a new Events object
func (e Events) Unset(event Event) Events {
	if i, exists := e.Exists(event); exists {
		e = append(e[0:i], e[i+1:]...)
	}
	return e
}

// Exists returns true if the given event
// exists in the Events
func (e Events) Exists(event Event) (int, bool) {
	for i, ev := range e {
		if ev == event {
			return i, true
		}
	}

	return -1, false
}

// convertInterfaceSliceToEvents will attempt to convert an
// array of interfaces to Events object. It panics if the the
// type assertion fails
func convertInterfaceSliceToEvents(ui []interface{}) Events {
	var ev Events
	for _, v := range ui {
		d, ok := v.(uint)
		if !ok {
			panic("Cannot convert interface to Event")
		}
		ev = ev.Set(Event(d))
	}

	return ev
}
