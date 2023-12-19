package parser

import (
	"errors"

	. "github.com/ASC8384/wen/src/ast"
	. "github.com/ASC8384/wen/src/lexer"
)

var _statEmpty = &EmptyStat{}

func parseStat(lexer *Lexer) ([]Stat, error) {
	lexer.LookAheadAndSkip(TOKEN_IGNORED) // skip if source code start with ignored token
	var stats []Stat
	var stat Stat
	var err error
	// var loops LoopStat

	switch lexer.LookAhead() {
	case TOKEN_PRINT:
		stat, err = parsePrint(lexer)
	case TOKEN_SCANF:
		stat, err = parseScanf(lexer)
	case TOKEN_VAR_PREFIX:
		stat, err = parseAssignStat(lexer)
	case TOKEN_INC_PTR:
		stat, err = parseIncPtr(lexer)
	case TOKEN_DEC_PTR:
		stat, err = parseDecPtr(lexer)
	case TOKEN_INC_MEM:
		stat, err = parseIncMem(lexer)
	case TOKEN_DEC_MEM:
		stat, err = parseDecMem(lexer)
	case TOKEN_LOOP_OPEN:
		stat, err = parseLoopOpen(lexer)
	case TOKEN_LOOP_CLOSE:
		stats, err = paeseLoopClose(lexer)
	default:
		stats, err = nil, errors.New("parseStat(): unknown Stat."+TokenNameMap[lexer.LookAhead()])
	}
	if stat != nil {
		stats = append(stats, stat)
	}
	// if loops.Stats != nil {
	// 	stats = append(stats, loops)
	// }
	return stats, err
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

func parseScanf(lexer *Lexer) (*Scanf, error) {
	var scanf Scanf
	var err error

	scanf.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_SCANF)
	if TOKEN_LEFT_PAREN == lexer.LookAhead() {
		lexer.NextTokenIs(TOKEN_LEFT_PAREN)
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
		if scanf.Variable, err = parseVariable(lexer); err != nil {
			return nil, err
		}
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
		lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
	} else {
		scanf.Variable = nil
	}
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &scanf, nil
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

func parseIncMem(lexer *Lexer) (*CellStat, error) {
	var CellStat CellStat
	CellStat.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_INC_MEM)
	CellStat.Cell = 1
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &CellStat, nil
}

func parseDecMem(lexer *Lexer) (*CellStat, error) {
	var CellStat CellStat
	CellStat.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_DEC_MEM)
	CellStat.Cell = -1
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &CellStat, nil
}

func parseLoopOpen(lexer *Lexer) (*LoopStat, error) {
	var body LoopStat
	var err error

	lexer.NextTokenIs(TOKEN_LOOP_OPEN)
	// Parse statements inside the loop until we encounter TOKEN_LOOP_CLOSE
	for lexer.LookAhead() != TOKEN_LOOP_CLOSE {
		if lexer.LookAhead() == TOKEN_EOF {
			return nil, errors.New("parseLoopOpen(): loop never closes")
		}

		var stat []Stat
		stat, err = parseStat(lexer)
		if err != nil {
			return nil, err
		}
		body.Stats = append(body.Stats, stat...)
	}

	// Ensure we consume the ']' TOKEN_LOOP_CLOSE
	if lexer.LookAhead() == TOKEN_LOOP_CLOSE {
		lexer.NextToken()
	} else {
		return nil, errors.New("parseLoopOpen(): expected ']' at the end of the loop")
	}

	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &body, nil
}

func paeseLoopClose(lexer *Lexer) ([]Stat, error) {
	return nil, errors.New("parseLoopClose(): unexpected ']' without matching '['")
}
