package lexer

import (
	"fmt"
	"regexp"
	"strings"
)

var reIdentifier = regexp.MustCompile(`^[_\d\w]+`)
var reNewLine = regexp.MustCompile("\r\n|\n\r|\n|\r")

// var reNumber = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([+\-]?[0-9]+)?|^[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)
var reOpeningLongBracket = regexp.MustCompile(`^\[=*\[`)

type Lexer struct {
	chunk         string // 源代码
	chunkName     string // 源文件名
	line          int    // 当前行号
	nextToken     string // 下一个 token
	nextTokenKind int    // 下一个 token 的类型
	nextTokenLine int    // 下一个 token 的行号
}

// 创建 Lexer 结构体实例，设初始行号为 1
func NewLexer(chunk, chunkName string) *Lexer {
	return &Lexer{
		chunk:         chunk,
		chunkName:     chunk,
		line:          1,
		nextToken:     "",
		nextTokenKind: 0,
		nextTokenLine: 0,
	}
}

// 抛出错误信息
func (lexer *Lexer) error(f string, a ...interface{}) {
	err := fmt.Sprintf(f, a...)
	err = fmt.Sprintf("%s:%d: %s", lexer.chunkName, lexer.line, err)
	panic(err)
}

// 查看下一个 token 的类型，不跳过
func (lexer *Lexer) LookAhead() int {
	if lexer.nextTokenLine > 0 { // 缓存中有下一个 token 的信息
		return lexer.nextTokenKind // 直接返回
	}
	currentLine := lexer.line              // 保存当前行号
	line, kind, token := lexer.NextToken() // 提取下一个 token
	lexer.line = currentLine
	lexer.nextTokenLine = line
	lexer.nextTokenKind = kind
	lexer.nextToken = token
	return kind
}

// 查看下一个 token 的类型，并且跳过
func (lexer *Lexer) LookAheadAndSkip(expectedType int) {
	// get next token
	nowLine := lexer.line
	line, kind, token := lexer.NextToken()
	// not is expected type, reverse cursor
	if kind != expectedType {
		lexer.line = nowLine
		lexer.nextTokenLine = line
		lexer.nextTokenKind = kind
		lexer.nextToken = token
	}
}

// 选择性跳过空白符、注释，返回下一个 token 的行号和类型
func (lexer *Lexer) NextToken() (line, kind int, token string) {
	if lexer.nextTokenLine > 0 { // 缓存中有下一个 token 的信息，直接返回
		line = lexer.nextTokenLine
		kind = lexer.nextTokenKind
		token = lexer.nextToken
		lexer.line = lexer.nextTokenLine
		lexer.nextTokenLine = 0
		return
	}
	// 否则，开始匹配
	lexer.skipWhitespace()     // 跳过空白符与注释
	if len(lexer.chunk) == 0 { // 返回 EOF
		return lexer.line, TOKEN_EOF, "EOF"
	}
	switch lexer.chunk[0] {
	case '$':
		lexer.next(1)
		return lexer.line, TOKEN_VAR_PREFIX, "$"
	case '(':
		lexer.next(1)
		return lexer.line, TOKEN_LEFT_PAREN, "("
	case ')':
		lexer.next(1)
		return lexer.line, TOKEN_RIGHT_PAREN, ")"
	case '=':
		lexer.next(1)
		return lexer.line, TOKEN_EQUAL, "="
	case '"':
		if lexer.nextCodeIs("\"\"") {
			lexer.next(2)
			return lexer.line, TOKEN_DUOQUOTE, "\"\""
		}
		lexer.next(1)
		return lexer.line, TOKEN_QUOTE, "\""
	case 'i':
		lexer.next(1)
		return lexer.line, TOKEN_INC_PTR, "i"
	case 'a':
		lexer.next(1)
		return lexer.line, TOKEN_DEC_PTR, "a"
	case 'l':
		lexer.next(1)
		return lexer.line, TOKEN_INC_MEM, "l"
	case 'e':
		lexer.next(1)
		return lexer.line, TOKEN_DEC_MEM, "e"
	case 'b':
		lexer.next(1)
		return lexer.line, TOKEN_LOOP_OPEN, "b"
	case 'u':
		lexer.next(1)
		return lexer.line, TOKEN_LOOP_CLOSE, "u"
	case 'v':
		lexer.next(1)
		return lexer.line, TOKEN_SCANF, "v"
	case 'o':
		lexer.next(1)
		return lexer.line, TOKEN_PRINT, "o"
	case '&':
		lexer.next(1)
		return lexer.line, TOKEN_INIT_DATA, "&"
	}
	c := lexer.chunk[0]
	// check multiple character token
	if c == '_' || isLetter(c) {
		token := lexer.scanIdentifier()
		if tokenType, isMatch := keywords[token]; isMatch {
			return lexer.line, tokenType, token
		} else {
			return lexer.line, TOKEN_NAME, token
		}
	}
	lexer.error("unexpected symbol near %q", c)
	return
}

