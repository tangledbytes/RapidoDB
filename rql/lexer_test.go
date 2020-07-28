package rql

import (
	"reflect"
	"testing"
)

func Test_lex(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		args    args
		want    []*token
		wantErr bool
	}{
		{
			"SET ANY",
			args{"SET data any 1"},
			[]*token{
				{"set", keywordType, location{0, 0}},
				{"data", identifierType, location{0, 4}},
				{"any", keywordType, location{0, 9}},
				{"1", numericType, location{0, 13}},
			},
			false,
		},
		{
			"SET STRING",
			args{`SET data string "Hello World"`},
			[]*token{
				{"set", keywordType, location{0, 0}},
				{"data", identifierType, location{0, 4}},
				{"string", keywordType, location{0, 9}},
				{"Hello World", stringType, location{0, 16}},
			},
			false,
		},
		{
			"GET",
			args{`GET data`},
			[]*token{
				{"get", keywordType, location{0, 0}},
				{"data", identifierType, location{0, 4}},
			},
			false,
		},
		{
			"DELETE",
			args{`DEL data`},
			[]*token{
				{"del", keywordType, location{0, 0}},
				{"data", identifierType, location{0, 4}},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lex(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("lex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lex() = %v, want %v", got, tt.want)
			}
		})
	}
}
