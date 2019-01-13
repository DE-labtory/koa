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

package translate

import (
	"bytes"
	"testing"

	"github.com/DE-labtory/koa/ast"
	"github.com/DE-labtory/koa/opcode"
)

func TestCompiler_emit(t *testing.T) {
	tests := []struct {
		operator       opcode.Type
		operands       [][]byte
		expectedByte   []byte
		expectedString []string
	}{
		{
			operator: opcode.Push,
			operands: [][]byte{
				{0x00, 0x00, 0x00, 01},
			},
			expectedByte:   []byte{0x21, 0x00, 0x00, 0x00, 0x01},
			expectedString: []string{"Push", "00000001"},
		},
		{
			operator:       opcode.Pop,
			operands:       nil,
			expectedByte:   []byte{0x20},
			expectedString: []string{"Pop"},
		},
		{
			operator: opcode.Add,
			operands: [][]byte{
				{0x00, 0x00, 0x00, 0x01},
				{0x00, 0x00, 0x00, 0x02},
			},
			expectedByte:   []byte{0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02},
			expectedString: []string{"Add", "00000001", "00000002"},
		},
	}

	for i, test := range tests {
		bytecode := &Bytecode{
			RawByte: make([]byte, 0),
			AsmCode: make([]string, 0),
			PC:      0,
		}

		emit(bytecode, test.operator, test.operands...)

		if !bytes.Equal(bytecode.RawByte, test.expectedByte) {
			t.Fatalf("test[%d] - emit() bytecode result is wrong. expected=%x, got=%x", i, test.expectedByte, bytecode.RawByte)
		}

		for n, s := range bytecode.AsmCode {
			expected := test.expectedString[n]
			if s != expected {
				t.Fatalf("test[%d] - emit() asmcode result is wrong. expected=%s, got=%s", i, expected, s)
			}
		}
	}
}

// TODO: implement test cases :-)
func TestCompiler_compileNode(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompiler_compileIdentifier(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompiler_compileString(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompiler_compileInteger(t *testing.T) {

}

func TestCompiler_compileBoolean(t *testing.T) {
	tests := []struct {
		node           ast.BooleanLiteral
		expectedByte   []byte
		expectedString []string
		err            error
	}{
		{
			node: ast.BooleanLiteral{
				Value: true,
			},
			expectedByte:   []byte{0x21, 0x00, 0x00, 0x00, 0x01},
			expectedString: []string{"Push", "00000001"},
			err:            nil,
		},
		{
			node: ast.BooleanLiteral{
				Value: false,
			},
			expectedByte:   []byte{0x21, 0x00, 0x00, 0x00, 0x00},
			expectedString: []string{"Push", "00000000"},
			err:            nil,
		},
	}

	for i, test := range tests {
		bytecode := &Bytecode{
			RawByte: make([]byte, 0),
			AsmCode: make([]string, 0),
			PC:      0,
		}

		n := test.node
		err := compileBoolean(n, bytecode)

		if !bytes.Equal(bytecode.RawByte, test.expectedByte) {
			t.Fatalf("test[%d] - compileBoolean() bytecode result is wrong. expected=%x, got=%x", i, test.expectedByte, bytecode.RawByte)
		}

		for n, s := range bytecode.AsmCode {
			expected := test.expectedString[n]
			if s != expected {
				t.Fatalf("test[%d] - compileBoolean() asmcode result is wrong. expected=%s, got=%s", i, expected, s)
			}
		}

		if err != nil {
			t.Fatalf("test[%d] - compileBoolean() had error. err=%v", i, err)
		}
	}
}

// TODO: implement test cases :-)
func TestCompiler_compilePrefixExpression(t *testing.T) {

}
