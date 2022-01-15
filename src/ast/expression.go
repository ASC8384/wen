package ast

/*
exp ::=  nil | false | true | Numeral | exp binop exp | unop exp | LiteralString | Variable
*/

type Exp interface{}

// LiteralString
type StringExp struct {
	Line int
	Str  string
}

type Variable struct {
	Line int
	Name string
}
