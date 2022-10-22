package rql

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name    string
		args    args
		want    *Ast
		wantErr bool
	}{
		{
			"AUTH STATEMENT",
			args{`AUTH user pass;`},
			&Ast{
				Statements: []*Statement{
					{
						AuthStatement: &AuthStatement{
							username: "user",
							password: "pass",
						},
						Typ: AuthType,
					},
				},
			},
			false,
		},
		{
			"REGUSER STATEMENT",
			args{`REGUSER user pass 5;`},
			&Ast{
				Statements: []*Statement{
					{
						RegUserStatement: &RegUserStatement{
							username: "user",
							password: "pass",
							access:   5,
						},
						Typ: RegUserType,
					},
				},
			},
			false,
		},
		{
			"REGUSER STATEMENT",
			args{`REGUSER user pass;`},
			&Ast{
				Statements: []*Statement{
					{
						RegUserStatement: &RegUserStatement{
							username: "user",
							password: "pass",
							access:   0,
						},
						Typ: RegUserType,
					},
				},
			},
			false,
		},
		{
			"SET STATEMENT",
			args{`SET data "Hello World";`},
			&Ast{
				Statements: []*Statement{
					{
						SetStatement: &SetStatement{
							key: "data",
							val: "Hello World",
						},
						Typ: SetType,
					},
				},
			},
			false,
		},
		{
			"SET STATEMENT WITH EXPIRY",
			args{`SET data "Hello World" 234;`},
			&Ast{
				Statements: []*Statement{
					{
						SetStatement: &SetStatement{
							key: "data",
							val: "Hello World",
							exp: 234,
						},
						Typ: SetType,
					},
				},
			},
			false,
		},
		{
			"MULTI SET STATEMENTS",
			args{`SET data "Hello World" 234; SET data1 3454 565;`},
			&Ast{
				Statements: []*Statement{
					{
						SetStatement: &SetStatement{
							key: "data",
							val: "Hello World",
							exp: 234,
						},
						Typ: SetType,
					},
					{
						SetStatement: &SetStatement{
							key: "data1",
							val: "3454",
							exp: 565,
						},
						Typ: SetType,
					},
				},
			},
			false,
		},
		{
			"GET STATEMENT",
			args{`GET data data1 data2 data3;`},
			&Ast{
				Statements: []*Statement{
					{
						GetStatement: &GetStatement{
							keys: []string{"data", "data1", "data2", "data3"},
						},
						Typ: GetType,
					},
				},
			},
			false,
		},
		{
			"DELETE STATEMENT",
			args{`DEL data data1 data2 data3;`},
			&Ast{
				Statements: []*Statement{
					{
						DeleteStatement: &DeleteStatement{
							keys: []string{"data", "data1", "data2", "data3"},
						},
						Typ: DeleteType,
					},
				},
			},
			false,
		},
		{
			"WIPE STATEMENT",
			args{`WIPE;`},
			&Ast{
				Statements: []*Statement{
					{
						WipeStatement: &WipeStatement{},
						Typ:           WipeType,
					},
				},
			},
			false,
		},
		{
			"VALID PING STATEMENT",
			args{`PING ON GET;`},
			&Ast{
				Statements: []*Statement{
					{
						PingStatement: &PingStatement{
							operation: "get",
							on:        true,
						},
						Typ: PingType,
					},
				},
			},
			false,
		},
		{
			"VALID UNPING STATEMENT",
			args{`PING OFF GET;`},
			&Ast{
				Statements: []*Statement{
					{
						PingStatement: &PingStatement{
							operation: "get",
							on:        false,
						},
						Typ: PingType,
					},
				},
			},
			false,
		},
		{
			"MIX STATEMENTS",
			args{`SET data "Hello World"; GET data data1 data2 data3; DEL data data1 data2 data3; GET data; DEL data;`},
			&Ast{
				Statements: []*Statement{
					{
						SetStatement: &SetStatement{
							key: "data",
							val: "Hello World",
						},
						Typ: SetType,
					},
					{
						GetStatement: &GetStatement{
							keys: []string{"data", "data1", "data2", "data3"},
						},
						Typ: GetType,
					},
					{
						DeleteStatement: &DeleteStatement{
							keys: []string{"data", "data1", "data2", "data3"},
						},
						Typ: DeleteType,
					},
					{
						GetStatement: &GetStatement{
							keys: []string{"data"},
						},
						Typ: GetType,
					},
					{
						DeleteStatement: &DeleteStatement{
							keys: []string{"data"},
						},
						Typ: DeleteType,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
