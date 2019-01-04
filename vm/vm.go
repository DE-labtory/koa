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
	"errors"

	"encoding/binary"

	"github.com/DE-labtory/koa/opcode"
)

var ErrInvalidData = errors.New("Invalid data")
var ErrInvalidOpcode = errors.New("invalid opcode")

// The Execute function assemble the rawByteCode into an assembly code,
// which in turn executes the assembly logic.
func Execute(rawByteCode []byte) (*stack, error) {

	s := newStack()
	asm, err := assemble(rawByteCode)
	if err != nil {
		return &stack{}, err
	}

	for h := asm.code[0]; h != nil; h = asm.next() {
		op, ok := h.(opCode)
		if !ok {
			return &stack{}, ErrInvalidOpcode
		}

		err := op.Do(s, asm)
		if err != nil {
			return s, err
		}
	}

	return s, nil
}

type opCode interface {
	Do(*stack, asmReader) error
	hexer
}

// Perform add opcode logic.
type add struct{}

// TODO: implement me w/ test cases :-)
func (add) Do(stack *stack, _ asmReader) error {
	rightValue := stack.pop()
	leftValue := stack.pop()
	result := rightValue + leftValue

	stack.push(result)

	return nil
}

// TODO: implement me w/ test cases :-)
func (add) hex() []uint8 {
	return []uint8{uint8(opcode.Add)}
}

type push struct{}

func (push) Do(stack *stack, asm asmReader) error {
	code := asm.next()
	data, ok := code.(Data)
	if !ok {
		return ErrInvalidData
	}
	item := item(bytesToUint32(data.hex()))
	stack.push(item)

	return nil
}

func (push) hex() []uint8 {
	return []uint8{uint8(opcode.Push)}
}

func uint32ToBytes(uint32 uint32) []byte {
	byteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(byteSlice, uint32)
	return byteSlice
}

func bytesToUint32(bytes []byte) uint32 {
	uint32 := binary.BigEndian.Uint32(bytes)
	return uint32
}
