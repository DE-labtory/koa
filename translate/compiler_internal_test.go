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
	"github.com/DE-labtory/koa/opcode"
)

type expressionCompileTestCase struct {
	expression ast.Expression
	expected   Bytecode
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

// TODO: implement test cases :-)
func TestCompileReturnStatement(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompileIfStatement(t *testing.T) {

}

func TestCompileBlockStatement(t *testing.T) {
	statements := makeTempStatements()

	tests := []struct {
		statements *ast.BlockStatement
		expected   Bytecode
		err        error
	}{
		{
			statements: &ast.BlockStatement{
				Statements: statements,
			},
			expected: Bytecode{
				RawByte: []byte{byte(opcode.Push), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xd2, byte(opcode.Push), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
				AsmCode: []string{"Push", "00000000000004d2", "Push", "0000000000000001"},
			},
			err: nil,
		},
	}

	for i, test := range tests {
		b := &Bytecode{
			RawByte: make([]byte, 0),
			AsmCode: make([]string, 0),
		}

		err := compileBlockStatement(test.statements, b)

		if err != nil && err != test.err {
			t.Fatalf("test[%d] - TestCompileBlockStatement() error wrong. expected=%v, got=%v", i, test.err, err)
		}

		if !bytes.Equal(b.RawByte, test.expected.RawByte) {
			t.Fatalf("test[%d] - TestCompileBlockStatement() result wrong for RawByte.\nexpected=%x, got=%x", i, test.expected.RawByte, b.RawByte)
		}

		for j, expected := range test.expected.AsmCode {
			if expected != b.AsmCode[j] {
				t.Fatalf("test[%d] - TestCompileBlockStatement() result wrong for RawByte.\nexpected=%v, got=%v", i, test.expected.AsmCode, b.AsmCode)
			}
		}
	}
}

func TestCompileExpressionStatement(t *testing.T) {
	tests := []struct {
		statement *ast.ExpressionStatement
		expected  Bytecode
		err       error
	}{
		{
			statement: &ast.ExpressionStatement{
				Expr: &ast.IntegerLiteral{
					Value: 12345678,
				},
			},
			expected: Bytecode{
				RawByte: []byte{byte(opcode.Push), 0x00, 0x00, 0x00, 0x00, 0x00, 0xbc, 0x61, 0x4e},
				AsmCode: []string{"Push", "0000000000bc614e"},
			},
			err: nil,
		},
		{
			statement: &ast.ExpressionStatement{
				Expr: &ast.BooleanLiteral{
					Value: true,
				},
			},
			expected: Bytecode{
				RawByte: []byte{byte(opcode.Push), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
				AsmCode: []string{"Push", "0000000000000001"},
			},
			err: nil,
		},
	}

	for i, test := range tests {
		b := &Bytecode{
			RawByte: make([]byte, 0),
			AsmCode: make([]string, 0),
		}

		err := compileExpressionStatement(test.statement, b)

		if err != nil && err != test.err {
			t.Fatalf("test[%d] - TestCompileExpressionStatement() error wrong. expected=%v, got=%v", i, test.err, err)
		}

		if !bytes.Equal(b.RawByte, test.expected.RawByte) {
			t.Fatalf("test[%d] - TestCompileExpressionStatement() result wrong for RawByte.\nexpected=%x, got=%x", i, test.expected.RawByte, b.RawByte)
		}

		for j, expected := range test.expected.AsmCode {
			if expected != b.AsmCode[j] {
				t.Fatalf("test[%d] - TestCompileExpressionStatement() result wrong for RawByte.\nexpected=%v, got=%v", i, test.expected.AsmCode, b.AsmCode)
			}
		}
	}

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

func TestCompilePrefixExpression(t *testing.T) {
	tests := []expressionCompileTestCase{
		// simple prefix expression case
		{
			expression: &ast.PrefixExpression{
				Operator: ast.Bang,
				Right:    &ast.BooleanLiteral{Value: true},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x15,
				},
				AsmCode: []string{"Push", "0000000000000001", "NOT"},
			},
		},
		{
			expression: &ast.PrefixExpression{
				Operator: ast.Minus,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x16,
				},
				AsmCode: []string{"Push", "0000000000000002", "Minus"},
			},
		},
		{
			expression: &ast.PrefixExpression{
				Operator: ast.Minus,
				Right:    &ast.BooleanLiteral{Value: true},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x16,
				},
				AsmCode: []string{"Push", "0000000000000001", "Minus"},
			},
		},
		// rather complex cases
		{
			expression: &ast.PrefixExpression{
				Operator: ast.Minus,
				Right: &ast.PrefixExpression{
					Operator: ast.Minus,
					Right: &ast.IntegerLiteral{
						Value: 1,
					},
				},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x16,
					0x16,
				},
				AsmCode: []string{"Push", "0000000000000001", "Minus", "Minus"},
			},
		},
		{
			expression: &ast.PrefixExpression{
				Operator: ast.Bang,
				Right: &ast.PrefixExpression{
					Operator: ast.Bang,
					Right: &ast.BooleanLiteral{
						Value: false,
					},
				},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x15,
					0x15,
				},
				AsmCode: []string{"Push", "0000000000000000", "NOT", "NOT"},
			},
		},
	}

	runExpressionCompileTests(t, tests)
}

func TestCompileIntegerLiteral(t *testing.T) {
	tests := []expressionCompileTestCase{
		{
			expression: &ast.IntegerLiteral{
				Value: 10,
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a},
				AsmCode: []string{"Push", "000000000000000a"},
			},
		},
		{
			expression: &ast.IntegerLiteral{
				Value: 20,
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14},
				AsmCode: []string{"Push", "0000000000000014"},
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
				RawByte: []byte{0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				AsmCode: []string{"Push", "0000000000000000"},
			},
		},
		{
			expression: &ast.BooleanLiteral{
				Value: true,
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
				AsmCode: []string{"Push", "0000000000000001"},
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
		case *ast.PrefixExpression:
			testFuncName = "compilePrefixExpression()"
			err = compilePrefixExpression(expr, bytecode)
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

func makeTempStatements() []ast.Statement {
	statements := make([]ast.Statement, 0)
	statements = append(statements, &ast.ExpressionStatement{
		Expr: &ast.IntegerLiteral{
			Value: 1234,
		},
	})
	statements = append(statements, &ast.ExpressionStatement{
		Expr: &ast.BooleanLiteral{
			Value: true,
		},
	})

	return statements
}
