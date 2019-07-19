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

package translate_test

import (
	"testing"

	"bytes"

	"github.com/DE-labtory/koa/opcode"
	"github.com/DE-labtory/koa/translate"
)

func TestBytecodeEmerge(t *testing.T) {
	tests := []struct {
		op       opcode.Type
		operands [][]byte
		result   translate.Asm
	}{
		{
			op: opcode.Add,
			operands: [][]byte{
				{0x01},
				{0x02},
			},
			result: translate.Asm{
				AsmCodes: []translate.AsmCode{
					translate.AsmCode{
						RawByte: []byte{byte(opcode.Add)},
						Value:   "Add",
					},
					{
						RawByte: []byte{0x01},
						Value:   "01",
					},
					{
						RawByte: []byte{0x02},
						Value:   "02",
					},
				},
			},
		},
		{
			op: opcode.Mul,
			operands: [][]byte{
				{0x11},
				{0x22},
			},
			result: translate.Asm{
				AsmCodes: []translate.AsmCode{
					translate.AsmCode{
						RawByte: []byte{byte(opcode.Mul)},
						Value:   "Mul",
					},
					{
						RawByte: []byte{0x11},
						Value:   "11",
					},
					{
						RawByte: []byte{0x22},
						Value:   "22",
					},
				},
			},
		},
		{
			op: opcode.Mload,
			operands: [][]byte{
				{0xff},
				{0xff},
			},
			result: translate.Asm{
				AsmCodes: []translate.AsmCode{
					translate.AsmCode{
						RawByte: []byte{byte(opcode.Mload)},
						Value:   "Mload",
					},
					{
						RawByte: []byte{0xff},
						Value:   "ff",
					},
					{
						RawByte: []byte{0xff},
						Value:   "ff",
					},
				},
			},
		},
		// How to deal with byte overflow
		{
			op: opcode.LoadFunc,
			operands: [][]byte{
				{0xff, 0xff},
				{0xff, 0xff},
			},
			result: translate.Asm{
				AsmCodes: []translate.AsmCode{
					translate.AsmCode{
						RawByte: []byte{byte(opcode.LoadFunc)},
						Value:   "LoadFunc",
					},
					{
						RawByte: []byte{0xff, 0xff},
						Value:   "ffff",
					},
					{
						RawByte: []byte{0xff, 0xff},
						Value:   "ffff",
					},
				},
			},
		},
		{
			op: opcode.LoadFunc,
			operands: [][]byte{
				{0x12, 0x34},
				{0x56, 0x78},
			},
			result: translate.Asm{
				AsmCodes: []translate.AsmCode{
					translate.AsmCode{
						RawByte: []byte{byte(opcode.LoadFunc)},
						Value:   "LoadFunc",
					},
					{
						RawByte: []byte{0x12, 0x34},
						Value:   "1234",
					},
					{
						RawByte: []byte{0x56, 0x78},
						Value:   "5678",
					},
				},
			},
		},
		{
			op: opcode.LoadFunc,
			operands: [][]byte{
				{0x12, 0x34, 0x56},
				{0xab, 0xcd, 0xef},
			},
			result: translate.Asm{
				AsmCodes: []translate.AsmCode{
					translate.AsmCode{
						RawByte: []byte{byte(opcode.LoadFunc)},
						Value:   "LoadFunc",
					},
					{
						RawByte: []byte{0x12, 0x34, 0x56},
						Value:   "123456",
					},
					{
						RawByte: []byte{0xab, 0xcd, 0xef},
						Value:   "abcdef",
					},
				},
			},
		},
		{
			op: opcode.LoadFunc,
			operands: [][]byte{
				{0x12, 0x34, 0x56, 0x78, 0x9a},
				{0xab, 0xcd, 0xef},
			},
			result: translate.Asm{
				AsmCodes: []translate.AsmCode{
					translate.AsmCode{
						RawByte: []byte{byte(opcode.LoadFunc)},
						Value:   "LoadFunc",
					},
					{
						RawByte: []byte{0x12, 0x34, 0x56, 0x78, 0x9a},
						Value:   "123456789a",
					},
					{
						RawByte: []byte{0xab, 0xcd, 0xef},
						Value:   "abcdef",
					},
				},
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
			result: translate.Asm{
				AsmCodes: []translate.AsmCode{
					translate.AsmCode{
						RawByte: []byte{byte(opcode.LoadFunc)},
						Value:   "LoadFunc",
					},
					{
						RawByte: []byte{0x12, 0x34, 0x56},
						Value:   "123456",
					},
					{
						RawByte: []byte{0xab, 0xcd, 0xef},
						Value:   "abcdef",
					},
					{
						RawByte: []byte{0x12, 0x34, 0x56},
						Value:   "123456",
					},
					{
						RawByte: []byte{0xab, 0xcd, 0xef},
						Value:   "abcdef",
					},
				},
			},
		},
	}

	for i, tt := range tests {
		a := translate.Asm{}
		a.Emerge(tt.op, tt.operands...)
		for j, code := range a.AsmCodes {
			if code.Value != tt.result.AsmCodes[j].Value {
				t.Errorf("test[%d] - bytecode.Emerge generate wrong AsmCode at [%d], expected=%v, got=%v",
					i, j, tt.result.AsmCodes[j], code.Value)
			}

			if !bytes.Equal(code.RawByte, tt.result.AsmCodes[j].RawByte) {
				t.Errorf("test[%d] - bytecode.Emerge generate wrong RawByte expected=%v, got=%v",
					i, tt.result.AsmCodes[j].RawByte, code.RawByte)
			}
		}
	}
}

func TestAsm_ToRawByteCode(t *testing.T) {
	tests := []struct {
		asm    translate.Asm
		expect []byte
	}{
		{
			asm: translate.Asm{
				AsmCodes: []translate.AsmCode{
					{
						RawByte: []byte{byte(opcode.Add)},
						Value:   "Add",
					},
					{
						RawByte: []byte{0x01},
						Value:   "01",
					},
					{
						RawByte: []byte{0x02},
						Value:   "02",
					},
				},
			},
			expect: []byte{byte(opcode.Add), 0x01, 0x02},
		},
		{
			asm: translate.Asm{
				AsmCodes: []translate.AsmCode{
					{
						RawByte: []byte{byte(opcode.Mul)},
						Value:   "Mul",
					},
					{
						RawByte: []byte{0x11},
						Value:   "11",
					},
					{
						RawByte: []byte{0x22},
						Value:   "22",
					},
				},
			},
			expect: []byte{byte(opcode.Mul), 0x11, 0x22},
		},
		{
			asm: translate.Asm{
				AsmCodes: []translate.AsmCode{
					{
						RawByte: []byte{byte(opcode.LoadFunc)},
						Value:   "LoadFunc",
					},
					{
						RawByte: []byte{0xff, 0xff},
						Value:   "ffff",
					},
					{
						RawByte: []byte{0xff, 0xff},
						Value:   "ffff",
					},
				},
			},
			expect: []byte{byte(opcode.LoadFunc), 0xff, 0xff, 0xff, 0xff},
		},
		{
			asm: translate.Asm{
				AsmCodes: []translate.AsmCode{
					{
						RawByte: []byte{byte(opcode.LoadFunc)},
						Value:   "LoadFunc",
					},
					{
						RawByte: []byte{0x12, 0x34, 0x56, 0x78, 0x9a},
						Value:   "123456789a",
					},
					{
						RawByte: []byte{0xab, 0xcd, 0xef},
						Value:   "abcdef",
					},
					{
						RawByte: []byte{byte(opcode.Mul)},
						Value:   "Mul",
					},
					{
						RawByte: []byte{0x11},
						Value:   "11",
					},
					{
						RawByte: []byte{0x22},
						Value:   "22",
					},
				},
			},
			expect: []byte{
				byte(opcode.LoadFunc), 0x12, 0x34, 0x56, 0x78, 0x9a, 0xab, 0xcd, 0xef,
				byte(opcode.Mul), 0x11, 0x22,
			},
		},
	}

	for i, tt := range tests {
		result := tt.asm.ToRawByteCode()

		if !bytes.Equal(result, tt.expect) {
			t.Errorf("test[%d] - got unexpected result from ToRawByteCode(), expected=%v, got=%v",
				i, tt.expect, result)
		}
	}
}
