package db

import (
	"sync"
	"time"
)

const (
	// DefaultExpiry of the item is set to "never expire"
	defaultExpiry time.Duration = NeverExpire
	// NeverExpire is set to -1
	NeverExpire = -1
)

// DB struct encapsualtes the database and its methods
type DB struct {
	sync.RWMutex
	// DefaultExpiry is the default time when
	// an item in the database will expire
	DefaultExpiry time.Duration
	data          map[string]Item
}

// New creates a new database and returns a pointer to it
func New(defaultExpiry time.Duration) *DB {
	return &DB{
		DefaultExpiry: defaultExpiry,
		data:          make(map[string]Item),
	}
}

// Set adds the value in the database
func (db *DB) Set(key string, data interface{}, expireIn time.Duration) {
	// Lock the map
	db.Lock()
	db.data[key] = newItem(data, expireIn)
	// Unlock
	db.Unlock()
}

// Get returns the data corresponding to the key
// If the data doesn't exists then it returns "nil"
func (db *DB) Get(key string) interface{} {
	db.RLock()

	item, ok := db.data[key]
	if !ok || item.isExpired() {
		db.RUnlock()
		return nil
	}

	db.RUnlock()
	return item.data
}
