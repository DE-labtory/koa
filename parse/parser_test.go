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
	"html/template"
	"testing"

	"github.com/DE-labtory/koa/ast"
	"github.com/DE-labtory/koa/parse"
)

type parserTestCase struct {
	assignTestCase []struct {
		ds    ast.DataStructure
		ident string
		value string
	}
	conditionTestCase []struct {
		condition   string
		consequence []string
		alternative []string
	}
	expressionStmtTestCase []struct {
		expression string
	}
	inlineReturnStmtTestCase struct {
		returnValue string
	}
	inlineAssignStmtTestCase struct {
		ds    ast.DataStructure
		ident string
		value string
	}
	inlineConditionTestCase struct {
		condition   string
		consequence []string
		alternative []string
	}
	inlineExpressionStmtTestCase struct {
		expression string
	}
}

type returnStmtExpect struct {
	returnValue ast.Expression
}

type assertFnHeader func(fn *ast.FunctionLiteral)

type testContractTmpl struct {
	FuncName string
	Args     string
	RetType  string
	Stmts    []string
}

const singleFnContractTmpl = `
contract {
    func {{.FuncName}}({{.Args}}) {{.RetType}} {
        {{range .Stmts}}
		{{.}}
        {{end}}
    }
}
`

var tmplInstance, _ = template.New("ContractTemplate").Parse(singleFnContractTmpl)

func genTestContractCode(c testContractTmpl) string {
	out := bytes.NewBufferString("")
	tmplInstance.Execute(out, c)
	return out.String()
}

func parseTestContract(input string) (*ast.Contract, error) {
	l := parse.NewLexer(input)
	buf := parse.NewTokenBuffer(l)
	return parse.Parse(buf)
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		contractTmpl testContractTmpl
		chkFnHeader  assertFnHeader
		expected     returnStmtExpect
	}{
		{
			contractTmpl: testContractTmpl{
				FuncName: "returnStatement1",
				Args:     "",
				RetType:  "",
				Stmts: []string{
					"return 1",
				},
			},
			chkFnHeader: func(fn *ast.FunctionLiteral) {
				if fn.ReturnType != ast.VoidType {
					t.Errorf("testReturn() has wrong return type expected=%v, got=%v",
						ast.VoidType.String(), fn.ReturnType.String())
				}

				if len(fn.Parameters) != 0 {
					t.Errorf("testReturn() has wrong parameters length got=%v",
						len(fn.Parameters))
				}
			},
			expected: returnStmtExpect{
				returnValue: &ast.IntegerLiteral{Value: 1},
			},
		},
	}

	for _, tt := range tests {
		input := genTestContractCode(tt.contractTmpl)
		contract, err := parseTestContract(input)

		if err != nil {
			t.Errorf("parser error: %q", err)
			t.FailNow()
		}

		testReturnStatement(t, tt.chkFnHeader, contract.Functions[0], tt.expected)
	}
}

func testReturnStatement(t *testing.T, chkFnHeader assertFnHeader, fn *ast.FunctionLiteral, tt returnStmtExpect) {
	t.Logf("test return statement - [%s]", fn.Name)

	chkFnHeader(fn)

	// test testReturn's return statements
	for _, stmt := range fn.Body.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.ReturnStatement. got=%T", stmt)
		}

		testExpression2(t, returnStmt.ReturnValue, tt.returnValue)
	}
}

func testAssignStatementFunc(t *testing.T, fn *ast.FunctionLiteral, tt parserTestCase) {
	t.Log("test assign statement")

	// test testAssign() function body, parameters, return type
	if fn.ReturnType != ast.StringType {
		t.Errorf("testAssign() has wrong return type expected=%v, got=%v",
			ast.StringType.String(), fn.ReturnType.String())
	}

	if len(fn.Parameters) != 1 {
		t.Errorf("testAssign() has wrong parameters length got=%v",
			len(fn.Parameters))
	}

	testFnParameters(t, fn.Parameters[0], ast.IntType, "foo")

	// test testAssign()'s assign statements
	for i, stmt := range fn.Body.Statements {
		assignStmt, ok := stmt.(*ast.AssignStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.AssignStatement. got=%T", stmt)
		}
		tc := tt.assignTestCase[i]
		testAssignStatement(t, assignStmt, tc.ds, tc.ident, tc.value)
	}
}

