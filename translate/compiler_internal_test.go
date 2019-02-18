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
	"errors"
	"testing"

	"github.com/DE-labtory/koa/abi"
	"github.com/DE-labtory/koa/ast"
	"github.com/DE-labtory/koa/opcode"
)

type setupTracer func() MemTracer

func defaultSetupTracer() MemTracer {
	return NewMemEntryTable()
}

type expressionCompileTestCase struct {
	setupTracer
	expression  ast.Expression
	expected    Asm
	expectedErr error
}

type statementCompileTestCase struct {
	setupTracer
	statement   ast.Statement
	expected    Asm
	expectedErr error
}

func TestCreateFuncJmprPlaceholder(t *testing.T) {
	tests := []struct {
		contract      ast.Contract
		expectAsm     *Asm
		expectFuncMap FuncMap
		err           error
	}{
		{
			contract: ast.Contract{
				Functions: []*ast.FunctionLiteral{
					{
						Name: &ast.Identifier{
							Value: "foo",
						},
					},
				},
			},
			expectAsm: &Asm{
				AsmCodes: []AsmCode{
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					// LoadFunc
					{
						RawByte: []byte{0x24},
						Value:   "LoadFunc",
					},
					// DUP
					{
						RawByte: []byte{0x30},
						Value:   "DUP",
					},
					// Push c5d2460100000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0xc5, 0xd2, 0x46, 0x01, 0x00, 0x00, 0x00, 0x00},
						Value:   "c5d2460100000000",
					},
					// EQ
					{
						RawByte: []byte{0x14},
						Value:   "EQ",
					},
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					// Jumpi
					{
						RawByte: []byte{0x29},
						Value:   "Jumpi",
					},
					// Returning
					{
						RawByte: []byte{0x26},
						Value:   "Returning",
					},
				},
			},
			expectFuncMap: FuncMap{
				string(abi.Selector("FuncJmpr")): 3,
				string(abi.Selector("Revert")):   10,
			},
			err: nil,
		},
	}

	for i, test := range tests {
		asm := &Asm{
			AsmCodes: []AsmCode{},
		}
		funcMap := FuncMap{}

		err := createFuncJmprPlaceholder(test.contract, asm, funcMap)

		if !asm.Equal(*test.expectAsm) {
			t.Fatalf("test[%d] - createFuncJmprPlaceholder() bytecode result wrong.\nexpected=%v,\ngot=%v", i, test.expectAsm, asm)
		}

		if !compareFuncMap(funcMap, test.expectFuncMap) {
			t.Fatalf("test[%d] - createFuncJmprPlaceholder() FuncMap result wrong.\nexpected=%v,\ngot=%v", i, test.expectFuncMap, funcMap)
		}

		if err != nil && err != test.err {
			t.Fatalf("test[%d] - createFuncJmprPlaceholder() error wrong.\nexpected=%v,\ngot=%v", i, test.err, err)
		}
	}
}

func TestCompileProgramEndPoint(t *testing.T) {
	tests := []struct {
		revertDst int
		expect    Asm
		err       error
	}{
		{
			revertDst: 0,
			expect: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
				},
			},
			err: nil,
		},
		{
			revertDst: 1234,
			expect: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xd2},
						Value:   "00000000000004d2",
					},
				},
			},
			err: nil,
		},
	}

	for i, test := range tests {
		asm := &Asm{
			AsmCodes: make([]AsmCode, 0),
		}

		err := compileProgramEndPoint(asm, test.revertDst)

		if !asm.Equal(test.expect) {
			t.Fatalf("test[%d] - compileProgramEndPoint() result wrong. expected=%v, got=%v", i, test.expect, asm)
		}

		if err != nil && err != test.err {
			t.Fatalf("test[%d] - compileProgramEndPoint() error wrong. expected=%v, got=%v", i, test.err, err)
		}
	}
}

func TestCompileRevert(t *testing.T) {
	tests := []struct {
		expect Asm
	}{
		{
			expect: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{0x26},
						Value:   "Returning",
					},
				},
			},
		},
	}

	asm := &Asm{
		AsmCodes: make([]AsmCode, 0),
	}

	for i, test := range tests {
		compileRevert(asm)

		if !asm.Equal(test.expect) {
			t.Fatalf("test[%d] - compileRevert() result wrong. expected=%v, got=%v", i, test.expect, asm)
		}
	}
}

