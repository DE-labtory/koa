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
	"bytes"
)

func makeTestABI() abi.ABI {
	abiJSON := `[
	{
		"name" : "foo",
		"arguments" : [
			{
				"name" : "first",
				"type" : "int64"
			},
			{
				"name" : "second",
				"type" : "string"
			},
			{
				"name" : "true",
				"type" : "bool"
			}
		],
		"output" : {
			"name" : "returnValue",
			"type" : "int64"
		}
	},
	{
		"name" : "var",
		"arguments" : [
			{
				"name" : "first",
				"type" : "int64"
			},
			{
				"name" : "second",
				"type" : "string"
			},
			{
				"name" : "false",
				"type" : "bool"
			}
		],
		"output" : {
			"name" : "returnValue",
			"type" : "int64"
		}
	}
]
`

	ABI, err := abi.New(abiJSON)
	if err != nil {
		return abi.ABI{}
	}

	return ABI
}

func TestMethod_Signature(t *testing.T) {
	ABI := makeTestABI()

	testExpected := []struct {
		signature string
	} {
		{
			signature: "foo(int64,string,bool)",
		},
		{
			signature: "var(int64,string,bool)",
		},
	}

	for index, method := range ABI.Methods {
		if method.Signature() != testExpected[index].signature {
			t.Errorf("Invalid Signature - expected = %s, got = %s", testExpected[index].signature, method.Signature())
		}
	}
}

func TestMethod_ID(t *testing.T) {
	ABI := makeTestABI()

	testExpected := []struct {
		id []byte
	} {
		{
			id: []byte{0x94, 0x1f, 0xe5, 0xf2},
		},
		{
			id: []byte{0x7a, 0x89, 0xe4, 0x4b},
		},
	}

	for index, method := range ABI.Methods {
		if !bytes.Equal(method.ID(), testExpected[index].id) {
			t.Errorf("Invalid ID - expected = %x, got = %x", testExpected[index].id, method.ID())
		}
	}
}