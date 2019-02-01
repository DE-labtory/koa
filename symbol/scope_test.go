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
	"errors"
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
		expectedErr error
	}{
		{
			Scope{
				map[string]Symbol{
					"a": &Integer{&ast.Identifier{Value: "a"}},
					"b": &Integer{&ast.Identifier{Value: "b"}},
				},
				nil,
			},
			"a",
			&Integer{&ast.Identifier{Value: "a"}},
			nil,
		},
		{
			Scope{
				map[string]Symbol{
					"a": &Integer{&ast.Identifier{Value: "a"}},
					"b": &Integer{&ast.Identifier{Value: "b"}},
				},
				&Scope{
					map[string]Symbol{
						"c": &String{&ast.Identifier{Value: "c"}},
					},
					nil,
				},
			},
			"c",
			&String{&ast.Identifier{Value: "c"}},
			nil,
		},
		{
			Scope{
				map[string]Symbol{
					"a": &Integer{&ast.Identifier{Value: "a"}},
					"b": &Integer{&ast.Identifier{Value: "b"}},
				},
				nil,
			},
			"c",
			nil,
			errors.New("c is not defined"),
		},
	}

	for i, test := range tests {
		sym, err := test.scope.Get(test.want)
		if sym != nil && test.expectedSym.String() != sym.String() {
			t.Fatalf("test[%d] testScopeGetter() returns invalid symbol.\n"+
				"expected=%s\n"+
				"got=%s", i, test.expectedSym.String(), sym.String())
		}
		if err != nil && err.Error() != test.expectedErr.Error() {
			t.Fatalf("test[%d] testScopeGetter() returns invalid error.\n"+
				"expected=%s\n"+
				"got=%s", i, test.expectedErr.Error(), err.Error())
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
				map[string]Symbol{},
				&Scope{},
			},
			"testInt",
			&Integer{&ast.Identifier{Value: "testInt"}},
		},
		{
			&Scope{
				map[string]Symbol{},
				&Scope{},
			},
			"testBool",
			&Boolean{&ast.Identifier{Value: "testBool"}},
		},
		{
			&Scope{
				map[string]Symbol{},
				&Scope{},
			},
			"testString",
			&String{&ast.Identifier{Value: "testString"}},
		},
	}

	for i, test := range tests {
		symbol := test.Scope.Set(test.Name, test.Symbol)
		if symbol != test.Symbol {
			t.Fatalf("test[%d] - TestScopeSetter() wrong result.\n"+
				"expected=%s\n"+
				"got=%s", i, test.Symbol.String(), symbol.String())
		}

		if _, err := test.Scope.Get(test.Name); err != nil {
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
				map[string]Symbol{
					"a": &Integer{&ast.Identifier{Value: "a"}},
					"b": &Integer{&ast.Identifier{Value: "b"}},
				},
				nil,
			},
			"[ Scope ]\nkey:a, symbol:a\nkey:b, symbol:b\n",
		},
		{
			Scope{
				map[string]Symbol{
					"a": &Integer{&ast.Identifier{Value: "a"}},
					"b": &Integer{&ast.Identifier{Value: "b"}},
				},
				&Scope{
					map[string]Symbol{
						"c": &String{&ast.Identifier{Value: "c"}},
					},
					nil,
				},
			},
			"[ Scope ]\nkey:a, symbol:a\nkey:b, symbol:b\n[ Scope ]\nkey:c, symbol:c\n",
		},
		{
			Scope{
				map[string]Symbol{
					"c": &Integer{&ast.Identifier{Value: "c"}},
					"b": &Integer{&ast.Identifier{Value: "b"}},
					"a": &Integer{&ast.Identifier{Value: "a"}},
				},
				&Scope{
					map[string]Symbol{
						"d": &String{&ast.Identifier{Value: "d"}},
					},
					nil,
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