func testIfElseStatementFunc(t *testing.T, fn *ast.FunctionLiteral, tt parserTestCase) {
	t.Log("test if-else statement")

	// test testIfElse() function body, parameters, return type
	if fn.ReturnType != ast.IntType {
		t.Errorf("testIfElse() has wrong return type expected=%v, got=%v",
			ast.StringType.String(), fn.ReturnType.String())
	}

	if len(fn.Parameters) != 3 {
		t.Errorf("testIfElse() has wrong parameters length got=%v",
			len(fn.Parameters))
	}

	testFnParameters(t, fn.Parameters[0], ast.IntType, "foo")
	testFnParameters(t, fn.Parameters[1], ast.StringType, "bar")
	testFnParameters(t, fn.Parameters[2], ast.StringType, "baz")

	// test testIfElse()'s assign statements
	for i, stmt := range fn.Body.Statements {
		ifStmt, ok := stmt.(*ast.IfStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.IfStatement. got=%T", stmt)
		}
		tc := tt.conditionTestCase[i]
		testIfStatement(t, ifStmt, tc.condition, tc.consequence, tc.alternative)
	}
}

func testExpressionStatementFunc(t *testing.T, fn *ast.FunctionLiteral, tt parserTestCase) {
	t.Log("test expression-statement statement")

	// test testExpressionStatement() function body, parameters, return type
	if fn.ReturnType != ast.BoolType {
		t.Errorf("testExpressionS	tatement() has wrong return type expected=%v, got=%v",
			ast.BoolType.String(), fn.ReturnType.String())
	}

	if len(fn.Parameters) != 1 {
		t.Errorf("testExpressionStatement() has wrong parameters length got=%v",
			len(fn.Parameters))
	}

	testFnParameters(t, fn.Parameters[0], ast.BoolType, "foo")

	// test testExpressionStatement()'s return statements
	for i, stmt := range fn.Body.Statements {
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.ExpressionStatement. got=%T", expStmt)
		}

		tc := tt.expressionStmtTestCase[i]
		testExpression(t, expStmt.Expr, tc.expression)
	}
}

func testInlineReturnStatementFunc(t *testing.T, fn *ast.FunctionLiteral, tt parserTestCase) {
	t.Log("test inline-function with return statement")

	if fn.ReturnType != ast.StringType {
		t.Errorf("testInlineReturnStatement() has wrong return type expected=%v, got=%v",
			ast.StringType.String(), fn.ReturnType.String())
	}

	if len(fn.Parameters) != 0 {
		t.Errorf("testInlineReturnStatement() has wrong parameters length got=%v",
			len(fn.Parameters))
	}

	if len(fn.Body.Statements) != 1 {
		t.Errorf("testInlineReturnStatement() has wrong body statements length got=%v",
			len(fn.Parameters))
	}

	returnStmt, ok := fn.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Errorf("function body stmt is not *ast.ReturnStatement. got=%T", fn.Body.Statements[0])
	}

	testExpression(t, returnStmt.ReturnValue, tt.inlineReturnStmtTestCase.returnValue)
}

func testInlineAssignStatementFunc(t *testing.T, fn *ast.FunctionLiteral, tt parserTestCase) {
	t.Log("test inline-function with assign statement")

	if fn.ReturnType != ast.VoidType {
		t.Errorf("testInlineAssignStatement() has wrong return type expected=%v, got=%v",
			ast.StringType.String(), fn.ReturnType.String())
	}

	if len(fn.Parameters) != 0 {
		t.Errorf("testInlineAssignStatement() has wrong parameters length got=%v",
			len(fn.Parameters))
	}

	if len(fn.Body.Statements) != 1 {
		t.Errorf("testInlineAssignStatement() has wrong body statements length got=%v",
			len(fn.Parameters))
	}

	stmt, ok := fn.Body.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Errorf("function body stmt is not *ast.AssignStatement. got=%T", fn.Body.Statements[0])
	}

	tc := tt.inlineAssignStmtTestCase
	testAssignStatement(t, stmt, tc.ds, tc.ident, tc.value)
}

