package dirp

import "unicode/utf8"

// TokenType は DSL を構成する最小単位の種類。
// Parser はこの種類だけを見て文法を判定する。
type TokenType int

const (
	TokenString TokenType = iota
	TokenLBrace
	TokenRBrace
	TokenLParen
	TokenRParen
	TokenComma
	TokenPipe
	TokenHash
	TokenAt
	TokenEOF
)

type Token struct {
	Type    TokenType
	Literal string
	Pos     int // 元入力文字列におけるバイト位置（エラー表示に使用）
}

// Lexer は UTF-8 文字列を1文字ずつ読み、Token へ分解する。
type Lexer struct {
	input        string
	position     int // 現在の rune が始まるバイト位置
	readPosition int // 次に読むバイト位置
	ch           rune
}

// NewLexer は先頭1文字を読み込んだ状態で初期化する。
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar は UTF-8 を壊さないよう rune 単位で1文字進める。
// EOF 到達時も position を末尾へ合わせ、最後の文字欠落を防ぐ。
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
		l.position = len(l.input)
		return
	}
	r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
	l.ch = r
	l.position = l.readPosition
	l.readPosition += size
}

func (l *Lexer) NextToken() Token {
	var tok Token

	// トークン境界に不要な空白はここで吸収する。
	l.skipWhitespace()

	switch l.ch {
	case '{':
		tok = Token{Type: TokenLBrace, Literal: "{", Pos: l.position}
	case '}':
		tok = Token{Type: TokenRBrace, Literal: "}", Pos: l.position}
	case '(':
		tok = Token{Type: TokenLParen, Literal: "(", Pos: l.position}
	case ')':
		tok = Token{Type: TokenRParen, Literal: ")", Pos: l.position}
	case ',':
		tok = Token{Type: TokenComma, Literal: ",", Pos: l.position}
	case '|':
		tok = Token{Type: TokenPipe, Literal: "|", Pos: l.position}
	case '#':
		tok = Token{Type: TokenHash, Literal: "#", Pos: l.position}
	case '@':
		tok = Token{Type: TokenAt, Literal: "@", Pos: l.position}
	case 0:
		tok = Token{Type: TokenEOF, Literal: "", Pos: l.position}
	default:
		// 予約記号以外は一塊の文字列として読む。
		// 例: "api_#(1,3,1)" は parser 側で必要に応じて再構成する。
		pos := l.position
		literal := l.readString()
		return Token{Type: TokenString, Literal: literal, Pos: pos}
	}

	l.readChar()
	return tok
}

// readString は予約記号または空白までを TokenString として切り出す。
func (l *Lexer) readString() string {
	start := l.position
	for !l.isReserved(l.ch) && l.ch != 0 && !isWhitespace(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// isReserved は文法で特別な意味を持つ文字の判定。
func (l *Lexer) isReserved(ch rune) bool {
	switch ch {
	case '{', '}', '(', ')', ',', '|', '#', '@':
		return true
	default:
		return false
	}
}

// skipWhitespace はトークンとして不要な空白をまとめて飛ばす。
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

// isWhitespace は lexer 内での空白定義を集中管理する。
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
