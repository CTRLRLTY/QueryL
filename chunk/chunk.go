package chunk

import (
	"fmt"
	"math"
)

const (
	OpConstant byte = iota
	OpAnd
	OpOr
	OpEqual
	OpNotEqual
	OpGreater
	OpLesser
	OpGreaterEqual
	OpLesserEqual
	OpNot
	OpPop
	OpJumpIfFalse
)

type Field string

type Value interface{}

type Chunk struct {
	Code    []byte
	Values  []Value
	Offsets []uint32
}

func (c *Chunk) Write(b byte, ofs uint32) {
	c.Code = append(c.Code, b)
	c.Offsets = append(c.Offsets, ofs)
}

func (c *Chunk) WriteConstant(v Value, ofs uint32) {
	c.Write(OpConstant, ofs)
	c.Values = append(c.Values, v)
	c.Write(byte(len(c.Values)-1), ofs)
}

func (c *Chunk) WriteJump(code byte, ofs uint32) int {
	c.Write(code, ofs)

	// Temporary 16 bytes operand
	c.Write(0xff, ofs)
	c.Write(0xff, ofs)

	// returns the little-end boundary index of the 16-bytes operand
	return len(c.Code) - 2
}

func (c *Chunk) PatchJump(chunkIndex uint16) error {
	jump := len(c.Code) - int(chunkIndex) - 2

	if jump > math.MaxUint16 {
		return fmt.Errorf("jump exceeds 16 bytes")
	}

	c.Code[chunkIndex] = byte(jump >> 8 & 0xff)
	c.Code[chunkIndex+1] = byte(jump & 0xff)

	return nil
}
