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

package opcode_test

import (
	"testing"

	"github.com/DE-labtory/koa/opcode"
)

func TestType_String(t *testing.T) {
	tests := []struct {
		input    opcode.Type
		expected string
	}{
		{
			opcode.Add,
			"Add",
		},
		{
			opcode.Mul,
			"Mul",
		},
		{
			opcode.Sub,
			"Sub",
		},
		{
			opcode.Div,
			"Div",
		},
		{
			opcode.Mod,
			"Mod",
		},
		{
			opcode.LT,
			"LT",
		},
		{
			opcode.GT,
			"GT",
		},
		{
			opcode.EQ,
			"EQ",
		},
		{
			opcode.NOT,
			"NOT",
		},
		{
			opcode.Pop,
			"Pop",
		},
		{
			opcode.Push,
			"Push",
		},
		{
			opcode.Mload,
			"Mload",
		},
		{
			opcode.Mstore,
			"Mstore",
		},
		{
			0x24,
			"String() error - Not defined opcode",
		},
	}

	for i, test := range tests {
		str, err := test.input.String()
		if str != "" && str != test.expected {
			t.Fatalf("test[%d] - TestType_String() wrong result.\n"+
				"expected: %s\n"+
				"got: %s", i, test.expected, str)
		}
		if err != nil && err.Error() != test.expected {
			t.Fatalf("test[%d] - TestType_String() wrong error.\n"+
				"expected: %s\n"+
				"got: %s", i, test.expected, err.Error())
		}
	}
}
