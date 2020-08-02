package manage

import (
	"reflect"
	"testing"
	"time"
)

// Mock database
type MockDB struct {
	db map[string]interface{}
}

// Mock Set
func (db *MockDB) Set(key string, data interface{}, _ time.Duration) {
	db.db[key] = data
}

// Mock Get
func (db *MockDB) Get(key string) (interface{}, bool) {
	item, ok := db.db[key]
	if !ok {
		return nil, false
	}

	return item, true
}

// Mock Delete
func (db *MockDB) Delete(key string) (interface{}, bool) {
	item, ok := db.db[key]
	if !ok {
		return nil, false
	}

	delete(db.db, key)
	return item, true
}

// Mock Wipe
func (db *MockDB) Wipe() {
	db.db = make(map[string]interface{})
}

// Mock default expiry
func (db *MockDB) DefaultExpiry() time.Duration {
	return 0
}

func TestUserDB_FindUserByUsername(t *testing.T) {
	type fields struct {
		UnsecureStore UnsecureStore
	}
	type args struct {
		username string
	}

	db := &MockDB{make(map[string]interface{})}

	// Add a user to the database
	db.Set("utkarsh", NewDBUser("utkarsh", "test", 5), db.DefaultExpiry())

	tests := []struct {
		name   string
		fields fields
		args   args
		want   DBUser
		want1  bool
	}{
		{
			"FIND A USER THAT EXISTS",
			fields{db},
			args{"utkarsh"},
			NewDBUser("utkarsh", "test", 5),
			true,
		},
		{
			"FIND A USER THAT DOESN'T EXISTS",
			fields{db},
			args{"utkarsh2"},
			NewDBUser("", "", 0),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udb := &UserDB{
				UnsecureStore: tt.fields.UnsecureStore,
			}
			got, got1 := udb.FindUserByUsername(tt.args.username)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserDB.FindUserByUsername() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UserDB.FindUserByUsername() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
