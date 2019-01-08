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

func TestAdd_negative(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(0xFFFFFFEC), // 0xFFFFFFEC : -20
		uint8(opcode.Push), uint32ToBytes(30),
		uint8(opcode.Add),
	)
	testExpected := item(10)

	stack, err := Execute(testByteCode)

	if err != nil {
		t.Error(err)
	}
	result := stack.pop()
	if testExpected != result {
		t.Errorf("stack.pop() result wrong - expected=%d, got=%d", testExpected, result)
	}
}

func TestMul(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(3),
		uint8(opcode.Push), uint32ToBytes(5),
		uint8(opcode.Mul),
	)
	testExpected := item(15)

	stack, err := Execute(testByteCode)

	if err != nil {
		t.Error(err)
	}
	result := stack.pop()
	if testExpected != result {
		t.Errorf("stack.pop() result wrong - expected=%d, got=%d", testExpected, result)
	}
}

func TestMul_negative(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(0XFFFFFFFD), // FFFFFFFD : -3
		uint8(opcode.Push), uint32ToBytes(5),
		uint8(opcode.Mul),
	)
	testExpected := item(0xFFFFFFF1) // FFFFFFF1 : -15

	stack, err := Execute(testByteCode)

	if err != nil {
		t.Error(err)
	}
	result := stack.pop()
	if testExpected != result {
		t.Errorf("stack.pop() result wrong - expected=%d, got=%d", testExpected, result)
	}
}

func TestSub(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(50),
		uint8(opcode.Push), uint32ToBytes(20),
		uint8(opcode.Sub),
	)
	testExpected := item(30)

	stack, err := Execute(testByteCode)

	if err != nil {
		t.Error(err)
	}
	result := stack.pop()
	if testExpected != result {
		t.Errorf("stack.pop() result wrong - expected=%d, got=%d", testExpected, result)
	}
}

func TestSub_negative(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(0xFFFFFFEC), // 0xFFFFFFEC : -20
		uint8(opcode.Push), uint32ToBytes(50),
		uint8(opcode.Sub),
	)
	testExpected := item(0xFFFFFFBA) // 0xFFFFFFBA : -70

	stack, err := Execute(testByteCode)

	if err != nil {
		t.Error(err)
	}
	result := stack.pop()
	if testExpected != result {
		t.Errorf("stack.pop() result wrong - expected=%d, got=%d", testExpected, result)
	}
}

func TestDiv(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(14),
		uint8(opcode.Push), uint32ToBytes(5),
		uint8(opcode.Div),
	)
	testExpected := item(2)

	stack, err := Execute(testByteCode)

	if err != nil {
		t.Error(err)
	}
	result := stack.pop()
	if testExpected != result {
		t.Errorf("stack.pop() result wrong - expected=%d, got=%d", testExpected, result)
	}
}

// Be careful! int.Div and int.Quo is different
func TestDiv_negative(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(0xFFFFFFEC), // 0xFFFFFFEC : -20
		uint8(opcode.Push), uint32ToBytes(6),
		uint8(opcode.Div),
	)
	testExpected := item(0xFFFFFFFC) // 0xFFFFFFFC : -4

	stack, err := Execute(testByteCode)

	if err != nil {
		t.Error(err)
	}
	result := stack.pop()
	if testExpected != result {
		t.Errorf("stack.pop() result wrong - expected=%d, got=%d", testExpected, result)
	}
}

func TestMod(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(14),
		uint8(opcode.Push), uint32ToBytes(5),
		uint8(opcode.Mod),
	)
	testExpected := item(4)

	stack, err := Execute(testByteCode)

	if err != nil {
		t.Error(err)
	}
	result := stack.pop()
	if testExpected != result {
		t.Errorf("stack.pop() result wrong - expected=%d, got=%d", testExpected, result)
	}
}

func TestMod_negative(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint32ToBytes(0xFFFFFFEC), // 0xFFFFFFEC : -20
		uint8(opcode.Push), uint32ToBytes(6),
		uint8(opcode.Mod),
	)
	testExpected := item(4)

	stack, err := Execute(testByteCode)

	if err != nil {
		t.Error(err)
	}
	result := stack.pop()
	if testExpected != result {
		t.Errorf("stack.pop() result wrong - expected=%d, got=%d", testExpected, result)
	}
}

// TODO: implement test cases :-)
func TestLT(t *testing.T) {

}

// TODO: implement test cases :-)
func TestGT(t *testing.T) {

}

// TODO: implement test cases :-)
func TestEQ(t *testing.T) {

}

// TODO: implement test cases :-)
func TestNOT(t *testing.T) {

}

// TODO: implement test cases :-)
func TestPop(t *testing.T) {

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

// TODO: implement test cases :-)
func TestMload(t *testing.T) {

}

// TODO: implement test cases :-)
func TestMstore(t *testing.T) {

}
