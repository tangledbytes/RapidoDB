package security

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

	// CreateUserAccess access type indicates that a
	// user can create new users
	CreateUserAccess

	// WipeAccess access type indicated that a user
	// can wipe out the database
	WipeAccess

	// AdminAccess access type indicated that a user
	// is an admin and hence can perform all of the above tasks
	AdminAccess
)
