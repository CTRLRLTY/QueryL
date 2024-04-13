package vm

import "github.com/CTRLRLTY/QueryL/chunk"

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
