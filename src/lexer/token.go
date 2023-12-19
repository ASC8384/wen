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
	TOKEN_INC_PTR            // >
	TOKEN_DEC_PTR            // <
	TOKEN_INC_MEM            // +
	TOKEN_DEC_MEM            // -
	TOKEN_SCANF              // ,
	TOKEN_LOOP_OPEN          // [
	TOKEN_LOOP_CLOSE         // ]
	TOKEN_INIT_DATA          // &
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
	TOKEN_PRINT:       "o",
	TOKEN_INC_PTR:     "i",
	TOKEN_DEC_PTR:     "a",
	TOKEN_INC_MEM:     "l",
	TOKEN_DEC_MEM:     "e",
	TOKEN_LOOP_OPEN:   "b",
	TOKEN_LOOP_CLOSE:  "u",
	TOKEN_SCANF:       "v",
	TOKEN_INIT_DATA:   "&",
}

var keywords = map[string]int{
	"o": TOKEN_PRINT,
	"i": TOKEN_INC_PTR,
	"a": TOKEN_DEC_PTR,
	"l": TOKEN_INC_MEM,
	"e": TOKEN_DEC_MEM,
	"b": TOKEN_LOOP_OPEN,
	"u": TOKEN_LOOP_CLOSE,
	"v": TOKEN_SCANF,
	"&": TOKEN_INIT_DATA,
	// "if": TOKEN_KW_IF,
}
