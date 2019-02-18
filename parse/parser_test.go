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

package parse_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/DE-labtory/koa/ast"
	"github.com/DE-labtory/koa/parse"
)

// expectedFnArg is used to verifing parsed function args data
type expectedFnArg struct {
	argIdent string
	argType  ast.DataStructure
}

// expectedFnArg is used to verifing parsed function header data
type expectedFnHeader struct {
	retType ast.DataStructure
	args    []expectedFnArg
}

// fnTmplData represents smart contract function which is going to
// be injected into contractTmpl
type fnTmplData struct {
	FuncName string
	Args     string
	RetType  string
	Stmts    []string
}

// contractTmplData represents smart contract which is going to
// be injected into contractTmpl
type contractTmplData struct {
	Fns []fnTmplData
}

// contractTmpl is template for creating smart contract code
// contractTmplData injects data to this template
const contractTmpl = `
contract {
	{{range .Fns -}}
		func {{.FuncName}}({{.Args}}) {{.RetType}} {
        {{range .Stmts -}}
			{{.}}
        {{end}}
    }
	{{end}}
}
`

// createTestContractCode creates and returns string code, which is made by
// contractTmpl and contractTmplData
func createTestContractCode(c contractTmplData) string {
	out := bytes.NewBufferString("")
	instance, _ := template.New("ContractTemplate").Parse(contractTmpl)
	instance.Execute(out, c)

	return out.String()
}

func parseTestContract(input string) (*ast.Contract, error) {
	l := parse.NewLexer(input)
	buf := parse.NewTokenBuffer(l)
	return parse.Parse(buf)
}

