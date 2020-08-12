package manage

// DBUser represents a user of the database
type DBUser struct {
	// Username is the Username that the user of the
	// database is supposed to use to authenticate themselves
	Username string

	// Password is the Password that the user of the database
	// is supposed to use to authenticate themselves
	Password string

	// Access determines what kind of permissions were allocated
	// for a certain user. A user will be able to do only the
	// tasks which are feasible with the Access level assigned
	// to them during user creation.
	Access Access

	// Events determines all the Events to which a database
	// user has subscribed
	Events Events
}

// NewDBUser creates a new database user object and return it
// It does not create an entry in the user's database for the user
func NewDBUser(username, pass string, access Access, events Events) DBUser {
	return DBUser{username, pass, access, events}
}

// ToDBUser converts an interface{} to DBUser type
// if the passed interface{} is not a DBUser then the
// function panics
func ToDBUser(data interface{}) DBUser {
	v, ok := data.(DBUser)
	if !ok {
		// Check if the data is of type "map", if posible then create a dbuser from the map
		mp, ok := data.(map[string]interface{})
		if !ok {
			panic("Invalid user exists in the DBUser store")
		}

		// Check if "Username" is available
		iun, ok := mp["Username"]
		un, ok := iun.(string)
		if !ok {
			panic("Invalid user exists in the DBUser store: Invalid username")
		}

		// Check if "Password" is available
		ips, ok := mp["Password"]
		ps, ok := ips.(string)
		if !ok {
			panic("Invalid user exists in the DBUser store: Invalid password")
		}

		// Check if "Access" is available
		iac, ok := mp["Access"]
		ac, ok := iac.(float64)
		if !ok {
			panic("Invalid user exists in the DBUser store: Invalid access type")
		}

		// Check if "Events" is available
		iev, ok := mp["Events"]
		ev, ok := iev.([]interface{})
		if !ok {
			panic("Invalid user exists in the DBUser store: Invalid events")
		}

		v = NewDBUser(un, ps, Access(ac), convertInterfaceSliceToEvents(ev))
	}

	return v
}
