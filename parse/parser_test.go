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
	"testing"

	"github.com/DE-labtory/koa/ast"
	"github.com/DE-labtory/koa/parse"
)

type parserTestCase struct {
	returnTestCase []struct {
		returnValue string
	}
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

func TestParse(t *testing.T) {
	input := `
	contract {
		// testReturn must have return statements, this function is for 
		// testing return statements.
		func testReturnStatement() {
			return 1
			return a
			return a * b + 1
			return a * (b + 1)
			return (
				a)
			return add(1, 2)
			return add(
				1,
				2)
			return add(
				1,
				2 /* at this point there's SEMICOLON, should consume */
			)
		}
		
		/* testIntType  
		is for testing int assign statements */
		func testAssignStatement(foo int) string {
			int a = 1
			int a = 1 + 2
			int a =
				1 + 2
			int a = (foo + 1) * 2
			int a = add(1 + 2, 3) + 4
			int a = add(1 % 2, 3) / 4
			
			// This is not working code according to go-spec
			//
			// int a = add(1 + 2, 3) 
			//	 + 4
			
			string a = "hello"
			string a = 
				"hello"
			string tabbed_string = "hello, \t world"

			bool a = true
			bool b= false
		}
		
		// testIfElse is for testing if-else statements and should only contain
		// if-else statements
		func testIfStatement(foo int, bar string, baz string) int {
			if (true) {}
		
			if (1 != 1 + 2) {
				int a = 1
				string a = "hello"
			}
		
			if (foo) {} else {}
			
			if (foo) {
				int a = 1
				string a = "hello"
			} else {
				int a = 1
				string a = "hello"
			}
		}
		
		// testExpressionStatement is for testing expression statement and should
		// only contain expression statement
		func testExpressionStatement(foo bool) bool {
			add(1, 2)
			add(add(1, 2), 3)
			add(add(1,
				2), 3)
			add(add(1, 2),
				3)
			add(
				add(
					1, 
					2
				),
				3
			)
		}
		
		// testInlineReturnStatement is for testing inline-function, in the case when
		// return statement does not change the line, then do not insert semicolon
		// related with issue #228
		func testInlineReturnStatement() string { return "hello" }

		// testInlineAssignStatement is for testing inline-function, in the case when
		// assign statement does not change the line, then do not insert semicolon
		// related with issue #228
		func testInlineAssignStatement() { int a = 1 }

		// testInlineConditionStatement is for testing inline-function, in the case when
		// condition statement does not change the line, then do not insert semicolon
		// related with issue #228
		func testInlineConditionStatement() { if(true) {} else {} }

		// testInlineExpressionStatement is for testing inline-function, in the case when
		// expression statement does not change the line, then do not insert semicolon
		// related with issue #228
		func testInlineExpressionStatement() { add(1, 2) }
	}	
`
	tcs := parserTestCase{
		returnTestCase: []struct {
			returnValue string
		}{
			{"1"},
			{"a"},
			{"((a * b) + 1)"},
			{"(a * (b + 1))"},
			{"a"},
			{"function add( 1, 2 )"},
			{"function add( 1, 2 )"},
			{"function add( 1, 2 )"},
		},
		assignTestCase: []struct {
			ds    ast.DataStructure
			ident string
			value string
		}{
			{ast.IntType, "a", "1"},
			{ast.IntType, "a", "(1 + 2)"},
			{ast.IntType, "a", "(1 + 2)"},
			{ast.IntType, "a", "((foo + 1) * 2)"},
			{ast.IntType, "a", "(function add( (1 + 2), 3 ) + 4)"},
			{ast.IntType, "a", "(function add( (1 % 2), 3 ) / 4)"},
			{ast.StringType, "a", "\"hello\""},
			{ast.StringType, "a", "\"hello\""},
			{ast.StringType, "tabbed_string", "\"hello, \\t world\""},
			{ast.BoolType, "a", "true"},
			{ast.BoolType, "b", "false"},
		},
		conditionTestCase: []struct {
			condition   string
			consequence []string
			alternative []string
		}{
			{
				condition: "true",
			},
			{
				condition: "(1 != (1 + 2))",
				consequence: []string{
					"int a = 1",
					"string a = \"hello\"",
				},
			},
			{
				condition: "foo",
			},
			{
				condition: "foo",
				consequence: []string{
					"int a = 1",
					"string a = \"hello\"",
				},
				alternative: []string{
					"int a = 1",
					"string a = \"hello\"",
				},
			},
		},
		expressionStmtTestCase: []struct {
			expression string
		}{
			{"function add( 1, 2 )"},
			{"function add( function add( 1, 2 ), 3 )"},
			{"function add( function add( 1, 2 ), 3 )"},
			{"function add( function add( 1, 2 ), 3 )"},
			{"function add( function add( 1, 2 ), 3 )"},
		},
		inlineReturnStmtTestCase: struct {
			returnValue string
		}{"\"hello\""},
		inlineAssignStmtTestCase: struct {
			ds    ast.DataStructure
			ident string
			value string
		}{ast.IntType, "a", "1"},
		inlineConditionTestCase: struct {
			condition   string
			consequence []string
			alternative []string
		}{
			condition:   "true",
			consequence: []string{},
			alternative: []string{},
		},
		inlineExpressionStmtTestCase: struct {
			expression string
		}{"function add( 1, 2 )"},
	}

	l := parse.NewLexer(input)
	buf := parse.NewTokenBuffer(l)
	contract, err := parse.Parse(buf)

	if err != nil {
		t.Errorf("parser error: %q", err)
		t.FailNow()
	}
	for _, fn := range contract.Functions {
		testFunctionLiteral(t, fn, tcs)
	}
}

func testFunctionLiteral(t *testing.T, fn *ast.FunctionLiteral, tt parserTestCase) {
	t.Helper()
	switch fn.Name.Value {
	case "testReturnStatement":
		testReturnStatementFunc(t, fn, tt)
	case "testAssignStatement":
		testAssignStatementFunc(t, fn, tt)
	case "testIfStatement":
		testIfElseStatementFunc(t, fn, tt)
	case "testExpressionStatement":
		testExpressionStatementFunc(t, fn, tt)
	case "testInlineReturnStatement":
		testInlineReturnStatementFunc(t, fn, tt)
	case "testInlineAssignStatement":
		testInlineAssignStatementFunc(t, fn, tt)
	case "testInlineConditionStatement":
		testInlineConditionStatement(t, fn, tt)
	case "testInlineExpressionStatement":
		testInlineExpressionStatement(t, fn, tt)
	}
}

func testReturnStatementFunc(t *testing.T, fn *ast.FunctionLiteral, tt parserTestCase) {
	t.Log("test return statement")

	// test testReturn() function body, parameters, return type
	if fn.ReturnType != ast.VoidType {
		t.Errorf("testReturn() has wrong return type expected=%v, got=%v",
			ast.VoidType.String(), fn.ReturnType.String())
	}

	if len(fn.Parameters) != 0 {
		t.Errorf("testReturn() has wrong parameters length got=%v",
			len(fn.Parameters))
	}

	// test testReturn's return statements
	for i, stmt := range fn.Body.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.ReturnStatement. got=%T", stmt)
		}

		tc := tt.returnTestCase[i]
		testExpression(t, returnStmt.ReturnValue, tc.returnValue)
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
