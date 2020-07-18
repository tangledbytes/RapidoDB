package security

// Access type indicates the available access types
type Access uint

const (
	// READ access type indicates that a user
	// has only read access to the database
	READ Access = iota

	// WRITE access type indicates that a user
	// has only write access to the database
	WRITE

	// WIPE access type indicated that a user
	// can wipe out the database
	WIPE

	// ADMIN access type indicated that a user
	// is an admin and hence can perform all of the above tasks
	ADMIN
)
