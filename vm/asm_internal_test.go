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
	"testing"

	"bytes"

	"github.com/DE-labtory/koa/opcode"
)

func TestAssemble(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(1), // 1 byte allocation
		uint8(opcode.Push), uint32ToBytes(300), // 2 byte allocation
		uint8(opcode.Add),
	)
	testAsmExpected := asm{
		code: []hexer{
			push{}, Data{Body: uint32ToBytes(1)},
			push{}, Data{Body: uint32ToBytes(300)},
			add{},
		},
	}

	asm, err := assemble(testByteCode)
	if err != nil {
		t.Error(err)
	}

	//opCode is a byte
	//Data is 4 bytes
	for i, code := range asm.code {
		switch code.(type) {
		case opCode:
			if !bytes.Equal(testAsmExpected.code[i].(opCode).hex(), code.(opCode).hex()) {
				t.Fatal("There is an error in Opcode during analysis")
			}
		case Data:
			if !bytes.Equal(testAsmExpected.code[i].(Data).hex(), code.(Data).hex()) {
				t.Fatal("There is an error in Data during analysis")
			}
		}

	}
}

func TestAssemble_invalid(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(1),
		uint8(255), uint32ToBytes(300),
		uint8(opcode.Add),
	)

	_, err := assemble(testByteCode)
	if err != ErrInvalidOpcode {
		t.Error("The desired error was not found")
	}
}

func TestNext(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(1),
		uint8(opcode.Push), uint32ToBytes(2),
		uint8(opcode.Add),
	)

	asm, err := assemble(testByteCode)
	if err != nil {
		t.Error(err)
	}

	for asm.pc+1 < uint64(len(asm.code)) {
		prevPc := asm.pc
		code := asm.next()
		if prevPc+1 != asm.pc {
			t.Error("Error in moving pc")
		}

		_, ok := code.(hexer)
		if !ok {
			t.Errorf("Error at %d", asm.pc-1)
		}
	}

	if asm.next() != nil {
		t.Error("Error at last index")
	}
}

func TestJump(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(1),
		uint8(opcode.Push), uint32ToBytes(2),
		uint8(opcode.Add),
	)

	asm, err := assemble(testByteCode)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Access to invalid program counter!")
		}
	}()
	asm.jump(1)
	if asm.pc != 1 {
		t.Errorf("Invalid pc - expected=1, got=%d", asm.pc)
	}

	asm.jump(2)
	if asm.pc != 2 {
		t.Errorf("Invalid pc - expected=2, got=%d", asm.pc)
	}

	asm.jump(4)
	if asm.pc != 4 {
		t.Errorf("Invalid pc - expected=4, got=%d", asm.pc)
	}
}

func TestJump_invalid(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(1),
		uint8(opcode.Push), uint32ToBytes(2),
		uint8(opcode.Add),
	)

	asm, err := assemble(testByteCode)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("The desired error was not found")
		}
	}()
	asm.jump(15)
	asm.jump(16)
	asm.jump(17)
}
