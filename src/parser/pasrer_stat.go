package parser

import (
	"errors"

	. "github.com/ASC8384/wen/src/ast"
	. "github.com/ASC8384/wen/src/lexer"
)

var _statEmpty = &EmptyStat{}

func parseStat(lexer *Lexer) (Stat, error) {
	lexer.LookAheadAndSkip(TOKEN_IGNORED) // skip if source code start with ignored token
	switch lexer.LookAhead() {
	case TOKEN_PRINT:
		return parsePrint(lexer)
	case TOKEN_VAR_PREFIX:
		return parseAssignStat(lexer)
	default:
		return nil, errors.New("parseStat(): unknown Stat.")
	}
}

// Print ::= "print" "(" Ignored Variable Ignored ")" Ignored
func parsePrint(lexer *Lexer) (*Print, error) {
	var print Print
	var err error

	print.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_PRINT)
	lexer.NextTokenIs(TOKEN_LEFT_PAREN)
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	if print.Variable, err = parseVariable(lexer); err != nil {
		return nil, err
	}
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &print, nil
}

// Assignment      ::= Variable Ignored "=" Ignored LiteralString Ignored
func parseAssignStat(lexer *Lexer) (*AssignStat, error) {
	var assignment AssignStat
	var err error

	assignment.Line = lexer.Line()
	if assignment.Variable, err = parseVariable(lexer); err != nil {
		return nil, err
	}
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	lexer.NextTokenIs(TOKEN_EQUAL)
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	if assignment.String, err = parseString(lexer); err != nil {
		return nil, err
	}
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &assignment, nil
}
