package rql

import (
	"fmt"
	"strings"
)

// ========================= TYPES =============================

// keyword represents the "keywords" of RQL
type keyword string

// symbol represents the "symbols" of RQL
type symbol string

// tokenType describes the type of the token
type tokenType uint

// lexer represents lexing functions which are used internally
// to analyse different tokens
type lexer func(string, cursor) (*token, cursor, bool)

// location describes the location of a token
// in the source code (RQL)
type location struct {
	line uint
	col  uint
}

type token struct {
	val string
	typ tokenType
	loc location
}

type cursor struct {
	ptr uint
	loc location
}

// ===================== CONSTANTS ============================

// RQL Keyword
const (
	// Commands
	authKeyword    keyword = "auth"
	getKeyword     keyword = "get"
	setKeyword     keyword = "set"
	delKeyword     keyword = "del"
	wipeKeyword    keyword = "wipe"
	reguserKeyword keyword = "reguser"
	pingKeyword    keyword = "ping"
	onKeyword      keyword = "on"
	offKeyword     keyword = "off"
	// Data types
	// numberKeyword keyword = "number"
	// stringKeyword keyword = "string"
	// boolKeyword   keyword = "bool"
	// jsonKeyword   keyword = "json"
	// anyKeyword    keyword = "any"

	// Conditionals
	ifKeyword  keyword = "if"
	andKeyword keyword = "and"
	orKeyword  keyword = "or"

	// Meta
	expireinKeyword keyword = "expirein"
)

// RQL Symbol
const (
	semicolonSymbol  symbol = ";"
	asteriskSymbol   symbol = "*"
	commaSymbol      symbol = ","
	leftParenSymbol  symbol = "("
	rightParenSymbol symbol = ")"
	eqSymbol         symbol = "=="
	neqSymbol        symbol = "!="
	ltSymbol         symbol = "<"
	lteSymbol        symbol = "<="
	gtSymbol         symbol = ">"
	gteSymbol        symbol = ">="
)

// RQL Token type
const (
	keywordType tokenType = iota
	symbolType
	identifierType
	stringType
	numericType
	boolType
)

// ============================================================

// lex splits an input string into a list of tokens. This
// process can be divided into following tasks:
//
// 1. Instantiating a cursor with pointing to the start of the string
//
// 2. Execute all the lexers in the series
//
// 3. If any of the lexer generates a token then add the token to the
// token slice, update the cursor and restart the process from the new
// cursor location
func lex(src string) ([]*token, error) {
	var tokens []*token
	cur := cursor{}

lex:
	for cur.ptr < uint(len(src)) {
		lexers := []lexer{lexKeyword, lexSymbol, lexString, lexNumeric, lexIdentifier}
		for _, l := range lexers {
			if token, newCursor, ok := l(src, cur); ok {
				cur = newCursor

				if token != nil {
					tokens = append(tokens, token)
				}

				continue lex
			}
		}

		// hint describes the hint to be printed
		// in case the lexing process fails at a point
		hint := ""

		if len(tokens) > 0 {
			// Store the last token from the tokens array
			hint = "after " + tokens[len(tokens)-1].val
		}

		// Print all the tokens extracted uptil that point
		for _, t := range tokens {
			fmt.Println(t)
		}

		return nil, fmt.Errorf("Unable to lex token %s, at %d %d", hint, cur.loc.line, cur.loc.col)
	}

	return tokens, nil
}

// lexKeyword analysis the source code for keywords
func lexKeyword(source string, ic cursor) (*token, cursor, bool) {
	cur := ic
	keywords := []keyword{
		// Commands
		authKeyword,
		getKeyword,
		setKeyword,
		delKeyword,
		wipeKeyword,
		reguserKeyword,
		pingKeyword,
		onKeyword,
		offKeyword,
		// Data types
		// numberKeyword,
		// stringKeyword,
		// boolKeyword,
		// jsonKeyword,
		// anyKeyword,

		// Conditionals
		ifKeyword,
		andKeyword,
		orKeyword,

		// Meta
		expireinKeyword,
	}

	var options []string
	for _, k := range keywords {
		options = append(options, string(k))
	}

	match := longestMatch(source, ic, options)
	if match == "" {
		return nil, ic, false
	}

	cur.ptr = ic.ptr + uint(len(match))
	cur.loc.col = ic.loc.col + uint(len(match))

	ttype := keywordType

	return &token{
		val: match,
		typ: ttype,
		loc: ic.loc,
	}, cur, true
}

// lexSymbol analysis the symbols in the source code
//
// lexSymbol will eat up the white spaces in the source string
func lexSymbol(source string, ic cursor) (*token, cursor, bool) {
	c := source[ic.ptr]
	cur := ic
	// Will get overwritten later if not an ignored symbol
	// but if the symbol is to be ignored then this increment
	// will account for the escaped/ignored symbol
	cur.ptr++
	cur.loc.col++

	switch c {
	// Symbols that should be thrown away
	case '\n':
		cur.loc.line++
		cur.loc.col = 0
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		return nil, cur, true
	}

	// Symbols that should be kept
	symbols := []symbol{
		commaSymbol,
		leftParenSymbol,
		rightParenSymbol,
		semicolonSymbol,
		eqSymbol,
		neqSymbol,
		ltSymbol,
		lteSymbol,
		gtSymbol,
		gteSymbol,
		asteriskSymbol,
	}

	var options []string
	for _, s := range symbols {
		options = append(options, string(s))
	}

	// Use `ic`, not `cur`
	match := longestMatch(source, ic, options)
	// Unknown character
	if match == "" {
		return nil, ic, false
	}

	cur.ptr = ic.ptr + uint(len(match))
	cur.loc.col = ic.loc.col + uint(len(match))

	return &token{
		val: match,
		loc: ic.loc,
		typ: symbolType,
	}, cur, true
}

