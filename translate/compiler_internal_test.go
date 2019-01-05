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
	"testing"

	"github.com/DE-labtory/koa/ast"
)

// TODO: implement test cases :-)
func TestCompiler_emit(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompiler_compileIdentifier(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompiler_compileString(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCompiler_compileInteger(t *testing.T) {

}

func TestCompiler_compileBoolean(t *testing.T) {
	c := Compiler{}

	tests := []struct {
		node     ast.BooleanLiteral
		expected error
	}{
		{
			node: ast.BooleanLiteral{
				Value: true,
			},
			expected: nil,
		},
		{
			node: ast.BooleanLiteral{
				Value: false,
			},
			expected: nil,
		},
	}

	for i, test := range tests {
		n := test.node
		err := c.compileBoolean(n)

		if err != nil {
			t.Fatalf("test[%d] - compileBoolean() had error. err=%v", i, err)
		}
	}
}

// TODO: implement test cases :-)
func TestCompiler_compilePrefixExpression(t *testing.T) {

}
