/*
// * Copyright 2018 De-labtory
// *
// * Licensed under the Apache License, Version 2.0 (the "License");
// * you may not use this file except in compliance with the License.
// * You may obtain a copy of the License at
// *
// * https://www.apache.org/licenses/LICENSE-2.0
// *
// * Unless required by applicable law or agreed to in writing, software
// * distributed under the License is distributed on an "AS IS" BASIS,
// * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// * See the License for the specific language governing permissions and
// * limitations under the License.
// */

package symbol

import (
	"reflect"
	"testing"

	"github.com/DE-labtory/koa/ast"
)

type setupResolverFn func() *Resolver

func defaultResolver() *Resolver {
	return NewResolver()
}

func postMakeScope(s *Scope) *Scope {
	for _, child := range s.child {
		child.parent = s
		postMakeScope(child)
	}
	return s
}

func TestNewScope(t *testing.T) {

}

func TestNewEnclosedScope(t *testing.T) {

}

func TestGet(t *testing.T) {
	tests := []struct {
		wanted   string
		input    *Scope
		expected Object
	}{
		{
			wanted: "a",
			input: &Scope{
				store: map[string]Object{
					"a": &Integer{Name: &ast.Identifier{Name: "a"}},
				},
			},
			expected: &Integer{Name: &ast.Identifier{Name: "a"}},
		},
	}

	for i, test := range tests {
		obj := test.input.Get(test.wanted)
		if !reflect.DeepEqual(obj, test.expected) {
			t.Fatalf("[test %d] - TestGet failed.\nexpected=%v\ngot=%v",
				i, test.expected, obj)
		}
	}
}

func TestResolveFunction(t *testing.T) {
	vars := []struct {
		fl *ast.FunctionLiteral
	}{
		{
			fl: &ast.FunctionLiteral{
				Name:       &ast.Identifier{Name: "testFunction"},
				Parameters: []*ast.ParameterLiteral{},
				Body:       &ast.BlockStatement{},
				ReturnType: ast.VoidType,
			},
		},
		{
			fl: &ast.FunctionLiteral{
				Name: &ast.Identifier{Name: "add"},
				Parameters: []*ast.ParameterLiteral{
					{
						Identifier: &ast.Identifier{Name: "a"},
						Type:       ast.IntType,
					},
					{
						Identifier: &ast.Identifier{Name: "b"},
						Type:       ast.IntType,
					},
				},
				Body: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.ReturnStatement{
							ReturnValue: &ast.InfixExpression{
								Left: &ast.Identifier{
									Name: "a",
								},
								Operator: ast.Plus,
								Right: &ast.Identifier{
									Name: "b",
								},
							},
						},
					},
				},
				ReturnType: ast.IntType,
			},
		},
	}

	tests := []struct {
		setupResolverFn
		input    *ast.FunctionLiteral
		expected *Resolver
		err      error
	}{
		{
			// test case 1
			// void testFunction()
			setupResolverFn: defaultResolver,
			input:           vars[0].fl,
			expected: &Resolver{
				scope: &Scope{
					store: map[string]Object{},
					child: []*Scope{
						{
							store: make(map[string]Object, 0),
							child: []*Scope{
								{
									store: make(map[string]Object, 0),
									child: []*Scope{},
								},
							},
						},
					},
				},
				types: map[ast.Expression]ObjectType{},
				defs:  map[*ast.Identifier]Object{},
				fns:   map[string]*ast.FunctionLiteral{},
			},
			err: nil,
		},
		{
			// test case 2
			// int add(int a, int b) { return a + b }
			setupResolverFn: defaultResolver,
			input:           vars[1].fl,
			expected: &Resolver{
				scope: &Scope{
					store: map[string]Object{},
					child: []*Scope{
						{
							store: map[string]Object{
								"a": &Integer{Name: &ast.Identifier{Name: "a"}},
								"b": &Integer{Name: &ast.Identifier{Name: "b"}},
							},
							child: []*Scope{
								{
									store: make(map[string]Object, 0),
									child: []*Scope{},
								},
							},
						},
					},
				},
				types: map[ast.Expression]ObjectType{
					vars[1].fl.Parameters[0]: IntegerObject,
					vars[1].fl.Parameters[1]: IntegerObject,
				},
				defs: map[*ast.Identifier]Object{
					vars[1].fl.Parameters[0].Identifier: &Integer{Name: &ast.Identifier{Name: "a"}},
					vars[1].fl.Parameters[1].Identifier: &Integer{Name: &ast.Identifier{Name: "b"}},
				},
				fns: map[string]*ast.FunctionLiteral{},
			},
		},
	}
	for i, test := range tests {
		r := test.setupResolverFn()
		postMakeScope(test.expected.scope)
		err := ResolveFunction(test.input, r)
		if err == nil && !reflect.DeepEqual(test.expected, r) {
			t.Fatalf("test [%d] - TestResolveFunction failed.\nexpected=%v\ngot=%v",
				i, test.expected, r)
		}
	}
}

