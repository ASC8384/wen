package parser

import (
	"errors"

	. "github.com/ASC8384/wen/src/ast"
	. "github.com/ASC8384/wen/src/lexer"
)

// var _statEmpty = &EmptyStat{}

// parse statement
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
	case TOKEN_INIT_DATA:
		stat, err = parseInitData(lexer)
	case TOKEN_REG:
		stat, err = parseReg(lexer)
	case TOKEN_IF_START:
		stat, err = parseIf(lexer)
	case TOKEN_IF_RUSH:
		stat, err = parseIf(lexer)
	case TOKEN_IF_END:
		stat, err = parseIf(lexer)
	case TOKEN_END:
		stat, err = parseIf(lexer)

	default:
		stats, err = nil, errors.New("parseStat(): unknown Stat."+TokenNameMap[lexer.LookAhead()]+" line: "+string(rune(lexer.Line())))
	}
	if stat != nil {
		stats = append(stats, stat)
	}

	return stats, err
}

// Init ::= "&" SourceCharacter
func parseInitData(lexer *Lexer) (*StringExp, error) {
	var InitData StringExp
	// var err error

	InitData.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_INIT_DATA)

	InitData.Str += lexer.GetChunk()

	// lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &InitData, nil
}

// Print ::= "o" "@"? "(" Ignored Variable Ignored ")" Ignored | "o" "@"? Ignored
func parsePrint(lexer *Lexer) (*Print, error) {
	var print Print
	var err error

	print.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_PRINT)
	if TOKEN_VAR_INT == lexer.LookAhead() {
		lexer.NextTokenIs(TOKEN_VAR_INT)
		print.Int = true
	} else {
		print.Int = false
	}
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

// Scanf ::= "v" "@"? "(" Ignored Variable Ignored ")" Ignored | "v" "@"? Ignored
func parseScanf(lexer *Lexer) (*Scanf, error) {
	var scanf Scanf
	var err error

	scanf.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_SCANF)
	if TOKEN_VAR_INT == lexer.LookAhead() {
		lexer.NextTokenIs(TOKEN_VAR_INT)
		scanf.Int = true
	} else {
		scanf.Int = false
	}
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

// Assignment      ::= Variable Ignored "=" "@"? Ignored LiteralString Ignored | Variable "(" Ignored LiteralNumber Ignored ")" Ignored | Variable "[" Ignored LiteralNumber Ignored "]" Ignored
func parseAssignStat(lexer *Lexer) (*AssignStat, error) {
	var assignment AssignStat
	var err error

	assignment.Line = lexer.Line()
	if assignment.Variable, err = parseVariable(lexer); err != nil {
		return nil, err
	}
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	if TOKEN_EQUAL == lexer.LookAhead() {
		lexer.NextTokenIs(TOKEN_EQUAL)
		if TOKEN_VAR_INT == lexer.LookAhead() {
			lexer.NextTokenIs(TOKEN_VAR_INT)
			assignment.Int = true
		} else {
			assignment.Int = false
		}
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
		if assignment.String, err = parseString(lexer); err != nil {
			return nil, err
		}
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
	} else if TOKEN_LEFT_PAREN == lexer.LookAhead() {
		lexer.NextTokenIs(TOKEN_LEFT_PAREN)
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
		if assignment.Length, err = parseInt(lexer); err != nil {
			return nil, err
		}
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
		lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
	} else if TOKEN_VAR_LEFT == lexer.LookAhead() {
		lexer.NextTokenIs(TOKEN_VAR_LEFT)
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
		if assignment.Start, err = parseInt(lexer); err != nil {
			return nil, err
		}
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
		lexer.NextTokenIs(TOKEN_VAR_RIGHT)
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
	} else if TOKEN_VAR_START == lexer.LookAhead() {
		lexer.NextTokenIs(TOKEN_VAR_START)
		lexer.LookAheadAndSkip(TOKEN_IGNORED)
		assignment.Start = 0
		assignment.Length = 0
	}
	return &assignment, nil
}

