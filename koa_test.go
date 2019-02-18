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

package koa

import (
	"errors"
	"os"
	"testing"

	"github.com/DE-labtory/koa/opcode"
	"github.com/DE-labtory/koa/translate"
)

type testData struct {
	fileName string
	asm      *translate.Asm
	err      error
}

func defineAsm() []testData {
	return []testData{
		{
			fileName: "test/hello.koa",
			asm: &translate.Asm{
				AsmCodes: []translate.AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x22, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x21, 0x22},
						Value:   "2268656c6c6f2122",
					},
					{
						RawByte: []byte{byte(opcode.Returning)},
						Value:   "Returning",
					},
				},
			},
			err: nil,
		},
		{
			fileName: "test/jun.koa",
			asm:      nil,
			err:      errors.New("[junbeomlee] definition doesn't exist"),
		},
		{
			fileName: "test/add1.koa",
			asm: &translate.Asm{
				AsmCodes: []translate.AsmCode{
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
						Value:   "0000000000000001",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
						Value:   "0000000000000002",
					},
					{
						RawByte: []byte{byte(opcode.Add)},
						Value:   "Add",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					{
						RawByte: []byte{byte(opcode.Mstore)},
						Value:   "Mstore",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03},
						Value:   "0000000000000003",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04},
						Value:   "0000000000000004",
					},
					{
						RawByte: []byte{byte(opcode.Add)},
						Value:   "Add",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Mstore)},
						Value:   "Mstore",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
						Value:   "0000000000000005",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
						Value:   "0000000000000006",
					},
					{
						RawByte: []byte{byte(opcode.Mul)},
						Value:   "Mul",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10},
						Value:   "0000000000000010",
					},
					{
						RawByte: []byte{byte(opcode.Mstore)},
						Value:   "Mstore",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04},
						Value:   "0000000000000004",
					},
					{
						RawByte: []byte{byte(opcode.Div)},
						Value:   "Div",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18},
						Value:   "0000000000000018",
					},
					{
						RawByte: []byte{byte(opcode.Mstore)},
						Value:   "Mstore",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
						Value:   "0000000000000000",
					},
					{
						RawByte: []byte{byte(opcode.Mload)},
						Value:   "Mload",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Mload)},
						Value:   "Mload",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10},
						Value:   "0000000000000010",
					},
					{
						RawByte: []byte{byte(opcode.Mload)},
						Value:   "Mload",
					},
					{
						RawByte: []byte{byte(opcode.Sub)},
						Value:   "Sub",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
						Value:   "0000000000000008",
					},
					{
						RawByte: []byte{byte(opcode.Push)},
						Value:   "Push",
					},
					{
						RawByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18},
						Value:   "0000000000000018",
					},
					{
						RawByte: []byte{byte(opcode.Mload)},
						Value:   "Mload",
					},
					{
						RawByte: []byte{byte(opcode.Mul)},
						Value:   "Mul",
					},
					{
						RawByte: []byte{byte(opcode.Add)},
						Value:   "Add",
					},
					{
						RawByte: []byte{byte(opcode.Returning)},
						Value:   "Returning",
					},
				},
			},
			err: nil,
		},
	}
}

func readFile(fileName string) (string, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, os.FileMode(644))
	if err != nil {
		return "", err
	}

	fs, err := file.Stat()
	if err != nil {
		return "", err
	}

	data := make([]byte, fs.Size())
	if err != nil {
		return "", err
	}

	_, err = file.Read(data)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func TestCompile(t *testing.T) {
	tests := defineAsm()

	for i, test := range tests {
		str, err := readFile(test.fileName)
		if err != nil {
			continue
		}

		asm, err := Compile(str)

		// TODO : after implements all functions in compile
		if err != nil && err.Error() != test.err.Error() {
			t.Fatalf("[test %d] - TestCompile() returns wrong error.\nexpected=%v\ngot=%v", i, test.err, err)
		}

		if test.asm != nil && !asm.Equal(*test.asm) {
			t.Fatalf("[test %d] - TestCompile() returns wrong asm.\nexpected=%v\ngot=%v", i, test.asm, asm)
		}
	}
}
