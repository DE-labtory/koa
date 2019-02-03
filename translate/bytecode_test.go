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

package translate_test

import (
	"bytes"
	"testing"

	"github.com/DE-labtory/koa/opcode"
	"github.com/DE-labtory/koa/translate"
)

func TestBytecodeEmerge(t *testing.T) {
	tests := []struct {
		op       opcode.Type
		operands [][]byte
		result   translate.Bytecode
	}{
		{
			op: opcode.Add,
			operands: [][]byte{
				{0x01},
				{0x02},
			},
			result: translate.Bytecode{
				RawByte: []byte{byte(opcode.Add), 0x01, 0x02},
				AsmCode: []string{"Add", "01", "02"},
			},
		},
		{
			op: opcode.Mul,
			operands: [][]byte{
				{0x11},
				{0x22},
			},
			result: translate.Bytecode{
				RawByte: []byte{byte(opcode.Mul), 0x11, 0x22},
				AsmCode: []string{"Mul", "11", "22"},
			},
		},
		{
			op: opcode.Mload,
			operands: [][]byte{
				{0xff},
				{0xff},
			},
			result: translate.Bytecode{
				RawByte: []byte{byte(opcode.Mload), 0xff, 0xff},
				AsmCode: []string{"Mload", "ff", "ff"},
			},
		},
		// How to deal with byte overflow
		{
			op: opcode.LoadFunc,
			operands: [][]byte{
				{0xff, 0xff},
				{0xff, 0xff},
			},
			result: translate.Bytecode{
				RawByte: []byte{byte(opcode.LoadFunc), 0xff, 0xff, 0xff, 0xff},
				AsmCode: []string{"LoadFunc", "ffff", "ffff"},
			},
		},
		{
			op: opcode.LoadFunc,
			operands: [][]byte{
				{0x12, 0x34},
				{0x56, 0x78},
			},
			result: translate.Bytecode{
				RawByte: []byte{byte(opcode.LoadFunc), 0x12, 0x34, 0x56, 0x78},
				AsmCode: []string{"LoadFunc", "1234", "5678"},
			},
		},
		{
			op: opcode.LoadFunc,
			operands: [][]byte{
				{0x12, 0x34, 0x56},
				{0xab, 0xcd, 0xef},
			},
			result: translate.Bytecode{
				RawByte: []byte{byte(opcode.LoadFunc), 0x12, 0x34, 0x56, 0xab, 0xcd, 0xef},
				AsmCode: []string{"LoadFunc", "123456", "abcdef"},
			},
		},
		{
			op: opcode.LoadFunc,
			operands: [][]byte{
				{0x12, 0x34, 0x56, 0x78, 0x9a},
				{0xab, 0xcd, 0xef},
			},
			result: translate.Bytecode{
				RawByte: []byte{byte(opcode.LoadFunc), 0x12, 0x34, 0x56, 0x78, 0x9a, 0xab, 0xcd, 0xef},
				AsmCode: []string{"LoadFunc", "123456789a", "abcdef"},
			},
		},
		{
			op: opcode.LoadFunc,
			operands: [][]byte{
				{0x12, 0x34, 0x56},
				{0xab, 0xcd, 0xef},
				{0x12, 0x34, 0x56},
				{0xab, 0xcd, 0xef},
			},
			result: translate.Bytecode{
				RawByte: []byte{byte(opcode.LoadFunc), 0x12, 0x34, 0x56, 0xab, 0xcd, 0xef, 0x12, 0x34, 0x56, 0xab, 0xcd, 0xef},
				AsmCode: []string{"LoadFunc", "123456", "abcdef", "123456", "abcdef"},
			},
		},
	}

	for i, tt := range tests {
		b := translate.Bytecode{}
		b.Emerge(tt.op, tt.operands...)

		if !bytes.Equal(b.RawByte, tt.result.RawByte) {
			t.Errorf("test[%d] - bytecode.Emerge generate wrong RawByte expected=%v, got=%v",
				i, tt.result.RawByte, b.RawByte)
		}

		for j, code := range b.AsmCode {
			if code != tt.result.AsmCode[j] {
				t.Errorf("test[%d] - bytecode.Emerge generate wrong AsmCode at [%d], expected=%v, got=%v",
					i, j, tt.result.AsmCode[j], code)
			}
		}
	}
}
