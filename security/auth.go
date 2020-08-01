package security

import "fmt"

// Auth holds the users' database and the access
// granted to the current user. By default it
// should be "NONE". Although it can be changed to
// anything in case of testing
type Auth struct {
	usersDB UnsecureDB
	Access  []Access
}

// RegisteredUser holds the data of a client that is
// using the database. It holds information like
// username, password and permitted access types
type RegisteredUser struct {
	Username string
	Password string
	Access   []Access
}

// Register function registers a new authentication detail and returns the auth object
func Register(username string, password string, access []Access) RegisteredUser {
	return RegisteredUser{username, password, access}
}

// Authenticate authenticates a user but does not handles authorization
// over the database resources!
func (auth *Auth) Authenticate(username string, password string) error {
	user, ok := auth.usersDB.Get(username)
	if ok {
		v, valid := user.(RegisteredUser)
		if !valid {
			panic("Invalid user exists in the DBUser store")
		}

		if v.Password == password {
			auth.Access = v.Access
			return nil
		}
	}

	return fmt.Errorf("Invalid credentials")
}

// Authorize method just authorizes a given action and doesn't handle
// authentication. For authentication Authenticate method should be used
func (auth Auth) Authorize(reqAccess Access) bool {
	for _, access := range auth.Access {
		if access >= reqAccess {
			return true
		}
	}

	return false
}

// RegisterUser creates a new user for the database
func (auth *Auth) RegisterUser(username string, password string, access []Access) error {
	if auth.Authorize(CREATE_USER_ACCESS) {
		ru := Register(username, password, access)
		auth.usersDB.Set(ru.Username, ru, 0)
		return nil
	}

	return deniedErr()
}
