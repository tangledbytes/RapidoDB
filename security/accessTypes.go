package security

// Access type indicates the available access types
type Access uint

const (
	// NONE access type indicates that a user
	// has no permissions at all
	NONE Access = iota

	// READ_ACCESS access type indicates that a user
	// has only read access to the database
	READ_ACCESS

	// WRITE_ACCESS access type indicates that a user
	// has only write access to the database
	WRITE_ACCESS

	// CREATE_USER_ACCESS access type indicates that a
	// user can create new users
	CREATE_USER_ACCESS

	// WIPE_ACCESS access type indicated that a user
	// can wipe out the database
	WIPE_ACCESS

	// ADMIN_ACCESS access type indicated that a user
	// is an admin and hence can perform all of the above tasks
	ADMIN_ACCESS
)
