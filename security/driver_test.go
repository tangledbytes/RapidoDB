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

func TestDriver_Set(t *testing.T) {
	type fields struct {
		db   UnsecureDB
		Auth *Auth
	}
	type args struct {
		key      string
		data     interface{}
		expireIn time.Duration
	}

	m := make(map[string]interface{})
	db := &MockDB{m}
	auth := &Auth{"admin", "pass", []Access{ADMIN_ACCESS}, true}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"SET STRING",
			fields{db, auth},
			args{
				"d1",
				"Hello World",
				time.Duration(100),
			},
		},
		{
			"SET NUMBERS",
			fields{db, auth},
			args{
				"d1",
				100,
				time.Duration(100),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Driver{
				db:   tt.fields.db,
				Auth: tt.fields.Auth,
			}
			d.Set(tt.args.key, tt.args.data, tt.args.expireIn)
		})
	}
}

func TestDriver_Get(t *testing.T) {
	type fields struct {
		db   UnsecureDB
		Auth *Auth
	}
	type args struct {
		key string
	}

	m := make(map[string]interface{})
	db := &MockDB{m}
	auth := &Auth{"admin", "pass", []Access{ADMIN_ACCESS}, true}

	// Add a key to the map -> key = 'd1'
	db.Set("d1", "Hello World", 0)

	// Add a key to the map -> key = 'd3'
	db.Set("d3", 2345, 0)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  bool
	}{
		{
			"GET STRING WHEN KEY IS PRESENT",
			fields{db, auth},
			args{
				"d1",
			},
			"Hello World",
			true,
		},
		{
			"GET NIL WHEN KEY IS ABSENT",
			fields{db, auth},
			args{
				"d2",
			},
			nil,
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Driver{
				db:   tt.fields.db,
				Auth: tt.fields.Auth,
			}
			got, got1 := d.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Driver.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Driver.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDriver_Delete(t *testing.T) {
	type fields struct {
		db   UnsecureDB
		Auth *Auth
	}
	type args struct {
		key string
	}

	m := make(map[string]interface{})
	db := &MockDB{m}
	auth := &Auth{"admin", "pass", []Access{ADMIN_ACCESS}, true}

	// Add a key to the map -> key = 'd1'
	db.Set("d1", "Hello World", 0)

	// Add a key to the map -> key = 'd3'
	db.Set("d3", 2345, 0)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  bool
	}{
		{
			"DELETE STRING WHEN KEY IS PRESENT",
			fields{db, auth},
			args{
				"d1",
			},
			"Hello World",
			true,
		},
		{
			"GET NIL WHEN KEY IS ABSENT",
			fields{db, auth},
			args{
				"d2",
			},
			nil,
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Driver{
				db:   tt.fields.db,
				Auth: tt.fields.Auth,
			}
			got, got1 := d.Delete(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Driver.Delete() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Driver.Delete() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
