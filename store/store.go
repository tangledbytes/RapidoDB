package store

import (
	"sync"
	"time"
)

const (
	// NeverExpire constant denotes the symbolic representation
	// for not removing an item from the database
	NeverExpire = 0

	// This is the interval at which the janitor would be triggered
	// It is hardcoded at the moment but shall be made configurable
	// in the future
	janitorInterval = 10 * time.Second
)

// Store struct encapsulates the store used by the database
type Store struct {
	sync.RWMutex
	defaultExpiry time.Duration
	data          map[string]Item
	janitor       *janitor
}

// New returns a new store
func New(defaultExpiry time.Duration) *Store {
	s := &Store{
		defaultExpiry: defaultExpiry,
		data:          make(map[string]Item),
		janitor:       newJanitor(janitorInterval),
	}

	// Setup janitor for this store
	setupJanitor(s)

	return s
}

// Set adds an entry to the map with the corresponding key and data
func (store *Store) Set(key string, data interface{}, expireIn time.Duration) {
	// Lock the map
	store.Lock()
	store.data[key] = newItem(data, expireIn)
	// Unlock the map
	store.Unlock()
}

// Get returns the data stored corresponding to the given key
// if the data is not found then it returns nil
func (store *Store) Get(key string) (interface{}, bool) {
	store.RLock()

	item, ok := store.data[key]
	if !ok || item.isExpired() {
		store.RUnlock()
		return nil, false
	}

	store.RUnlock()
	return item.data, true
}

// Delete method deletes a key from the store. If the key doesn't exists
// then it's a no-op
//
// Method returns the deleted item after the delete operation
// If nothing has been deleted then it returns nil and false
func (store *Store) Delete(key string) (interface{}, bool) {
	store.Lock()

	// Get the item using the native method only
	// Get method on the store was avoided to be used here
	// to avoid creating and removing locks twice which would
	// affect the performance of the store
	item, ok := store.data[key]

	if ok {
		// Delete the key from the map
		// With the current implementation of golang
		// delete function, the runtime doesn't crashes even
		// if the key doesn't exists in the map
		delete(store.data, key)
		store.Unlock()
		return item.data, ok
	}

	store.Unlock()

	return item.data, ok
}

// DeleteExpired loops through the store and deletes
// all the expired items
func (store *Store) DeleteExpired() {
	store.Lock()

	for k, v := range store.data {
		if v.isExpired() {
			// Delete the key from the map
			// With the current implementation of golang
			// delete function, the runtime doesn't crashes even
			// if the key doesn't exists in the map
			delete(store.data, k)
		}
	}

	store.Unlock()
}

// Wipe method clears the entire map by creating a new map
// and assigning a pointer to that map to the "data" attribute
// of the store. Clearing up of that memory is the responsibility
// of the garbage collector
func (store *Store) Wipe() {
	store.Lock()
	store.data = make(map[string]Item)
	store.Unlock()
}

// DefaultExpiry returns the default expiry of the store items
func (store *Store) DefaultExpiry() time.Duration {
	return store.defaultExpiry
}
