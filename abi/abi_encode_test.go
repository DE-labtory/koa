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
	"bytes"
	"testing"

	"github.com/DE-labtory/koa/abi"
)

func TestEncode(t *testing.T) {
	test := struct {
		args1 int
		args2 string
		args3 int
	}{
		args1: 50,
		args2: "HelloKOA",
		args3: 256,
	}

	testExpected := []byte{
		0x00, 0x00, 0x00, 0x0C,
		0x00, 0x00, 0x00, 0x11,
		0x00, 0x00, 0x00, 0x1D,
		0x00, 0x00, 0x00, 0x01,
		0x32,
		0x00, 0x00, 0x00, 0x08,
		0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x4b, 0x4f, 0x41,
		0x00, 0x00, 0x00, 0x02,
		0x01, 0x00,
	}

	encodedParams, err := abi.Encode(test.args1, test.args2, test.args3)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(testExpected, encodedParams) {
		t.Error("There is a problem with Encode")
	}
}

func TestSelector(t *testing.T) {
	tests := []struct {
		input  string
		output []byte
	}{
		{
			input:  "foo(uint128,bytes32,bool[2])",
			output: []byte{0x5C, 0x8D, 0x62, 0x96},
		},
		{
			input:  "foo(string,bytes32,bool[2])",
			output: []byte{0xCC, 0xBE, 0xD7, 0x8C},
		},
		{
			input:  "transfer(string,string,int)",
			output: []byte{0x82, 0x00, 0x55, 0x85},
		},
	}

	for i, test := range tests {
		selector := abi.Selector(test.input)
		if bytes.Compare(selector, test.output) != 0 {
			t.Fatalf("test[%d] - Selector() result wrong expected=%x, got=%x", i, test.output, selector)
		}
	}
}
