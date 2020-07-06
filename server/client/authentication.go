package client

import "errors"

// Auth struct holds the user
// and pass of the current server
type Auth struct {
	User string
	Pass string
}

const (
	// NotAuthenticated indicates that the user is not authenticated
	NotAuthenticated = false
	// Authenticated indicates that the user is authenticated
	Authenticated = true
)

// HandleAuth extracts the user and pass from the command
// and check it against the passed user and pass
func (a Auth) HandleAuth(cmdString []string) (bool, error) {
	// Check if the command length is 3
	// AUTH <user> <pass>
	if len(cmdString) != 3 {
		return false, errors.New("Invalid AUTH command")
	}

	// Extract the user and pass
	user := cmdString[1]
	pass := cmdString[2]

	if user == a.User && pass == a.Pass {
		return true, nil
	}

	return false, errors.New("Invalid Credentials")
}
