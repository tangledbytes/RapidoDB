package rql

// ================================ TYPES ================================

// Ast stands for the abstract syntax tree. Here it is implemented as
// a slice of statements instead of a tree
type Ast struct {
	Statements []*Statement
}

// Statement represents the statement structure inside the AST
type Statement struct {
	SetStatement    *SetStatement
	GetStatement    *GetStatement
	DeleteStatement *DeleteStatement
	AuthStatement   *AuthStatement
	WipeStatement   *WipeStatement
	Stype           AstType
}

// SetStatement contains the structure for a "SET" command
type SetStatement struct {
	key    token
	values *[]*expression
}

// GetStatement contains the structure for a "GET" command
type GetStatement struct {
	key token
}

// DeleteStatement contains the structure for a "DELETE" command
type DeleteStatement struct {
	key token
}

// AuthStatement contains the structure for a "AUTH" command
type AuthStatement struct {
	username token
	password token
}

// WipeStatement contains the structure for a "WIPE" command
type WipeStatement struct {
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
	etype   expressionType
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
)
