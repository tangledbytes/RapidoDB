package manage

import "time"

// UnsecureStore is the interface of the underlying store and supports
// every read, write operation other than authentication and authorization
type UnsecureStore interface {
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

	// DefaultExpiry returns the default expiration time for the specified store
	DefaultExpiry() time.Duration
}

// New function returns an instance of an UnsecureDB
// Both of the parameters can be any store that satisfies the
// UnsecureStore interface. The first store would be used to
// store the data provided by the users while the second store
// would be used internally to store the user's info
func New(unsecureStore UnsecureStore, userdb UnsecureStore) *SecureDB {
	// Here a new DBClient has no username, password and has no privileges
	// and are not associated with any events
	return &SecureDB{unsecureStore, &UserDB{userdb}, newDBClient("", "", NONE, Events{})}
}
