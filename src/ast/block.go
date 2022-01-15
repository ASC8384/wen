package ast

// block           ::= stat+
// retstat         ::= return exp+

// SourceCode      ::= block+

type Block struct {
	LastLine int
	Stats    []Stat
	// RetExps  []Exp
}