func TestGenerateFuncJmpr(t *testing.T) {
	tests := []struct {
		contract  ast.Contract
		asm       *Asm
		expectAsm *Asm
		funcMap   FuncMap
		err       error
	}{
		{
			contract: ast.Contract{
				Functions: []*ast.FunctionLiteral{
					{
						Name: &ast.Identifier{
							Value: "foo",
						},
					},
					{
						Name: &ast.Identifier{
							Value: "sam",
						},
						Parameters: nil,
					},
				},
			},
			asm: &Asm{
				AsmCodes: []AsmCode{
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					// LoadFunc
					{
						RawByte: []byte{0x24},
						Value:   "LoadFunc",
					},
					// DUP
					{
						RawByte: []byte{0x30},
						Value:   "DUP",
					},
					// Push c5d2460100000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0xc5, 0xd2, 0x46, 0x01, 0x00, 0x00, 0x00, 0x00},
						Value:   "c5d2460100000000",
					},
					// EQ
					{
						RawByte: []byte{0x14},
						Value:   "EQ",
					},
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					// Jumpi
					{
						RawByte: []byte{0x29},
						Value:   "Jumpi",
					},
					// DUP
					{
						RawByte: []byte{0x30},
						Value:   "DUP",
					},
					// Push c5d2460100000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0xc5, 0xd2, 0x46, 0x01, 0x00, 0x00, 0x00, 0x00},
						Value:   "c5d2460100000000",
					},
					// EQ
					{
						RawByte: []byte{0x14},
						Value:   "EQ",
					},
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					// Jumpi
					{
						RawByte: []byte{0x29},
						Value:   "Jumpi",
					},
					// Returning
					{
						RawByte: []byte{0x26},
						Value:   "Returning",
					},
				},
			},
			expectAsm: &Asm{
				AsmCodes: []AsmCode{
					// Push 000000000000000a
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a},
						Value:   "000000000000000a",
					},
					// LoadFunc
					{
						RawByte: []byte{0x24},
						Value:   "LoadFunc",
					},
					// DUP
					{
						RawByte: []byte{0x30},
						Value:   "DUP",
					},
					// Push 1b24aabc00000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x1b, 0x24, 0xaa, 0xbc, 0x00, 0x00, 0x00, 0x00},
						Value:   "1b24aabc00000000",
					},
					// EQ
					{
						RawByte: []byte{0x14},
						Value:   "EQ",
					},
					// Push 0000000000000013
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x13},
						Value:   "0000000000000013",
					},
					// Jumpi
					{
						RawByte: []byte{0x29},
						Value:   "Jumpi",
					},
					// DUP
					{
						RawByte: []byte{0x30},
						Value:   "DUP",
					},
					// Push 9f24b46700000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x9f, 0x24, 0xb4, 0x67, 0x00, 0x00, 0x00, 0x00},
						Value:   "9f24b46700000000",
					},
					// EQ
					{
						RawByte: []byte{0x14},
						Value:   "EQ",
					},
					// Push 0000000000000018
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18},
						Value:   "0000000000000018",
					},
					// Jumpi
					{
						RawByte: []byte{0x29},
						Value:   "Jumpi",
					},
					// Returning
					{
						RawByte: []byte{0x26},
						Value:   "Returning",
					},
				},
			},
			funcMap: FuncMap{
				string(abi.Selector("FuncJmpr")):   3,
				string(abi.Selector("Revert")):     10,
				string(abi.Selector("func foo()")): 19,
				string(abi.Selector("func sam()")): 24,
			},
			err: nil,
		},
	}

	for i, test := range tests {
		err := generateFuncJmpr(test.contract, test.asm, test.funcMap)

		if !test.asm.Equal(*test.expectAsm) {
			t.Fatalf("test[%d] - generateFuncJmpr() bytecode result wrong.\nexpected=%v,\ngot=%v", i, test.expectAsm, test.asm)
		}

		if err != nil && err != test.err {
			t.Fatalf("test[%d] - generateFuncJmpr() error wrong.\nexpected=%v,\ngot=%v", i, test.err, err)
		}
	}
}

