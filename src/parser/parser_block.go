package parser

import (
	. "github.com/ASC8384/wen/src/ast"
	. "github.com/ASC8384/wen/src/lexer"
)

// block ::= {stat} [retstat]
func parseBlock(lexer *Lexer) (*Block, error) {
	var stats []Stat
	var err error
	if stats, err = parseStats(lexer); err != nil {
		return nil, err
	}
	return &Block{
		LastLine: lexer.Line(),
		Stats:    stats,
	}, nil
}

// Stat ::= Print | Assignment | Register | If | Loop | Scanf | Pointer | Cell | Init
func parseStats(lexer *Lexer) ([]Stat, error) {
	var statements []Stat
	for !_isReturnOrBlockEnd(lexer.LookAhead()) {
		var statement []Stat
		var err error
		if statement, err = parseStat(lexer); err != nil {
			return nil, err
		}
		statements = append(statements, statement...)
	}
	return statements, nil
}

func _isReturnOrBlockEnd(tokenKind int) bool {
	switch tokenKind {
	case TOKEN_EOF:
		return true
	}
	return false
}
