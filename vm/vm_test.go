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

package vm

import (
	"testing"

	"reflect"

	"github.com/DE-labtory/koa/opcode"
	"github.com/stretchr/testify/assert"
)

func makeTestByteCode(slice ...interface{}) []byte {
	testByteCode := make([]byte, 0)
	for _, code := range slice {
		value := reflect.ValueOf(code)
		switch value.Kind() {
		case reflect.Slice:
			testByteCode = append(testByteCode, code.([]byte)...)
		case reflect.Uint8:
			testByteCode = append(testByteCode, code.(uint8))
		default:
		}
	}

	return testByteCode
}

func TestPush(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(1),
		uint8(opcode.Push), uint32ToBytes(2),
	)

	testExpected := []item{1, 2}

	stack, err := Execute(testByteCode)
	if err != nil {
		t.Error(err)
	}

	if len(stack.items) != 2 {
		t.Errorf("Invalid stack size - expected=%d, got =%d", len(testExpected), stack.len())
	}

	for i, item := range stack.items {
		assert.Equal(t, testExpected[i], item)
	}
}

func TestPush_invalid(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint8(opcode.Push),
		uint8(opcode.Push), uint32ToBytes(3),
	)

	_, err := Execute(testByteCode)

	if err != ErrInvalidOpcode {
		t.Error("The desired error was not found")
	}
}

func TestAdd(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(1),
		uint8(opcode.Push), uint32ToBytes(2),
		uint8(opcode.Add),
	)
	testExpected := item(3)

	stack, err := Execute(testByteCode)

	if err != nil {
		t.Error(err)
	}
	result := stack.pop()
	if testExpected != result {
		t.Errorf("stack.pop() result wrong - expected=%d, got=%d", testExpected, result)
	}
}
