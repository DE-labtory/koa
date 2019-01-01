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

	"github.com/DE-labtory/koa/opcode"
)

// The Execute function assemble the rawByteCode into an assembly code,
// which in turn executes the assembly logic.
func Execute(rawByteCode []byte) (*stack, error) {

	s := newStack()
	asm, err := assemble(rawByteCode)
	if err != nil {
		return &stack{}, err
	}

	for h := asm.next(); h != nil; {
		op, ok := h.(opCode)
		if !ok {
			return &stack{}, errors.New("invalid opcode err")
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

	return nil
}

// TODO: implement me w/ test cases :-)
func (add) hex() []uint8 {
	return []uint8{uint8(opcode.Add)}
}
