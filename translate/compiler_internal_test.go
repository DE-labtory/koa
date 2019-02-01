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

	"github.com/DE-labtory/koa/ast"
)

type expressionCompileTestCase struct {
	expression ast.Expression
	expected   Bytecode
}

type statementCompileTestCase struct {
	statement ast.Statement
	expected  Bytecode
}

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

// TODO: implement test cases after making other compile functions
func TestCompileReturnStatement(t *testing.T) {
	tests := []statementCompileTestCase{
		{
			statement: &ast.ReturnStatement{},
			expected: Bytecode{
				RawByte: []byte{0x26},
				AsmCode: []string{"Returning, "},
			},
		},
		{
			statement: &ast.ReturnStatement{
				ReturnValue: &ast.IntegerLiteral{Value: 10},
			},
			expected: Bytecode{
				RawByte: []byte{0x26},
				AsmCode: []string{"Returning, 10"},
			},
		},
		{
			statement: &ast.ReturnStatement{
				ReturnValue: &ast.BooleanLiteral{Value: true},
			},
			expected: Bytecode{
				RawByte: []byte{0x26},
				AsmCode: []string{"Returning, 1"},
			},
		},
	}

	for i, test := range tests {
		b := &Bytecode{}
		err := compileReturnStatement(test.statement.(*ast.ReturnStatement), b)

		if b.String() != test.expected.String() {
			t.Fatalf("test[%d] - TestCompileReturnStatement failed.\n"+
				"expected=%s\n"+
				"got=%s\n", i, test.expected.String(), b.String())
		}

		if err != nil {
			t.Fatalf("test[%d] - TestCompileReturnStatement has error.", i)
		}
	}
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

func TestCompileIntegerLiteral(t *testing.T) {
	tests := []expressionCompileTestCase{
		{
			expression: &ast.IntegerLiteral{
				Value: 10,
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0x00, 0x00, 0x00, 0x0a},
				AsmCode: []string{"Push", "0000000a"},
			},
		},
		{
			expression: &ast.IntegerLiteral{
				Value: 20,
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0x00, 0x00, 0x00, 0x14},
				AsmCode: []string{"Push", "00000014"},
			},
		},
	}

	runExpressionCompileTests(t, tests)
}

// TODO: implement test cases :-)
func TestCompileStringLiteral(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileBooleanLiteral(t *testing.T) {
	tests := []expressionCompileTestCase{
		{
			expression: &ast.BooleanLiteral{
				Value: false,
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0x00, 0x00, 0x00, 0x00},
				AsmCode: []string{"Push", "00000000"},
			},
		},
		{
			expression: &ast.BooleanLiteral{
				Value: true,
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0x00, 0x00, 0x00, 0x01},
				AsmCode: []string{"Push", "00000001"},
			},
		},
	}

	runExpressionCompileTests(t, tests)
}

// TODO: implement test cases :-)
func TestCompileIdentifier(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileParameterLiteral(t *testing.T) {

}

func runExpressionCompileTests(t *testing.T, tests []expressionCompileTestCase) {
	for i, test := range tests {
		bytecode := &Bytecode{
			RawByte: make([]byte, 0),
			AsmCode: make([]string, 0),
		}

		var err error
		var testFuncName string

		// add your test expression here with its function name
		switch expr := test.expression.(type) {
		case *ast.BooleanLiteral:
			testFuncName = "compileBooleanLiteral()"
			err = compileBooleanLiteral(expr, bytecode)
		case *ast.IntegerLiteral:
			testFuncName = "compileIntegerLiteral()"
			err = compileIntegerLiteral(expr, bytecode)
		default:
			t.Fatalf("%T type not support, abort.", expr)
			t.FailNow()
		}

		if err != nil {
			t.Fatalf("test[%d] - %s had error. err=%v",
				i, testFuncName, err)
		}

		expectedRawByte := test.expected.RawByte
		resultRawByte := bytecode.RawByte
		if !bytes.Equal(expectedRawByte, resultRawByte) {
			t.Fatalf("test[%d] - %s result wrong. expected %x, got=%x",
				i, testFuncName, expectedRawByte, resultRawByte)
		}

		expectedAsmCode := test.expected.AsmCode
		resultAsmCode := bytecode.AsmCode
		if !reflect.DeepEqual(expectedAsmCode, resultAsmCode) {
			t.Fatalf("test[%d] - %s result wrong. expected %v, got=%v",
				i, testFuncName, expectedAsmCode, resultAsmCode)
		}
	}
}
