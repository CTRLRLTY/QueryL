package filter

import (
	"fmt"

	"github.com/CTRLRLTY/QueryL/chunk"
	"github.com/CTRLRLTY/QueryL/parser"
	"github.com/CTRLRLTY/QueryL/vm"
)

func Filter(str string, record []map[string]any) (filtered []map[string]any, err error) {
	var (
		p    parser.Parser
		cnk  chunk.Chunk
		expr vm.VM
	)

	recordcopy := append(make([]map[string]any, 0, len(record)), record...)

	p = parser.Parser{}
	p.Init()
	cnk, err = p.ParseString(str)

	if err != nil {
		return
	}

	for i := 0; i < len(cnk.Code); i++ {
		code := cnk.Code[i]

		switch code {
		case chunk.OpResetFiltered:
			filtered = make([]map[string]any, 0, len(recordcopy))

		case chunk.OpResetCopy:
			recordcopy = append(make([]map[string]any, 0, len(record)), record...)

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
			a := expr.StackPop()
			b := expr.StackPop()

			field, ok := b.(chunk.Field)

			if !ok {
				err = fmt.Errorf("%v is not a valid Field OpCode", b)
				return
			}

			key := string(field)

			for _, doc := range recordcopy {
				if val, ok := doc[key]; ok {
					if vm.Equal(val, a) {
						filtered = append(filtered, doc)
					}
				}
			}

			expr.StackPush(len(filtered) > 0)
			recordcopy = filtered

		case chunk.OpNotEqual:
			a := expr.StackPop()
			b := expr.StackPop()

			field, ok := b.(chunk.Field)

			if !ok {
				err = fmt.Errorf("%v is not a valid Field OpCode", b)
				return
			}

			key := string(field)

			for _, doc := range recordcopy {
				if val, ok := doc[key]; ok {
					if !vm.Equal(val, a) {
						filtered = append(filtered, doc)
					}
				}
			}

			expr.StackPush(len(filtered) > 0)
			recordcopy = filtered

		case chunk.OpGreater:
			a := expr.StackPop()
			b := expr.StackPop()

			field, ok := b.(chunk.Field)

			if !ok {
				err = fmt.Errorf("%v is not a valid Field OpCode", b)
				return
			}

			key := string(field)

			for _, doc := range recordcopy {
				if val, ok := doc[key]; ok {
					if vm.GreaterThan(val, a) {
						filtered = append(filtered, doc)
					}
				}
			}

			expr.StackPush(len(filtered) > 0)
			recordcopy = filtered

		case chunk.OpLesser:
			a := expr.StackPop()
			b := expr.StackPop()

			field, ok := b.(chunk.Field)

			if !ok {
				err = fmt.Errorf("%v is not a valid Field OpCode", b)
				return
			}

			key := string(field)

			for _, doc := range recordcopy {
				if val, ok := doc[key]; ok {
					if vm.LesserThan(val, a) {
						filtered = append(filtered, doc)
					}
				}
			}

			expr.StackPush(len(filtered) > 0)
			recordcopy = filtered

		case chunk.OpGreaterEqual:
			a := expr.StackPop()
			b := expr.StackPop()

			field, ok := b.(chunk.Field)

			if !ok {
				err = fmt.Errorf("%v is not a valid Field OpCode", b)
				return
			}

			key := string(field)

			for _, doc := range recordcopy {
				if val, ok := doc[key]; ok {
					if vm.Equal(val, a) || vm.GreaterThan(val, a) {
						filtered = append(filtered, doc)
					}
				}
			}

			expr.StackPush(len(filtered) > 0)
			recordcopy = filtered

		case chunk.OpLesserEqual:
			a := expr.StackPop()
			b := expr.StackPop()

			field, ok := b.(chunk.Field)

			if !ok {
				err = fmt.Errorf("%v is not a valid Field OpCode", b)
				return
			}

			key := string(field)

			for _, doc := range recordcopy {
				if val, ok := doc[key]; ok {
					if vm.Equal(val, a) || vm.LesserThan(val, a) {
						filtered = append(filtered, doc)
					}
				}
			}

			expr.StackPush(len(filtered) > 0)
			recordcopy = filtered
		}
	}

	return
}
