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
	"github.com/DE-labtory/koa/opcode"
	"github.com/pkg/errors"
)

var opcodes = map[opcode.Type]Opcode{
	opcode.Add:  Add{},
	opcode.Push: Push{},
}

type Item uint32

func Execute(instructions []Item) (*Stack, error) {
	stack := NewStack()

	Instructions := NewInstructions(instructions)

	for ; Instructions.ip < len(Instructions.Buffer); Instructions.ip++ {
		hex := opcode.Type(Instructions.Buffer[Instructions.ip])
		opCode, ok := opcodes[hex]

		if !ok {
			return &Stack{}, errors.New("No opcode")
		}

		err := opCode.Do(stack, &Instructions)
		if err != nil {
			return &Stack{}, err
		}
	}

	return stack, nil
}

// Before execute vm, check some validation.
func ParseInstructions(instructions []Item) []Item {
	return nil
}

type Opcode interface {
	Do(*Stack, InstructionBuffer) error
	Code() uint32
}

type Add struct{}

func (Add) Do(stack *Stack, instructionBuffer InstructionBuffer) error {
	rightValue := stack.Pop()
	leftValue := stack.Pop()

	result := rightValue + leftValue

	stack.Push(result)

	return nil
}

func (Add) Code() uint32 {
	return uint32(opcode.Add)
}

type Push struct{}

func (Push) Do(stack *Stack, instructionBuffer InstructionBuffer) error {
	item := instructionBuffer.GetItems(1)
	stack.Push(item[0])
	return nil
}

func (Push) Code() uint32 {
	return uint32(opcode.Push)
}
