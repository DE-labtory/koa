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

import "testing"

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
		scope        Scope
		want         string
		expectedSym  Symbol
		expectedBool bool
	}{
		{
			Scope{
				map[string]Symbol{
					"a": &Integer{0},
					"b": &Integer{1},
				},
				nil,
			},
			"a",
			&Integer{0},
			true,
		},
		{
			Scope{
				map[string]Symbol{
					"a": &Integer{0},
					"b": &Integer{1},
				},
				&Scope{
					map[string]Symbol{
						"c": &String{"abc"},
					},
					nil,
				},
			},
			"c",
			&String{"abc"},
			true,
		},
		{
			Scope{
				map[string]Symbol{
					"a": &Integer{0},
					"b": &Integer{1},
				},
				nil,
			},
			"c",
			nil,
			false,
		},
	}

	for i, test := range tests {
		sym, ok := test.scope.Get(test.want)
		if sym != nil && test.expectedSym.String() != sym.String() {
			t.Fatalf("test[%d] testScopeGetter() returns invalid symbol.\n"+
				"expected=%s\n"+
				"got=%s", i, test.expectedSym.String(), sym.String())
		}
		if ok != test.expectedBool {
			t.Fatalf("test[%d] testScopeGetter() returns invalid ok.\n"+
				"expected=%v\n"+
				"got=%v", i, test.expectedBool, ok)
		}
	}
}

// TODO implement me w/ test cases :-)
func TestScopeSetter(t *testing.T) {

}

// TODO implement me w/ test cases :-)
func TestScopeString(t *testing.T) {

}
