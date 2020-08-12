package manage

import (
	"reflect"
	"testing"
)

func TestToDBUser(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want DBUser
	}{
		{
			"CONVERT A VALID INTERFACE TO DBUSER",
			args{map[string]interface{}{
				"Username": "utkarsh",
				"Password": "test",
				"Access":   float64(2),
				"Events":   []interface{}{uint(2), uint(3)},
			}},
			DBUser{"utkarsh", "test", WriteAccess, Events{2, 3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToDBUser(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDBUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