// Pointer ::= "i" Ignored | "a" Ignored
func parseIncPtr(lexer *Lexer) (*PointerStat, error) {
	var PointerStat PointerStat

	PointerStat.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_INC_PTR)
	PointerStat.Pointer = 1
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &PointerStat, nil
}

// Pointer ::= "i" Ignored | "a" Ignored
func parseDecPtr(lexer *Lexer) (*PointerStat, error) {
	var PointerStat PointerStat

	PointerStat.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_DEC_PTR)
	PointerStat.Pointer = -1
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &PointerStat, nil
}

// Cell ::= "l" Ignored | "e" Ignored
func parseIncMem(lexer *Lexer) (*CellStat, error) {
	var CellStat CellStat
	CellStat.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_INC_MEM)
	CellStat.Cell = 1
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &CellStat, nil
}

// Cell ::= "l" Ignored | "e" Ignored
func parseDecMem(lexer *Lexer) (*CellStat, error) {
	var CellStat CellStat
	CellStat.Line = lexer.Line()
	lexer.NextTokenIs(TOKEN_DEC_MEM)
	CellStat.Cell = -1
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &CellStat, nil
}

// Loop ::= "b" Ignored Stat* "u" Ignored
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

// Loop ::= "b" Ignored Stat* "u" Ignored
func paeseLoopClose(lexer *Lexer) ([]Stat, error) {
	return nil, errors.New("parseLoopClose(): unexpected ']' without matching '['")
}

// Reg ::= "A" [+-*/%] Ignored
func parseReg(lexer *Lexer) (*RegStat, error) {
	var reg RegStat
	// var err error

	lexer.NextTokenIs(TOKEN_REG)
	switch lexer.LookAhead() {
	case TOKEN_REG_STORE:
		lexer.NextTokenIs(TOKEN_REG_STORE)
		reg.Reg = TOKEN_REG_STORE
	case TOKEN_REG_PLUS:
		lexer.NextTokenIs(TOKEN_REG_PLUS)
		reg.Reg = TOKEN_REG_PLUS
	case TOKEN_REG_MINUS:
		lexer.NextTokenIs(TOKEN_REG_MINUS)
		reg.Reg = TOKEN_REG_MINUS
	case TOKEN_REG_MUL:
		lexer.NextTokenIs(TOKEN_REG_MUL)
		reg.Reg = TOKEN_REG_MUL
	case TOKEN_REG_DIV:
		lexer.NextTokenIs(TOKEN_REG_DIV)
		reg.Reg = TOKEN_REG_DIV
	case TOKEN_REG_MOD:
		lexer.NextTokenIs(TOKEN_REG_MOD)
		reg.Reg = TOKEN_REG_MOD
	case TOKEN_REG_READ:
		lexer.NextTokenIs(TOKEN_REG_READ)
		reg.Reg = TOKEN_REG_READ
	default:
		return nil, errors.New("parseReg(): unknown Reg." + TokenNameMap[lexer.LookAhead()])
	}
	return &reg, nil
}

// If ::= "B" | "V" | "#"
func parseIf(lexer *Lexer) (*IfStat, error) {
	var ifstat IfStat
	ifstat.Line = lexer.Line()
	switch lexer.LookAhead() {
	case TOKEN_IF_START:
		lexer.NextTokenIs(TOKEN_IF_START)
		ifstat.Type = TOKEN_IF_START
	case TOKEN_IF_RUSH:
		lexer.NextTokenIs(TOKEN_IF_RUSH)
		ifstat.Type = TOKEN_IF_RUSH
	case TOKEN_IF_END:
		lexer.NextTokenIs(TOKEN_IF_END)
		ifstat.Type = TOKEN_IF_END
	case TOKEN_END:
		lexer.NextTokenIs(TOKEN_END)
		ifstat.Type = TOKEN_END
	}
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &ifstat, nil
}
