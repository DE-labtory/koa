/*
 * Copyright 2018 De-labtory
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

	"github.com/DE-labtory/koa/opcode"
)

const (
	// PTRSIZE is size of arguments pointer
	PTRSIZE = 4

	// SIZEPTRSIZE is size of arguments size pointer
	SIZEPTRSIZE = 4
)

var ErrInvalidData = errors.New("Invalid data")
var ErrInvalidOpcode = errors.New("invalid opcode")

// The Execute function assemble the rawByteCode into an assembly code,
// which in turn executes the assembly logic.
func Execute(rawByteCode []byte, memory *Memory, callFunc *CallFunc) (*stack, error) {

	s := newStack()
	asm, err := disassemble(rawByteCode)
	if err != nil {
		return &stack{}, err
	}

	for h := asm.code[0]; h != nil; h = asm.next() {
		op, ok := h.(opCode)
		if !ok {
			return &stack{}, ErrInvalidOpcode
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
// Pointer : 4bytes
// Size : 4bytes
// Value : dynamic
// arguments(n) : if the number of arguments is 4, range of n is 0~3
// cf.Args[n:n+4] : Pointer which point to value's size
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

	sizePtr := binary.BigEndian.Uint32(cf.Args[ptr : ptr+PTRSIZE])
	sizeVal := binary.BigEndian.Uint32(cf.Args[sizePtr : sizePtr+SIZEPTRSIZE])

	return cf.Args[sizePtr+SIZEPTRSIZE : sizePtr+SIZEPTRSIZE+sizeVal]
}

type opCode interface {
	Do(*stack, asmReader, *Memory, *CallFunc) error
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
type loadfunc struct{}
type loadargs struct{}
type returning struct{}
type jump struct{}
type jumpDst struct{}
type jumpi struct{}

// 0x30 range
type dup struct{}
type swap struct{}

func (add) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.pop()
	x := stack.pop()

	stack.push(x + y)

	return nil
}

func (add) hex() []uint8 {
	return []uint8{uint8(opcode.Add)}
}

func (mul) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.pop()
	x := stack.pop()

	stack.push(x * y)

	return nil
}

func (mul) hex() []uint8 {
	return []uint8{uint8(opcode.Mul)}
}

func (sub) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.pop()
	x := stack.pop()

	stack.push(x - y)

	return nil
}

func (sub) hex() []uint8 {
	return []uint8{uint8(opcode.Sub)}
}

// Be careful! int.Div and int.Quo is different
func (div) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.pop()
	x := stack.pop()

	item, _ := euclidean_div(x, y)

	stack.push(item)

	return nil
}

func (div) hex() []uint8 {
	return []uint8{uint8(opcode.Div)}
}

func (mod) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.pop()
	x := stack.pop()

	_, item := euclidean_div(x, y)

	stack.push(item)

	return nil
}

func (mod) hex() []uint8 {
	return []uint8{uint8(opcode.Mod)}
}

func (and) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.pop()
	x := stack.pop()

	ret := x & y

	stack.push(ret)

	return nil
}

func (and) hex() []uint8 {
	return []uint8{uint8(opcode.And)}
}

func (or) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y := stack.pop()
	x := stack.pop()

	ret := x | y

	stack.push(ret)

	return nil
}

func (or) hex() []uint8 {
	return []uint8{uint8(opcode.Or)}
}

func (lt) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y, x := stack.pop(), stack.pop()

	if x < y { // x < y
		stack.push(item(1))
	} else {
		stack.push(item(0))
	}

	return nil
}

func (lt) hex() []uint8 {
	return []uint8{uint8(opcode.LT)}
}

func (lte) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y, x := stack.pop(), stack.pop()

	if x <= y { // x <= y
		stack.push(item(1))
	} else {
		stack.push(item(0))
	}

	return nil
}

func (lte) hex() []uint8 {
	return []uint8{uint8(opcode.LTE)}
}

func (gt) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y, x := stack.pop(), stack.pop()

	if x > y { // x > y
		stack.push(item(1))
	} else {
		stack.push(item(0))
	}

	return nil
}

func (gt) hex() []uint8 {
	return []uint8{uint8(opcode.GT)}
}

func (gte) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y, x := stack.pop(), stack.pop()

	if x >= y { // x >= y
		stack.push(item(1))
	} else {
		stack.push(item(0))
	}

	return nil
}

func (gte) hex() []uint8 {
	return []uint8{uint8(opcode.GTE)}
}

func (eq) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	y, x := stack.pop(), stack.pop()

	if x == y { // x == y
		stack.push(item(1))
	} else {
		stack.push(item(0))
	}

	return nil
}

func (eq) hex() []uint8 {
	return []uint8{uint8(opcode.EQ)}
}

func (not) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	x := stack.pop()

	stack.push(^x)
	return nil
}

func (not) hex() []uint8 {
	return []uint8{uint8(opcode.NOT)}
}

func (pop) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	_ = stack.pop()
	return nil
}

func (pop) hex() []uint8 {
	return []uint8{uint8(opcode.Pop)}
}

func (push) Do(stack *stack, asm asmReader, _ *Memory, contract *CallFunc) error {
	code := asm.next()
	data, ok := code.(Data)
	if !ok {
		return ErrInvalidData
	}
	item := item(bytesToInt32(data.hex()))
	stack.push(item)

	return nil
}

func (push) hex() []uint8 {
	return []uint8{uint8(opcode.Push)}
}

// TODO: implement me w/ test cases :-)
func (mload) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	return nil
}

func (mload) hex() []uint8 {
	return []uint8{uint8(opcode.Mload)}
}

// TODO: implement me w/ test cases :-)
func (mstore) Do(stack *stack, _ asmReader, _ *Memory, _ *CallFunc) error {
	return nil
}

func (mstore) hex() []uint8 {
	return []uint8{uint8(opcode.Mstore)}
}

// TODO: implement me w/ test cases :-)
func (loadfunc) Do(stack *stack, _ asmReader, _ *Memory, callfunc *CallFunc) error {
	return nil
}

func (loadfunc) hex() []uint8 {
	return []uint8{uint8(opcode.LoadFunc)}
}

// TODO: implement me w/ test cases :-)
func (loadargs) Do(stack *stack, _ asmReader, _ *Memory, callfunc *CallFunc) error {
	return nil
}

func (loadargs) hex() []uint8 {
	return []uint8{uint8(opcode.LoadArgs)}
}

// TODO: implement me w/ test cases :-)
func (returning) Do(stack *stack, _ asmReader, memory *Memory, _ *CallFunc) error {
	return nil
}

func (returning) hex() []uint8 {
	return []uint8{uint8(opcode.Returning)}
}

func (jump) Do(stack *stack, a asmReader, memory *Memory, _ *CallFunc) error {
	pos := stack.pop()
	if !a.validateJumpDst(uint64(pos)) {
		return errors.New("invalid jump target")
	}
	a.jump(uint64(pos))
	return nil
}

func (jump) hex() []uint8 {
	return []uint8{uint8(opcode.Jump)}
}

func (jumpDst) Do(stack *stack, a asmReader, memory *Memory, _ *CallFunc) error {
	return nil
}

func (jumpDst) hex() []uint8 {
	return []uint8{uint8(opcode.JumpDst)}
}

func (jumpi) Do(stack *stack, a asmReader, memory *Memory, _ *CallFunc) error {
	pos, cond := stack.pop(), stack.pop()
	if cond == item(0) { // cond == false
		if !a.validateJumpDst(uint64(pos)) {
			return errors.New("invalid jump target")
		}
		a.jump(uint64(pos))
	}
	return nil
}

func (jumpi) hex() []uint8 {
	return []uint8{uint8(opcode.Jumpi)}
}

func (dup) Do(stack *stack, _ asmReader, memory *Memory, _ *CallFunc) error {
	stack.dup()
	return nil
}

func (dup) hex() []uint8 {
	return []uint8{uint8(opcode.DUP)}
}

func (swap) Do(stack *stack, _ asmReader, memory *Memory, _ *CallFunc) error {
	stack.swap()
	return nil
}

func (swap) hex() []uint8 {
	return []uint8{uint8(opcode.SWAP)}
}

func int32ToBytes(int32 int32) []byte {
	byteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(byteSlice, uint32(int32))
	return byteSlice
}

func bytesToInt32(bytes []byte) int32 {
	int32 := int32(binary.BigEndian.Uint32(bytes))
	return int32
}

func euclidean_div(a item, b item) (item, item) {
	var q int32
	var r int32
	A := int32(a)
	B := int32(b)

	if A < 0 && B > 0 {
		q = int32(A/B) - 1
		r = A - (B * q)
	} else if A > 0 && B < 0 {
		q = int32(A / B)
		r = A - (B * q)
	} else if A > 0 && B > 0 {
		q = int32(A / B)
		r = A - (B * q)
	} else if A < 0 && B < 0 {
		q = int32((A + B) / B)
		r = A - (B * q)
	}

	return item(q), item(r)
}
