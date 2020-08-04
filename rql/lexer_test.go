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
			"SET NUMBER",
			args{"SET data 1"},
			[]*token{
				{"set", keywordType, location{0, 0}},
				{"data", identifierType, location{0, 4}},
				{"1", numericType, location{0, 9}},
			},
			false,
		},
		{
			"SET STRING",
			args{`SET data "Hello World"`},
			[]*token{
				{"set", keywordType, location{0, 0}},
				{"data", identifierType, location{0, 4}},
				{"Hello World", stringType, location{0, 9}},
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
		{
			"WIPE",
			args{`WIPE`},
			[]*token{
				{"wipe", keywordType, location{0, 0}},
			},
			false,
		},
		{
			"REGUSER",
			args{`REGUSER utkarsh safepass 5`},
			[]*token{
				{"reguser", keywordType, location{0, 0}},
				{"utkarsh", identifierType, location{0, 8}},
				{"safepass", identifierType, location{0, 16}},
				{"5", numericType, location{0, 25}},
			},
			false,
		},
		{
			"AUTH",
			args{`AUTH utkarsh safepass`},
			[]*token{
				{"auth", keywordType, location{0, 0}},
				{"utkarsh", identifierType, location{0, 5}},
				{"safepass", identifierType, location{0, 13}},
			},
			false,
		},
		{
			"PING ON",
			args{`PING ON GET`},
			[]*token{
				{"ping", keywordType, location{0, 0}},
				{"on", keywordType, location{0, 5}},
				{"get", keywordType, location{0, 8}},
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
