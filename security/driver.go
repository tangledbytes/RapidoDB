package security

import (
	"time"
)

// UnsecureDB is composed of the IRapidoDB interface
// and hence supports all the methods supported by it
type UnsecureDB interface {
	Set(key string, data interface{}, expireIn time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string) (interface{}, bool)
}

// Driver represents the object which deals with the security aspects
// of the database. It binds the lower and uppper layers of the database.
type Driver struct {
	db UnsecureDB
	*Auth
}

// New returns a new instance of the security driver
func New(db UnsecureDB) *Driver {
	return &Driver{
		db:   db,
		Auth: &Auth{"admin", "pass", []Access{ADMIN_ACCESS}, false},
	}
}

// Set method performs set operation on the database after checking
// the user permissions
func (d *Driver) Set(key string, data interface{}, expireIn time.Duration) {
	if d.IsAuthenticated && d.Authorize(WRITE_ACCESS) {
		d.db.Set(key, data, expireIn)
	}
}

// Get method performs get operation on the database after checking
// the user permissions
func (d *Driver) Get(key string) (interface{}, bool) {
	if d.IsAuthenticated && d.Authorize(READ_ACCESS) {
		return d.db.Get(key)
	}

	return nil, false
}

// Delete method performs delete operation on the database after
// checking the permissions
func (d *Driver) Delete(key string) (interface{}, bool) {
	if d.IsAuthenticated && d.Authorize(WRITE_ACCESS) {
		return d.db.Delete(key)
	}

	return nil, false
}
