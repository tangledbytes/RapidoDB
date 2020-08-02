package manage

import "fmt"

// Access type indicates the available access types
type Access uint

const (
	// NONE access type indicates that a user
	// has no permissions at all
	NONE Access = iota

	// ReadAccess access type indicates that a user
	// has only read access to the database
	ReadAccess

	// WriteAccess access type indicates that a user
	// has only write access to the database
	WriteAccess

	// ModifyUserAccess access type indicates that a
	// user can create new users
	ModifyUserAccess

	// WipeAccess access type indicated that a user
	// can wipe out the database
	WipeAccess

	// AdminAccess access type indicated that a user
	// is an admin and hence can perform all of the above tasks
	AdminAccess
)

// ConvertUintToAccess converts the uint type to Access type
func ConvertUintToAccess(access uint) (Access, error) {
	if access > uint(AdminAccess) {
		return 0, fmt.Errorf("Access parameter too high, max can be %d", AdminAccess)
	}

	if access < uint(NONE) {
		return 0, fmt.Errorf("Access parameter too low, min can be %d", NONE)
	}

	return Access(access), nil
}