func TestFillFuncJmpr(t *testing.T) {
	tests := []struct {
		asm      *Asm
		funcJmpr *Asm
		expect   *Asm
		err      error
	}{
		{
			asm: &Asm{
				AsmCodes: []AsmCode{
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					// LoadFunc
					{
						RawByte: []byte{0x24},
						Value:   "LoadFunc",
					},
					// DUP
					{
						RawByte: []byte{0x30},
						Value:   "DUP",
					},
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					// EQ
					{
						RawByte: []byte{0x14},
						Value:   "EQ",
					},
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					// Jumpi
					{
						RawByte: []byte{0x29},
						Value:   "Jumpi",
					},
					// Returning
					{
						RawByte: []byte{0x26},
						Value:   "Returning",
					},
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
				},
			},
			funcJmpr: &Asm{
				AsmCodes: []AsmCode{
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					// LoadFunc
					{
						RawByte: []byte{0x24},
						Value:   "LoadFunc",
					},
					// DUP
					{
						RawByte: []byte{0x30},
						Value:   "DUP",
					},
					// Push 0000000012345678
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x56, 0x78},
						Value:   "0000000012345678",
					},
					// EQ
					{
						RawByte: []byte{0x14},
						Value:   "EQ",
					},
					// Push 0000000000001234
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34},
						Value:   "0000000000001234",
					},
					// Jumpi
					{
						RawByte: []byte{0x29},
						Value:   "Jumpi",
					},
					// Returning
					{
						RawByte: []byte{0x26},
						Value:   "Returning",
					},
				},
			},
			expect: &Asm{
				AsmCodes: []AsmCode{
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					// LoadFunc
					{
						RawByte: []byte{0x24},
						Value:   "LoadFunc",
					},
					// DUP
					{
						RawByte: []byte{0x30},
						Value:   "DUP",
					},
					// Push 0000000012345678
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x56, 0x78},
						Value:   "0000000012345678",
					},
					// EQ
					{
						RawByte: []byte{0x14},
						Value:   "EQ",
					},
					// Push 0000000000001234
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34},
						Value:   "0000000000001234",
					},
					// Jumpi
					{
						RawByte: []byte{0x29},
						Value:   "Jumpi",
					},
					// Returning
					{
						RawByte: []byte{0x26},
						Value:   "Returning",
					},
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					// Push 0000000000000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
				},
			},
			err: nil,
		},
	}

	for i, test := range tests {
		err := fillFuncJmpr(test.asm, *test.funcJmpr)

		if !test.asm.Equal(*test.expect) {
			t.Fatalf("test[%d] - fillFuncJmpr() result wrong.\nexpected=%v,\ngot=%v", i, test.expect, test.asm)
		}

		if err != nil && err != test.err {
			t.Fatalf("test[%d] - fillFuncJmpr() error wrong.\nexpected=%v,\ngot=%v", i, test.err, err)
		}
	}
}

