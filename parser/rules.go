package parser

import (
	"fmt"
	"strconv"

	"github.com/CTRLRLTY/QueryL/chunk"
	"github.com/CTRLRLTY/QueryL/scanner"
)

var defaultRules = [...]ParseRule{
	scanner.TokenField:  {parseField, parseNothing, PrecNone},
	scanner.TokenString: {parseString, parseNothing, PrecNone},
	scanner.TokenNumber: {parseNumber, parseNothing, PrecNone},
	scanner.TokenEqual:  {parseNothing, parseBinary, PrecEquality},
}

func parseNothing(c *Parser, s *scanner.Scanner) error {
	return nil
}

func parseBinary(c *Parser, s *scanner.Scanner) error {
	rule := c.GetRule(c.Previous)
	prevToken := c.Previous

	if err := c.parsePrecedence(s, rule.Precedence+1); err != nil {
		return err
	}

	switch prevToken.Code {
	case scanner.TokenEqual:
		c.chunk.Write(chunk.OpEqual, prevToken.Offset)
	case scanner.TokenNotEqual:
		c.chunk.Write(chunk.OpNotEqual, prevToken.Offset)
	}

	return nil
}

func parseField(c *Parser, _ *scanner.Scanner) error {
	fieldName := chunk.Field(c.Previous.Lexeme[:len(c.Previous.Lexeme)])

	c.chunk.WriteConstant(fieldName, c.Previous.Offset)

	return nil
}

func parseString(c *Parser, _ *scanner.Scanner) error {
	lexeme := string(c.Previous.Lexeme[1 : len(c.Previous.Lexeme)-1])

	c.chunk.WriteConstant(lexeme, c.Previous.Offset)

	return nil
}

func parseNumber(c *Parser, _ *scanner.Scanner) error {
	num, err := strconv.ParseFloat(string(c.Previous.Lexeme), 64)

	if err != nil {
		return fmt.Errorf("unable to parse float; %w", err)
	}

	c.chunk.WriteConstant(num, c.Previous.Offset)

	return nil
}
