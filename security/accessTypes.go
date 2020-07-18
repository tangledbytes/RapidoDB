package security

// Access type indicates the available access types
type Access uint

const (
	// READ_ACCESS access type indicates that a user
	// has only read access to the database
	READ_ACCESS Access = iota

	// WRITE_ACCESS access type indicates that a user
	// has only write access to the database
	WRITE_ACCESS

	// WIPE_ACCESS access type indicated that a user
	// can wipe out the database
	WIPE_ACCESS

	// ADMIN_ACCESS access type indicated that a user
	// is an admin and hence can perform all of the above tasks
	ADMIN_ACCESS
)
