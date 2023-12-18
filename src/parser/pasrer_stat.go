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
	case TOKEN_INC_PTR:
		return parseIncPtr(lexer)
	case TOKEN_DEC_PTR:
		return parseDecPtr(lexer)
	default:
		return nil, errors.New("parseStat(): unknown Stat." + TokenNameMap[lexer.LookAhead()])
	}
}

// Print ::= "print" "(" Ignored Variable Ignored ")" Ignored
func parsePrint(lexer *Lexer) (*Print, error) {
	var print Print
	var err error

	print.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_PRINT)
	if TOKEN_LEFT_PAREN == lexer.LookAhead() {
		lexer.NextTokenIs(TOKEN_LEFT_PAREN)
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
		if print.Variable, err = parseVariable(lexer); err != nil {
			return nil, err
		}
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
		lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
	} else {
		print.Variable = nil
	}
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

func parseIncPtr(lexer *Lexer) (*PointerStat, error) {
	var PointerStat PointerStat

	PointerStat.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_INC_PTR)
	PointerStat.Pointer = 1
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &PointerStat, nil
}

func parseDecPtr(lexer *Lexer) (*PointerStat, error) {
	var PointerStat PointerStat

	PointerStat.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_DEC_PTR)
	PointerStat.Pointer = -1
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &PointerStat, nil
}
