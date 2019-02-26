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
	"github.com/DE-labtory/koa/ast"
	"reflect"
	"testing"
)

func TestNewScope(t *testing.T) {

}

func TestNewEnclosedScope(t *testing.T) {

}

func TestGenerateScope(t *testing.T) {
	tests := []struct {
		contract ast.Contract
		expected Scope
	}{
		{
			// contract which has one function named "foo", and it has one integer parameter named "a"
			contract: ast.Contract{
				Functions: []*ast.FunctionLiteral{
					{
						Name: &ast.Identifier{Name: "foo"},
						Parameters: []*ast.ParameterLiteral{
							{
								Type:       ast.IntType,
								Identifier: &ast.Identifier{Name: "a"},
							},
						},
					},
				},
			},
			expected: Scope{
				store: map[string]Symbol{
					"foo": &Function{Name: "foo"},
				},
				parent: nil,
				child: []*Scope{
					{
						store: map[string]Symbol{
							"a": &Integer{Name: &ast.Identifier{Name: "a"}},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		scope := GenerateScope(&test.contract)
		if !reflect.DeepEqual(test.expected, scope) {
			t.Fatalf("")
		}
	}
}
