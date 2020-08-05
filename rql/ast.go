package rql

import "fmt"

// ================================ TYPES ================================

// Ast stands for the abstract syntax tree. Here it is implemented as
// a slice of statements instead of a tree
type Ast struct {
	Statements []*Statement
}

// Statement represents the statement structure inside the AST
type Statement struct {
	SetStatement     *SetStatement
	GetStatement     *GetStatement
	DeleteStatement  *DeleteStatement
	AuthStatement    *AuthStatement
	WipeStatement    *WipeStatement
	RegUserStatement *RegUserStatement
	PingStatement    *PingStatement
	Typ              AstType
}

// SetStatement contains the structure for a "SET" command
type SetStatement struct {
	key string
	val interface{}
	exp uint
}

// GetStatement contains the structure for a "GET" command
type GetStatement struct {
	keys []string
}

// DeleteStatement contains the structure for a "DEL" command
type DeleteStatement struct {
	keys []string
}

// AuthStatement contains the structure for a "AUTH" command
type AuthStatement struct {
	username string
	password string
}

// WipeStatement contains the structure for a "WIPE" command
type WipeStatement struct {
}

// RegUserStatement contains the structure for a "REGUSER" command
type RegUserStatement struct {
	username string
	password string
	access   uint
}

// PingStatement contains the structure for a "PING ON" command
type PingStatement struct {
	operation string
}

// AstType represents the type of abstract syntax tree
type AstType uint

type binaryExpression struct {
	a  expression
	b  expression
	op token
}

type expression struct {
	literal *token
	binary  *binaryExpression
	typ     expressionType
}

type expressionType uint

// ================================ CONSTANTS ================================

const (
	literalType expressionType = iota
	binaryType
)

// Supported AST type
const (
	SetType AstType = iota
	GetType
	DeleteType
	AuthType
	WipeType
	RegUserType
	PingType
)

// ===========================================================================

func (a Ast) String() string {
	s := "[ "
	for _, stmt := range a.Statements {
		if stmt.SetStatement != nil {
			s += fmt.Sprintf("%+v", stmt.SetStatement)
		}
		if stmt.GetStatement != nil {
			s += fmt.Sprintf("%+v", stmt.GetStatement)
		}
		if stmt.DeleteStatement != nil {
			s += fmt.Sprintf("%+v", stmt.DeleteStatement)
		}
		if stmt.WipeStatement != nil {
			s += fmt.Sprintf("%+v", stmt.WipeStatement)
		}
		if stmt.AuthStatement != nil {
			s += fmt.Sprintf("%+v", stmt.AuthStatement)
		}
		if stmt.RegUserStatement != nil {
			s += fmt.Sprintf("%+v", stmt.RegUserStatement)
		}
		if stmt.PingStatement != nil {
			s += fmt.Sprintf("%+v", stmt.PingStatement)
		}
	}

	return s + " ]"
}
