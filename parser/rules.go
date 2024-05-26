package parser

import (
	"fmt"
	"strconv"

	"github.com/CTRLRLTY/QueryL/chunk"
	"github.com/CTRLRLTY/QueryL/scanner"
)

var defaultRules = [...]ParseRule{
	scanner.TokenOr:           {nil, &parseOr, PrecOr},
	scanner.TokenAnd:          {nil, &parseAnd, PrecAnd},
	scanner.TokenField:        {&parseField, nil, PrecNone},
	scanner.TokenString:       {&parseString, nil, PrecNone},
	scanner.TokenNumber:       {&parseNumber, nil, PrecNone},
	scanner.TokenEqual:        {nil, &parseBinary, PrecEquality},
	scanner.TokenGreater:      {nil, &parseBinary, PrecComparison},
	scanner.TokenLesser:       {nil, &parseBinary, PrecComparison},
	scanner.TokenGreaterEqual: {nil, &parseBinary, PrecComparison},
	scanner.TokenLesserEqual:  {nil, &parseBinary, PrecComparison},
}

var (
	parseBinary ParseFunc = func(c *Parser, s *scanner.Scanner) error {
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
		case scanner.TokenGreater:
			c.chunk.Write(chunk.OpGreater, prevToken.Offset)
		case scanner.TokenLesser:
			c.chunk.Write(chunk.OpLesser, prevToken.Offset)
		case scanner.TokenGreaterEqual:
			c.chunk.Write(chunk.OpGreaterEqual, prevToken.Offset)
		case scanner.TokenLesserEqual:
			c.chunk.Write(chunk.OpLesserEqual, prevToken.Offset)
		}

		return nil
	}

	parseField ParseFunc = func(c *Parser, s *scanner.Scanner) error {
		fieldName := chunk.Field(c.Previous.Lexeme[:len(c.Previous.Lexeme)])

		c.chunk.WriteConstant(fieldName, c.Previous.Offset)

		return nil
	}

	parseString ParseFunc = func(c *Parser, _ *scanner.Scanner) error {
		lexeme := string(c.Previous.Lexeme[1 : len(c.Previous.Lexeme)-1])

		c.chunk.WriteConstant(lexeme, c.Previous.Offset)

		return nil
	}

	parseNumber ParseFunc = func(c *Parser, _ *scanner.Scanner) error {
		num, err := strconv.ParseFloat(string(c.Previous.Lexeme), 64)

		if err != nil {
			return fmt.Errorf("unable to parse number; %w", err)
		}

		c.chunk.WriteConstant(num, c.Previous.Offset)

		return nil
	}

	parseAnd ParseFunc = func(c *Parser, s *scanner.Scanner) error {
		endJump := c.chunk.WriteJump(chunk.OpJumpIfFalse, c.Previous.Offset)

		c.chunk.Write(chunk.OpSetAndFlag, c.Previous.Offset)
		c.chunk.Write(chunk.OpPop, c.Previous.Offset)

		if err := c.parsePrecedence(s, PrecAnd); err != nil {
			return err
		}

		c.chunk.Write(chunk.OpClearAndFlag, c.Previous.Offset)

		if err := c.chunk.PatchJump(uint16(endJump)); err != nil {
			return err
		}

		return nil
	}

	parseOr ParseFunc = func(c *Parser, s *scanner.Scanner) error {
		c.chunk.Write(chunk.OpPop, c.Previous.Offset)

		if err := c.parsePrecedence(s, PrecOr); err != nil {
			return err
		}

		return nil
	}
)