func TestCompileFuncSel(t *testing.T) {
	tests := []struct {
		funcSel string
		funcDst int
		expect  *Asm
		err     error
	}{
		{
			funcSel: string(abi.Selector("func Foo()")),
			funcDst: 15,
			expect: &Asm{
				AsmCodes: []AsmCode{
					// DUP
					{
						RawByte: []byte{0x30},
						Value:   "DUP",
					},
					// Push e3170de100000000
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0xe3, 0x17, 0x0d, 0xe1, 0x00, 0x00, 0x00, 0x00},
						Value:   "e3170de100000000",
					},
					// EQ
					{
						RawByte: []byte{0x14},
						Value:   "EQ",
					},
					// Push 000000000000000f
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0f},
						Value:   "000000000000000f",
					},
					// Jumpi
					{
						RawByte: []byte{0x29},
						Value:   "Jumpi",
					},
				},
			},
			err: nil,
		},
	}

	for i, test := range tests {
		asm := &Asm{
			AsmCodes: make([]AsmCode, 0),
		}
		err := compileFuncSel(asm, test.funcSel, test.funcDst)

		if !asm.Equal(*test.expect) {
			t.Fatalf("test[%d] - compileFuncSel() result wrong.\nexpected=%v,\ngot=%v", i, test.expect, asm)
		}

		if err != nil && err != test.err {
			t.Fatalf("test[%d] - compileFuncSel() error wrong.\nexpected=%v,\ngot=%v", i, test.err, err)
		}

	}
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
		expected  *Asm
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
			expected: &Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					{
						RawByte: []byte{byte(opcode.Mstore)},
						Value:   "Mstore",
					},
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
			expected: &Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
						Value:   "0000000000000005",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					{
						RawByte: []byte{byte(opcode.Mstore)},
						Value:   "Mstore",
					},
				},
			},
		},
		{
			// string str = "str"
			statement: &ast.AssignStatement{
				Value: &ast.StringLiteral{
					Value: "str",
				},
				Variable: ast.Identifier{
					Value: "str",
				},
				Type: ast.StringType,
			},
			expected: &Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x73, 0x74, 0x72, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "7374720000000000",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					{
						RawByte: []byte{byte(opcode.Mstore)},
						Value:   "Mstore",
					},
				},
			},
		},
		{
			// string str = "abcdefgh"
			statement: &ast.AssignStatement{
				Value: &ast.StringLiteral{
					Value: "abcdefgh",
				},
				Variable: ast.Identifier{
					Value: "str",
				},
				Type: ast.StringType,
			},
			expected: &Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68},
						Value:   "6162636465666768",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					{
						RawByte: []byte{byte(opcode.Mstore)},
						Value:   "Mstore",
					},
				},
			},
		},
	}

	for i, test := range tests {
		a := &Asm{
			AsmCodes: make([]AsmCode, 0),
		}

		memTracer := NewMemEntryTable()

		err := compileAssignStatement(test.statement, a, memTracer)
		if err != nil {
			t.Fatalf("test[%d] - compileAssignStatement had error. err=%v",
				i, err)
		}

		if !a.Equal(*test.expected) {
			t.Fatalf("test[%d] - result wrong. \nexpected %x, \ngot=%x",
				i, test.expected, a)
		}
	}
}

func TestCompileReturnStatement(t *testing.T) {
	tests := []statementCompileTestCase{
		{
			statement: &ast.ReturnStatement{
				ReturnValue: &ast.IntegerLiteral{Value: 1},
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Returning)},
						Value:   "Returning",
					},
				},
			},
		},
		{
			statement: &ast.ReturnStatement{
				ReturnValue: &ast.StringLiteral{Value: "abc"},
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x61, 0x62, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "6162630000000000",
					},
					{
						RawByte: []byte{byte(opcode.Returning)},
						Value:   "Returning",
					},
				},
			},
		},
		// test return void
		{
			statement: &ast.ReturnStatement{},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					{
						RawByte: []byte{byte(opcode.Returning)},
						Value:   "Returning",
					},
				},
			},
		},
	}

	runStatementCompileTests(t, tests)
}

