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
	"testing"

	"github.com/DE-labtory/koa/ast"
)

func TestInteger(t *testing.T) {
	tests := []struct {
		input          Object
		expectedStr    string
		expectedSymbol ObjectType
	}{
		{&Integer{&ast.Identifier{Name: "testName"}}, "testName", IntegerObject},
		{&Integer{&ast.Identifier{Name: "a"}}, "a", IntegerObject},
		{&Integer{&ast.Identifier{Name: "b"}}, "b", IntegerObject},
	}

	for i, test := range tests {
		str := test.input.String()
		obj := test.input.Type()

		if str != test.expectedStr {
			t.Fatalf("test[%d] String() in Integer wrong result.\n"+
				"expected: %s\n"+
				"got: %s", i, test.expectedStr, str)
		}

		if obj != test.expectedSymbol {
			t.Fatalf("test[%d] Type() in Integer wrong result.\n"+
				"expected: %s\n"+
				"got: %s", i, test.expectedSymbol, obj)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		input          Object
		expectedStr    string
		expectedSymbol ObjectType
	}{
		{&String{&ast.Identifier{Name: "testName"}}, "testName", StringObject},
		{&String{&ast.Identifier{Name: "a"}}, "a", StringObject},
		{&String{&ast.Identifier{Name: "b"}}, "b", StringObject},
	}

	for i, test := range tests {
		str := test.input.String()
		obj := test.input.Type()

		if str != test.expectedStr {
			t.Fatalf("test[%d] String() in Integer wrong result.\n"+
				"expected: %s\n"+
				"got: %s", i, test.expectedStr, str)
		}

		if obj != test.expectedSymbol {
			t.Fatalf("test[%d] Type() in Integer wrong result.\n"+
				"expected: %s\n"+
				"got: %s", i, test.expectedSymbol, obj)
		}
	}
}

func TestBoolean(t *testing.T) {
	tests := []struct {
		input       Object
		expectedStr string
		expectedObj ObjectType
	}{
		{&Boolean{&ast.Identifier{Name: "testName"}}, "testName", BooleanObject},
		{&Boolean{&ast.Identifier{Name: "a"}}, "a", BooleanObject},
	}

	for i, test := range tests {
		str := test.input.String()
		obj := test.input.Type()

		if str != test.expectedStr {
			t.Fatalf("test[%d] String() in Integer wrong result.\n"+
				"expected: %s\n"+
				"got: %s", i, test.expectedStr, str)
		}

		if obj != test.expectedObj {
			t.Fatalf("test[%d] Type() in Integer wrong result.\n"+
				"expected: %s\n"+
				"got: %s", i, test.expectedObj, obj)
		}
	}
}

func TestFunction(t *testing.T) {
	tests := []struct {
		input          Object
		expectedStr    string
		expectedSymbol ObjectType
	}{
		{
			&Function{
				"add",
			},
			"add",
			FunctionObject,
		},
	}

	for i, test := range tests {
		str := test.input.String()
		obj := test.input.Type()

		if str != test.expectedStr {
			t.Fatalf("test[%d] String() in Integer wrong result.\n"+
				"expected: %s\n"+
				"got: %s", i, test.expectedStr, str)
		}

		if obj != test.expectedSymbol {
			t.Fatalf("test[%d] Type() in Integer wrong result.\n"+
				"expected: %s\n"+
				"got: %s", i, test.expectedSymbol, obj)
		}
	}
}