func TestResolveFunctionParameter(t *testing.T) {
	vars := []struct {
		pls []*ast.ParameterLiteral
	}{
		{
			pls: []*ast.ParameterLiteral{
				{
					Identifier: &ast.Identifier{
						Name: "a",
					},
					Type: ast.IntType,
				},
			},
		},
		{
			pls: []*ast.ParameterLiteral{
				{
					Identifier: &ast.Identifier{
						Name: "a",
					},
					Type: ast.IntType,
				},
				{
					Identifier: &ast.Identifier{
						Name: "b",
					},
					Type: ast.BoolType,
				},
				{
					Identifier: &ast.Identifier{
						Name: "c",
					},
					Type: ast.StringType,
				},
			},
		},
		{
			pls: []*ast.ParameterLiteral{
				{
					Identifier: &ast.Identifier{
						Name: "a",
					},
					Type: ast.IntType,
				},
				{
					Identifier: &ast.Identifier{
						Name: "a",
					},
					Type: ast.IntType,
				},
			},
		},
	}
	tests := []struct {
		setupResolverFn
		input    []*ast.ParameterLiteral
		expected *Resolver
		err      error
	}{
		{
			// test case 1
			// function ... (int a)
			setupResolverFn: defaultResolver,
			input:           vars[0].pls,
			expected: &Resolver{
				scope: &Scope{
					store: map[string]Object{
						"a": &Integer{Name: &ast.Identifier{Name: "a"}},
					},
				},
				types: map[ast.Expression]ObjectType{
					vars[0].pls[0]: IntegerObject,
				},
				defs: map[*ast.Identifier]Object{
					vars[0].pls[0].Identifier: &Integer{Name: &ast.Identifier{Name: "a"}},
				},
				fns: map[string]*ast.FunctionLiteral{},
			},
			err: nil,
		},
		{
			// test case 2
			// function ... (int a, bool b, string c)
			setupResolverFn: defaultResolver,
			input:           vars[1].pls,
			expected: &Resolver{
				scope: &Scope{
					store: map[string]Object{
						"a": &Integer{Name: &ast.Identifier{Name: "a"}},
						"b": &Boolean{Name: &ast.Identifier{Name: "b"}},
						"c": &String{Name: &ast.Identifier{Name: "c"}},
					},
				},
				types: map[ast.Expression]ObjectType{
					vars[1].pls[0]: IntegerObject,
					vars[1].pls[1]: BooleanObject,
					vars[1].pls[2]: StringObject,
				},
				defs: map[*ast.Identifier]Object{
					vars[1].pls[0].Identifier: &Integer{Name: &ast.Identifier{Name: "a"}},
					vars[1].pls[1].Identifier: &Boolean{Name: &ast.Identifier{Name: "b"}},
					vars[1].pls[2].Identifier: &String{Name: &ast.Identifier{Name: "c"}},
				},
				fns: map[string]*ast.FunctionLiteral{},
			},
			err: nil,
		},
		{
			setupResolverFn: defaultResolver,
			input:           vars[2].pls,
			expected:        nil,
			err: DupError{
				object: &Integer{
					Name: &ast.Identifier{
						Name: "a",
					},
				},
			},
		},
	}

	for i, test := range tests {
		r := test.setupResolverFn()
		err := ResolveFunctionParameter(test.input, r)
		if err == nil && !reflect.DeepEqual(test.expected, r) {
			t.Fatalf("test [%d] - TestScopingFunctionParameter failed.\nexpected=%v\ngot=%v",
				i, test.expected, r)
		}
		if err != nil && !reflect.DeepEqual(test.err, err) {
			t.Fatalf("test [%d] - TestScopingFunctionParameter failed (err case).\nexpected=%v\ngot=%v",
				i, test.err, err)
		}
	}
}

