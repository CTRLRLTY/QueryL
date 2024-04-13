package chunk

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
