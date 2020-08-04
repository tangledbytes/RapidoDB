package transportext

import (
	"github.com/utkarsh-pro/RapidoDB/eventbus"
)

// ClientConn interface describes the interface
// for sending messages to the clients
type ClientConn interface {
	Msg(string)
}

// PingClient takes in a net.Conn object and variadic number
// of events to which this will subscribe and will automatically
// send them to the client
func PingClient(c ClientConn, events ...string) {
	chs := subscribeToEvents(events...)

	// Process all the channels
	for _, ch := range chs {
		go func(ch eventbus.DataChannel) {
			for msg := range ch {
				c.Msg(msg.String())
			}
		}(ch)
	}
}

// subscribeToEvents subscribe to all the events passed
// as function. This is a variadic function and hence can take in
// any number of events.
//
// It returns DataChannelSlice which can be used to listen to messages
// passed on the channel by the event bus
func subscribeToEvents(events ...string) eventbus.DataChannelSlice {
	var ch eventbus.DataChannelSlice
	for _, e := range events {
		ch = append(ch, eventbus.Instance.Subscribe(e, 0))
	}

	return ch
}
