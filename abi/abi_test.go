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

package abi_test

import (
	"reflect"
	"testing"

	"github.com/DE-labtory/koa/abi"
	"github.com/DE-labtory/koa/ast"
)

func TestNew(t *testing.T) {
	abiJSON := `[
	{
		"name" : "foo",
		"arguments" : [
			{
				"name" : "first",
				"type" : "int64"
			},
			{
				"name" : "second",
				"type" : "string"
			},
			{
				"name" : "true",
				"type" : "bool"
			}
		],
		"output" : {
			"name" : "returnValue",
			"type" : "int64"
		}
	},
	{
		"name" : "var",
		"arguments" : [
			{
				"name" : "first",
				"type" : "int64"
			},
			{
				"name" : "second",
				"type" : "string"
			},
			{
				"name" : "false",
				"type" : "bool"
			}
		],
		"output" : {
			"name" : "returnValue",
			"type" : "int64"
		}
	}
]
`

	ABI, err := abi.New(abiJSON)
	if err != nil {
		t.Error(err)
	}

	if len(ABI.Methods) != 2 {
		t.Error("Invalid JSON parsing!")
	}
}

func TestExtractAbiFromFunction(t *testing.T) {
	tests := []struct {
		f      ast.FunctionLiteral
		expect abi.Method
		err    error
	}{
		{
			f: ast.FunctionLiteral{
				Name: &ast.Identifier{
					Name: "add",
				},
				Parameters: []*ast.ParameterLiteral{
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
						Type: ast.IntType,
					},
				},
				ReturnType: ast.IntType,
			},
			expect: abi.Method{
				Name: "add",
				Arguments: []abi.Argument{
					{
						Name: "a",
						Type: abi.Type{
							Type: "int",
						},
					},
					{
						Name: "b",
						Type: abi.Type{
							Type: "int",
						},
					},
				},
				Output: abi.Argument{
					Name: "",
					Type: abi.Type{
						Type: "int",
					},
				},
			},
			err: nil,
		},
		// test void return function
		{
			f: ast.FunctionLiteral{
				Name: &ast.Identifier{
					Name: "add",
				},
				Parameters: []*ast.ParameterLiteral{},
				ReturnType: ast.VoidType,
			},
			expect: abi.Method{
				Name:      "add",
				Arguments: []abi.Argument{},
				Output: abi.Argument{
					Name: "",
					Type: abi.Type{
						Type: "void",
					},
				},
			},
			err: nil,
		},
	}

	for i, test := range tests {
		m, err := abi.ExtractAbiFromFunction(test.f)

		if !reflect.DeepEqual(m, test.expect) {
			t.Fatalf("test[%d] - TestExtractAbiFromFunction() result wrong.\nexpected=%v,\ngot=%v", i, test.expect, m)
		}

		if err != nil && err != test.err {
			t.Fatalf("test[%d] - TestExtractAbiFromFunction() error wrong.\nexpected=%v,\ngot=%v", i, test.err, err)
		}
	}
}
