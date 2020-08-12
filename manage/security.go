package manage

import (
	"fmt"
	"time"
)

// SecureDB has the unsecure store and adds the user info to it
// this abstraction adds relevant methods to modify database user's
// and client's information
//
// This addition of user's info along with the store itself makes this
// an "SecureDB"
type SecureDB struct {
	// UnsecureStore is used to store the data
	ust UnsecureStore

	// userdb is used internally to store the information
	// of the users of the database
	userdb *UserDB

	// active client holds the information of the currently
	// active client using the layer
	activeClient *DBClient
}

////////////// DATABASE SPECIFIC COMMANDS //////////////////

// Set method performs set operation on the database after checking
// the user permissions
func (sdb *SecureDB) Set(key string, data interface{}, expireIn time.Duration) error {
	if sdb.Authorize(WriteAccess) {
		sdb.ust.Set(key, data, expireIn)
		return nil
	}

	return deniedErr()
}

// Get method performs get operation on the database after checking
// the user permissions
func (sdb *SecureDB) Get(key string) (interface{}, bool, error) {
	if sdb.Authorize(ReadAccess) {
		i, b := sdb.ust.Get(key)
		return i, b, nil
	}

	return nil, false, deniedErr()
}

// Delete method performs delete operation on the database after
// checking the permissions
func (sdb *SecureDB) Delete(key string) (interface{}, bool, error) {
	if sdb.Authorize(WriteAccess) {
		i, b := sdb.ust.Delete(key)
		return i, b, nil
	}

	return nil, false, deniedErr()
}

// Wipe method performs wipe operation on the database after
// checking the permissions
func (sdb *SecureDB) Wipe() error {
	if sdb.Authorize(WipeAccess) {
		sdb.ust.Wipe()
		return nil
	}

	return deniedErr()
}

////////////// ACTIVE CLIENT SPECIFIC COMMANDS //////////////////

// RegisterUser registers a new user with specified username, password and access level
// it does not check for the already existing user with the same username. If a user with
// same username exists then it will overwrite that user's data
func (sdb *SecureDB) RegisterUser(username, password string, access uint) error {
	if sdb.Authorize(ModifyUserAccess) {
		a, err := ConvertUintToAccess(access)
		if err != nil {
			return err
		}

		// Add a new user to the userdb
		sdb.userdb.New(username, password, a, Events{})

		return nil
	}

	return deniedErr()
}

// Authenticate authenticates a client and changes the permissions for the
// active client to that allocated to it
func (sdb *SecureDB) Authenticate(username, password string) error {
	user, ok := sdb.userdb.FindUserByUsername(username)
	if !ok || user.Password != password {
		return fmt.Errorf("Invalid Credentials")
	}

	sdb.ChangeActiveClient(user.Username, user.Password, user.Access, user.Events)
	return nil
}

// Ping will subscribe the currently active user to the
// passed in event. Only admins can use this method
func (sdb *SecureDB) Ping(event string) error {
	if !sdb.Authorize(AdminAccess) {
		return deniedErr()
	}

	ev, err := ConvertStringToEvent(event)
	if err != nil {
		return err
	}

	username := sdb.activeClient.Username
	password := sdb.activeClient.Password
	access := sdb.activeClient.Access
	events := sdb.activeClient.Events.Set(ev)

	// Change the currently active client
	sdb.ChangeActiveClient(username, password, access, events)

	// Update the same in the users database
	sdb.userdb.New(username, password, access, events)

	return nil
}

// Authorize authorizes the requests and returns true if a client
// is permitted to perform a certain action
func (sdb *SecureDB) Authorize(reqAccess Access) bool {
	return sdb.activeClient.Access >= reqAccess
}

// ChangeActiveClient changes the active client of the database by assigning new username
// password and access levels to the activeClient attribute
func (sdb *SecureDB) ChangeActiveClient(username, password string, access Access, events Events) {
	sdb.activeClient = newDBClient(username, password, access, events)
}

// IsSubscribed returns true if the active client has subscribed
// to the mentioned event
func (sdb *SecureDB) IsSubscribed(event Event) bool {
	return sdb.activeClient.Events.Exists(event)
}

// ========================= HELPER FUNCTIONS =============================

// deniedErr returns a pre formatted error
func deniedErr() error {
	return fmt.Errorf("Access denied")
}