//// TODO: implement test cases :-)
func TestCompileIfStatement(t *testing.T) {
	tests := []struct {
		statement *ast.IfStatement
		expected  Asm
		err       error
	}{
		{
			statement: &ast.IfStatement{
				Condition: &ast.BooleanLiteral{
					Value: true,
				},
				Consequence: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.ExpressionStatement{
							Expr: &ast.IntegerLiteral{
								Value: 12345678,
							},
						},
					},
				},
				Alternative: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.ExpressionStatement{
							Expr: &ast.IntegerLiteral{
								Value: 12345678,
							},
						},
					},
				},
			},

			// [Push 0000000000000001 Push 000000000000000a Jumpi Push 0000000000bc614e Push 000000000000000c Jump Push 0000000000bc614e]
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a},
						Value:   "000000000000000a",
					},
					{
						RawByte: []byte{byte(opcode.Jumpi)},
						Value:   "Jumpi",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0xbc, 0x61, 0x4e},
						Value:   "0000000000bc614e",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c},
						Value:   "000000000000000c",
					},
					{
						RawByte: []byte{byte(opcode.Jump)},
						Value:   "Jump",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0xbc, 0x61, 0x4e},
						Value:   "0000000000bc614e",
					},
				},
			},
			err: nil,
		},
		{
			statement: &ast.IfStatement{
				Condition: &ast.BooleanLiteral{
					Value: true,
				},
				Consequence: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.ExpressionStatement{
							Expr: &ast.IntegerLiteral{
								Value: 12345678,
							},
						},
					},
				},
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07},
						Value:   "0000000000000007",
					},
					{
						RawByte: []byte{byte(opcode.Jumpi)},
						Value:   "Jumpi",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0xbc, 0x61, 0x4e},
						Value:   "0000000000bc614e",
					},
				},
			},
			err: nil,
		},
		{
			statement: &ast.IfStatement{
				Condition: &ast.BooleanLiteral{
					Value: true,
				},
				Consequence: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.IfStatement{
							Condition: &ast.BooleanLiteral{
								Value: true,
							},
							Consequence: &ast.BlockStatement{
								Statements: []ast.Statement{
									&ast.ExpressionStatement{
										Expr: &ast.IntegerLiteral{
											Value: 12345678,
										},
									},
								},
							},
						},
					},
				},
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c},
						Value:   "000000000000000c",
					},
					{
						RawByte: []byte{byte(opcode.Jumpi)},
						Value:   "Jumpi",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c},
						Value:   "000000000000000c",
					},
					{
						RawByte: []byte{byte(opcode.Jumpi)},
						Value:   "Jumpi",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0xbc, 0x61, 0x4e},
						Value:   "0000000000bc614e",
					},
				},
			},
			err: nil,
		},
	}

	for i, test := range tests {
		asm := &Asm{
			AsmCodes: make([]AsmCode, 0),
		}

		memTracer := NewMemEntryTable()
		err := compileIfStatement(test.statement, asm, memTracer)
		if err != nil && err != test.err {
			t.Fatalf("test[%d] - TestCompileIfStatement() error wrong. expected=%v, got=%v", i, test.err, err)
		}

		if !asm.Equal(test.expected) {
			t.Fatalf("test[%d] - result wrong. \n expected %x, \n got=%x",
				i, test.expected, asm)
		}
	}
}

func TestCompileBlockStatement(t *testing.T) {
	statements := makeTempStatements()

	tests := []struct {
		statements *ast.BlockStatement
		expected   Asm
		err        error
	}{
		{
			statements: &ast.BlockStatement{
				Statements: statements,
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xd2},
						Value:   "00000000000004d2",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
				},
			},
			err: nil,
		},
	}

	for i, test := range tests {
		a := &Asm{
			AsmCodes: make([]AsmCode, 0),
		}

		memTracer := NewMemEntryTable()

		err := compileBlockStatement(test.statements, a, memTracer)

		if err != nil && err != test.err {
			t.Fatalf("test[%d] - TestCompileBlockStatement() error wrong. expected=%v, got=%v", i, test.err, err)
		}

		if !a.Equal(test.expected) {
			t.Fatalf("test[%d] - result wrong. expected %x, got=%x",
				i, test.expected, a)
		}
	}
}

