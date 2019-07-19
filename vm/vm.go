/*
 * Copyright 2018-2019 De-labtory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vm

import (
	"encoding/binary"
	"errors"

	"github.com/DE-labtory/koa/encoding"
	"github.com/DE-labtory/koa/opcode"
)

const (
	// PTRSIZE is size of arguments pointer
	PTRSIZE = 8

	// SIZEPTRSIZE is size of arguments size pointer
	SIZEPTRSIZE = 8
)

var ErrInvalidData = errors.New("Invalid data")
var ErrInvalidOpcode = errors.New("invalid opcode")

// The Execute function assemble the rawByteCode into an assembly code,
// which in turn executes the assembly logic.
func Execute(rawByteCode []byte, memory *Memory, callFunc *CallFunc) (*Stack, error) {

	s := newStack()
	asm, err := disassemble(rawByteCode)
	if err != nil {
		return &Stack{}, err
	}

	for h := asm.code[0]; h != nil; h = asm.next() {
		op, ok := h.(opCode)
		if !ok {
			return &Stack{}, ErrInvalidOpcode
		}

		err := op.Do(s, asm, memory, callFunc)
		if err != nil {
			return s, err
		}
	}

	return s, nil
}

type CallFunc struct {
	Func []byte
	Args []byte
}

// function return the Func in CallFunc
// TODO: Implements test case :-)
func (cf CallFunc) function() []byte {
	return cf.Func
}

// Example)
// Pointer : 8bytes
// Size : 8bytes
// Name : 8bytes
// arguments(n) : if the number of arguments is 4, range of n is 0~3
// cf.Args[n:n+8] : Pointer which point to value's size
// after we know size, next to size is value.
//
// CallFunc's Args
// -----------------------------------------------------------------
//  ptr1 | ptr2 | ... | size1 | value1 | size2 | value2 | ...
// -----------------------------------------------------------------
//
// arguments retrieve nth value from CallFunc Args
func (cf CallFunc) arguments(n int) []byte {
	if n < 0 {
		panic("CallFunc.arguments receive minus value as parameters")
	}

	ptr := n * PTRSIZE

	sizePtr := binary.BigEndian.Uint64(cf.Args[ptr : ptr+PTRSIZE])
	sizeVal := binary.BigEndian.Uint64(cf.Args[sizePtr : sizePtr+SIZEPTRSIZE])

	return cf.Args[sizePtr+SIZEPTRSIZE : sizePtr+SIZEPTRSIZE+sizeVal]
}

type opCode interface {
	Do(*Stack, asmReader, *Memory, *CallFunc) error
	hexer
}

// Perform opcodes logic.
// 0x0 range
type add struct{}
type mul struct{}
type sub struct{}
type div struct{}
type mod struct{}
type and struct{}
type or struct{}

// 0x10 range
type lt struct{}
type lte struct{}
type gt struct{}
type gte struct{}
type eq struct{}
type not struct{}

// 0x20 range
type pop struct{}
type push struct{}
type mload struct{}
type mstore struct{}
type msize struct{}
type loadfunc struct{}
type loadargs struct{}
type returning struct{}
type jump struct{}
type jumpDst struct{}

// 0x30 range
type jumpi struct{}
type dup struct{}
type swap struct{}
type exit struct{}

func (add) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.Pop()
	x := stack.Pop()

	stack.Push(x + y)

	return nil
}

func (add) hex() []uint8 {
	return []uint8{uint8(opcode.Add)}
}

func (mul) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.Pop()
	x := stack.Pop()

	stack.Push(x * y)

	return nil
}

func (mul) hex() []uint8 {
	return []uint8{uint8(opcode.Mul)}
}

func (sub) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.Pop()
	x := stack.Pop()

	stack.Push(x - y)

	return nil
}

func (sub) hex() []uint8 {
	return []uint8{uint8(opcode.Sub)}
}

// Be careful! int.Div and int.Quo is different
func (div) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.Pop()
	x := stack.Pop()

	item, _ := euclidean_div(x, y)

	stack.Push(item)

	return nil
}

func (div) hex() []uint8 {
	return []uint8{uint8(opcode.Div)}
}

func (mod) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.Pop()
	x := stack.Pop()

	_, item := euclidean_div(x, y)

	stack.Push(item)

	return nil
}

func (mod) hex() []uint8 {
	return []uint8{uint8(opcode.Mod)}
}

func (and) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.Pop()
	x := stack.Pop()

	ret := x & y

	stack.Push(ret)

	return nil
}

func (and) hex() []uint8 {
	return []uint8{uint8(opcode.And)}
}

func (or) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.Pop()
	x := stack.Pop()

	ret := x | y

	stack.Push(ret)

	return nil
}

func (or) hex() []uint8 {
	return []uint8{uint8(opcode.Or)}
}

func (lt) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y, x := stack.Pop(), stack.Pop()

	if x < y { // x < y
		stack.Push(item(1))
	} else {
		stack.Push(item(0))
	}

	return nil
}

func (lt) hex() []uint8 {
	return []uint8{uint8(opcode.LT)}
}

func (lte) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y, x := stack.Pop(), stack.Pop()

	if x <= y { // x <= y
		stack.Push(item(1))
	} else {
		stack.Push(item(0))
	}

	return nil
}

func (lte) hex() []uint8 {
	return []uint8{uint8(opcode.LTE)}
}

func (gt) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y, x := stack.Pop(), stack.Pop()

	if x > y { // x > y
		stack.Push(item(1))
	} else {
		stack.Push(item(0))
	}

	return nil
}

func (gt) hex() []uint8 {
	return []uint8{uint8(opcode.GT)}
}

func (gte) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y, x := stack.Pop(), stack.Pop()

	if x >= y { // x >= y
		stack.Push(item(1))
	} else {
		stack.Push(item(0))
	}

	return nil
}

func (gte) hex() []uint8 {
	return []uint8{uint8(opcode.GTE)}
}

func (eq) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y, x := stack.Pop(), stack.Pop()

	if x == y { // x == y
		stack.Push(item(1))
	} else {
		stack.Push(item(0))
	}

	return nil
}

func (eq) hex() []uint8 {
	return []uint8{uint8(opcode.EQ)}
}

func (not) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	x := stack.Pop()

	if x == 1 {
		x = 0
	} else {
		x = 1
	}

	stack.Push(x)
	return nil
}

func (not) hex() []uint8 {
	return []uint8{uint8(opcode.NOT)}
}

func (pop) Do(stack *Stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	_ = stack.Pop()
	return nil
}

func (pop) hex() []uint8 {
	return []uint8{uint8(opcode.Pop)}
}

func (push) Do(stack *Stack, asm asmReader, _ *Memory, contract *CallFunc) error {
	code := asm.next()
	data, ok := code.(Data)
	if !ok {
		return ErrInvalidData
	}
	item := bytesToItem(data.hex())
	stack.Push(item)

	return nil
}

func (push) hex() []uint8 {
	return []uint8{uint8(opcode.Push)}
}

func (mload) Do(stack *Stack, _ asmReader, memory *Memory, _ *CallFunc) error {
	offset, size := stack.Pop(), stack.Pop()
	value := memory.GetVal(uint64(offset), uint64(size))

	stack.Push(bytesToItem(value))
	return nil
}

func (mload) hex() []uint8 {
	return []uint8{uint8(opcode.Mload)}
}

func (mstore) Do(stack *Stack, _ asmReader, memory *Memory, _ *CallFunc) error {
	offset, size, value := stack.Pop(), stack.Pop(), stack.Pop()

	//memSize := uint64(memory.Len()) + uint64(size)
	//memory.Resize(memSize)

	convertedValue := int64ToBytes(int64(value))
	memory.Sets(uint64(offset), uint64(size), convertedValue)
	return nil
}

func (mstore) hex() []uint8 {
	return []uint8{uint8(opcode.Mstore)}
}

func (msize) Do(stack *Stack, _ asmReader, memory *Memory, _ *CallFunc) error {
	size := stack.Pop()
	memory.Resize(uint64(size))
	return nil
}

func (msize) hex() []uint8 {
	return []uint8{uint8(opcode.Mstore)}
}

func (loadfunc) Do(stack *Stack, _ asmReader, _ *Memory, callfunc *CallFunc) error {
	function := callfunc.function()

	convertedFunc, err := encoding.EncodeOperand(function)
	if err != nil {
		return err
	}

	stack.Push(bytesToItem(convertedFunc))
	return nil
}

func (loadfunc) hex() []uint8 {
	return []uint8{uint8(opcode.LoadFunc)}
}

func (loadargs) Do(stack *Stack, _ asmReader, _ *Memory, callfunc *CallFunc) error {
	index := stack.Pop()
	argument := callfunc.arguments(int(index))

	stack.Push(bytesToItem(argument))

	return nil
}

func (loadargs) hex() []uint8 {
	return []uint8{uint8(opcode.LoadArgs)}
}

func (returning) Do(stack *Stack, asm asmReader, memory *Memory, _ *CallFunc) error {
	value, _, pos := stack.Pop(), stack.Pop(), stack.Pop()

	asm.jump(uint64(pos - 1))

	stack.Push(value)
	return nil
}

func (returning) hex() []uint8 {
	return []uint8{uint8(opcode.Returning)}
}

func (jump) Do(stack *Stack, asm asmReader, memory *Memory, _ *CallFunc) error {
	pos := stack.Pop()
	asm.jump(uint64(pos - 1))
	return nil
}

func (jump) hex() []uint8 {
	return []uint8{uint8(opcode.Jump)}
}

func (jumpDst) Do(stack *Stack, asm asmReader, memory *Memory, _ *CallFunc) error {
	return nil
}

func (jumpDst) hex() []uint8 {
	return []uint8{uint8(opcode.JumpDst)}
}

func (jumpi) Do(stack *Stack, asm asmReader, memory *Memory, _ *CallFunc) error {
	pos, cond := stack.Pop(), stack.Pop()
	if cond == item(0) { // cond == false
		asm.jump(uint64(pos - 1))
	}
	return nil
}

func (jumpi) hex() []uint8 {
	return []uint8{uint8(opcode.Jumpi)}
}

func (dup) Do(stack *Stack, _ asmReader, memory *Memory, _ *CallFunc) error {
	stack.Dup()
	return nil
}

func (dup) hex() []uint8 {
	return []uint8{uint8(opcode.DUP)}
}

func (swap) Do(stack *Stack, _ asmReader, memory *Memory, _ *CallFunc) error {
	stack.Swap()
	return nil
}

func (swap) hex() []uint8 {
	return []uint8{uint8(opcode.SWAP)}
}

func (exit) Do(stack *Stack, asm asmReader, memory *Memory, _ *CallFunc) error {
	for asm.next() != nil {
	}
	return nil
}

func (exit) hex() []uint8 {
	return []uint8{uint8(opcode.Exit)}
}

func int64ToBytes(int64 int64) []byte {
	byteSlice := make([]byte, 8)
	binary.BigEndian.PutUint64(byteSlice, uint64(int64))
	return byteSlice
}

func bytesToItem(bytes []byte) item {
	item := item(binary.BigEndian.Uint64(bytes))
	return item
}

func euclidean_div(a item, b item) (item, item) {
	var q int64
	var r int64
	A := int64(a)
	B := int64(b)

	if A < 0 && B > 0 {
		q = int64(A/B) - 1
		r = A - (B * q)
	} else if A > 0 && B < 0 {
		q = int64(A / B)
		r = A - (B * q)
	} else if A > 0 && B > 0 {
		q = int64(A / B)
		r = A - (B * q)
	} else if A < 0 && B < 0 {
		q = int64((A + B) / B)
		r = A - (B * q)
	}

	return item(q), item(r)
}
