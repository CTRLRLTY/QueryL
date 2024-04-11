package scanner

import (
	"math"
	"unicode"
)

type Scanner struct {
	data   string
	lexeme []rune
	cursor int
	offset uint32
}

func (s *Scanner) advance() (c rune) {
	if s.IsAtEnd() {
		return
	}

	s.cursor++
	s.offset++
	c = rune(s.data[s.cursor-1])

	if len(s.lexeme) < math.MaxInt {
		s.lexeme = append(s.lexeme, c)
	}

	return
}

func (s *Scanner) peek() (c rune) {
	if s.IsAtEnd() {
		return
	}

	c = rune(s.data[s.cursor])

	return
}

func (s *Scanner) peekNext() (c rune) {
	if s.IsAtEnd() {
		return
	}

	c = rune(s.data[s.cursor+1])

	return
}

func numberToken(s *Scanner) Token {
	isDigit := func(x rune) bool {
		return x > 47 && x < 58
	}

	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	return makeToken(s, TokenNumber)
}

func stringToken(s *Scanner) Token {
	for s.peek() != '"' && !s.IsAtEnd() {
		s.advance()
	}

	if s.IsAtEnd() {
		return makeToken(s, TokenError)
	}

	s.advance() // consume closing "
	return makeToken(s, TokenString)
}

func fieldToken(s *Scanner) Token {
	for !s.IsAtEnd() && !unicode.IsSpace(s.peek()) {
		s.advance()
	}

	return makeToken(s, TokenField)
}

func makeToken(s *Scanner, t TokenType) Token {
	return Token{Code: t, Lexeme: s.lexeme, Offset: s.offset - 1}
}

func (s *Scanner) IsAtEnd() bool {
	return s.cursor >= len(s.data)
}

func (s *Scanner) ScanToken() (tkn Token) {
	s.lexeme = nil
	c := s.advance()

	switch c {
	case '&':
		if s.advance() == '&' {
			return makeToken(s, TokenAnd)
		}
	case '=':
		if s.advance() == '=' {
			return makeToken(s, TokenEqual)
		}
	case '!':
		if s.advance() == '=' {
			return makeToken(s, TokenNotEqual)
		} else {
			return makeToken(s, TokenNot)
		}
	case '|':
		if s.advance() == '|' {
			return makeToken(s, TokenOr)
		}

	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		return numberToken(s)

	case '"':
		return stringToken(s)

	// Ignore limited white space
	case ' ', '\t', '\r', '\n':
		return s.ScanToken()

	case rune(0):
		return makeToken(s, TokenEof)

	default:
		if unicode.IsLetter(c) {
			return fieldToken(s)
		}
	}

	return makeToken(s, TokenError)
}

func FromString(str string) Scanner {
	return Scanner{data: str, cursor: 0, offset: 0}
}
