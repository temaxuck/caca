package lexer

import "unicode"

type TokenKind int

const (
	TK_Digit TokenKind = iota
	TK_Comment
	TK_NewLine
	TK_Invalid
	TK_EOF
)

var TokenKindRepr = map[TokenKind]string{
	TK_Digit:   "Digit",
	TK_Comment: "Comment",
	TK_NewLine: "New Line",
	TK_Invalid: "Invalid",
	TK_EOF:     "EOF",
}

type Token struct {
	Str  string
	Kind TokenKind
}

type Lexer struct {
	Input string
	Pos   int

	curLine         int
	curLineStartPos int

	Token *Token
}

func (l *Lexer) NextToken() bool {
	l.skipWS()

	if l.Pos >= len(l.Input) {
		l.Token = &Token{"", TK_EOF}
		return false
	}

	if unicode.IsDigit(rune(l.Input[l.Pos])) {
		l.Token = &Token{l.Input[l.Pos : l.Pos+1], TK_Digit}
		l.advancePos()
		return true
	}

	if l.expectRune('\n') {
		l.Token = &Token{l.Input[l.Pos : l.Pos+1], TK_NewLine}
		l.advancePos()
		return true
	}

	if l.expectRune('#') {
		startPos := l.Pos
		l.skipToEOL()
		l.Token = &Token{l.Input[startPos:l.Pos], TK_Comment}
		return true
	}

	startPos := l.Pos
	for l.Pos < len(l.Input) {
		if unicode.IsSpace(rune(l.Input[l.Pos])) {
			break
		}
		l.Pos++
	}
	l.Token = &Token{l.Input[startPos:l.Pos], TK_Invalid}
	return false
}

func (l *Lexer) GetCursorFilePos() (line int, col int) {
	line = l.curLine
	col = l.Pos - l.curLineStartPos
	return line, col
}

func (l *Lexer) skipToEOL() {
	for l.Pos < len(l.Input) && l.Input[l.Pos] != '\n' {
		l.Pos++
	}
}

// Returns true if current cursor is a line break, false otherwise
func (l *Lexer) skipWS() {
	// I like it C-way
	for l.isSpaceButNotNewLine() && l.advancePos() {
	}
}

func (l *Lexer) advancePos() bool {
	if l.Pos >= len(l.Input) {
		return false
	}

	if l.expectRune('\n') {
		l.curLine++
		l.curLineStartPos = l.Pos + 1
	}

	l.Pos++
	return true
}

func (l *Lexer) isSpaceButNotNewLine() bool {
	return l.Pos < len(l.Input) && !l.expectRune('\n') && unicode.IsSpace(rune(l.Input[l.Pos]))
}

func (l *Lexer) expectRune(r rune) bool {
	return l.Pos < len(l.Input) && rune(l.Input[l.Pos]) == r
}
