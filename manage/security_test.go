package manage

import (
	"reflect"
	"testing"
	"time"
)

func TestSecureDB_Set(t *testing.T) {
	type fields struct {
		ust          UnsecureStore
		userdb       *UserDB
		activeClient *DBClient
	}
	type args struct {
		key      string
		data     interface{}
		expireIn time.Duration
	}

	db := &MockDB{make(map[string]interface{})}
	udb := &MockDB{make(map[string]interface{})}
	ac := newDBClient("admin", "pass", AdminAccess)
	ac2 := newDBClient("test", "test", WriteAccess)
	ac3 := newDBClient("test2", "test", ReadAccess)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"SET DATA WITH ADMIN ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac},
			args{"k1", 123, db.DefaultExpiry()},
			false,
		},
		{
			"SET DATA WITH WRITE ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac2},
			args{"k2", 123, db.DefaultExpiry()},
			false,
		},
		{
			"SET DATA WITH READ ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac3},
			args{"k3", 123, db.DefaultExpiry()},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdb := &SecureDB{
				ust:          tt.fields.ust,
				userdb:       tt.fields.userdb,
				activeClient: tt.fields.activeClient,
			}
			if err := sdb.Set(tt.args.key, tt.args.data, tt.args.expireIn); (err != nil) != tt.wantErr {
				t.Errorf("SecureDB.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecureDB_Get(t *testing.T) {
	type fields struct {
		ust          UnsecureStore
		userdb       *UserDB
		activeClient *DBClient
	}
	type args struct {
		key string
	}

	db := &MockDB{make(map[string]interface{})}
	udb := &MockDB{make(map[string]interface{})}
	ac := newDBClient("admin", "pass", AdminAccess)
	ac2 := newDBClient("test", "test", ReadAccess)
	ac3 := newDBClient("test2", "test", NONE)

	// Set data
	db.Set("k1", 1234, db.DefaultExpiry())
	db.Set("k2", "t1", db.DefaultExpiry())
	db.Set("k3", 123456, db.DefaultExpiry())

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		want1   bool
		wantErr bool
	}{
		{
			"GET DATA THAT EXISTS WITH ADMIN ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac},
			args{"k1"},
			1234,
			true,
			false,
		},
		{
			"GET DATA THAT DOES NOT EXISTS WITH ADMIN ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac},
			args{"k4"},
			nil,
			false,
			false,
		},
		{
			"GET DATA WITH READ ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac2},
			args{"k2"},
			"t1",
			true,
			false,
		},
		{
			"GET DATA THAT DOES NOT EXISTS WITH READ ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac2},
			args{"k4"},
			nil,
			false,
			false,
		},
		{
			"GET DATA WITH NONE ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac3},
			args{"k3"},
			nil,
			false,
			true,
		},
		{
			"GET DATA THAT DOES NOT EXISTS WITH NONE ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac3},
			args{"k4"},
			nil,
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdb := &SecureDB{
				ust:          tt.fields.ust,
				userdb:       tt.fields.userdb,
				activeClient: tt.fields.activeClient,
			}
			got, got1, err := sdb.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("SecureDB.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecureDB.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SecureDB.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSecureDB_Delete(t *testing.T) {
	type fields struct {
		ust          UnsecureStore
		userdb       *UserDB
		activeClient *DBClient
	}
	type args struct {
		key string
	}

	db := &MockDB{make(map[string]interface{})}
	udb := &MockDB{make(map[string]interface{})}
	ac := newDBClient("admin", "pass", AdminAccess)
	ac2 := newDBClient("test", "test", WriteAccess)
	ac3 := newDBClient("test2", "test", NONE)

	// Set data
	db.Set("k1", 1234, db.DefaultExpiry())
	db.Set("k2", "t1", db.DefaultExpiry())
	db.Set("k3", 123456, db.DefaultExpiry())

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		want1   bool
		wantErr bool
	}{
		{
			"DELETE DATA THAT EXISTS WITH ADMIN ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac},
			args{"k1"},
			1234,
			true,
			false,
		},
		{
			"DELETE DATA THAT DOES NOT EXISTS WITH ADMIN ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac},
			args{"k4"},
			nil,
			false,
			false,
		},
		{
			"DELETE DATA WITH WRITE ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac2},
			args{"k2"},
			"t1",
			true,
			false,
		},
		{
			"DELETE DATA THAT DOES NOT EXISTS WITH READ ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac2},
			args{"k4"},
			nil,
			false,
			false,
		},
		{
			"DELETE DATA WITH NONE ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac3},
			args{"k3"},
			nil,
			false,
			true,
		},
		{
			"DELETE DATA THAT DOES NOT EXISTS WITH NONE ACCESS LEVEL",
			fields{db, &UserDB{udb}, ac3},
			args{"k4"},
			nil,
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdb := &SecureDB{
				ust:          tt.fields.ust,
				userdb:       tt.fields.userdb,
				activeClient: tt.fields.activeClient,
			}
			got, got1, err := sdb.Delete(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("SecureDB.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecureDB.Delete() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SecureDB.Delete() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSecureDB_Wipe(t *testing.T) {
	type fields struct {
		ust          UnsecureStore
		userdb       *UserDB
		activeClient *DBClient
	}

	db := &MockDB{make(map[string]interface{})}
	udb := &MockDB{make(map[string]interface{})}
	ac := newDBClient("admin", "pass", AdminAccess)
	ac2 := newDBClient("test", "test", WipeAccess)
	ac3 := newDBClient("test2", "test", NONE)

	// Set data
	db.Set("k1", 1234, db.DefaultExpiry())
	db.Set("k2", "t1", db.DefaultExpiry())
	db.Set("k3", 123456, db.DefaultExpiry())

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"WIPE DATA WITH ADMIN ACCESS",
			fields{db, &UserDB{udb}, ac},
			false,
		},
		{
			"WIPE DATA WITH WIPE ACCESS",
			fields{db, &UserDB{udb}, ac2},
			false,
		},
		{
			"WIPE DATA WITH NONE ACCESS",
			fields{db, &UserDB{udb}, ac3},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdb := &SecureDB{
				ust:          tt.fields.ust,
				userdb:       tt.fields.userdb,
				activeClient: tt.fields.activeClient,
			}
			if err := sdb.Wipe(); (err != nil) != tt.wantErr {
				t.Errorf("SecureDB.Wipe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecureDB_RegisterUser(t *testing.T) {
	type fields struct {
		ust          UnsecureStore
		userdb       *UserDB
		activeClient *DBClient
	}
	type args struct {
		username string
		password string
		access   uint
	}

	db := &MockDB{make(map[string]interface{})}
	udb := &MockDB{make(map[string]interface{})}
	ac := newDBClient("admin", "pass", AdminAccess)
	ac2 := newDBClient("test", "test", ModifyUserAccess)
	ac3 := newDBClient("test2", "test", NONE)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"ADD A USER WITH VALID ACCESS USING ADMIN ACCESS",
			fields{db, &UserDB{udb}, ac},
			args{"utkarsh", "test", 5},
			false,
		},
		{
			"ADD A USER WITH INVALID ACCESS USING ADMIN ACCESS",
			fields{db, &UserDB{udb}, ac},
			args{"utkarsh", "test", 50},
			true,
		},
		{
			"ADD A USER WITH VALID ACCESS USING MODIFY USER ACCESS",
			fields{db, &UserDB{udb}, ac2},
			args{"utkarsh", "test", 5},
			false,
		},
		{
			"ADD A USER WITH VALID ACCESS USING MODIFY USER ACCESS",
			fields{db, &UserDB{udb}, ac},
			args{"utkarsh", "test", 500},
			true,
		},
		{
			"ADD A USER WITH VALID ACCESS USING NONE ACCESS",
			fields{db, &UserDB{udb}, ac3},
			args{"utkarsh", "test", 5},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdb := &SecureDB{
				ust:          tt.fields.ust,
				userdb:       tt.fields.userdb,
				activeClient: tt.fields.activeClient,
			}
			if err := sdb.RegisterUser(tt.args.username, tt.args.password, tt.args.access); (err != nil) != tt.wantErr {
				t.Errorf("SecureDB.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecureDB_Authenticate(t *testing.T) {
	type fields struct {
		ust          UnsecureStore
		userdb       *UserDB
		activeClient *DBClient
	}
	type args struct {
		username string
		password string
	}

	db := &MockDB{make(map[string]interface{})}
	udb := &MockDB{make(map[string]interface{})}
	ac := newDBClient("test", "test", NONE) // Simulate the behaviour of a normal client

	// Add new users to the database
	udb.Set("test2", NewDBUser("test2", "test2", AdminAccess), udb.DefaultExpiry())

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"AUTHENTICATE USER WITH VALID USERNAME AND PASSWORD",
			fields{db, &UserDB{udb}, ac},
			args{"test2", "test2"},
			false,
		},
		{
			"AUTHENTICATE USER WITH VALID USERNAME AND INVALID PASSWORD",
			fields{db, &UserDB{udb}, ac},
			args{"test2", "test21"},
			true,
		},
		{
			"AUTHENTICATE USER WITH INVALID USERNAME AND PASSWORD",
			fields{db, &UserDB{udb}, ac},
			args{"test21", "test21"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdb := &SecureDB{
				ust:          tt.fields.ust,
				userdb:       tt.fields.userdb,
				activeClient: tt.fields.activeClient,
			}
			if err := sdb.Authenticate(tt.args.username, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("SecureDB.Authenticate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecureDB_Authorize(t *testing.T) {
	type fields struct {
		ust          UnsecureStore
		userdb       *UserDB
		activeClient *DBClient
	}
	type args struct {
		reqAccess Access
	}

	db := &MockDB{make(map[string]interface{})}
	udb := &MockDB{make(map[string]interface{})}
	ac := newDBClient("test", "test", ModifyUserAccess)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"REQUEST ADMIN ACCESS WITH MODIFY USER ACCESS",
			fields{db, &UserDB{udb}, ac},
			args{AdminAccess},
			false,
		},
		{
			"REQUEST WIPE WITH MODIFY USER ACCESS",
			fields{db, &UserDB{udb}, ac},
			args{WipeAccess},
			false,
		},
		{
			"REQUEST WRITE ACCESS WITH MODIFY USER ACCESS",
			fields{db, &UserDB{udb}, ac},
			args{WriteAccess},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdb := &SecureDB{
				ust:          tt.fields.ust,
				userdb:       tt.fields.userdb,
				activeClient: tt.fields.activeClient,
			}
			if got := sdb.Authorize(tt.args.reqAccess); got != tt.want {
				t.Errorf("SecureDB.Authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}
