package store

import (
	"sync"
	"time"
)

const (
	// NeverExpire constant denotes the symbolic representation
	// for not removing an item from the database
	NeverExpire = -1
)

// Store struct encapsulates the store used by the database
type Store struct {
	sync.RWMutex
	DefaultExpiry time.Duration
	data          map[string]Item
}

// New returns a new store
func New(defaultExpiry time.Duration) *Store {
	return &Store{
		DefaultExpiry: defaultExpiry,
		data:          make(map[string]Item),
	}
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
