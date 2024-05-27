package parser

import (
	"fmt"

	"github.com/CTRLRLTY/QueryL/chunk"
	"github.com/CTRLRLTY/QueryL/scanner"
)

type Parser struct {
	Current  scanner.Token
	Previous scanner.Token
	rules    []ParseRule
	chunk    chunk.Chunk
}

type ParseFunc func(*Parser, *scanner.Scanner) error

type Precedence int

type ParseRule struct {
	Prefix     *ParseFunc
	Infix      *ParseFunc
	Precedence Precedence
}

const (
	PrecNone       Precedence = iota
	PrecOr                    // ||
	PrecAnd                   // &&
	PrecEquality              // == !=
	PrecComparison            // < > <= >=
)

func (c *Parser) advance(s *scanner.Scanner) error {
	c.Previous = c.Current

	if c.Current.Code != scanner.TokenEof {
		c.Current = s.ScanToken()

		if c.Current.Code == scanner.TokenError {
			return fmt.Errorf("error token %s at %d", string(c.Current.Lexeme), c.Current.Offset)
		}
	}

	return nil
}

func (c *Parser) GetRule(tkn scanner.Token) *ParseRule {
	if int(tkn.Code) >= len(c.rules) {
		return nil
	}

	return &c.rules[tkn.Code]
}

func (c *Parser) parsePrecedence(s *scanner.Scanner, precedence Precedence) error {
	if err := c.advance(s); err != nil {
		return err
	}

	prefixFunc := *c.GetRule(c.Previous).Prefix

	if prefixFunc == nil {
		return fmt.Errorf("token(%v) rule not found", c.Previous.Code)
	}

	if err := prefixFunc(c, s); err != nil {
		return err
	}

	// parse next token if the current rule precedence
	// is lower than the next token's rule precedence.
	for precedence <= c.GetRule(c.Current).Precedence {
		if err := c.advance(s); err != nil {
			return err
		}

		infixFunc := *c.GetRule(c.Previous).Infix

		if err := infixFunc(c, s); err != nil {
			return err
		}
	}

	return nil
}

func (c *Parser) Parse(s *scanner.Scanner) (cnk chunk.Chunk, err error) {
	// forward the compiler so it moves the current token to previous
	if err = c.advance(s); err != nil {
		return
	}

	for c.Current.Code != scanner.TokenEof {
		// Parse expression
		if err = c.parsePrecedence(s, 1); err != nil {
			return
		}
	}

	// Consume the Eof token
	if err = c.advance(s); err != nil {
		return
	}

	return c.chunk, nil
}

func (c *Parser) ParseString(str string) (cnk chunk.Chunk, err error) {
	s := scanner.FromString(str)

	return c.Parse(&s)
}

func (c *Parser) Init() {
	c.rules = defaultRules[:]
}
