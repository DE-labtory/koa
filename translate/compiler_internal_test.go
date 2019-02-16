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

func TestCompileAssignStatement(t *testing.T) {
	tests := []struct {
		statement *ast.AssignStatement
		expected  *Bytecode
	}{
		{
			// int a = true
			statement: &ast.AssignStatement{
				Value: &ast.BooleanLiteral{
					Value: true,
				},
				Variable: ast.Identifier{
					Value: "a",
				},
				Type: ast.BoolType,
			},
			expected: &Bytecode{
				RawByte: []byte{
					byte(opcode.Push), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					byte(opcode.Push), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08,
					byte(opcode.Push), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					byte(opcode.Mstore),
				},
				AsmCode: []string{
					"Push", "0000000000000001",
					"Push", "0000000000000008",
					"Push", "0000000000000000",
					"Mstore",
				},
			},
		},
		{
			// int sum = 5
			statement: &ast.AssignStatement{
				Value: &ast.IntegerLiteral{
					Value: 5,
				},
				Variable: ast.Identifier{
					Value: "sum",
				},
				Type: ast.IntType,
			},
			expected: &Bytecode{
				RawByte: []byte{
					byte(opcode.Push), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05,
					byte(opcode.Push), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08,
					byte(opcode.Push), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					byte(opcode.Mstore),
				},
				AsmCode: []string{
					"Push", "0000000000000005",
					"Push", "0000000000000008",
					"Push", "0000000000000000",
					"Mstore",
				},
			},
		},
	}

	for i, test := range tests {
		b := &Bytecode{
			RawByte: make([]byte, 0),
			AsmCode: make([]string, 0),
		}

		memTracer := NewMemEntryTable()

		err := compileAssignStatement(test.statement, b, memTracer)
		if err != nil {
			t.Fatalf("test[%d] - compileAssignStatement had error. err=%v",
				i, err)
		}

		if !compareByteCode(*b, *test.expected) {
			t.Fatalf("test[%d] - result wrong. expected %x, got=%x",
				i, test.expected, b)
		}
	}
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

		memTracer := NewMemEntryTable()

		err := compileBlockStatement(test.statements, b, memTracer)

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
		memTracer := NewMemEntryTable()
		err := compileExpressionStatement(test.statement, b, memTracer)

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

// TODO: after implement compileIdentifier, add test cases for compiling
// TODO: identifier contained infix expression
//
// TestCompileInfixExpression tests compileInfixExpression and test cases
// consists of three parts
//
// 1. test simple expression
// 2. test rather complex expression
// 3. test edge cases
//
func TestCompileInfixExpression(t *testing.T) {
	tests := []expressionCompileTestCase{
		//
		// 1. test simple infix expression
		//
		// Add
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.Plus,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x01,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"Add",
				},
			},
		},
		// Sub
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.Minus,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x03,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"Sub",
				},
			},
		},
		// Mul
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.Asterisk,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x02,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"Mul",
				},
			},
		},
		// Div
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.Slash,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x04,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"Div",
				},
			},
		},
		// Mod
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.Mod,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x05,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"Mod",
				},
			},
		},
		// LT
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.LT,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x10,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"LT",
				},
			},
		},
		// GT
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.GT,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x12,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"GT",
				},
			},
		},
		// LTE
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.LTE,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x11,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"LTE",
				},
			},
		},
		// GTE
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.GTE,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x13,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"GTE",
				},
			},
		},
		// EQ
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.EQ,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x14,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"EQ",
				},
			},
		},
		// NOT_EQ
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.NOT_EQ,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x14,
					0x15,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"EQ",
					"NOT",
				},
			},
		},
		// LAND
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.LAND,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x06,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"And",
				},
			},
		},
		// LOR
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.LOR,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x07,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"Or",
				},
			},
		},
		//
		// 2. test rather complex expression
		//
		// Plus Minus
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.Plus,
				Right: &ast.InfixExpression{
					Left:     &ast.IntegerLiteral{Value: 2},
					Operator: ast.Minus,
					Right:    &ast.IntegerLiteral{Value: 3},
				},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
					0x03,
					0x01,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000002",
					"Push",
					"0000000000000003",
					"Sub",
					"Add",
				},
			},
		},
		// LT Plus
		{
			expression: &ast.InfixExpression{
				Left: &ast.InfixExpression{
					Left:     &ast.IntegerLiteral{Value: 0},
					Operator: ast.Plus,
					Right:    &ast.IntegerLiteral{Value: 1},
				},
				Operator: ast.LT,
				Right: &ast.InfixExpression{
					Left:     &ast.IntegerLiteral{Value: 2},
					Operator: ast.Asterisk,
					Right:    &ast.IntegerLiteral{Value: 3},
				},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
					0x02,
					0x10,
				},
				AsmCode: []string{
					"Push",
					"0000000000000000",
					"Push",
					"0000000000000001",
					"Add",
					"Push",
					"0000000000000002",
					"Push",
					"0000000000000003",
					"Mul",
					"LT",
				},
			},
		},
		//
		// 3. test edge cases - type mismatching, etc
		//
		// Add Integer with Boolean
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.Plus,
				Right:    &ast.BooleanLiteral{Value: true},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x01,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000001",
					"Add",
				},
			},
		},
		// EQ Integer with Boolean
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.EQ,
				Right:    &ast.BooleanLiteral{Value: true},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x14,
				},
				AsmCode: []string{
					"Push",
					"0000000000000001",
					"Push",
					"0000000000000001",
					"EQ",
				},
			},
		},
		// Mul negative integer
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: -1},
				Operator: ast.Asterisk,
				Right:    &ast.IntegerLiteral{Value: 1},
			},
			expected: Bytecode{
				RawByte: []byte{
					0x21,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0x21,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
					0x02,
				},
				AsmCode: []string{
					"Push",
					"ffffffffffffffff",
					"Push",
					"0000000000000001",
					"Mul",
				},
			},
		},
	}

	runExpressionCompileTests(t, tests)
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

