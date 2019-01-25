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
	"reflect"
	"testing"

	"github.com/DE-labtory/koa/abi"
	"github.com/DE-labtory/koa/ast"
)

// TODO: implement test cases :-)
func TestGenerateFuncJumper(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileAbi(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileFunction(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileStatement(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileAssignStatement(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileReturnStatement(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileIfStatement(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileBlockStatement(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileExpressionStatement(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileFunctionLiteral(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileExpression(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileCallExpression(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileInfixExpression(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompilePrefixExpression(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileIntegerLiteral(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileStringLiteral(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileBooleanLiteral(t *testing.T) {
	tests := []struct {
		expression *ast.BooleanLiteral
		expected   Bytecode
	}{
		{
			expression: &ast.BooleanLiteral{
				Value: false,
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0x00, 0x00, 0x00, 0x00},
				AsmCode: []string{"Push", "00000000"},
				Abi:     abi.ABI{},
			},
		},
		{
			expression: &ast.BooleanLiteral{
				Value: true,
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0x00, 0x00, 0x00, 0x01},
				AsmCode: []string{"Push", "00000001"},
				Abi:     abi.ABI{},
			},
		},
	}

	for i, test := range tests {
		bytecode := &Bytecode{
			RawByte: make([]byte, 0),
			AsmCode: make([]string, 0),
			Abi:     abi.ABI{},
		}
		err := compileBooleanLiteral(test.expression, bytecode)

		if err != nil {
			t.Fatalf("test[%d] - compileBooleanLiteral() had error. err=%v", i, err)
		}

		expectedRawByte := test.expected.RawByte
		resultRawByte := bytecode.RawByte
		if !bytes.Equal(expectedRawByte, resultRawByte) {
			t.Fatalf("test[%d] - compileBooleanLiteral() result wrong. expected %x, got=%x", i, expectedRawByte, resultRawByte)
		}

		expectedAsmCode := test.expected.AsmCode
		resultAsmCode := bytecode.AsmCode
		if !reflect.DeepEqual(expectedAsmCode, resultAsmCode) {
			t.Fatalf("test[%d] - compileBooleanLiteral() result wrong. expected %v, got=%v", i, expectedAsmCode, resultAsmCode)
		}
	}
}

// TODO: implement test cases :-)
func TestCompileIdentifier(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileParameterLiteral(t *testing.T) {

}
