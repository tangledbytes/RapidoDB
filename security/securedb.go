package security

import (
	"fmt"
	"time"
)

// UnsecureDB is the interface of the underlying store and supports
// every read, write operation other than authentication and authorization
type UnsecureDB interface {
	Set(key string, data interface{}, expireIn time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string) (interface{}, bool)
	Wipe()
}

// SecureDB represents the object which deals with the security aspects
// of the database. It binds the lower and uppper layers of the database.
type SecureDB struct {
	db UnsecureDB
	*Auth
}

// New returns a new instance of the security driver
func New(db UnsecureDB) *SecureDB {
	return &SecureDB{
		db:   db,
		Auth: &Auth{"admin", "pass", []Access{ADMIN_ACCESS}, false},
	}
}

// Set method performs set operation on the database after checking
// the user permissions
func (d *SecureDB) Set(key string, data interface{}, expireIn time.Duration) error {
	if d.IsAuthenticated && d.Authorize(WRITE_ACCESS) {
		d.db.Set(key, data, expireIn)
		return nil
	}

	return deniedErr()
}

// Get method performs get operation on the database after checking
// the user permissions
func (d *SecureDB) Get(key string) (interface{}, bool, error) {
	if d.IsAuthenticated && d.Authorize(READ_ACCESS) {
		i, b := d.db.Get(key)
		return i, b, nil
	}

	return nil, false, deniedErr()
}

// Delete method performs delete operation on the database after
// checking the permissions
func (d *SecureDB) Delete(key string) (interface{}, bool, error) {
	if d.IsAuthenticated && d.Authorize(WRITE_ACCESS) {
		i, b := d.db.Delete(key)
		return i, b, nil
	}

	return nil, false, deniedErr()
}

// Wipe method performs wipe operation on the database after
// checking the permissions
func (d *SecureDB) Wipe() error {
	if d.IsAuthenticated && d.Authorize(WIPE_ACCESS) {
		d.db.Wipe()
		return nil
	}

	return deniedErr()
}

func deniedErr() error {
	return fmt.Errorf("Access denied")
}
