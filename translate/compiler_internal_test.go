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
		operator opcode.Type
		operands [][]byte
		expected []byte
	}{
		{
			operator: opcode.Push,
			operands: [][]byte{
				{0x00, 0x00, 0x00, 01},
			},
			expected: []byte{0x21, 0x00, 0x00, 0x00, 0x01},
		},
		{
			operator: opcode.Pop,
			operands: [][]byte{
				{},
			},
			expected: []byte{0x20},
		},
		{
			operator: opcode.Add,
			operands: [][]byte{
				{0x00, 0x00, 0x00, 0x01},
				{0x00, 0x00, 0x00, 0x02},
			},
			expected: []byte{0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02},
		},
	}

	for i, test := range tests {
		b := emit(test.operator, test.operands...)

		if !bytes.Equal(b, test.expected) {
			t.Fatalf("test[%d] - emit() result is wrong. expected=%x, got=%x", i, test.expected, b)
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
		node     ast.BooleanLiteral
		expected []byte
		err      error
	}{
		{
			node: ast.BooleanLiteral{
				Value: true,
			},
			expected: []byte{0x21, 0x00, 0x00, 0x00, 0x01},
			err:      nil,
		},
		{
			node: ast.BooleanLiteral{
				Value: false,
			},
			expected: []byte{0x21, 0x00, 0x00, 0x00, 0x00},
			err:      nil,
		},
	}

	for i, test := range tests {
		n := test.node
		b, err := compileBoolean(n)

		if !bytes.Equal(b, test.expected) {
			t.Fatalf("test[%d] - compileBoolean() result is wrong. expected=%x, got=%x", i, test.expected, b)
		}

		if err != nil {
			t.Fatalf("test[%d] - compileBoolean() had error. err=%v", i, err)
		}
	}
}

// TODO: implement test cases :-)
func TestCompiler_compilePrefixExpression(t *testing.T) {

}