func TestResolveFunctionBody(t *testing.T) {

}

func TestResolveAssignStatement(t *testing.T) {
	vars := []struct {
		as *ast.AssignStatement
	}{
		{
			as: &ast.AssignStatement{
				Type:     ast.IntType,
				Variable: ast.Identifier{Name: "a"},
				Value: &ast.IntegerLiteral{
					Value: 0,
				},
			},
		},
		{
			as: &ast.AssignStatement{
				Type:     ast.StringType,
				Variable: ast.Identifier{Name: "str"},
				Value: &ast.StringLiteral{
					Value: "testString",
				},
			},
		},
		{},
	}

	tests := []struct {
		setupResolverFn
		input    *ast.AssignStatement
		expected *Resolver
		err      error
	}{
		{
			// test case 1
			// int a = 0
			setupResolverFn: defaultResolver,
			input:           vars[0].as,
			expected: &Resolver{
				scope: &Scope{
					store: map[string]Object{
						"a": &Integer{Name: &ast.Identifier{Name: "a"}},
					},
				},
				types: map[ast.Expression]ObjectType{
					vars[0].as.Value: IntegerObject,
				},
				defs: map[*ast.Identifier]Object{
					&vars[0].as.Variable: &Integer{Name: &vars[0].as.Variable},
				},
			},
		},
		{
			// test case 2
			// string str = "str"
			setupResolverFn: defaultResolver,
			input:           vars[1].as,
			expected: &Resolver{
				scope: &Scope{
					store: map[string]Object{
						"str": &String{Name: &ast.Identifier{Name: "str"}},
					},
				},
				types: map[ast.Expression]ObjectType{
					vars[1].as.Value: StringObject,
				},
				defs: map[*ast.Identifier]Object{
					&vars[1].as.Variable: &String{Name: &vars[1].as.Variable},
				},
			},
		},
		{
			// test case 3
			// string str = "str"
			// but val named str is already defined
			setupResolverFn: func() *Resolver {
				return &Resolver{
					scope: &Scope{
						store: map[string]Object{
							"str": &String{Name: &ast.Identifier{Name: "str"}},
						},
					},
				}
			},
			input: &ast.AssignStatement{
				Type:     ast.StringType,
				Variable: ast.Identifier{Name: "str"},
				Value: &ast.StringLiteral{
					Value: "testString",
				},
			},
			expected: nil,
			err: DupError{
				object: &String{
					Name: &ast.Identifier{
						Name: "str",
					},
				},
			},
		},
		{
			// test case 4
			// int a = "str"
			setupResolverFn: defaultResolver,
			input: &ast.AssignStatement{
				Type:     ast.IntType,
				Variable: ast.Identifier{Name: "a"},
				Value: &ast.StringLiteral{
					Value: "str",
				},
			},
			expected: nil,
			err: &TypeError{
				target: IntegerObject,
				object: StringObject,
			},
		},
	}

	for i, test := range tests {
		r := test.setupResolverFn()
		err := ResolveAssignStatement(test.input, r)

		if err == nil && !reflect.DeepEqual(r.defs, test.expected.defs) {
			t.Fatalf("test [%d] - TestScopingAssignStatement failed.\nexpected=%v\ngot=%v",
				i, test.expected, r)
		}
		if err != nil && err.Error() != test.err.Error() {
			t.Fatalf("test [%d] - TestScopingAssignStatement failed (err case).\nexpected=%v\ngot=%v",
				i, test.err, err)
		}
	}
}

