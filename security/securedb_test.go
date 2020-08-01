package security

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

func TestSecureDB_Set(t *testing.T) {
	type fields struct {
		db   UnsecureDB
		Auth *Auth
	}
	type args struct {
		key      string
		data     interface{}
		expireIn time.Duration
	}

	m1 := make(map[string]interface{})
	m2 := make(map[string]interface{})
	db := &MockDB{m1}
	udb := &MockDB{m2}
	auth := &Auth{udb, []Access{ADMIN_ACCESS}}
	auth2 := &Auth{udb, []Access{READ_ACCESS}}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"SET STRING",
			fields{db, auth},
			args{
				"d1",
				"Hello World",
				time.Duration(100),
			},
			false,
		},
		{
			"SET NUMBERS",
			fields{db, auth},
			args{
				"d1",
				100,
				time.Duration(100),
			},
			false,
		},
		{
			"UNAUTHORIZED SET OPERATION",
			fields{db, auth2},
			args{
				"d3",
				100,
				time.Duration(100),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &SecureDB{
				db:   tt.fields.db,
				Auth: tt.fields.Auth,
			}
			if err := d.Set(tt.args.key, tt.args.data, tt.args.expireIn); (err != nil) != tt.wantErr {
				t.Errorf("Driver.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecureDB_Get(t *testing.T) {
	type fields struct {
		db   UnsecureDB
		Auth *Auth
	}
	type args struct {
		key string
	}

	m1 := make(map[string]interface{})
	m2 := make(map[string]interface{})
	db := &MockDB{m1}
	udb := &MockDB{m2}
	auth := &Auth{udb, []Access{ADMIN_ACCESS}}
	auth2 := &Auth{udb, []Access{NONE}}

	// Add a key to the map -> key = 'd1'
	db.Set("d1", "Hello World", 0)

	// Add a key to the map -> key = 'd3'
	db.Set("d3", 2345, 0)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		want1   bool
		wantErr bool
	}{
		{
			"GET STRING WHEN KEY IS PRESENT",
			fields{db, auth},
			args{
				"d1",
			},
			"Hello World",
			true,
			false,
		},
		{
			"GET NIL WHEN KEY IS ABSENT",
			fields{db, auth},
			args{
				"d2",
			},
			nil,
			false,
			false,
		},
		{
			"GET NUMBER WHEN KEY IS PRESENT",
			fields{db, auth},
			args{
				"d3",
			},
			2345,
			true,
			false,
		},
		{
			"UNAUTHORIZED GET NUMBER WHEN KEY IS PRESENT",
			fields{db, auth2},
			args{
				"d3",
			},
			nil,
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &SecureDB{
				db:   tt.fields.db,
				Auth: tt.fields.Auth,
			}
			got, got1, err := d.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Driver.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Driver.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Driver.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSecureDB_Delete(t *testing.T) {
	type fields struct {
		db   UnsecureDB
		Auth *Auth
	}
	type args struct {
		key string
	}

	m1 := make(map[string]interface{})
	m2 := make(map[string]interface{})
	db := &MockDB{m1}
	udb := &MockDB{m2}
	auth := &Auth{udb, []Access{ADMIN_ACCESS}}
	auth2 := &Auth{udb, []Access{READ_ACCESS}}

	// Add a key to the map -> key = 'd1'
	db.Set("d1", "Hello World", 0)

	// Add a key to the map -> key = 'd3'
	db.Set("d3", 2345, 0)

	// Add a key to the map -> key = 'd4'
	db.Set("d4", 2345, 0)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		want1   bool
		wantErr bool
	}{
		{
			"DELETE STRING WHEN KEY IS PRESENT",
			fields{db, auth},
			args{
				"d1",
			},
			"Hello World",
			true,
			false,
		},
		{
			"GET NIL WHEN KEY IS ABSENT",
			fields{db, auth},
			args{
				"d2",
			},
			nil,
			false,
			false,
		},
		{
			"DELETE NUMBER WHEN KEY IS PRESENT",
			fields{db, auth},
			args{
				"d3",
			},
			2345,
			true,
			false,
		},
		{
			"UNAUTHORIZED DELETE STRING WHEN KEY IS PRESENT",
			fields{db, auth2},
			args{
				"d4",
			},
			nil,
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &SecureDB{
				db:   tt.fields.db,
				Auth: tt.fields.Auth,
			}
			got, got1, err := d.Delete(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Driver.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Driver.Delete() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Driver.Delete() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSecureDB_Wipe(t *testing.T) {
	type fields struct {
		db   UnsecureDB
		Auth *Auth
	}

	m1 := make(map[string]interface{})
	m2 := make(map[string]interface{})
	db := &MockDB{m1}
	udb := &MockDB{m2}
	auth := &Auth{udb, []Access{ADMIN_ACCESS}}
	auth2 := &Auth{udb, []Access{NONE}}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"WIPE DB",
			fields{db, auth},
			false,
		},
		{
			"WIPE DB UNAUTH",
			fields{db, auth2},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &SecureDB{
				db:   tt.fields.db,
				Auth: tt.fields.Auth,
			}
			if err := d.Wipe(); (err != nil) != tt.wantErr {
				t.Errorf("Driver.Wipe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
