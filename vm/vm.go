package vm

import (
	"fmt"

	"github.com/CTRLRLTY/QueryL/chunk"
)

type VM struct {
	stack []chunk.Value
}

func (vm *VM) StackPush(v chunk.Value) {
	vm.stack = append(vm.stack, v)
}

func (vm *VM) StackPop() (v chunk.Value) {
	l := len(vm.stack)

	if l == 0 {
		return nil
	}

	v = vm.stack[l-1]
	vm.stack = vm.stack[:l-1]

	return
}

func (vm *VM) StackPeek(distance int) chunk.Value {
	return vm.stack[len(vm.stack)-1-distance]
}

func Number2Float(val chunk.Value) (f float64, err error) {
	switch v := val.(type) {
	case int:
		f = float64(v)
	case int32:
		f = float64(v)
	case int64:
		f = float64(v)
	case float32:
		f = float64(v)
	case float64:
		f = v
	default:
		err = fmt.Errorf("unsupported typecast")
		return
	}

	return
}

func Equal(a chunk.Value, b chunk.Value) bool {
	var (
		num1 float64
		num2 float64
		err1 error
		err2 error
	)

	num1, err1 = Number2Float(a)
	num2, err2 = Number2Float(b)

	if err1 == nil && err2 == nil {
		return num1 == num2
	}

	return a == b
}

func LesserThan(a chunk.Value, b chunk.Value) bool {
	var (
		num1 float64
		num2 float64
		err  error
	)

	if num1, err = Number2Float(a); err != nil {
		return false
	}

	if num2, err = Number2Float(b); err != nil {
		return false
	}

	return num1 < num2
}

func GreaterThan(a chunk.Value, b chunk.Value) bool {
	var (
		num1 float64
		num2 float64
		err  error
	)

	if num1, err = Number2Float(a); err != nil {
		return false
	}

	if num2, err = Number2Float(b); err != nil {
		return false
	}

	return num1 > num2
}
