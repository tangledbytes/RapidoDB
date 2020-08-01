package main

import (
	"os"
	"testing"
)

func Test_getEnv(t *testing.T) {
	type args struct {
		env      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"ENV IS NOT EMPTY",
			args{"USER", "root"},
			os.Getenv("USER"),
		},
		{
			"ENV IS EMPTY",
			args{"RANDOM_XYZ", "0000"},
			"0000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnv(tt.args.env, tt.args.fallback); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