// 跳过空白符、注释
func (lexer *Lexer) skipWhitespace() {
	isWhiteSpace := func(c byte) bool { // 空白
		switch c {
		case '\t', '\n', '\v', '\f', '\r', ' ':
			return true
		}
		return false
	}
	for len(lexer.chunk) > 0 {
		if lexer.nextCodeIs("--") { // 123
			lexer.skipComment()
		} else if lexer.nextCodeIs("\r\n") || lexer.nextCodeIs("\n\r") {
			lexer.next(2)
			lexer.line += 1
		} else if isNewLine(lexer.chunk[0]) {
			lexer.next(1)
			lexer.line += 1
		} else if isWhiteSpace(lexer.chunk[0]) {
			lexer.next(1)
		} else {
			break
		}
	}
}

// 字符串开头
func (lexer *Lexer) nextCodeIs(s string) bool {
	return strings.HasPrefix(lexer.chunk, s)
}

// 跳过 n 个字符
func (lexer *Lexer) next(n int) {
	lexer.chunk = lexer.chunk[n:]
}

// 回车或换行
func isNewLine(c byte) bool {
	return c == '\r' || c == '\n'
}

// 跳过注释
func (lexer *Lexer) skipComment() {
	lexer.next(2) // skip --
	// 块注释 123
	/*
		-- 这是行注释
		--[[
		这是块注释，
		块注释可以
		注释多行！
		--]]
	*/
	if lexer.nextCodeIs("[") {
		if reOpeningLongBracket.FindString(lexer.chunk) != "" {
			lexer.scanLongString()
			return
		}
	}
	// 行注释
	for len(lexer.chunk) > 0 && !isNewLine(lexer.chunk[0]) {
		lexer.next(1)
	}
}

// 长字符串字面量
func (lexer *Lexer) scanLongString() string {
	openingLongBracket := reOpeningLongBracket.FindString(lexer.chunk)
	if openingLongBracket == "" { // 找不到匹配的右边部分
		lexer.error("invalid long string delimiter near '%s'",
			lexer.chunk[0:2])
	}

	closingLongBracket := strings.Replace(openingLongBracket, "[", "]", -1)
	closingLongBracketIdx := strings.Index(lexer.chunk, closingLongBracket)
	if closingLongBracketIdx < 0 {
		lexer.error("unfinished long string or comment")
	}

	str := lexer.chunk[len(openingLongBracket):closingLongBracketIdx]
	lexer.next(closingLongBracketIdx + len(closingLongBracket))

	str = reNewLine.ReplaceAllString(str, "\n")
	lexer.line += strings.Count(str, "\n")
	if len(str) > 0 && str[0] == '\n' {
		str = str[1:]
	}

	return str
}

func (lexer *Lexer) scanIdentifier() string {
	return lexer.scan(reIdentifier)
}

// func (lexer *Lexer) scanNumber() string {
// 	return lexer.scan(reNumber)
// }

func (lexer *Lexer) scan(re *regexp.Regexp) string {
	if token := re.FindString(lexer.chunk); token != "" {
		lexer.next(len(token))
		return token
	}
	panic("unreachable!")
}

// 提取标识符
func (lexer *Lexer) NextIdentifier() (line int, token string) {
	return lexer.NextTokenIs(TOKEN_IDENTIFIER)
}

// 提取制定类型的 token
func (lexer *Lexer) NextTokenIs(kind int) (line int, token string) {
	line, _kind, token := lexer.NextToken()
	if kind != _kind {
		lexer.error("syntax error near '%s'", token)
	}
	return line, token
}

// 返回当前行号
func (lexer *Lexer) Line() int {
	return lexer.line
}

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// 返回 token 之前的内容
func (lexer *Lexer) ScanBeforeToken(token string) string {
	s := strings.Split(lexer.chunk, token)
	if len(s) < 2 {
		panic("unreachable!")
	}
	lexer.next(len(s[0]))
	return s[0]
}

// 返回 lexer.chunk[0] 的内容，并跳过1个字符
func (lexer *Lexer) NextChar() byte {
	c := lexer.chunk[0]
	lexer.next(1)
	return c
}

// 返回 lexer.chunk 的内容，并跳过所有
func (lexer *Lexer) GetChunk() string {
	ret := lexer.chunk
	lexer.next(len(lexer.chunk))
	return ret
}
