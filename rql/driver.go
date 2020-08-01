package rql

import (
	"fmt"
	"time"
)

// SecureDB interface defines the set of functions that RQL
// driver expects.
type SecureDB interface {
	Set(key string, data interface{}, expireIn time.Duration) error
	Get(key string) (interface{}, bool, error)
	Delete(key string) (interface{}, bool, error)
	Wipe() error
	Authenticate(username string, password string) error
	RegisterUser(username string, password string, access uint) error
}

// Driver is the RQL driver which acts as an interface between a database client and
// the querying language parser
//
// Driver takes in a database which conforms to the RQL DB interaface.
// This ensures that the RQL driver isn't tied to a single implementation
// of the database. Any database API that conforms this interface will work
type Driver struct {
	db SecureDB
}

// New function returns a pointer to an instance of RQL driver
func New(db SecureDB) *Driver {
	return &Driver{db}
}

// Operate method can take in any RQL query and perform action
// based on query.
//
// This method doesn't returns anything even if the query is invalid
// instead it will use the io.Writer to write the response to the
// specified stream
func (d *Driver) Operate(src string) (string, error) {
	// Parse the src
	ast, err := Parse(src)
	if err != nil {
		return "", err
	}
	if ast == nil {
		return "", nil
	}

	var result string

	for _, stmt := range ast.Statements {
		switch stmt.Typ {
		case SetType:
			res, err := d.set(stmt.SetStatement)
			if err != nil {
				return result, err
			}
			result = prepareResponse(result, res)
		case GetType:
			res, err := d.get(stmt.GetStatement)
			if err != nil {
				return result, err
			}
			result = prepareResponse(result, res)
		case DeleteType:
			res, err := d.delete(stmt.DeleteStatement)
			if err != nil {
				return result, err
			}
			result = prepareResponse(result, res)
		case WipeType:
			res, err := d.wipe(stmt.WipeStatement)
			if err != nil {
				return result, err
			}
			result = prepareResponse(result, res)
		case AuthType:
			res, err := d.auth(stmt.AuthStatement)
			if err != nil {
				return result, err
			}
			result = prepareResponse(result, res)
		case RegUserType:
			res, err := d.reguser(stmt.RegUserStatement)
			if err != nil {
				return result, err
			}
			result = prepareResponse(result, res)
		}
	}

	return result, nil
}

// set method calls the set method on the database by providing
// appropriate parameters
func (d *Driver) set(stmt *SetStatement) (string, error) {
	err := d.db.Set(stmt.key, stmt.val, convertToDuration(stmt.exp))

	if err != nil {
		return "", err
	}
	return "Success", nil
}

// get method calls the get method on the database by providing
// appropriate parameters
// it ignores the "keys" which do not exists in the database and places
// nil in the slice for them
//
// It returns the stringified slice
func (d *Driver) get(stmt *GetStatement) (string, error) {
	var res []interface{}

	for _, key := range stmt.keys {
		val, _, err := d.db.Get(key)
		if err != nil {
			return "", err
		}
		res = append(res, val)
	}

	return stringify(res), nil
}

// delete method calls the delete method on the database by providing
// appropriate parameters
// it ignores the "keys" which do not exists in the database and places
// nil in the slice for them
//
// It returns the stringified slice
func (d *Driver) delete(stmt *DeleteStatement) (string, error) {
	var res []interface{}

	for _, key := range stmt.keys {
		val, _, err := d.db.Delete(key)
		if err != nil {
			return "", err
		}
		res = append(res, val)
	}

	return stringify(res), nil
}

// wipe method call the wipe method on the secure database
// if any error occurs in the process then that error is passed
// on to the client
func (d *Driver) wipe(stmt *WipeStatement) (string, error) {
	if err := d.db.Wipe(); err != nil {
		return "", err
	}

	return "Success", nil
}

// auth takes in the authStatement and executes Authenticate method on the database
func (d *Driver) auth(stmt *AuthStatement) (string, error) {
	if err := d.db.Authenticate(stmt.username, stmt.password); err != nil {
		return "", err
	}

	return "Successfully Authenticated", nil
}

// reguser takes username, password and access level for the user and creates a newuser
// by invoking the RegiseterUser method on the SecureDB
func (d *Driver) reguser(stmt *RegUserStatement) (string, error) {
	if err := d.db.RegisterUser(stmt.username, stmt.password, stmt.access); err != nil {
		return "", err
	}

	return "Created user " + stmt.username, nil
}

// ============================ HELPER FUNCTIONS ===================================

// convertToDuration converts uint to time.Duration object.
// This uint is supposed to be in MILLISECONDS.
// It's internally converted into nanoseconds and is then casted into
// time.Duration object
func convertToDuration(t uint) time.Duration {
	return time.Duration(time.Duration(t) * time.Millisecond)
}

// stringify function can be used to stringify any data type
// It internally uses fmt.Sprintf("%v", ...) to perform the conversion
// which internally uses the String() method on the objects to perform the conversion
func stringify(any interface{}) string {
	return fmt.Sprintf("%v", any)
}

// prepareResponse just concatanates the passed string and separates
// them with a newline character
func prepareResponse(str1, str2 string) string {
	if str1 == "" {
		return str2
	}
	return str1 + "\n" + str2
}