func TestCompileExpressionStatement(t *testing.T) {
	tests := []struct {
		setupTracer
		statement *ast.ExpressionStatement
		expected  Asm
		err       error
	}{
		{
			setupTracer: defaultSetupTracer,
			statement: &ast.ExpressionStatement{
				Expr: &ast.IntegerLiteral{
					Value: 12345678,
				},
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0xbc, 0x61, 0x4e},
						Value:   "0000000000bc614e",
					},
				},
			},
			err: nil,
		},
		{
			setupTracer: defaultSetupTracer,
			statement: &ast.ExpressionStatement{
				Expr: &ast.BooleanLiteral{
					Value: true,
				},
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
				},
			},
			err: nil,
		},
	}

	for i, test := range tests {
		a := &Asm{
			AsmCodes: make([]AsmCode, 0),
		}

		err := compileExpressionStatement(test.statement, a, test.setupTracer())
		if err != nil && err != test.err {
			t.Fatalf("test[%d] - TestCompileExpressionStatement() error wrong. expected=%v, got=%v", i, test.err, err)
		}

		if !a.Equal(test.expected) {
			t.Fatalf("test[%d] - result wrong. expected %x, got=%x",
				i, test.expected, a)
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.Add)},
						Value:   "Add",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.Sub)},
						Value:   "Sub",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.Mul)},
						Value:   "Mul",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.Div)},
						Value:   "Div",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.Mod)},
						Value:   "Mod",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.LT)},
						Value:   "LT",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.GT)},
						Value:   "GT",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.LTE)},
						Value:   "LTE",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.GTE)},
						Value:   "GTE",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.EQ)},
						Value:   "EQ",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.EQ)},
						Value:   "EQ",
					},
					{
						RawByte: []byte{byte(opcode.NOT)},
						Value:   "NOT",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.And)},
						Value:   "And",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.Or)},
						Value:   "Or",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03},
						Value:   "0000000000000003",
					},
					{
						RawByte: []byte{byte(opcode.Sub)},
						Value:   "Sub",
					},
					{
						RawByte: []byte{byte(opcode.Add)},
						Value:   "Add",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Add)},
						Value:   "Add",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03},
						Value:   "0000000000000003",
					},
					{
						RawByte: []byte{byte(opcode.Mul)},
						Value:   "Mul",
					},
					{
						RawByte: []byte{byte(opcode.LT)},
						Value:   "LT",
					},
				},
			},
		},

		// 3. test edge cases - type mismatching, etc

		// Add Integer with Boolean
		{
			expression: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: ast.Plus,
				Right:    &ast.BooleanLiteral{Value: true},
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Add)},
						Value:   "Add",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.EQ)},
						Value:   "EQ",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
						Value:   "ffffffffffffffff",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Mul)},
						Value:   "Mul",
					},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.NOT)},
						Value:   "NOT",
					},
				},
			},
		},
		{
			expression: &ast.PrefixExpression{
				Operator: ast.Minus,
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.Minus)},
						Value:   "Minus",
					},
				},
			},
		},
		{
			expression: &ast.PrefixExpression{
				Operator: ast.Minus,
				Right:    &ast.BooleanLiteral{Value: true},
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Minus)},
						Value:   "Minus",
					},
				},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Minus)},
						Value:   "Minus",
					},
					{
						RawByte: []byte{byte(opcode.Minus)},
						Value:   "Minus",
					},
				},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					{
						RawByte: []byte{byte(opcode.NOT)},
						Value:   "NOT",
					},
					{
						RawByte: []byte{byte(opcode.NOT)},
						Value:   "NOT",
					},
				},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a},
						Value:   "000000000000000a",
					},
				},
			},
		},
		{
			expression: &ast.IntegerLiteral{
				Value: 20,
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14},
						Value:   "0000000000000014",
					},
				},
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF6},
						Value:   "fffffffffffffff6",
					},
				},
			},
		},
		{
			expression: &ast.IntegerLiteral{
				Value: -20,
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xEC},
						Value:   "ffffffffffffffec",
					},
				},
			},
		},
	}

	runExpressionCompileTests(t, tests)
}

func TestCompileStringLiteral(t *testing.T) {
	tests := []expressionCompileTestCase{
		{
			expression: &ast.StringLiteral{
				Value: "a",
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x61, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "6100000000000000",
					},
				},
			},
		},
		{
			expression: &ast.StringLiteral{
				Value: "ab",
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x61, 0x62, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "6162000000000000",
					},
				},
			},
		},
		{
			expression: &ast.StringLiteral{
				Value: "ab,c",
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x61, 0x62, 0x2c, 0x63, 0x00, 0x00, 0x00, 0x00},
						Value:   "61622c6300000000",
					},
				},
			},
		},
		{
			expression: &ast.StringLiteral{
				Value: "ababababababababababababababababab",
			},
			expected: Asm{
				AsmCodes: []AsmCode{},
			},
			expectedErr: errors.New("Length of string must shorter than 8"),
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
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
				},
			},
		},
		{
			expression: &ast.BooleanLiteral{
				Value: true,
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
				},
			},
		},
	}

	runExpressionCompileTests(t, tests)
}