func testInlineConditionStatement(t *testing.T, fn *ast.FunctionLiteral, tt parserTestCase) {
	t.Log("test inline-function with condition statement")

	if fn.ReturnType != ast.VoidType {
		t.Errorf("testInlineConditionStatement() has wrong return type expected=%v, got=%v",
			ast.StringType.String(), fn.ReturnType.String())
	}

	if len(fn.Parameters) != 0 {
		t.Errorf("testInlineConditionStatement() has wrong parameters length got=%v",
			len(fn.Parameters))
	}

	if len(fn.Body.Statements) != 1 {
		t.Errorf("testInlineConditionStatement() has wrong body statements length got=%v",
			len(fn.Parameters))
	}

	stmt, ok := fn.Body.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Errorf("function body stmt is not *ast.IfStatement. got=%T", fn.Body.Statements[0])
	}

	tc := tt.inlineConditionTestCase
	testIfStatement(t, stmt, tc.condition, tc.consequence, tc.alternative)
}

func testInlineExpressionStatement(t *testing.T, fn *ast.FunctionLiteral, tt parserTestCase) {
	t.Log("test inline-function with expression statement")

	if fn.ReturnType != ast.VoidType {
		t.Errorf("testInlineExpressionStatement() has wrong return type expected=%v, got=%v",
			ast.StringType.String(), fn.ReturnType.String())
	}

	if len(fn.Parameters) != 0 {
		t.Errorf("testInlineExpressionStatement() has wrong parameters length got=%v",
			len(fn.Parameters))
	}

	if len(fn.Body.Statements) != 1 {
		t.Errorf("testInlineExpressionStatement() has wrong body statements length got=%v",
			len(fn.Parameters))
	}

	stmt, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("function body stmt is not *ast.ExpressionStatement. got=%T", stmt)
	}

	tc := tt.inlineExpressionStmtTestCase
	testExpression(t, stmt.Expr, tc.expression)
}

func testExpression(t *testing.T, exp ast.Expression, expected string) {
	if exp.String() != expected {
		t.Errorf(`expression is not "%s", bot got "%s"`, exp.String(), expected)
	}
}

func testExpression2(t *testing.T, exp ast.Expression, expected ast.Expression) {
	if exp.String() != expected.String() {
		t.Errorf(`expression is not "%s", bot got "%s"`, exp.String(), expected.String())
	}
}

func testFnParameters(t *testing.T, p *ast.ParameterLiteral, ds ast.DataStructure, id string) {
	if p.Type != ds {
		t.Errorf("wrong parameter type expected=%s, got=%s",
			p.Type.String(), ds.String())
	}
	if p.Identifier.Value != id {
		t.Errorf("wrong parameter identifier expected=%T, got=%T",
			p.Type, ds)
	}
}

func testAssignStatement(t *testing.T, stmt *ast.AssignStatement, ds ast.DataStructure, ident string, value string) {
	if stmt.Type != ds {
		t.Errorf("wrong assign statement type expected=%T, got=%T",
			ds, stmt.Type)
	}
	if stmt.Variable.Value != ident {
		t.Errorf("wrong assign statement variable expected=%s, got=%s",
			ident, stmt.Variable.Value)
	}
	if stmt.Value.String() != value {
		t.Errorf("wrong assign statement value expected=%s, got=%s",
			value, stmt.Value.String())
	}
}

func testIfStatement(t *testing.T, stmt *ast.IfStatement, condition string, consequences []string, alternatives []string) {
	if stmt.Condition.String() != condition {
		t.Errorf("wrong condition statement type expected=%s, got=%s",
			condition, stmt.Condition.String())
	}

	if len(stmt.Consequence.Statements) != len(consequences) {
		t.Errorf("wrong condition statement consequences length expected=%d, got=%d",
			len(consequences), len(stmt.Consequence.Statements))
	}
	for i, csq := range stmt.Consequence.Statements {
		if csq.String() != consequences[i] {
			t.Errorf("wrong condition statement consequences literal expected=%s, got=%s",
				csq.String(), consequences[i])
		}
	}

	if stmt.Alternative == nil {
		return
	}
	if len(stmt.Alternative.Statements) != len(alternatives) {
		t.Errorf("wrong condition statement alternatives length expected=%d, got=%d",
			len(alternatives), len(stmt.Alternative.Statements))
	}
	for i, alt := range stmt.Alternative.Statements {
		if alt.String() != alternatives[i] {
			t.Errorf("wrong condition statement alternative literal expected=%s, got=%s",
				alt.String(), alternatives[i])
		}
	}
}
