package db

import (
	"sync"
	"time"
)

const (
	// Default expiry of the item is set to "never expire"
	defaultExpiry time.Duration = neverExpire
	// Value of never expire
	neverExpire = -1
)

// DB struct encapsualtes the database and its methods
type DB struct {
	*sync.RWMutex
	defaultExpiry time.Duration
	data          map[string]Item
}

// Set adds the value in the database
func (db *DB) Set(key string, value Item) {
	// Lock the map
	db.Lock()
	db.data[key] = value
	// Unlock
	db.Unlock()
}

// Get returns the data corresponding to the key
// If the data doesn't exists then it returns "nil"
func (db *DB) Get(key string) interface{} {
	db.RLock()
	item, ok := db.data[key]
	if !ok || item.IsExpired() {
		db.RUnlock()
		return nil
	}

	db.RUnlock()
	return item.data
}
