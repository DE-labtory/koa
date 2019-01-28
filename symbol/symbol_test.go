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
)

func TestInteger(t *testing.T) {
	tests := []struct {
		input          Symbol
		expectedStr    string
		expectedSymbol SymbolType
	}{
		{&Integer{10}, "10", IntegerSymbol},
		{&Integer{20}, "20", IntegerSymbol},
		{&Integer{30}, "30", IntegerSymbol},
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
		input          Symbol
		expectedStr    string
		expectedSymbol SymbolType
	}{
		{&String{"abcd"}, "abcd", StringSymbol},
		{&String{"Hi"}, "Hi", StringSymbol},
		{&String{"@Aa1231"}, "@Aa1231", StringSymbol},
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
		input       Symbol
		expectedStr string
		expectedObj SymbolType
	}{
		{&Boolean{true}, "true", BooleanSymbol},
		{&Boolean{false}, "false", BooleanSymbol},
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
		input          Symbol
		expectedStr    string
		expectedSymbol SymbolType
	}{
		{
			&Function{
				"add",
				&Scope{},
			},
			"func add",
			FunctionSymbol,
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
