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
func PingClient(c ClientConn, eb *eventbus.EventBus, events ...string) {
	muxcd := eventbus.ChannelMultiplexer(eb, 0, events...)

	go func(ch eventbus.DataChannel) {
		for msg := range muxcd {
			c.Msg(msg.String())
		}
	}(muxcd)
}