func TestCompileIntegerLiteral_negative(t *testing.T) {
	tests := []expressionCompileTestCase{
		{
			expression: &ast.IntegerLiteral{
				Value: -10,
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF6},
				AsmCode: []string{"Push", "fffffffffffffff6"},
			},
		},
		{
			expression: &ast.IntegerLiteral{
				Value: -20,
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xEC},
				AsmCode: []string{"Push", "ffffffffffffffec"},
			},
		},
	}

	runExpressionCompileTests(t, tests)
}

// TODO: implement test cases :-)
func TestCompileStringLiteral(t *testing.T) {
	tests := []expressionCompileTestCase{
		{
			expression: &ast.StringLiteral{
				Value: "string",
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x00, 0x00, // Push 0x737472696e6700
					0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, // Push 0x0000000000000008 (size)
					0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Push 0x0000000000000000 (offset)
					0x23}, // Memory store
				AsmCode: []string{"Push", "737472696e670000",
					"Push", "0000000000000008",
					"Push", "0000000000000000",
					"Mstore"},
			},
		},
		{
			expression: &ast.StringLiteral{
				Value: "12345678",
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, // Push 0x3132333435363738
					0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, // Push 0x0000000000000008 (size)
					0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Push 0x0000000000000000 (offset)
					0x23}, // Memory store
				AsmCode: []string{"Push", "3132333435363738",
					"Push", "0000000000000008",
					"Push", "0000000000000000",
					"Mstore"},
			},
		},
		{
			expression: &ast.StringLiteral{
				Value: "ssssssss1111111",
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, // Push 0x7373737373737373
					0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, // Push 0x0000000000000008 (size)
					0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Push 0x0000000000000000 (offset)
					0x21, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x00, // Push 0x3131313131313100
					0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, // Push 0x0000000000000008 (size)
					0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Push 0x0000000000000000 (offset)},
					0x23, // Memory store
				},
				AsmCode: []string{"Push", "7373737373737373",
					"Push", "0000000000000008",
					"Push", "0000000000000000",
					"Push", "3131313131313100",
					"Push", "0000000000000008",
					"Push", "0000000000000000",
					"Mstore"},
			},
		},
		{
			expression: &ast.StringLiteral{
				Value: "ssssssss11111111",
			},
			expected: Bytecode{
				RawByte: []byte{0x21, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, 0x73, // Push 0x7373737373737373
					0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, // Push 0x0000000000000008 (size)
					0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Push 0x0000000000000000 (offset)
					0x21, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, // Push 0x3131313131313131
					0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, // Push 0x0000000000000008 (size)
					0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Push 0x0000000000000000 (offset)},
					0x23, // Memory store
				},
				AsmCode: []string{"Push", "7373737373737373",
					"Push", "0000000000000008",
					"Push", "0000000000000000",
					"Push", "3131313131313131",
					"Push", "0000000000000008",
					"Push", "0000000000000000",
					"Mstore"},
			},
		},
	}

	runExpressionCompileTests(t, tests)
}

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
		memTracer := NewMemEntryTable()

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
		case *ast.StringLiteral:
			testFuncName = "compileStringLiteral()"
			err = compileStringLiteral(expr, bytecode, memTracer)
		case *ast.PrefixExpression:
			testFuncName = "compilePrefixExpression()"
			err = compilePrefixExpression(expr, bytecode, memTracer)
		case *ast.InfixExpression:
			testFuncName = "compileInfixExpression()"
			err = compileInfixExpression(expr, bytecode, memTracer)
		default:
			t.Fatalf("%T type not support, abort.", expr)
			t.FailNow()
		}

		if err != nil {
			t.Fatalf("test[%d] - %s had error. err=%v",
				i, testFuncName, err)
		}

		if !compareByteCode(*bytecode, test.expected) {
			t.Fatalf("test[%d] - %s result wrong. "+
				"\nexpected %x, "+
				"\ngot %x",
				i, testFuncName, test.expected, bytecode)
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

func compareByteCode(b1 Bytecode, b2 Bytecode) bool {

	if !bytes.Equal(b1.RawByte, b2.RawByte) {
		return false
	}

	if !reflect.DeepEqual(b1.AsmCode, b2.AsmCode) {
		return false
	}

	return true
}
