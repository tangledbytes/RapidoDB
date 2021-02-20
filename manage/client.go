package manage

// DBClient is composed of DBUser
// They both are very similar in nature but are
// intended for different use cases
//
// DBUser is meant to indicate a user in the database
// whereas a DBClient is meant to represent the current
// active client which could be anyone of the users mentioned
// in the user's database
type DBClient struct {
	DBUser
}

// newDBClient creates an instance of the DBClient and returns a pointer
// to the instance
func newDBClient(username, password string, access Access, events Events) *DBClient {
	return &DBClient{NewDBUser(username, password, access, events)}
}
