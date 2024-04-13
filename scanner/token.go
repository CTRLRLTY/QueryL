package scanner

type TokenType int

const (
	TokenAnd TokenType = iota
	TokenOr
	TokenNot
	TokenEqual
	TokenNotEqual
	TokenGreater
	TokenLesser
	TokenGreaterEqual
	TokenLesserEqual
	TokenError
	TokenField
	TokenEof
	TokenNumber
	TokenString
)

type Token struct {
	Code   TokenType
	Lexeme []rune
	Offset uint32
}