func TestResolveBlockStatement(t *testing.T) {
	vars := []struct {
		bs *ast.BlockStatement
	}{
		{
			bs: &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.AssignStatement{
						Type:     ast.IntType,
						Variable: ast.Identifier{Name: "a"},
						Value: &ast.IntegerLiteral{
							Value: 0,
						},
					},
					&ast.AssignStatement{
						Type:     ast.IntType,
						Variable: ast.Identifier{Name: "b"},
						Value: &ast.IntegerLiteral{
							Value: 1,
						},
					},
				},
			},
		},
		{
			&ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.AssignStatement{
						Type:     ast.IntType,
						Variable: ast.Identifier{Name: "a"},
						Value: &ast.IntegerLiteral{
							Value: 0,
						},
					},
					&ast.AssignStatement{
						Type:     ast.IntType,
						Variable: ast.Identifier{Name: "a"},
						Value: &ast.IntegerLiteral{
							Value: 1,
						},
					},
				},
			},
		},
	}

	tests := []struct {
		setupResolverFn
		input    *ast.BlockStatement
		expected *Resolver
		err      error
	}{
		{
			// test case 1
			// int a = 0
			// int b = 1
			setupResolverFn: defaultResolver,
			input:           vars[0].bs,
			expected: &Resolver{
				scope: &Scope{
					store: make(map[string]Object, 0),
					child: []*Scope{
						{
							store: map[string]Object{
								"a": &Integer{Name: &ast.Identifier{Name: "a"}},
								"b": &Integer{Name: &ast.Identifier{Name: "b"}},
							},
							child: []*Scope{},
						},
					},
				},
				types: map[ast.Expression]ObjectType{
					vars[0].bs.Statements[0].(*ast.AssignStatement).Value: IntegerObject,
					vars[0].bs.Statements[1].(*ast.AssignStatement).Value: IntegerObject,
				},
				defs: map[*ast.Identifier]Object{
					&vars[0].bs.Statements[0].(*ast.AssignStatement).Variable: &Integer{Name: &ast.Identifier{Name: "a"}},
					&vars[0].bs.Statements[1].(*ast.AssignStatement).Variable: &Integer{Name: &ast.Identifier{Name: "b"}},
				},
				fns: map[string]*ast.FunctionLiteral{},
			},
			err: nil,
		},
		{
			// test case 2
			// int a = 0
			// int a = 1
			// This is an error case.
			setupResolverFn: defaultResolver,
			input:           vars[1].bs,
			expected:        nil,
			err: DupError{
				object: &Integer{
					Name: &ast.Identifier{
						Name: "a",
					},
				},
			},
		},
	}

	for i, test := range tests {
		if test.expected != nil {
			postMakeScope(test.expected.scope)
		}

		r := test.setupResolverFn()
		err := ResolveBlockStatement(test.input, r)
		if err == nil && !reflect.DeepEqual(r, test.expected) {
			t.Fatalf("test [%d] - TestScopingBlockStatement failed.\nexpected=%v\ngot=%v",
				i, test.expected, r)
		}

		if err != nil && !reflect.DeepEqual(test.err, err) {
			t.Fatalf("test [%d] - TestScopingAssignStatement failed (err case).\nexpected=%v\ngot=%v",
				i, test.err, err)
		}
	}
}

