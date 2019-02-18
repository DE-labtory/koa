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

package translate_test

import (
	"testing"

	"github.com/DE-labtory/koa/abi"
	"github.com/DE-labtory/koa/translate"
)

// TODO: implement test cases :-)
func TestCompileContract(t *testing.T) {

}

func TestFuncMap_Declare(t *testing.T) {
	tests := []struct {
		signature string
		asm       translate.Asm
		expectPC  int
		expectLen int
	}{
		{
			signature: "foo()",
			asm: translate.Asm{
				AsmCodes: []translate.AsmCode{
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
						Value:   "0001020304050607",
					},
				},
			},
			expectPC:  2,
			expectLen: 1,
		},
		{
			signature: "bar()",
			asm: translate.Asm{
				AsmCodes: []translate.AsmCode{
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
						Value:   "0001020304050607",
					},
					{
						RawByte: []byte{0x21},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
						Value:   "0001020304050607",
					},
					{
						RawByte: []byte{0x01},
						Value:   "Add",
					},
				},
			},
			expectPC:  5,
			expectLen: 2,
		},
	}

	funcMap := translate.FuncMap{}
	for i, test := range tests {
		funcMap.Declare(test.signature, test.asm)

		result := funcMap[string(abi.Selector(test.signature))]

		if test.expectPC != result {
			t.Fatalf("test[%d] - Declare() pc result wrong. expected=%d, got=%d", i, test.expectPC, result)
		}

		if test.expectLen != len(funcMap) {
			t.Fatalf("test[%d] - Declare() FuncMap result wrong. expected=%d, got=%d", i, test.expectLen, len(funcMap))
		}
	}
}