// chkFnHeader verify smart contract's function header
func chkFnHeader(t *testing.T, fn *ast.FunctionLiteral, efh expectedFnHeader) {
	if fn.ReturnType != efh.retType {
		t.Errorf("function has wrong return type Expected=%v, got=%v",
			ast.VoidType.String(), fn.ReturnType.String())
	}

	if len(fn.Parameters) != len(efh.args) {
		t.Errorf("function has wrong parameters length got=%v",
			len(fn.Parameters))
	}

	for i, a := range fn.Parameters {
		testFnParameters(t, a, efh.args[i].argType, efh.args[i].argIdent)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		contractTmpl      contractTmplData
		expectedFnHeaders []expectedFnHeader
		expected          []ast.ReturnStatement
		expectedErr       error
	}{
		/*
			func returnStatement1() int {
				return 1
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "returnStatement1",
						Args:     "",
						RetType:  "int",
						Stmts: []string{
							"return 1",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.IntType,
				},
			},
			expected: []ast.ReturnStatement{
				{ReturnValue: &ast.IntegerLiteral{Value: 1}},
			},
			expectedErr: nil,
		},
		/*
			func returnStatement2(a int) {
				return a
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "returnStatement2",
						Args:     "a int",
						RetType:  "int",
						Stmts: []string{
							"return a",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{"a", ast.IntType},
					},
				},
			},
			expected: []ast.ReturnStatement{
				{
					ReturnValue: &ast.Identifier{Name: "a"},
				},
			},
			expectedErr: nil,
		},
		/*
			func returnStatement3(a int, b int) int {
				return a + b + 1
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "returnStatement3",
						Args:     "a int, b int",
						RetType:  "int",
						Stmts: []string{
							"return a + b + 1",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{"a", ast.IntType},
						{"b", ast.IntType},
					},
				},
			},
			expected: []ast.ReturnStatement{
				{
					ReturnValue: &ast.InfixExpression{
						Left: &ast.InfixExpression{
							Left:     &ast.Identifier{Name: "a"},
							Operator: ast.Plus,
							Right:    &ast.Identifier{Name: "b"},
						},
						Operator: ast.Plus,
						Right:    &ast.IntegerLiteral{Value: 1},
					},
				},
			},
			expectedErr: nil,
		},
		/*
			func returnStatement4(a int) int {
				return (
					a)
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "returnStatement4",
						Args:     "a int",
						RetType:  "int",
						Stmts: []string{
							"return (",
							"a)",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{"a", ast.IntType},
					},
				},
			},
			expected: []ast.ReturnStatement{
				{
					ReturnValue: &ast.Identifier{Name: "a"},
				},
			},
			expectedErr: nil,
		},
		/*
			func returnStatement5(a int) int {
				return (
					a)
			}
			func returnStatement5_1(b int) int {
				return (
					b)
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "returnStatement5",
						Args:     "a int",
						RetType:  "int",
						Stmts: []string{
							"return (",
							"a)",
						},
					},
					{
						FuncName: "returnStatement5_1",
						Args:     "b int",
						RetType:  "int",
						Stmts: []string{
							"return (",
							"b)",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{"a", ast.IntType},
					},
				},
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{"b", ast.IntType},
					},
				},
			},
			expected: []ast.ReturnStatement{
				{
					ReturnValue: &ast.Identifier{Name: "a"},
				},
				{
					ReturnValue: &ast.Identifier{Name: "b"},
				},
			},
			expectedErr: nil,
		},
		/*
			// type mismatch - invalid return type
			func returnStatement6() string {
				return 1
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "returnStatement6",
						Args:     "",
						RetType:  "string",
						Stmts: []string{
							"return 1",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.StringType,
					args:    []expectedFnArg{},
				},
			},
			expected: []ast.ReturnStatement{
				{
					ReturnValue: &ast.IntegerLiteral{Value: 1},
				},
			},
			expectedErr: nil,
		},
		/*
			// symbol resolution - return undefined variable
			func returnStatement7() string {
				return a
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "returnStatement7",
						Args:     "",
						RetType:  "string",
						Stmts: []string{
							"return a",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.StringType,
					args:    []expectedFnArg{},
				},
			},
			expected: []ast.ReturnStatement{
				{
					ReturnValue: &ast.Identifier{Name: "a"},
				},
			},
			expectedErr: parse.NotExistSymError{Source: parse.Token{Type: parse.Ident, Val: "a", Line: 3, Column: 17}},
		},
		/*
			// void return
			func returnStatement7() string {
				return
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "returnStatement7",
						Args:     "",
						RetType:  "string",
						Stmts: []string{
							"return",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.StringType,
					args:    []expectedFnArg{},
				},
			},
			expected: []ast.ReturnStatement{
				{},
			},
		},
	}

	for _, tt := range tests {
		input := createTestContractCode(tt.contractTmpl)
		contract, err := parseTestContract(input)

		if err != nil && err == tt.expectedErr {
			continue
		}

		if err != nil {
			t.Errorf("parser error: %q", err)
			t.FailNow()
		}

		for i, fn := range contract.Functions {
			runReturnStatementTestCases(t, fn, tt.expectedFnHeaders[i], tt.expected[i])
		}
	}
}

func runReturnStatementTestCases(t *testing.T, fn *ast.FunctionLiteral, efhs expectedFnHeader, tt ast.ReturnStatement) {
	t.Logf("test ReturnStatement - [%s]", fn.Name)

	chkFnHeader(t, fn, efhs)

	for _, stmt := range fn.Body.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.ReturnStatement. got=%T", stmt)
		}

		testExpression(t, returnStmt.ReturnValue, tt.ReturnValue)
	}
}

func TestAssignStatement(t *testing.T) {
	tests := []struct {
		contractTmpl      contractTmplData
		expectedFnHeaders []expectedFnHeader
		expected          []ast.AssignStatement
		expectedErr       error
	}{
		/*
			func assignStatement1() {
				int a = 1
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "assignStatement1",
						Args:     "",
						RetType:  "",
						Stmts: []string{
							"int a = 1",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
			},
			expected: []ast.AssignStatement{
				{
					Type:     ast.IntType,
					Variable: ast.Identifier{Name: "a"},
					Value:    &ast.IntegerLiteral{Value: 1},
				},
			},
			expectedErr: nil,
		},
		/*
			func assignStatement2() {
				int a = 1 + 2 * 3
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "assignStatement2",
						Args:     "",
						RetType:  "",
						Stmts: []string{
							"int a = 1 + 2 * 3",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
			},
			expected: []ast.AssignStatement{
				{
					Type:     ast.IntType,
					Variable: ast.Identifier{Name: "a"},
					Value: &ast.InfixExpression{
						Left:     &ast.IntegerLiteral{Value: 1},
						Operator: ast.Plus,
						Right: &ast.InfixExpression{
							Left:     &ast.IntegerLiteral{Value: 2},
							Operator: ast.Asterisk,
							Right:    &ast.IntegerLiteral{Value: 3},
						},
					},
				},
			},
			expectedErr: nil,
		},
		/*
			func assignStatement3() {
				int a =
						1 + 2 * 3
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "assignStatement3",
						Args:     "",
						RetType:  "",
						Stmts: []string{
							"int a = ",
							"		1 + 2 * 3",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
			},
			expected: []ast.AssignStatement{
				{
					Type:     ast.IntType,
					Variable: ast.Identifier{Name: "a"},
					Value: &ast.InfixExpression{
						Left:     &ast.IntegerLiteral{Value: 1},
						Operator: ast.Plus,
						Right: &ast.InfixExpression{
							Left:     &ast.IntegerLiteral{Value: 2},
							Operator: ast.Asterisk,
							Right:    &ast.IntegerLiteral{Value: 3},
						},
					},
				},
			},
			expectedErr: nil,
		},
		/*
			func assignStatement4(foo int) {
				int a = (foo + 2) * 3
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "assignStatement4",
						Args:     "foo int",
						RetType:  "",
						Stmts: []string{
							"int a = ((foo + 2) * 3)",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args: []expectedFnArg{
						{argIdent: "foo", argType: ast.IntType},
					},
				},
			},
			expected: []ast.AssignStatement{
				{
					Type:     ast.IntType,
					Variable: ast.Identifier{Name: "a"},
					Value: &ast.InfixExpression{
						Left: &ast.InfixExpression{
							Left:     &ast.Identifier{Name: "foo"},
							Operator: ast.Plus,
							Right:    &ast.IntegerLiteral{Value: 2},
						},
						Operator: ast.Asterisk,
						Right:    &ast.IntegerLiteral{Value: 3},
					},
				},
			},
			expectedErr: nil,
		},
		/*
			func add(aa int, b int) int {
			}
			func assignStatement5(foo int) {
				int a = add(foo, 1)
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "add",
						Args:     "aa int, b int",
						RetType:  "int",
						Stmts:    []string{},
					},
					{
						FuncName: "assignStatement5",
						Args:     "foo int",
						RetType:  "",
						Stmts: []string{
							"int a = add(foo, 1)",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{argIdent: "aa", argType: ast.IntType},
						{argIdent: "b", argType: ast.IntType},
					},
				},
				{
					retType: ast.VoidType,
					args: []expectedFnArg{
						{argIdent: "foo", argType: ast.IntType},
					},
				},
			},
			expected: []ast.AssignStatement{
				{},
				{
					Type:     ast.IntType,
					Variable: ast.Identifier{Name: "a"},
					Value: &ast.CallExpression{
						Function: &ast.Identifier{Name: "add"},
						Arguments: []ast.Expression{
							&ast.Identifier{Name: "foo"},
							&ast.IntegerLiteral{Value: 1},
						},
					},
				},
			},
		},
		/*
			// as in 0.1.0 version, add symbol in order
			// so by the time parser parsing "int a = add(foo, 1)",
			// foo function still not defined
			//
			func assignStatement6(foo int) {
				int a = add(foo, 1)
			}
			func add(aa int, b int) int {
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "assignStatement6",
						Args:     "foo int",
						RetType:  "",
						Stmts: []string{
							"int a = add(foo, 1)",
						},
					},
					{
						FuncName: "add",
						Args:     "aa int, b int",
						RetType:  "int",
						Stmts:    []string{},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args: []expectedFnArg{
						{argIdent: "foo", argType: ast.IntType},
					},
				},
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{argIdent: "aa", argType: ast.IntType},
						{argIdent: "b", argType: ast.IntType},
					},
				},
			},
			expected: []ast.AssignStatement{
				{
					Type:     ast.IntType,
					Variable: ast.Identifier{Name: "a"},
					Value: &ast.CallExpression{
						Function: &ast.Identifier{Name: "add"},
						Arguments: []ast.Expression{
							&ast.Identifier{Name: "foo"},
							&ast.IntegerLiteral{Value: 1},
						},
					},
				},
				{},
			},
			expectedErr: parse.NotExistSymError{
				Source: parse.Token{Type: parse.Ident, Val: "add", Line: 3, Column: 19}},
		},
		/*
			// type mismatch - missing return type
			func add(aa int, b int) {
			}
			func assignStatement7(foo int) {
				int a = add(foo, 1)
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "add",
						Args:     "aa int, b int",
						RetType:  "",
						Stmts:    []string{},
					},
					{
						FuncName: "assignStatement7",
						Args:     "foo int",
						RetType:  "",
						Stmts: []string{
							"int a = add(foo, 1)",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args: []expectedFnArg{
						{argIdent: "aa", argType: ast.IntType},
						{argIdent: "b", argType: ast.IntType},
					},
				},
				{
					retType: ast.VoidType,
					args: []expectedFnArg{
						{argIdent: "foo", argType: ast.IntType},
					},
				},
			},
			expected: []ast.AssignStatement{
				{},
				{
					Type:     ast.IntType,
					Variable: ast.Identifier{Name: "a"},
					Value: &ast.CallExpression{
						Function: &ast.Identifier{Name: "add"},
						Arguments: []ast.Expression{
							&ast.Identifier{Name: "foo"},
							&ast.IntegerLiteral{Value: 1},
						},
					},
				},
			},
		},
		/*
			// type mismatch - missing parameters
			func add() {
			}
			func assignStatement8(foo int) {
				int a = add(foo, 1)
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "add",
						Args:     "",
						RetType:  "",
						Stmts:    []string{},
					},
					{
						FuncName: "assignStatement8",
						Args:     "foo int",
						RetType:  "",
						Stmts: []string{
							"int a = add(foo, 1)",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
				{
					retType: ast.VoidType,
					args: []expectedFnArg{
						{argIdent: "foo", argType: ast.IntType},
					},
				},
			},
			expected: []ast.AssignStatement{
				{},
				{
					Type:     ast.IntType,
					Variable: ast.Identifier{Name: "a"},
					Value: &ast.CallExpression{
						Function: &ast.Identifier{Name: "add"},
						Arguments: []ast.Expression{
							&ast.Identifier{Name: "foo"},
							&ast.IntegerLiteral{Value: 1},
						},
					},
				},
			},
		},
		/*
			func assignStatement9() {
				string a = "hello, world"
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "assignStatement9",
						Args:     "",
						RetType:  "",
						Stmts: []string{
							`string a = "hello, world"`,
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
			},
			expected: []ast.AssignStatement{
				{
					Type:     ast.StringType,
					Variable: ast.Identifier{Name: "a"},
					Value:    &ast.StringLiteral{Value: "\"hello, world\""},
				},
			},
		},
		/*
			func assignStatement10() {
				string a =
						"hello, world"
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "assignStatement10",
						Args:     "",
						RetType:  "",
						Stmts: []string{
							`string a = `,
							`			"hello, world"`,
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
			},
			expected: []ast.AssignStatement{
				{
					Type:     ast.StringType,
					Variable: ast.Identifier{Name: "a"},
					Value:    &ast.StringLiteral{Value: "\"hello, world\""},
				},
			},
		},
		/*
			func assignStatement11() {
				bool a = true
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "assignStatement11",
						Args:     "",
						RetType:  "",
						Stmts: []string{
							`bool a = true`,
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
			},
			expected: []ast.AssignStatement{
				{
					Type:     ast.BoolType,
					Variable: ast.Identifier{Name: "a"},
					Value:    &ast.BooleanLiteral{Value: true},
				},
			},
		},
	}

	for _, tt := range tests {
		input := createTestContractCode(tt.contractTmpl)

		contract, err := parseTestContract(input)

		if err != nil && err == tt.expectedErr {
			continue
		}

		if err != nil {
			t.Errorf("parser error: %q", err)
			t.FailNow()
		}

		for i, fn := range contract.Functions {
			runAssignStatementTestCases(t, fn, tt.expectedFnHeaders[i], tt.expected[i])
		}

	}
}

func runAssignStatementTestCases(t *testing.T, fn *ast.FunctionLiteral, efh expectedFnHeader, tt ast.AssignStatement) {
	t.Logf("test AssignStatement - [%s]", fn.Name)

	chkFnHeader(t, fn, efh)

	for _, stmt := range fn.Body.Statements {
		assignStmt, ok := stmt.(*ast.AssignStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.ReturnStatement. got=%T", stmt)
		}

		testAssignStatement(t, assignStmt, tt.Type, tt.Variable, tt.Value)
	}
}

func TestIfElseStatement(t *testing.T) {
	tests := []struct {
		contractTmpl      contractTmplData
		expectedFnHeaders []expectedFnHeader
		expected          []ast.IfStatement
		expectedErr       error
	}{
		/*
			func ifStatement1() {
				if (true) {}
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "ifStatement1",
						Args:     "",
						RetType:  "",
						Stmts: []string{
							"if (true) {}",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
			},
			expected: []ast.IfStatement{
				{
					Condition:   &ast.BooleanLiteral{Value: true},
					Consequence: &ast.BlockStatement{},
					Alternative: &ast.BlockStatement{},
				},
			},
			expectedErr: nil,
		},
		/*
			func ifStatement2(a int, b int) int {
				if (1 != 1 + 2) {
					return a
				} else {
					return b
				}
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "ifStatement2",
						Args:     "a int, b int",
						RetType:  "int",
						Stmts: []string{
							"if (1 != 1 + 2) {",
							"	return a",
							"} else {",
							"	return b",
							"}",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{argIdent: "a", argType: ast.IntType},
						{argIdent: "b", argType: ast.IntType},
					},
				},
			},
			expected: []ast.IfStatement{
				{
					Condition: &ast.InfixExpression{
						Left:     &ast.IntegerLiteral{Value: 1},
						Operator: ast.NOT_EQ,
						Right: &ast.InfixExpression{
							Left:     &ast.IntegerLiteral{Value: 1},
							Operator: ast.Plus,
							Right:    &ast.IntegerLiteral{Value: 2},
						},
					},
					Consequence: &ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ReturnStatement{
								ReturnValue: &ast.Identifier{Name: "a"},
							},
						},
					},
					Alternative: &ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ReturnStatement{
								ReturnValue: &ast.Identifier{Name: "b"},
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
		/*
			func ifStatement3(a int, b int) int {
				if (1 < 1 + 2) {
					return a
				} else {
				}
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "ifStatement3",
						Args:     "a int, b int",
						RetType:  "int",
						Stmts: []string{
							"if (1 < 1 + 2) {",
							"	return a",
							"} else {",
							"}",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{argIdent: "a", argType: ast.IntType},
						{argIdent: "b", argType: ast.IntType},
					},
				},
			},
			expected: []ast.IfStatement{
				{
					Condition: &ast.InfixExpression{
						Left:     &ast.IntegerLiteral{Value: 1},
						Operator: ast.LT,
						Right: &ast.InfixExpression{
							Left:     &ast.IntegerLiteral{Value: 1},
							Operator: ast.Plus,
							Right:    &ast.IntegerLiteral{Value: 2},
						},
					},
					Consequence: &ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ReturnStatement{
								ReturnValue: &ast.Identifier{Name: "a"},
							},
						},
					},
					Alternative: &ast.BlockStatement{
						Statements: []ast.Statement{},
					},
				},
			},
			expectedErr: nil,
		},
		/*
			func ifStatement4(a int, b int) int {
				if (1 != 1 + 2) {
				} else {
					return b
				}
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "ifStatement4",
						Args:     "a int, b int",
						RetType:  "int",
						Stmts: []string{
							"if (1 != 1 + 2) {",
							"} else {",
							"	return b",
							"}",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{argIdent: "a", argType: ast.IntType},
						{argIdent: "b", argType: ast.IntType},
					},
				},
			},
			expected: []ast.IfStatement{
				{
					Condition: &ast.InfixExpression{
						Left:     &ast.IntegerLiteral{Value: 1},
						Operator: ast.NOT_EQ,
						Right: &ast.InfixExpression{
							Left:     &ast.IntegerLiteral{Value: 1},
							Operator: ast.Plus,
							Right:    &ast.IntegerLiteral{Value: 2},
						},
					},
					Consequence: &ast.BlockStatement{
						Statements: []ast.Statement{},
					},
					Alternative: &ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ReturnStatement{
								ReturnValue: &ast.Identifier{Name: "b"},
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
		/*
			func ifStatement5(foo bool) {
				if (foo) {
					int a = 1
				} else {
					int a = 2
				}
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "ifStatement5",
						Args:     "foo bool",
						RetType:  "",
						Stmts: []string{
							"if (foo) {",
							"	int a = 1",
							"} else {",
							"	int a = 2",
							"}",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args: []expectedFnArg{
						{argIdent: "foo", argType: ast.BoolType},
					},
				},
			},
			expected: []ast.IfStatement{
				{
					Condition: &ast.Identifier{Name: "foo"},
					Consequence: &ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.AssignStatement{
								Type:     ast.IntType,
								Variable: ast.Identifier{Name: "a"},
								Value:    &ast.IntegerLiteral{Value: 1},
							},
						},
					},
					Alternative: &ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.AssignStatement{
								Type:     ast.IntType,
								Variable: ast.Identifier{Name: "a"},
								Value:    &ast.IntegerLiteral{Value: 2},
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
		/*
			func ifStatement5(foo bool) {
				if (foo) {
					bool a = true
					if (a) {
						int a = 1
					}
				} else {
					int a = 2
				}
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "ifStatement5",
						Args:     "foo bool",
						RetType:  "",
						Stmts: []string{
							"if (foo) {",
							"	bool a = true",
							"	if (a) {",
							"		int a = 1",
							"	}",
							"} else {",
							"	int a = 2",
							"}",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args: []expectedFnArg{
						{argIdent: "foo", argType: ast.BoolType},
					},
				},
			},
			expected: []ast.IfStatement{
				{
					Condition: &ast.Identifier{Name: "foo"},
					Consequence: &ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.IfStatement{
								Condition: &ast.Identifier{Name: "a"},
								Consequence: &ast.BlockStatement{
									Statements: []ast.Statement{
										&ast.AssignStatement{
											Type:     ast.IntType,
											Variable: ast.Identifier{Name: "a"},
											Value:    &ast.IntegerLiteral{Value: 1},
										},
									},
								},
							},
						},
					},
					Alternative: &ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.AssignStatement{
								Type:     ast.IntType,
								Variable: ast.Identifier{Name: "a"},
								Value:    &ast.IntegerLiteral{Value: 2},
							},
						},
					},
				},
			},
			expectedErr: parse.DupSymError{
				Source: parse.Token{Type: parse.Ident, Val: "a", Line: 6, Column: 15}},
		},
	}

	for _, tt := range tests {
		input := createTestContractCode(tt.contractTmpl)
		contract, err := parseTestContract(input)

		if err != nil && err == tt.expectedErr {
			continue
		}

		if err != nil {
			t.Errorf("parser error: %q", err)
			t.FailNow()
		}

		for i, fn := range contract.Functions {
			runIfStatementTestCases(t, fn, tt.expectedFnHeaders[i], tt.expected[i])
		}
	}
}

func runIfStatementTestCases(t *testing.T, fn *ast.FunctionLiteral, efh expectedFnHeader, tt ast.IfStatement) {
	t.Logf("test IfStatement - [%s]", fn.Name)

	chkFnHeader(t, fn, efh)

	for _, stmt := range fn.Body.Statements {
		ifStmt, ok := stmt.(*ast.IfStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.IfStatement. got=%T", stmt)
		}
		testIfStatement(t, ifStmt, tt.Condition, tt.Consequence, tt.Alternative)
	}
}

func TestExpressionStatement(t *testing.T) {
	tests := []struct {
		contractTmpl      contractTmplData
		expectedFnHeaders []expectedFnHeader
		expected          []ast.ExpressionStatement
		expectedErr       error
	}{
		/*
			func add() {
			}
			func expressionStatement1() {
				add()
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "add",
						Args:     "",
						RetType:  "",
						Stmts:    []string{},
					},
					{
						FuncName: "expressionStatement",
						Args:     "",
						RetType:  "",
						Stmts: []string{
							"add()",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
			},
			expected: []ast.ExpressionStatement{
				{},
				{
					Expr: &ast.CallExpression{
						Function:  &ast.Identifier{Name: "add"},
						Arguments: []ast.Expression{},
					},
				},
			},
			expectedErr: nil,
		},
		/*
			func add() {
			}
			func expressionStatement1() {
				add(1, 2)
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "add",
						Args:     "",
						RetType:  "",
						Stmts:    []string{},
					},
					{
						FuncName: "expressionStatement",
						Args:     "",
						RetType:  "",
						Stmts: []string{
							"add(1, 2)",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
			},
			expected: []ast.ExpressionStatement{
				{},
				{
					Expr: &ast.CallExpression{
						Function: &ast.Identifier{Name: "add"},
						Arguments: []ast.Expression{
							&ast.IntegerLiteral{Value: 1},
							&ast.IntegerLiteral{Value: 2},
						},
					},
				},
			},
			expectedErr: nil,
		},
		/*
			func add() {
			}
			func expressionStatement1(foo int) {
				add(foo, 2)
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "add",
						Args:     "",
						RetType:  "",
						Stmts:    []string{},
					},
					{
						FuncName: "expressionStatement",
						Args:     "foo int",
						RetType:  "",
						Stmts: []string{
							"add(foo, 2)",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
				{
					retType: ast.VoidType,
					args: []expectedFnArg{
						{
							argIdent: "foo",
							argType:  ast.IntType,
						},
					},
				},
			},
			expected: []ast.ExpressionStatement{
				{},
				{
					Expr: &ast.CallExpression{
						Function: &ast.Identifier{Name: "add"},
						Arguments: []ast.Expression{
							&ast.Identifier{Name: "foo"},
							&ast.IntegerLiteral{Value: 2},
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		input := createTestContractCode(tt.contractTmpl)
		contract, err := parseTestContract(input)

		if err != nil {
			t.Errorf("parser error: %q", err)
			t.FailNow()
		}

		for i, fn := range contract.Functions {
			runExpressionStatementTestCases(t, fn, tt.expectedFnHeaders[i], tt.expected[i])
		}
	}
}

func runExpressionStatementTestCases(t *testing.T, fn *ast.FunctionLiteral, efh expectedFnHeader, tt ast.ExpressionStatement) {
	t.Logf("test ExpressionStatement - [%s]", fn.Name)

	chkFnHeader(t, fn, efh)

	for _, stmt := range fn.Body.Statements {
		exprStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.ExpressionStatement. got=%T", stmt)
		}
		testExpression(t, exprStmt.Expr, tt.Expr)
	}
}

func TestReassignStatement(t *testing.T) {
	tests := []struct {
		contractTmpl      contractTmplData
		expectedFnHeaders []expectedFnHeader
		expected          []ast.ReassignStatement
		expectedErr       error
	}{
		/*
			func reassignStatement1(a string) {
				a = "hello, world"
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "reassignStatement1",
						Args:     "a string",
						RetType:  "",
						Stmts: []string{
							`a = "hello, world"`,
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args: []expectedFnArg{
						{
							argIdent: "a",
							argType:  ast.StringType,
						},
					},
				},
			},
			expected: []ast.ReassignStatement{
				{
					Variable: &ast.Identifier{Name: "a"},
					Value:    &ast.StringLiteral{Value: "\"hello, world\""},
				},
			},
		},
		/*
			func reassignStatement2() {
				a = "hello, world"
			}
		*/
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "reassignStatement2",
						Args:     "",
						RetType:  "",
						Stmts: []string{
							`a = "hello, world"`,
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
					args:    []expectedFnArg{},
				},
			},
			expected: []ast.ReassignStatement{
				{
					Variable: &ast.Identifier{Name: "a"},
					Value:    &ast.StringLiteral{Value: "\"hello, world\""},
				},
			},
			expectedErr: parse.NotExistSymError{
				Source: parse.Token{Type: parse.Ident, Val: "a", Line: 3, Column: 9},
			},
		},
	}

	for i, tt := range tests {
		input := createTestContractCode(tt.contractTmpl)
		contract, err := parseTestContract(input)

		if err != nil && err.Error() == tt.expectedErr.Error() {
			continue
		}

		if err != nil && err.Error() != tt.expectedErr.Error() {
			t.Errorf("test[%d] - unexpected parser error expected=%s, got=%s",
				i, err.Error(), tt.expectedErr.Error())
			t.FailNow()
		}

		for i, fn := range contract.Functions {
			runReassignStatementTestCases(t, fn, tt.expectedFnHeaders[i], tt.expected[i])
		}
	}
}