func TestResolveIfStatement(t *testing.T) {
	vars := []*ast.IfStatement{
		{
			// case 1
			Condition: &ast.BooleanLiteral{
				Value: true,
			},
			Consequence: &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.AssignStatement{
						Type:     ast.IntType,
						Variable: ast.Identifier{Name: "a"},
						Value: &ast.IntegerLiteral{
							Value: 0,
						},
					},
				},
			},
			Alternative: nil,
		},
		{
			// case 2
			Condition: &ast.BooleanLiteral{
				Value: true,
			},
			Consequence: &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.AssignStatement{
						Type:     ast.IntType,
						Variable: ast.Identifier{Name: "a"},
						Value: &ast.IntegerLiteral{
							Value: 0,
						},
					},
					&ast.AssignStatement{
						Type:     ast.IntType,
						Variable: ast.Identifier{Name: "b"},
						Value: &ast.IntegerLiteral{
							Value: 0,
						},
					},
				},
			},
		},
		{
			// case 3
			Condition: &ast.BooleanLiteral{
				Value: true,
			},
			Consequence: &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.AssignStatement{
						Type:     ast.IntType,
						Variable: ast.Identifier{Name: "a"},
						Value: &ast.IntegerLiteral{
							Value: 0,
						},
					},
					&ast.AssignStatement{
						Type:     ast.IntType,
						Variable: ast.Identifier{Name: "b"},
						Value: &ast.IntegerLiteral{
							Value: 0,
						},
					},
				},
			},
			Alternative: &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.AssignStatement{
						Type:     ast.IntType,
						Variable: ast.Identifier{Name: "a"},
						Value: &ast.IntegerLiteral{
							Value: 0,
						},
					},
				},
			},
		},
	}
	tests := []struct {
		setupResolverFn
		input    *ast.IfStatement
		expected *Resolver
		err      error
	}{
		{
			// Test case 1
			// if (true) {
			//    int a = 0
			// }
			setupResolverFn: defaultResolver,
			input:           vars[0],
			expected: &Resolver{
				scope: &Scope{
					store: make(map[string]Object, 0),
					child: []*Scope{
						{
							store: map[string]Object{
								"a": &Integer{Name: &ast.Identifier{Name: "a"}},
							},
							child: []*Scope{},
						},
					},
				},
				types: map[ast.Expression]ObjectType{
					vars[0].Consequence.Statements[0].(*ast.AssignStatement).Value: IntegerObject,
				},
				defs: map[*ast.Identifier]Object{
					&vars[0].Consequence.Statements[0].(*ast.AssignStatement).Variable: &Integer{Name: &ast.Identifier{Name: "a"}},
				},
				fns: map[string]*ast.FunctionLiteral{},
			},
		},
		{
			// Test case 2
			// if (true) {
			//   int a = 0
			//   int b = 0
			// }
			setupResolverFn: defaultResolver,
			input:           vars[1],
			expected: &Resolver{
				scope: &Scope{
					store: make(map[string]Object, 0),
					child: []*Scope{
						{
							store: map[string]Object{
								"a": &Integer{Name: &ast.Identifier{Name: "a"}},
								"b": &Integer{Name: &ast.Identifier{Name: "b"}},
							},
							child: []*Scope{},
						},
					},
				},
				types: map[ast.Expression]ObjectType{
					vars[1].Consequence.Statements[0].(*ast.AssignStatement).Value: IntegerObject,
					vars[1].Consequence.Statements[1].(*ast.AssignStatement).Value: IntegerObject,
				},
				defs: map[*ast.Identifier]Object{
					&vars[1].Consequence.Statements[0].(*ast.AssignStatement).Variable: &Integer{Name: &ast.Identifier{Name: "a"}},
					&vars[1].Consequence.Statements[1].(*ast.AssignStatement).Variable: &Integer{Name: &ast.Identifier{Name: "b"}},
				},
				fns: map[string]*ast.FunctionLiteral{},
			},
		},
		{
			// Test case 3
			// if (true) {
			//    int a = 0
			//    int b = 0
			// } else {
			//    int a = 0
			// }
			setupResolverFn: defaultResolver,
			input:           vars[2],
			expected: &Resolver{
				scope: &Scope{
					store: make(map[string]Object, 0),
					child: []*Scope{
						{
							store: map[string]Object{
								"a": &Integer{Name: &ast.Identifier{Name: "a"}},
								"b": &Integer{Name: &ast.Identifier{Name: "b"}},
							},
							child: []*Scope{},
						},
						{
							store: map[string]Object{
								"a": &Integer{Name: &ast.Identifier{Name: "a"}},
							},
							child: []*Scope{},
						},
					},
				},
				types: map[ast.Expression]ObjectType{
					vars[2].Consequence.Statements[0].(*ast.AssignStatement).Value: IntegerObject,
					vars[2].Consequence.Statements[1].(*ast.AssignStatement).Value: IntegerObject,
					vars[2].Alternative.Statements[0].(*ast.AssignStatement).Value: IntegerObject,
				},
				defs: map[*ast.Identifier]Object{
					&vars[2].Consequence.Statements[0].(*ast.AssignStatement).Variable: &Integer{Name: &ast.Identifier{Name: "a"}},
					&vars[2].Consequence.Statements[1].(*ast.AssignStatement).Variable: &Integer{Name: &ast.Identifier{Name: "b"}},
					&vars[2].Consequence.Statements[0].(*ast.AssignStatement).Variable: &Integer{Name: &ast.Identifier{Name: "a"}},
				},
				fns: map[string]*ast.FunctionLiteral{},
			},
		},
	}

	for i, test := range tests {
		if test.expected != nil {
			postMakeScope(test.expected.scope)
		}
		r := test.setupResolverFn()
		err := ResolveIfStatement(test.input, r)
		if err == nil && !reflect.DeepEqual(test.expected.scope, r.scope) {
			t.Fatalf("test [%d] - TestScopingIfStatement failed.\nexpected=%v\ngot=%v",
				i, test.expected, r)
		}

		if err != nil && !reflect.DeepEqual(test.err, err) {
			t.Fatalf("test [%d] - TestScopingIfStatement failed (err case).\nexpected=%v\ngot=%v",
				i, test.err, err)
		}
	}
}

