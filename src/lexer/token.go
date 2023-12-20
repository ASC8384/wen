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
	TOKEN_VAR_LEFT           // 【
	TOKEN_VAR_RIGHT          // 】
	TOKEN_VAR_START          // var start
	TOKEN_NUMBER             // number
	TOKEN_VAR_INT            // not acill
	TOKEN_REG                //register
	TOKEN_REG_STORE          // store register
	TOKEN_REG_PLUS           // plus register
	TOKEN_REG_MINUS
	TOKEN_REG_MUL
	TOKEN_REG_DIV
	TOKEN_REG_MOD
	TOKEN_REG_READ
	TOKEN_IF_START
	TOKEN_IF_RUSH
	TOKEN_IF_END
	TOKEN_END
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
	TOKEN_VAR_LEFT:    "[",
	TOKEN_VAR_RIGHT:   "]",
	TOKEN_VAR_START:   "^",
	TOKEN_VAR_INT:     "@",
	TOKEN_REG:         "A",
	TOKEN_REG_STORE:   "!",
	TOKEN_REG_PLUS:    "+",
	TOKEN_REG_MINUS:   "-",
	TOKEN_REG_MUL:     "*",
	TOKEN_REG_DIV:     "/",
	TOKEN_REG_MOD:     "%",
	TOKEN_REG_READ:    "?",
	TOKEN_IF_START:    "B",
	TOKEN_IF_RUSH:     "R",
	TOKEN_IF_END:      "#",
	TOKEN_END:         "E",
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
	"A": TOKEN_REG,
	"E": TOKEN_END,
	// "&": TOKEN_INIT_DATA,
	// "if": TOKEN_KW_IF,
}