func runReassignStatementTestCases(t *testing.T, fn *ast.FunctionLiteral, efh expectedFnHeader, tt ast.ReassignStatement) {
	t.Logf("test ExpressionStatement - [%s]", fn.Name)

	chkFnHeader(t, fn, efh)

	for _, stmt := range fn.Body.Statements {
		reasnStmt, ok := stmt.(*ast.ReassignStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.ExpressionStatement. got=%T", stmt)
		}
		testReassignStatement(t, reasnStmt, *tt.Variable, tt.Value)
	}
}

func testExpression(t *testing.T, exp ast.Expression, expected ast.Expression) {
	if exp == nil && expected == nil {
		return
	}

	if exp.String() != expected.String() {
		t.Errorf(`expression is not "%s", bot got "%s"`, exp.String(), expected.String())
	}
}

func testFnParameters(t *testing.T, p *ast.ParameterLiteral, ds ast.DataStructure, id string) {
	if p.Type != ds {
		t.Errorf("wrong parameter type Expected=%s, got=%s",
			p.Type.String(), ds.String())
	}
	if p.Identifier.Name != id {
		t.Errorf("wrong parameter identifier expected=%T, got=%T",
			p.Type, ds)
	}
}

func testAssignStatement(t *testing.T, stmt *ast.AssignStatement, ds ast.DataStructure, ident ast.Identifier, value ast.Expression) {
	if stmt.Type != ds {
		t.Errorf("wrong assign statement type Expected=%s, got=%s",
			ds.String(), stmt.Type.String())
	}
	if stmt.Variable.String() != ident.String() {
		t.Errorf("wrong assign statement variable Expected=%s, got=%s",
			ident.String(), stmt.Variable.String())
	}
	if stmt.Value.String() != value.String() {
		t.Errorf("wrong assign statement value Expected=%s, got=%s",
			value.String(), stmt.Value.String())
	}
}

