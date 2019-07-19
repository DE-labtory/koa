/*
 * Copyright 2018-2019 De-labtory
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
	"testing"

	"github.com/DE-labtory/koa/ast"
)

func TestNewEnclosedScope(t *testing.T) {
	outer := NewScope()
	s := NewEnclosedScope(outer)

	if s.outer != outer {
		t.Fatalf("testNewEnclosedScope() failed. outer must be set")
	}

	if len(s.store) > 0 {
		t.Fatalf("testNewEnclosedScope() failed. store's size must be 0")
	}
}

func TestNewScope(t *testing.T) {
	s := NewScope()
	if s.outer != nil {
		t.Fatalf("testNewScope() failed. outer must be nil")
	}

	if len(s.store) > 0 {
		t.Fatalf("testNewScope() failed. store's size must be 0")
	}
}

func TestScopeGetter(t *testing.T) {
	tests := []struct {
		scope       Scope
		want        string
		expectedSym Symbol
	}{
		{
			Scope{
				store: map[string]Symbol{
					"a": &Integer{&ast.Identifier{Name: "a"}},
					"b": &Integer{&ast.Identifier{Name: "b"}},
				},
			},
			"a",
			&Integer{&ast.Identifier{Name: "a"}},
		},
		{
			Scope{
				store: map[string]Symbol{
					"a": &Integer{&ast.Identifier{Name: "a"}},
					"b": &Integer{&ast.Identifier{Name: "b"}},
				},
				outer: &Scope{
					store: map[string]Symbol{
						"c": &String{&ast.Identifier{Name: "c"}},
					},
				},
			},
			"c",
			&String{&ast.Identifier{Name: "c"}},
		},
		{
			Scope{
				store: map[string]Symbol{
					"a": &Integer{&ast.Identifier{Name: "a"}},
					"b": &Integer{&ast.Identifier{Name: "b"}},
				},
			},
			"c",
			nil,
		},
	}

	for i, test := range tests {
		sym := test.scope.Get(test.want)
		if sym != nil && test.expectedSym.String() != sym.String() {
			t.Fatalf("test[%d] testScopeGetter() returns invalid symbol.\n"+
				"expected=%s\n"+
				"got=%s", i, test.expectedSym.String(), sym.String())
		}

		if sym != nil && test.expectedSym == nil {
			t.Fatalf("test[%d] testScopeGetter() returns invalid symbol.\n"+
				"expected=nil\n"+
				"got=%s", i, sym.String())
		}
	}
}

func TestScopeSetter(t *testing.T) {
	tests := []struct {
		Scope  *Scope
		Name   string
		Symbol Symbol
	}{
		{
			&Scope{
				store: map[string]Symbol{},
				outer: &Scope{},
			},
			"testInt",
			&Integer{&ast.Identifier{Name: "testInt"}},
		},
		{
			&Scope{
				store: map[string]Symbol{},
				outer: &Scope{},
			},
			"testBool",
			&Boolean{&ast.Identifier{Name: "testBool"}},
		},
		{
			&Scope{
				store: map[string]Symbol{},
				outer: &Scope{},
			},
			"testString",
			&String{&ast.Identifier{Name: "testString"}},
		},
	}

	for i, test := range tests {
		symbol := test.Scope.Set(test.Name, test.Symbol)
		if symbol.String() != test.Symbol.String() {
			t.Fatalf("test[%d] - TestScopeSetter() wrong result.\n"+
				"expected=%s\n"+
				"got=%s", i, test.Symbol.String(), symbol.String())
		}

		if sym := test.Scope.Get(test.Name); sym.String() != test.Symbol.String() {
			t.Fatalf("test[%d] - TestScopeSetter() must set in scope store", i)
		}
	}

}

func TestScopeString(t *testing.T) {
	tests := []struct {
		scope    Scope
		expected string
	}{
		{
			Scope{
				store: map[string]Symbol{
					"a": &Integer{&ast.Identifier{Name: "a"}},
					"b": &Integer{&ast.Identifier{Name: "b"}},
				},
			},
			"[ Scope ]\nkey:a, symbol:a\nkey:b, symbol:b\n",
		},
		{
			Scope{
				store: map[string]Symbol{
					"a": &Integer{&ast.Identifier{Name: "a"}},
					"b": &Integer{&ast.Identifier{Name: "b"}},
				},
				outer: &Scope{
					store: map[string]Symbol{
						"c": &String{&ast.Identifier{Name: "c"}},
					},
					outer: nil,
				},
			},
			"[ Scope ]\nkey:a, symbol:a\nkey:b, symbol:b\n[ Scope ]\nkey:c, symbol:c\n",
		},
		{
			Scope{
				store: map[string]Symbol{
					"c": &Integer{&ast.Identifier{Name: "c"}},
					"b": &Integer{&ast.Identifier{Name: "b"}},
					"a": &Integer{&ast.Identifier{Name: "a"}},
				},
				outer: &Scope{
					store: map[string]Symbol{
						"d": &String{&ast.Identifier{Name: "d"}},
					},
				},
			},
			"[ Scope ]\nkey:a, symbol:a\nkey:b, symbol:b\nkey:c, symbol:c\n[ Scope ]\nkey:d, symbol:d\n",
		},
	}

	for i, test := range tests {
		str := test.scope.String()
		if str != test.expected {
			t.Fatalf("test[%d] - TestScopeString() wrong result.\n"+
				"expected :\n%s\n"+
				"got :\n%s\n", i, test.expected, str)
		}
	}
}
