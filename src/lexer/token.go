package lexer

// token const
const (
	TOKEN_EOF         = iota // end-of-file
	TOKEN_VAR_PREFIX         // $
	TOKEN_LEFT_PAREN         // (
	TOKEN_RIGHT_PAREN        // )
	TOKEN_EQUAL              // =
	TOKEN_QUOTE              // "
	TOKEN_DUOQUOTE           // ""
	TOKEN_NAME               // Name ::= [_A-Za-z][_0-9A-Za-z]*
	TOKEN_IGNORED            // Ignored
	TOKEN_IDENTIFIER         // identifier
	TOKEN_PRINT              // print
	// TOKEN_KW_IF           // if
)

var TokenNameMap = map[int]string{
	TOKEN_EOF:         "EOF",
	TOKEN_VAR_PREFIX:  "$",
	TOKEN_LEFT_PAREN:  "(",
	TOKEN_RIGHT_PAREN: ")",
	TOKEN_EQUAL:       "=",
	TOKEN_QUOTE:       "\"",
	TOKEN_DUOQUOTE:    "\"\"",
	TOKEN_NAME:        "Name",
	TOKEN_IGNORED:     "Ignored",
	TOKEN_PRINT:       "print",
}

var keywords = map[string]int{
	"print": TOKEN_PRINT,
	// "if": TOKEN_KW_IF,
}