func testIfStatement(t *testing.T, stmt *ast.IfStatement, condition ast.Expression, consequences *ast.BlockStatement, alternatives *ast.BlockStatement) {
	if stmt.Condition.String() != condition.String() {
		t.Errorf("wrong condition statement type Expected=%s, got=%s",
			condition, stmt.Condition.String())
	}

	if len(stmt.Consequence.Statements) != len(consequences.Statements) {
		t.Errorf("wrong condition statement consequences length Expected=%d, got=%d",
			len(consequences.Statements), len(stmt.Consequence.Statements))
	}
	for i, csq := range stmt.Consequence.Statements {
		if csq.String() != consequences.Statements[i].String() {
			t.Errorf("wrong condition statement consequences literal Expected=%s, got=%s",
				csq.String(), consequences.Statements[i].String())
		}
	}

	if stmt.Alternative == nil {
		return
	}
	if len(stmt.Alternative.Statements) != len(alternatives.Statements) {
		t.Errorf("wrong condition statement alternatives length Expected=%d, got=%d",
			len(alternatives.Statements), len(stmt.Alternative.Statements))
	}
	for i, alt := range stmt.Alternative.Statements {
		if alt.String() != alternatives.Statements[i].String() {
			t.Errorf("wrong condition statement alternative literal Expected=%s, got=%s",
				alt.String(), alternatives.Statements[i].String())
		}
	}
}

func testReassignStatement(t *testing.T, stmt *ast.ReassignStatement, ident ast.Identifier, value ast.Expression) {
	if stmt.Variable.String() != ident.String() {
		t.Errorf("wrong re-assign statement variable Expected=%s, got=%s",
			ident.String(), stmt.Variable.String())
	}
	if stmt.Value.String() != value.String() {
		t.Errorf("wrong re-assign statement value Expected=%s, got=%s",
			value.String(), stmt.Value.String())
	}
}
