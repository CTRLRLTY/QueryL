package filter

import (
	"fmt"
	"slices"

	"github.com/CTRLRLTY/QueryL/chunk"
	"github.com/CTRLRLTY/QueryL/parser"
	"github.com/CTRLRLTY/QueryL/vm"
)

func binaryOperator(expr *vm.VM, predicate func(v1, v2 chunk.Value) bool, record *[]map[string]any, filtered *[]map[string]any) (err error) {
	src := *record
	dst := *filtered

	a := expr.StackPop()
	b := expr.StackPop()

	field, ok := b.(chunk.Field)

	if !ok {
		err = fmt.Errorf("%v is not a valid Field OpCode", b)
		return
	}

	if expr.RegFlag&chunk.RfAnd == chunk.RfAnd {
		src = dst
		dst = make([]map[string]any, 0, len(src))
	}

	key := string(field)

	for _, doc := range src {
		if val, ok := doc[key]; ok {
			if predicate(val, a) {
				// If not already filtered
				if !slices.ContainsFunc(dst, func(fDoc map[string]any) bool {
					// Check by id
					if id, hasProp := fDoc["id"]; hasProp {
						return id == doc["id"]
					}

					return false
				}) {
					dst = append(dst, doc)
				}
			}
		}
	}

	expr.StackPush(len(dst) > 0)
	*filtered = dst

	return nil
}

func Filter(str string, record []map[string]any) (filtered []map[string]any, err error) {
	var (
		p    parser.Parser
		cnk  chunk.Chunk
		expr vm.VM
	)

	p = parser.Parser{}
	p.Init()
	cnk, err = p.ParseString(str)

	if err != nil {
		return
	}

	for i := 0; i < len(cnk.Code); i++ {
		code := cnk.Code[i]

		switch code {
		case chunk.OpSetAndFlag:
			expr.RegFlag |= chunk.RfAnd
		case chunk.OpClearAndFlag:
			expr.RegFlag &= 0b1111_1110

		case chunk.OpPop:
			expr.StackPop()
		case chunk.OpJump:
			i += 2
			offset := uint16(cnk.Code[i-1])<<8 | uint16(cnk.Code[i])
			i += int(offset)
		case chunk.OpJumpIfFalse:
			i += 2
			offset := uint16(cnk.Code[i-1])<<8 | uint16(cnk.Code[i])

			if isTrue, ok := expr.StackPeek(0).(bool); ok && !isTrue {
				i += int(offset)
			}

		case chunk.OpConstant:
			i += 1
			index := cnk.Code[i]
			val := cnk.Values[index]
			expr.StackPush(val)

		case chunk.OpEqual:
			binaryOperator(&expr, vm.Equal, &record, &filtered)

		case chunk.OpNotEqual:
			binaryOperator(&expr, vm.NotEqual, &record, &filtered)

		case chunk.OpGreater:
			binaryOperator(&expr, vm.GreaterThan, &record, &filtered)

		case chunk.OpLesser:
			binaryOperator(&expr, vm.LesserThan, &record, &filtered)

		case chunk.OpLesserEqual:
			binaryOperator(&expr, vm.LesserThanEqual, &record, &filtered)

		case chunk.OpGreaterEqual:
			binaryOperator(&expr, vm.GreaterThanEqual, &record, &filtered)
		}
	}

	return
}
