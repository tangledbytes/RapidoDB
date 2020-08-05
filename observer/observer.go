package observer

import (
	"time"

	"github.com/utkarsh-pro/RapidoDB/eventbus"
	"github.com/utkarsh-pro/RapidoDB/manage"
)

// ObservedDB adds a very minor layer over the Client Management
// layer and publishes the events that are happening on the database
// to the event bus which later sends those events to the subscribers
//
// ObserverDB is very tightly tied to the Client Management layer and
// the singleton instance of the event bus
type ObservedDB struct {
	*manage.SecureDB
}

// New returns a new observed store
func New(db *manage.SecureDB) (*ObservedDB, *eventbus.EventBus) {
	odb := &ObservedDB{db}
	eb := eventbus.New()

	go setupListenersAndDispatcher(
		odb,
		eb,
		string(opGet),
		string(opSet),
		string(opDel),
		string(opWipe),
	)

	return odb, eb
}

// Set is a thin wrapper over the native set method which adds an observer
// on the set operation.
//
// Whenever a set operation is completed, this publishes a "op_set" event
func (ost *ObservedDB) Set(key string, data interface{}, expireIn time.Duration) error {
	// perform the action
	err := ost.SecureDB.Set(key, data, expireIn)
	// publish the event
	publish(opSet, key, data)

	return err
}

// Get is a thin wapper over the native get method adds an observer
// on the get operation
//
// Whenever a get operation is completed, this published a "op_get" event
func (ost *ObservedDB) Get(key string) (interface{}, bool, error) {
	// perform the action
	v, ok, err := ost.SecureDB.Get(key)
	// publish the event
	publish(opGet, key, v)

	return v, ok, err
}

// Delete is a thin wapper over the native delete method adds an observer
// on the delete operation
//
// Whenever a delete operation is completed, this published a "op_delete" event
func (ost *ObservedDB) Delete(key string) (interface{}, bool, error) {
	// perform the action
	v, ok, err := ost.SecureDB.Delete(key)
	// publish the event
	publish(opDel, key, v)

	return v, ok, err
}

// Wipe is a thin wapper over the native wipe method adds an observer
// on the wipe operation
//
// Whenever a wipe operation is completed, this published a "op_wipe" event
func (ost *ObservedDB) Wipe() error {
	// perform the action
	err := ost.SecureDB.Wipe()
	// publish the event
	publish(opWipe, "wipe", true)

	return err
}

// publish publishes the event to the event bus to be consumed by the subscribers
func publish(event event, key string, value interface{}) {
	eventbus.Instance.Publish(string(event), eventbus.NewDataEvent(string(event), key, value))
}

// setupListenerAndDispatcher sets up the listeners on the multiplexed channel
// it publishes "verified_event" if an event is subscribed by the current client
func setupListenersAndDispatcher(odb *ObservedDB, eb *eventbus.EventBus, events ...string) {
	muxcd := eventbus.ChannelMultiplexer(eventbus.Instance, 0, events...)

	for msg := range muxcd {
		if odb.IsSubscribed(eventToClientEvent(event(msg.Event()))) {
			eb.Publish(
				string(verifiedEvent),
				eventbus.NewDataEvent(string(msg.Event()),
					msg.Key(),
					msg.Value(),
				),
			)
		}
	}
}

// eventToClientEvent converts the local events to the
// events valid in the client management layer
func eventToClientEvent(event event) manage.Event {
	switch event {
	case opGet:
		return manage.GET
	case opSet:
		return manage.SET
	case opDel:
		return manage.DEL
	case opWipe:
		return manage.WIPE
	default:
		return manage.NULL
	}
}