// lexString analysis strings in the source code. It internally
// uses the lexCharacterDelimited function to do so
func lexString(src string, ic cursor) (*token, cursor, bool) {
	return lexCharacterDelimited(src, ic, '"')
}

// lexNumeric analysis the numbers in the source code
func lexNumeric(src string, ic cursor) (*token, cursor, bool) {
	cur := ic

	// Notes if a period has been found
	periodFound := false
	// Notes if an exponent has been found
	expMarkerFound := false

	for ; cur.ptr < uint(len(src)); cur.ptr++ {
		ch := src[cur.ptr]
		// Move ahead in the source code
		cur.loc.col++

		isDigit := ch >= '0' && ch <= '9'
		isPeriod := ch == '.'
		isExpMarker := ch == 'e'

		// A number must start with a digit or period
		if cur.ptr == ic.ptr {
			if !isDigit && !isPeriod {
				return nil, ic, false
			}

			periodFound = isPeriod
			continue
		}

		if isPeriod {
			// If this is also true then it means that 2 consecutive
			// periods were found and hence it's not a number
			if periodFound {
				return nil, ic, false
			}

			periodFound = isPeriod
			continue
		}

		if isExpMarker {
			// If this is also true then it means that 2 consecutive
			// exp markers were found and hence it's not a number
			if expMarkerFound {
				return nil, ic, false
			}

			expMarkerFound = isExpMarker

			// No periods should be after the marker hence
			periodFound = true

			// Marker must be followed by digits
			if cur.ptr == uint(len(src)-1) {
				return nil, ic, false
			}

			chNext := src[cur.ptr+1]
			if chNext == '-' || chNext == '+' {
				cur.ptr++
				cur.loc.col++
			}

			continue
		}

		if !isDigit {
			break
		}
	}

	// No characters accumulated
	if cur.ptr == ic.ptr {
		return nil, ic, false
	}

	return &token{
		val: src[ic.ptr:cur.ptr],
		loc: ic.loc,
		typ: numericType,
	}, cur, true
}

// lexIdentifier analysis the source code for identifiers
func lexIdentifier(source string, ic cursor) (*token, cursor, bool) {
	cur := ic

	c := source[cur.ptr]
	// Other characters count too, big ignoring non-ascii for now
	isAlphabetical := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
	if !isAlphabetical {
		return nil, ic, false
	}
	cur.ptr++
	cur.loc.col++

	value := []byte{c}
	for ; cur.ptr < uint(len(source)); cur.ptr++ {
		c = source[cur.ptr]

		// Other characters count too, big ignoring non-ascii for now
		isAlphabetical := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
		isNumeric := c >= '0' && c <= '9'
		if isAlphabetical || isNumeric || c == '$' || c == '_' {
			value = append(value, c)
			cur.loc.col++
			continue
		}

		break
	}

	if len(value) == 0 {
		return nil, ic, false
	}

	return &token{
		// all identifiers are case sensitive
		val: string(value),
		loc: ic.loc,
		typ: identifierType,
	}, cur, true
}

// lexCharacterDelimited analysis the source code for string with custom delimiter
func lexCharacterDelimited(src string, ic cursor, delimiter byte) (*token, cursor, bool) {
	cur := ic

	if len(src[cur.ptr:]) == 0 {
		return nil, ic, false
	}

	if src[cur.ptr] != delimiter {
		return nil, ic, false
	}

	cur.loc.col++
	cur.ptr++

	var val []byte

	for ; cur.ptr < uint(len(src)); cur.ptr++ {
		ch := src[cur.ptr]

		if ch == delimiter {
			// Unlike SQL, RQL will escape characters with backslash
			if cur.ptr+1 >= uint(len(src)) || src[cur.ptr+1] != '\\' {
				cur.ptr++
				cur.loc.col++
				return &token{
					val: string(val),
					loc: ic.loc,
					typ: stringType,
				}, cur, true
			}

			val = append(val, ch)
			cur.ptr++
			cur.loc.col++
		}

		val = append(val, ch)
		cur.loc.col++
	}

	return nil, ic, false
}

// longestMatch iterates through a source string starting at the given
// cursor to find the longest matching substring among the provided
// options
func longestMatch(src string, ic cursor, opts []string) string {
	var val []byte
	var skipList []int
	var match string

	cur := ic

	for cur.ptr < uint(len(src)) {

		val = append(val, strings.ToLower(string(src[cur.ptr]))...)
		cur.ptr++

	match:
		for i, opt := range opts {
			for _, skip := range skipList {
				if i == skip {
					continue match
				}
			}

			// Deal with cases like INT vs INTO
			if opt == string(val) {
				skipList = append(skipList, i)
				if len(opt) > len(match) {
					match = opt
				}

				continue
			}

			sharesPrefix := string(val) == opt[:cur.ptr-ic.ptr]
			tooLong := len(val) > len(opt)
			if tooLong || !sharesPrefix {
				skipList = append(skipList, i)
			}
		}

		if len(skipList) == len(opts) {
			break
		}
	}

	return match
}

// equals returns true if passed token has the same value
// and token type as of the current token
func (t *token) equals(other *token) bool {
	return t.val == other.val && t.typ == other.typ
}

// ================ HELPER FUNCTIONS (ONLY FOR DEBUGGING) ======================

func (t token) String() string {
	return fmt.Sprintf("val: %s, type: %d, location: %v\n", t.val, t.typ, t.loc)
}

func (l location) String() string {
	return fmt.Sprintf("[ line: %d, col: %d ]", l.line, l.col)
}
