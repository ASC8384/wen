package parser

import (
	"errors"

	. "github.com/ASC8384/wen/src/ast"
	. "github.com/ASC8384/wen/src/lexer"
)

// Name ::= [_A-Za-z][_0-9A-Za-z]*
func parseName(lexer *Lexer) (string, error) {
	_, name := lexer.NextTokenIs(TOKEN_NAME)
	return name, nil
}

// LiteralString   ::= '"' '"' Ignored | '"' StringCharacter '"' Ignored
func parseString(lexer *Lexer) (string, error) {
	str := ""
	switch lexer.LookAhead() {
	case TOKEN_DUOQUOTE:
		lexer.NextTokenIs(TOKEN_DUOQUOTE)
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
		return str, nil
	case TOKEN_QUOTE:
		lexer.NextTokenIs(TOKEN_QUOTE)
		str = lexer.ScanBeforeToken(TokenNameMap[TOKEN_QUOTE])
		lexer.NextTokenIs(TOKEN_QUOTE)
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
		return str, nil
	default:
		return "", errors.New("parseString(): not a string.")
	}
}

// Variable ::= "$" Name Ignored
func parseVariable(lexer *Lexer) (*Variable, error) {
	var variable Variable
	var err error

	variable.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_VAR_PREFIX)
	if variable.Name, err = parseName(lexer); err != nil {
		return nil, err
	}
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &variable, nil
}
