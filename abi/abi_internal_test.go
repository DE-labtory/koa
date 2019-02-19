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

package abi

import (
	"testing"

	"github.com/DE-labtory/koa/ast"
)

func TestConvertAstTypeToAbi(t *testing.T) {
	tests := []struct {
		p      ast.ParameterLiteral
		expect Type
		err    error
	}{
		{
			p: ast.ParameterLiteral{
				Identifier: &ast.Identifier{
					Name: "a",
				},
				Type: ast.IntType,
			},
			expect: Type{
				Type: "int",
			},
			err: nil,
		},
		{
			p: ast.ParameterLiteral{
				Identifier: &ast.Identifier{
					Name: "b",
				},
				Type: ast.BoolType,
			},
			expect: Type{
				Type: "bool",
			},
			err: nil,
		},
		{
			p: ast.ParameterLiteral{
				Identifier: &ast.Identifier{
					Name: "c",
				},
				Type: ast.StringType,
			},
			expect: Type{
				Type: "string",
			},
			err: nil,
		},
	}

	for i, test := range tests {
		typeResult, err := convertAstTypeToAbi(test.p.Type)

		if typeResult.Type != test.expect.Type {
			t.Fatalf("test[%d] - TestConvertAstTypeToAbi() result wrong. expected=%s, got=%s", i, test.expect.Type, typeResult.Type)
		}

		if err != nil && err != test.err {
			t.Fatalf("test[%d] - TestConvertAstTypeToAbi() error wrong. expected=%v, got=%v", i, test.err, err)
		}
	}

}
