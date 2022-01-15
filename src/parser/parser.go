package parser

import (
	. "github.com/ASC8384/wen/src/ast"
	. "github.com/ASC8384/wen/src/lexer"
)

/* recursive descent parser */
func Parse(chunk, chunkName string) (*Block, error) {
	var block *Block
	var err error
	lexer := NewLexer(chunk, chunkName)
	if block, err = parseBlock(lexer); err != nil {
		return nil, err
	}
	lexer.NextTokenIs(TOKEN_EOF)
	return block, err
}
