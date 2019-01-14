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

// TODO: implement pc position test case :-)
func TestCompiler_Emit(t *testing.T) {
	tests := []struct {
		operator        opcode.Type
		operands        [][]byte
		expectedRawByte []byte
		expectedAsmCode []string
	}{
		{
			operator: opcode.Push,
			operands: [][]byte{
				{0x00, 0x00, 0x00, 01},
			},
			expectedRawByte: []byte{0x21, 0x00, 0x00, 0x00, 0x01},
			expectedAsmCode: []string{"Push", "00000001"},
		},
		{
			operator:        opcode.Pop,
			operands:        nil,
			expectedRawByte: []byte{0x20},
			expectedAsmCode: []string{"Pop"},
		},
		{
			operator: opcode.Add,
			operands: [][]byte{
				{0x00, 0x00, 0x00, 0x01},
				{0x00, 0x00, 0x00, 0x02},
			},
			expectedRawByte: []byte{0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02},
			expectedAsmCode: []string{"Add", "00000001", "00000002"},
		},
	}

	for i, test := range tests {
		bytecode := &translate.Bytecode{
			RawByte: make([]byte, 0),
			AsmCode: make([]string, 0),
		}

		_, err := bytecode.Emit(test.operator, test.operands...)

		if err != nil {
			t.Fatalf("test[%d] - Emit() has error. err=%v", i, err)
		}

		if !bytes.Equal(bytecode.RawByte, test.expectedRawByte) {
			t.Fatalf("test[%d] - Emit() bytecode result is wrong. expected=%x, got=%x", i, test.expectedRawByte, bytecode.RawByte)
		}

		for n, s := range bytecode.AsmCode {
			expected := test.expectedAsmCode[n]
			if s != expected {
				t.Fatalf("test[%d] - Emit() asmcode result is wrong. expected=%s, got=%s", i, expected, s)
			}
		}
	}
}
