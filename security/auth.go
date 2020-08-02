package security

// UserDB holds the users' database
type UserDB struct {
	UnsecureDB
}

// ActiveUser represents the current active user
// using the database. Each active user has its
// own security layer hence this instance is unique
// to each of the connected client.
//
// By default the Access for the user should be "NONE"
// and should be changed only once the user is authenticated
type ActiveUser struct {
	RegisteredUser
}

// RegisteredUser holds the data of a client that is
// using the database. It holds information like
// username, password and permitted access types
type RegisteredUser struct {
	Username string
	Password string
	Access   Access
}

// newUserDB returns an instance of UserDB struct
func newUserDB(udb UnsecureDB) *UserDB {
	return &UserDB{udb}
}

// newActiveUser returns an instance of ActiveUser
func newActiveUser(username, password string, access Access) *ActiveUser {
	return &ActiveUser{NewRegisteredUser(username, password, access)}
}

// setAccess sets the access level of the current active user
func (au *ActiveUser) setAccess(a Access) {
	au.Access = a
}

// authorize returns true if the current active user
// has the permissions to perform a certain action
func (au *ActiveUser) authorize(reqAccess Access) bool {
	return au.Access >= reqAccess
}

// NewRegisteredUser function registers a new authentication detail and returns the auth object
func NewRegisteredUser(username string, password string, access Access) RegisteredUser {
	return RegisteredUser{username, password, access}
}
