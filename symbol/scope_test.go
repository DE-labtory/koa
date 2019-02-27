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

package symbol

import (
	"reflect"
	"testing"

	"github.com/DE-labtory/koa/ast"
)

type setupScopeFn func() *Scope

func defaultScope() *Scope {
	return NewScope()
}

func TestNewScope(t *testing.T) {

}

func TestNewEnclosedScope(t *testing.T) {

}

func TestScopingFunction(t *testing.T) {

}

func TestScopingFunctionParameter(t *testing.T) {
	tests := []struct {
		setupScopeFn
		input    []*ast.ParameterLiteral
		expected *Scope
		err      error
	}{
		{
			setupScopeFn: defaultScope,
			input: []*ast.ParameterLiteral{
				{
					Identifier: &ast.Identifier{
						Name: "a",
					},
					Type: ast.IntType,
				},
			},
			expected: &Scope{
				store: map[string]Symbol{
					"a": &Integer{Name: &ast.Identifier{Name: "a"}},
				},
			},
			err: nil,
		},
		{
			setupScopeFn: defaultScope,
			input: []*ast.ParameterLiteral{
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
			expected: &Scope{
				store: map[string]Symbol{
					"a": &Integer{Name: &ast.Identifier{Name: "a"}},
					"b": &Boolean{Name: &ast.Identifier{Name: "b"}},
					"c": &String{Name: &ast.Identifier{Name: "c"}},
				},
			},
			err: nil,
		},
		{
			setupScopeFn: defaultScope,
			input: []*ast.ParameterLiteral{
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
			expected: nil,
			err:      DupError{Identifier: ast.Identifier{Name: "a"}},
		},
	}

	for i, test := range tests {
		s := test.setupScopeFn()
		err := ScopingFunctionParameter(test.input, s)

		if err == nil && !reflect.DeepEqual(test.expected, s) {
			t.Fatalf("test [%d] - TestScopingAssignStatement failed.\nexpected=%v\ngot=%v",
				i, test.expected, s)
		}
		if err != nil && reflect.DeepEqual(test.err, err) {
			t.Fatalf("test [%d] - TestScopingAssignStatement failed (err case).\nexpected=%v\ngot=%v",
				i, test.err, err)
		}
	}
}

func TestScopingFunctionBody(t *testing.T) {

}

func TestScopingAssignStatement(t *testing.T) {
	tests := []struct {
		setupScopeFn
		input    *ast.AssignStatement
		expected *Scope
		err      error
	}{
		{
			setupScopeFn: defaultScope,
			input: &ast.AssignStatement{
				Type:     ast.IntType,
				Variable: ast.Identifier{Name: "a"},
				Value: &ast.IntegerLiteral{
					Value: 0,
				},
			},
			expected: &Scope{
				store: map[string]Symbol{
					"a": &Integer{Name: &ast.Identifier{Name: "a"}},
				},
			},
		},
		{
			setupScopeFn: defaultScope,
			input: &ast.AssignStatement{
				Type:     ast.StringType,
				Variable: ast.Identifier{Name: "str"},
				Value: &ast.StringLiteral{
					Value: "testString",
				},
			},
			expected: &Scope{
				store: map[string]Symbol{
					"str": &String{Name: &ast.Identifier{Name: "str"}},
				},
			},
		},
		{
			setupScopeFn: func() *Scope {
				return &Scope{
					store: map[string]Symbol{
						"str": &String{Name: &ast.Identifier{Name: "str"}},
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
				Identifier: ast.Identifier{Name: "str"},
			},
		},
	}

	for i, test := range tests {
		s := test.setupScopeFn()
		err := ScopingAssignStatement(test.input, s)

		if err == nil && !reflect.DeepEqual(s, test.expected) {
			t.Fatalf("test [%d] - TestScopingAssignStatement failed.\nexpected=%v\ngot=%v",
				i, test.expected, s)
		}
		if err != nil && !reflect.DeepEqual(err, test.err) {
			t.Fatalf("test [%d] - TestScopingAssignStatement failed (err case).\nexpected=%v\ngot=%v",
				i, test.err, err)
		}
	}
}

func TestScopingIfStatement(t *testing.T) {
	tests := []struct {
		setupScopeFn
		input    *ast.IfStatement
		expected *Scope
		err      error
	}{
		{
			setupScopeFn: defaultScope,
			input: &ast.IfStatement{
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
			},
			expected: &Scope{
				store: map[string]Symbol{
					"a": &Integer{Name: &ast.Identifier{Name: "a"}},
				},
			},
		},
		{
			setupScopeFn: defaultScope,
			input: &ast.IfStatement{
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
			expected: &Scope{
				store: map[string]Symbol{
					"a": &Integer{Name: &ast.Identifier{Name: "a"}},
					"b": &Integer{Name: &ast.Identifier{Name: "b"}},
				},
			},
		},
		{
			setupScopeFn: defaultScope,
			input: &ast.IfStatement{
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
			expected: &Scope{
				store: map[string]Symbol{
					"a": &Integer{Name: &ast.Identifier{Name: "a"}},
					"b": &Integer{Name: &ast.Identifier{Name: "b"}},
				},
			},
		},
	}

	for i, test := range tests {
		s := test.setupScopeFn()
		ScopingIfStatement(test.input, s)
		if test.expected != s {
			t.Fatalf("test [%d] - TestScopingAssignStatement failed", i)
		}
	}
}

func TestScopingIdentifier(t *testing.T) {
	tests := []struct {
		setupScopeFn
		idf      ast.Identifier
		ds       ast.DataStructure
		expected *Scope
		err      error
	}{
		{
			setupScopeFn: defaultScope,
			idf: ast.Identifier{
				Name: "a",
			},
			ds: ast.IntType,
			expected: &Scope{
				store: map[string]Symbol{
					"a": &Integer{Name: &ast.Identifier{Name: "a"}},
				},
			},
			err: nil,
		},
		{
			setupScopeFn: func() *Scope {
				return &Scope{
					store: map[string]Symbol{
						"a": &Boolean{Name: &ast.Identifier{Name: "a"}},
					},
				}
			},
			idf: ast.Identifier{
				Name: "a",
			},
			ds:       ast.IntType,
			expected: nil,
			err: DupError{
				Identifier: ast.Identifier{Name: "a"},
			},
		},
	}

	for i, test := range tests {
		s := test.setupScopeFn()
		err := ScopingIdentifier(test.idf, test.ds, s)
		if err == nil && !reflect.DeepEqual(s, test.expected) {
			t.Fatalf("test [%d] - TestScopingAssignStatement failed.\nexpected=%v\ngot=%v",
				i, test.expected, s)
		}
		if err != nil && !reflect.DeepEqual(err, test.err) {
			t.Fatalf("test [%d] - TestScopingAssignStatement failed (err case).\nexpected=%v\ngot=%v",
				i, test.err, err)
		}
	}
}