func TestResolveCallExpression(t *testing.T) {
	vars := []struct {
		fl *ast.FunctionLiteral
		ce *ast.CallExpression
	}{
		{
			// case 1
			fl: &ast.FunctionLiteral{
				Name: &ast.Identifier{Name: "add"},
				Parameters: []*ast.ParameterLiteral{
					{
						Identifier: &ast.Identifier{Name: "a"},
						Type:       ast.IntType,
					},
					{
						Identifier: &ast.Identifier{Name: "b"},
						Type:       ast.IntType,
					},
				},
				Body:       nil,
				ReturnType: ast.IntType,
			},
			ce: &ast.CallExpression{
				Function: &ast.Identifier{
					Name: "add",
				},
				Arguments: []ast.Expression{
					&ast.Identifier{Name: "a"},
					&ast.Identifier{Name: "b"},
				},
			},
		},
	}

	tests := []struct {
		setupResolverFn
		input        *ast.CallExpression
		expected     *Resolver
		expectedType ObjectType
		err          error
	}{
		{
			// test case 1
			// function add(int a, int b)
			// ...
			// add(a, b)
			setupResolverFn: func() *Resolver {
				return &Resolver{
					scope: &Scope{
						store: map[string]Object{
							"add": &Function{Name: "add"},
						},
					},
					types: map[ast.Expression]ObjectType{
						&ast.Identifier{Name: "a"}: IntegerObject,
						&ast.Identifier{Name: "b"}: IntegerObject,
					},
					defs: map[*ast.Identifier]Object{
						vars[0].fl.Name: &Function{Name: "add"},
					},
					fns: map[string]*ast.FunctionLiteral{
						"add": vars[0].fl,
					},
				}
			},
			input: vars[0].ce,
			expected: &Resolver{
				scope: &Scope{
					store: map[string]Object{
						"add": &Function{Name: "add"},
					},
					child: []*Scope{
						{
							store: map[string]Object{
								"a": &Integer{Name: &ast.Identifier{Name: "a"}},
								"b": &Integer{Name: &ast.Identifier{Name: "b"}},
							},
							child: []*Scope{
								{
									store: map[string]Object{},
									child: []*Scope{},
								},
							},
						},
					},
				},
				types: map[ast.Expression]ObjectType{
					vars[0].ce.Arguments[0].(*ast.Identifier): IntegerObject,
					vars[0].ce.Arguments[1].(*ast.Identifier): IntegerObject,
				},
				defs: map[*ast.Identifier]Object{
					vars[0].ce.Arguments[0].(*ast.Identifier): &Integer{Name: &ast.Identifier{Name: "a"}},
					vars[0].ce.Arguments[1].(*ast.Identifier): &Integer{Name: &ast.Identifier{Name: "b"}},
				},
				fns: map[string]*ast.FunctionLiteral{
					"add": vars[0].fl,
				},
			},
			expectedType: IntegerObject,
			err:          nil,
		},
	}

	for i, test := range tests {
		r := test.setupResolverFn()
		if test.expected != nil {
			postMakeScope(test.expected.scope)
		}
		objType, err := ResolveCallExpression(test.input, r)
		if err == nil && objType != test.expectedType {
			t.Fatalf("[test %d] - TestResolveCallExpression failed.(wrong obj type).\n"+
				"expected: %v\ngot:%v", i, test.expectedType, objType)
		}
		if err == nil && !reflect.DeepEqual(r, test.expected) {
			t.Fatalf("[test %d] - TestResolveCallExpression failed.(wrong resolver).\n"+
				"expected: %v\ngot:%v", i, test.expected, r)
		}
		if err != nil && !reflect.DeepEqual(err, test.err) {
			t.Fatalf("[test %d] - TestResolveCallExpression failed.(wrong err).\n"+
				"expected: %v\ngot:%v", i, test.err, err)
		}

	}
}
