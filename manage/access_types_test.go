package manage

import "testing"

func TestConvertUintToAccess(t *testing.T) {
	type args struct {
		access uint
	}
	tests := []struct {
		name    string
		args    args
		want    Access
		wantErr bool
	}{
		{
			"CONVERT A MAX VALID UINT TO ACCESS",
			args{5},
			AdminAccess,
			false,
		},
		{
			"CONVERT A MIN VALID UINT TO ACCESS",
			args{0},
			NONE,
			false,
		},
		{
			"CONVERT AN INVALID UINT TO ACCESS",
			args{50},
			NONE,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertUintToAccess(tt.args.access)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertUintToAccess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertUintToAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
