package security

// Auth holds the data of a client that is
// using the database. It holds information like
// username, password and permitted access types
type Auth struct {
	username string
	password string
	access   []Access
}

// Register function registers a new authentication detail and returns the auth object
func Register(username string, password string, access []Access) Auth {
	return Auth{username, password, access}
}

// Authenticate authenticates a user but does not handles authorization
// over the database resources!
func (auth Auth) Authenticate(username string, password string) bool {
	return username == auth.username && password == auth.password
}

// Authorize method just authorizes a given action and doesn't handle
// authentication. For authentication Authenticate method should be used
func (auth Auth) Authorize(reqAccess Access) bool {
	for _, access := range auth.access {
		if access == reqAccess {
			return true
		}
	}

	return false
}
