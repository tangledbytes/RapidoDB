package rql

import (
	"fmt"
	"strings"
)

// location describes the location of a token
// in the source code (RQL)
type location struct {
	line uint
	col  uint
}

// keyword represents the "keywords" of RQL
type keyword string

const (
	getKeyword      keyword = "get"
	setKeyword      keyword = "set"
	delKeyword      keyword = "del"
	ifKeyword       keyword = "if"
	andKeyword      keyword = "and"
	orKeyword       keyword = "or"
	expireinKeyword keyword = "expirein"
	rfieldKeyword   keyword = "r_field"
)

// symbol represents the "symbols" of RQL
type symbol string

const (
	semicolonSymbol  symbol = ";"
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

// tokenType describes the type of the token
type tokenType uint

const (
	keywordType tokenType = iota
	symbolType
	identifierType
	stringType
	numericType
)

type token struct {
	val   string
	ttype tokenType
	loc   location
}

type cursor struct {
	ptr uint
	loc location
}

func (t *token) equals(other *token) bool {
	return t.val == other.val && t.ttype == other.ttype
}

type lexer func(string, cursor) (*token, cursor, bool)

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
		val:   src[ic.ptr:cur.ptr],
		loc:   ic.loc,
		ttype: numericType,
	}, cur, true
}

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
				return &token{
					val:   string(val),
					loc:   ic.loc,
					ttype: stringType,
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

func lexString(src string, ic cursor) (*token, cursor, bool) {
	return lexCharacterDelimited(src, ic, '"')
}

func lexSymbol(source string, ic cursor) (*token, cursor, bool) {
	c := source[ic.ptr]
	cur := ic
	// Will get overwritten later if not an ignored symbol
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
		val:   match,
		loc:   ic.loc,
		ttype: symbolType,
	}, cur, true
}

func lexKeyword(source string, ic cursor) (*token, cursor, bool) {
	cur := ic
	keywords := []keyword{
		getKeyword,
		setKeyword,
		delKeyword,
		ifKeyword,
		andKeyword,
		orKeyword,
		expireinKeyword,
		rfieldKeyword,
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
		val:   match,
		ttype: ttype,
		loc:   ic.loc,
	}, cur, true
}

// longestMatch iterates through a source string starting at the given
// cursor to find the longest matching substring among the provided
// options
func longestMatch(source string, ic cursor, options []string) string {
	var value []byte
	var skipList []int
	var match string

	cur := ic

	for cur.ptr < uint(len(source)) {

		value = append(value, strings.ToLower(string(source[cur.ptr]))...)
		cur.ptr++

	match:
		for i, option := range options {
			for _, skip := range skipList {
				if i == skip {
					continue match
				}
			}

			// Deal with cases like INT vs INTO
			if option == string(value) {
				skipList = append(skipList, i)
				if len(option) > len(match) {
					match = option
				}

				continue
			}

			sharesPrefix := string(value) == option[:cur.ptr-ic.ptr]
			tooLong := len(value) > len(option)
			if tooLong || !sharesPrefix {
				skipList = append(skipList, i)
			}
		}

		if len(skipList) == len(options) {
			break
		}
	}

	return match
}

func lexIdentifier(source string, ic cursor) (*token, cursor, bool) {
	// Handle separately if is a double-quoted identifier
	if token, newCursor, ok := lexCharacterDelimited(source, ic, '"'); ok {
		return token, newCursor, true
	}

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
		// Unquoted dentifiers are case-insensitive
		val:   strings.ToLower(string(value)),
		loc:   ic.loc,
		ttype: identifierType,
	}, cur, true
}

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
			fmt.Println(t.val)
		}

		return nil, fmt.Errorf("Unable to lex token %s, at %d %d", hint, cur.loc.line, cur.loc.col)
	}

	return tokens, nil
}
