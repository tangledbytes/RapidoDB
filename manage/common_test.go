package manage

import "time"

// Mock database
type MockDB struct {
	db map[string]interface{}
}

// Mock Set
func (db *MockDB) Set(key string, data interface{}, _ time.Duration) {
	db.db[key] = data
}

// Mock Get
func (db *MockDB) Get(key string) (interface{}, bool) {
	item, ok := db.db[key]
	if !ok {
		return nil, false
	}

	return item, true
}

// Mock Delete
func (db *MockDB) Delete(key string) (interface{}, bool) {
	item, ok := db.db[key]
	if !ok {
		return nil, false
	}

	delete(db.db, key)
	return item, true
}

// Mock Wipe
func (db *MockDB) Wipe() {
	db.db = make(map[string]interface{})
}

// Mock default expiry
func (db *MockDB) DefaultExpiry() time.Duration {
	return 0
}