func TestCompileIdentifier(t *testing.T) {
	tests := []expressionCompileTestCase{
		{
			setupTracer: func() MemTracer {
				tracer := NewMemEntryTable()
				tracer.Define("a")
				return tracer
			},
			expression: &ast.Identifier{
				Value: "a",
			},
			expected: Asm{
				AsmCodes: []AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					{
						RawByte: []byte{byte(opcode.Mload)},
						Value:   "Mload",
					},
				},
			},
		},
		{
			setupTracer: defaultSetupTracer,
			expression: &ast.Identifier{
				Value: "a",
			},
			expected: Asm{
				AsmCodes: []AsmCode{},
			},
			expectedErr: EntryError{Id: "a"},
		},
	}

	runExpressionCompileTests(t, tests)
}

// TODO: implement test cases :-)
func TestCompileParameterLiteral(t *testing.T) {

}

func runExpressionCompileTests(t *testing.T, tests []expressionCompileTestCase) {
	t.Helper()

	for i, test := range tests {
		asm := &Asm{
			AsmCodes: make([]AsmCode, 0),
		}

		var err error
		var testFuncName string
		var tracer MemTracer

		if test.setupTracer != nil {
			tracer = test.setupTracer()
		}

		// add your test expression here with its function name
		switch expr := test.expression.(type) {
		case *ast.BooleanLiteral:
			testFuncName = "compileBooleanLiteral()"
			err = compilePrimitive(expr.Value, asm)
		case *ast.IntegerLiteral:
			testFuncName = "compileIntegerLiteral()"
			err = compilePrimitive(expr.Value, asm)
		case *ast.StringLiteral:
			testFuncName = "compileStringLiteral()"
			err = compilePrimitive(expr.Value, asm)
		case *ast.PrefixExpression:
			testFuncName = "compilePrefixExpression()"
			err = compilePrefixExpression(expr, asm, tracer)
		case *ast.InfixExpression:
			testFuncName = "compileInfixExpression()"
			err = compileInfixExpression(expr, asm, tracer)
		case *ast.Identifier:
			testFuncName = "compileIdentifier()"
			err = compileIdentifier(expr, asm, tracer)
		default:
			t.Fatalf("%T type not support, abort.", expr)
			t.FailNow()
		}

		if err != nil && err.Error() != test.expectedErr.Error() {
			t.Fatalf("test[%d] - [%s] got unexpected error, expected=%s, got=%s",
				i, testFuncName, test.expectedErr.Error(), err.Error())
		}

		if !asm.Equal(test.expected) {
			t.Fatalf("test[%d] - %s result wrong. \n expected %x, \n got=%x",
				i, testFuncName, test.expected, asm)
		}
	}
}

func runStatementCompileTests(t *testing.T, tests []statementCompileTestCase) {
	t.Helper()

	for i, test := range tests {
		asm := &Asm{
			AsmCodes: make([]AsmCode, 0),
		}

		var err error
		var testFuncName string
		var tracer MemTracer

		if test.setupTracer != nil {
			tracer = test.setupTracer()
		}

		switch stmt := test.statement.(type) {
		case *ast.ReturnStatement:
			testFuncName = "compileReturnStatement()"
			err = compileReturnStatement(stmt, asm, tracer)
		default:
			t.Fatalf("%T type not support, abort.", stmt)
			t.FailNow()
		}

		if err != nil && err.Error() != test.expectedErr.Error() {
			t.Fatalf("test[%d] - [%s] got unexpected error, expected=%s, got=%s",
				i, testFuncName, test.expectedErr.Error(), err.Error())
		}

		if !asm.Equal(test.expected) {
			t.Fatalf("test[%d] - %s result wrong. expected %x, got=%x",
				i, testFuncName, test.expected, asm)
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

func compareFuncMap(funcMap1 FuncMap, funcMap2 FuncMap) bool {
	if len(funcMap1) != len(funcMap2) {
		return false
	}

	for funcSel, item := range funcMap1 {
		if item != funcMap2[funcSel] {
			return false
		}
	}

	return true
}
