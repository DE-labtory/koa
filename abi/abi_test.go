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
	"testing"

	"github.com/DE-labtory/koa/abi"
)

func TestNew(t *testing.T) {
	abiJSON := `[
	{
		"name" : "foo",
		"arguments" : [
			{
				"name" : "first",
				"type" : "int256"
			},
			{
				"name" : "second",
				"type" : "string"
			},
			{
				"name" : "third",
				"type" : "byte"
			}
		],
		"output" : {
			"name" : "returnValue",
			"type" : "int256"
		}
	},
	{
		"name" : "var",
		"arguments" : [
			{
				"name" : "first",
				"type" : "int256"
			},
			{
				"name" : "second",
				"type" : "string"
			},
			{
				"name" : "third",
				"type" : "byte"
			}
		],
		"output" : {
			"name" : "returnValue",
			"type" : "int256"
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

// TODO: implement test cases :-)
func TestExtractAbiFromFunction(t *testing.T) {

}
