package manage

// UserDB is just an abstraction over the UnsecureStore
// it is meant to be used to store the database user's info
// so the authentication and authorization are not required
// by this store
type UserDB struct {
	UnsecureStore
}

// New adds a new db user to the database
func (udb *UserDB) New(username, password string, access Access, events Events) {
	udb.Set(username, NewDBUser(username, password, access, events), udb.DefaultExpiry())
}

// FindUserByUsername finds a user by its username in the user database
// It returns the DBUser and true if the user exists or an empty
// DBUser object and false
func (udb *UserDB) FindUserByUsername(username string) (DBUser, bool) {
	user, ok := udb.Get(username)
	if !ok {
		return NewDBUser("", "", NONE, Events{}), false
	}

	return ToDBUser(user), true
}
