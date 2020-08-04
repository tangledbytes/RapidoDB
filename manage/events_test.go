package manage

import (
	"reflect"
	"testing"
)

func TestEvents_Set(t *testing.T) {
	type args struct {
		event Event
	}
	tests := []struct {
		name string
		e    Events
		args args
		want Events
	}{
		{
			"ADD AN EVENT WHICH DOESN'T EXISTS",
			Events{1, 2},
			args{3},
			Events{1, 2, 3},
		},
		{
			"ADD AN EVENT WHICH DOES EXISTS",
			Events{1, 2},
			args{1},
			Events{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Set(tt.args.event); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Events.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvents_Exists(t *testing.T) {
	type args struct {
		event Event
	}
	tests := []struct {
		name string
		e    Events
		args args
		want bool
	}{
		{
			"WHEN EVENT EXISTS",
			Events{1, 2, 3},
			args{2},
			true,
		},
		{
			"WHEN EVENT DOESN'T EXISTS",
			Events{1, 2, 3},
			args{4},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Exists(tt.args.event); got != tt.want {
				t.Errorf("Events.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertStringToEvent(t *testing.T) {
	type args struct {
		event string
	}
	tests := []struct {
		name    string
		args    args
		want    Event
		wantErr bool
	}{
		{
			"CONVERT A VALID GET STRING TO EVENT",
			args{"GET"},
			GET,
			false,
		},
		{
			"CONVERT A VALID MIXED CASE GET STRING TO EVENT",
			args{"GeT"},
			GET,
			false,
		},
		{
			"CONVERT AN INVALID STRING TO EVENT",
			args{"GeTi"},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertStringToEvent(tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertStringToEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertStringToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
