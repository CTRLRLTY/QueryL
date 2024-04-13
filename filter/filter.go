package filter

import (
	"fmt"

	"github.com/CTRLRLTY/QueryL/chunk"
	"github.com/CTRLRLTY/QueryL/parser"
	"github.com/CTRLRLTY/QueryL/vm"
)

func Filter(str string, record []map[string]any) (filtered []map[string]any, err error) {
	var (
		p   parser.Parser
		cnk chunk.Chunk
		vm  vm.VM
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
		case chunk.OpConstant:
			i += 1
			index := cnk.Code[i]
			val := cnk.Values[index]
			vm.StackPush(val)

		case chunk.OpEqual:
			a := vm.StackPop()
			b := vm.StackPop()

			field, ok := b.(chunk.Field)

			if !ok {
				err = fmt.Errorf("%v is not a valid Field OpCode", b)
				return
			}

			key := string(field)

			for _, doc := range record {
				if val, ok := doc[key]; ok {
					if val == a {
						filtered = append(filtered, doc)
					}
				}
			}

			vm.StackPush(len(filtered) > 0)

		case chunk.OpNotEqual:
			a := vm.StackPop()
			b := vm.StackPop()

			field, ok := b.(chunk.Field)

			if !ok {
				err = fmt.Errorf("%v is not a valid Field OpCode", b)
				return
			}

			key := string(field)

			for _, doc := range record {
				if val, ok := doc[key]; ok {
					if val != a {
						filtered = append(filtered, doc)
					}
				}
			}

			vm.StackPush(len(filtered) > 0)
		}
	}

	return
}
