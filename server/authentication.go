package server

import "errors"

// Auth struct holds the user
// and pass of the current server
type Auth struct {
	user string
	pass string
}

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

	if user == a.user && pass == a.pass {
		return true, nil
	}

	return false, errors.New("Invalid Credentials")
}
