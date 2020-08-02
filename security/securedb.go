package security

import (
	"fmt"
	"time"
)

// UnsecureDB is the interface of the underlying store and supports
// every read, write operation other than authentication and authorization
type UnsecureDB interface {
	// Set method should add the passed key into the store with the provided data
	Set(key string, data interface{}, expireIn time.Duration)

	// Get method should return the value corresponding to the provided key
	// if the value doesn't exist in the store then the bool should be false
	Get(key string) (interface{}, bool)

	// Delete method should delete the key value pair for the provided key
	// it should return the deleted value. The second returned value should
	// be true if the item was found in the store and was successfully removed
	// if not found then this value should be false
	Delete(key string) (interface{}, bool)

	// Wipe method should just wipe the entire underlying database
	// and should start with a fresh store again
	Wipe()
}

// SecureDB represents the object which deals with the security aspects
// of the database. It binds the lower and uppper layers of the database.
type SecureDB struct {
	// Original store for the database
	db UnsecureDB

	// Metadata for the users
	userdb *UserDB

	// Metadata for the current active user
	activeUser *ActiveUser
}

// New returns a new instance of the security driver
func New(db UnsecureDB, userDB UnsecureDB) *SecureDB {
	return &SecureDB{
		db:         db,
		userdb:     newUserDB(userDB),
		activeUser: newActiveUser("", "", NONE),
	}
}

// Set method performs set operation on the database after checking
// the user permissions
func (d *SecureDB) Set(key string, data interface{}, expireIn time.Duration) error {
	if d.Authorize(WriteAccess) {
		d.db.Set(key, data, expireIn)
		return nil
	}

	return deniedErr()
}

// Get method performs get operation on the database after checking
// the user permissions
func (d *SecureDB) Get(key string) (interface{}, bool, error) {
	if d.Authorize(ReadAccess) {
		i, b := d.db.Get(key)
		return i, b, nil
	}

	return nil, false, deniedErr()
}

// Delete method performs delete operation on the database after
// checking the permissions
func (d *SecureDB) Delete(key string) (interface{}, bool, error) {
	if d.Authorize(WriteAccess) {
		i, b := d.db.Delete(key)
		return i, b, nil
	}

	return nil, false, deniedErr()
}

// Wipe method performs wipe operation on the database after
// checking the permissions
func (d *SecureDB) Wipe() error {
	if d.Authorize(WipeAccess) {
		d.db.Wipe()
		return nil
	}

	return deniedErr()
}

// Authorize method just authorizes a given action and doesn't handle
// authentication. For authentication Authenticate method should be used
func (d *SecureDB) Authorize(reqAccess Access) bool {
	return d.activeUser.authorize(reqAccess)
}

// Authenticate authenticates a user but does not handles authorization
// over the database resources
func (d *SecureDB) Authenticate(username, password string) error {
	user, ok := d.userdb.Get(username)
	if ok {
		v, valid := user.(RegisteredUser)
		if !valid {
			panic("Invalid user exists in the DBUser store")
		}

		if v.Password == password {
			d.activeUser = newActiveUser(username, password, v.Access)
			return nil
		}
	}

	return fmt.Errorf("Invalid credentials")
}

// RegisterUser creates a new user for the database
func (d *SecureDB) RegisterUser(username string, password string, access uint) error {
	if d.Authorize(CreateUserAccess) {
		// Convert uint to Access
		if access > uint(AdminAccess) {
			return fmt.Errorf("Invalid access level, max can be %d for admins", AdminAccess)
		}
		a := Access(access)

		ru := NewRegisteredUser(username, password, a)
		d.userdb.Set(ru.Username, ru, 0)
		return nil
	}

	return deniedErr()
}

// deniedErr returns a pre formatted error
func deniedErr() error {
	return fmt.Errorf("Access denied")
}
