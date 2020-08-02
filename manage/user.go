package manage

// DBUser represents a user of the database
type DBUser struct {
	// username is the username that the user of the
	// database is supposed to use to authenticate themselves
	username string

	// password is the password that the user of the database
	// is supposed to use to authenticate themselves
	password string

	// access determines what kind of permissions were allocated
	// for a certain user. A user will be able to do only the
	// tasks which are feasible with the access level assigned
	// to them during user creation.
	access Access
}

// NewDBUser creates a new database user object and return it
// It does not create an entry in the user's database for the user
func NewDBUser(username, pass string, access Access) DBUser {
	return DBUser{username, pass, access}
}

// ToDBUser converts an interface{} to DBUser type
// if the passed interface{} is not a DBUser then the
// function panics
func ToDBUser(data interface{}) DBUser {
	v, ok := data.(DBUser)
	if !ok {
		panic("Invalid user exists in the DBUser store")
	}

	return v
}
