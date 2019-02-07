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
		uint8(opcode.Push), int32ToBytes(1),
		uint8(opcode.Push), int32ToBytes(2),
		uint8(opcode.Add),
	)
	testExpected := item(3)

	stack, err := Execute(testByteCode, nil, nil)
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
		uint8(opcode.Push), int32ToBytes(-20),
		uint8(opcode.Push), int32ToBytes(30),
		uint8(opcode.Add),
	)
	testExpected := item(10)

	stack, err := Execute(testByteCode, nil, nil)
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
		uint8(opcode.Push), int32ToBytes(3),
		uint8(opcode.Push), int32ToBytes(5),
		uint8(opcode.Mul),
	)
	testExpected := item(15)

	stack, err := Execute(testByteCode, nil, nil)
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
		uint8(opcode.Push), int32ToBytes(-3),
		uint8(opcode.Push), int32ToBytes(5),
		uint8(opcode.Mul),
	)
	testExpected := item(-15)

	stack, err := Execute(testByteCode, nil, nil)
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
		uint8(opcode.Push), int32ToBytes(50),
		uint8(opcode.Push), int32ToBytes(20),
		uint8(opcode.Sub),
	)
	testExpected := item(30)

	stack, err := Execute(testByteCode, nil, nil)
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
		uint8(opcode.Push), int32ToBytes(-20),
		uint8(opcode.Push), int32ToBytes(50),
		uint8(opcode.Sub),
	)
	testExpected := item(-70)

	stack, err := Execute(testByteCode, nil, nil)
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
		uint8(opcode.Push), int32ToBytes(14),
		uint8(opcode.Push), int32ToBytes(5),
		uint8(opcode.Div),
	)
	testExpected := item(2)

	stack, err := Execute(testByteCode, nil, nil)

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
		uint8(opcode.Push), int32ToBytes(-20),
		uint8(opcode.Push), int32ToBytes(6),
		uint8(opcode.Div),
	)
	testExpected := item(-4)

	stack, err := Execute(testByteCode, nil, nil)
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
		uint8(opcode.Push), int32ToBytes(14),
		uint8(opcode.Push), int32ToBytes(5),
		uint8(opcode.Mod),
	)
	testExpected := item(4)

	stack, err := Execute(testByteCode, nil, nil)
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
		uint8(opcode.Push), int32ToBytes(-20),
		uint8(opcode.Push), int32ToBytes(6),
		uint8(opcode.Mod),
	)
	testExpected := item(4)

	stack, err := Execute(testByteCode, nil, nil)
	if err != nil {
		t.Error(err)
	}
	result := stack.pop()
	if testExpected != result {
		t.Errorf("stack.pop() result wrong - expected=%d, got=%d", testExpected, result)
	}
}

func TestAnd(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), int32ToBytes(0xAC), // 000...10101100
		uint8(opcode.Push), int32ToBytes(0xF0), // 000...11110000
		uint8(opcode.And),
	)
	testExpected := item(0xA0) // 000...10100000

	stack, err := Execute(testByteCode, nil, nil)
	if err != nil {
		t.Error(err)
	}
	result := stack.pop()
	if testExpected != result {
		t.Errorf("stack.pop() result wrong - expected=%d, got=%d", testExpected, result)
	}
}

func TestOr(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), int32ToBytes(0xAC), //  000...10101100
		uint8(opcode.Push), int32ToBytes(0xF0), //  000...11110000
		uint8(opcode.Or),
	)
	testExpected := item(0xFC) // 000...11111100

	stack, err := Execute(testByteCode, nil, nil)
	if err != nil {
		t.Error(err)
	}
	result := stack.pop()
	if testExpected != result {
		t.Errorf("stack.pop() result wrong - expected=%d, got=%d", testExpected, result)
	}
}

func TestLT(t *testing.T) {
	tests := []struct {
		x      int32
		y      int32
		answer int
	}{
		{1, 2, 1}, // true = 1 , false = 0
		{2, 1, 0},
		{-20, -21, 0},
		{-21, -20, 1},
	}

	for i, test := range tests {
		testByteCode := makeTestByteCode(
			uint8(opcode.Push), int32ToBytes(test.x),
			uint8(opcode.Push), int32ToBytes(test.y),
			uint8(opcode.LT),
		)
		testExpected := item(test.answer)

		stack, err := Execute(testByteCode, nil, nil)
		if err != nil {
			t.Error(err)
		}
		result := stack.pop()
		if testExpected != result {
			t.Errorf("test[%d]:stack.pop() result wrong - expected=%d, got=%d", i, testExpected, result)
		}
	}
}

func TestLTE(t *testing.T) {
	tests := []struct {
		x      int32
		y      int32
		answer int
	}{
		{1, 2, 1}, // true = 1 , false = 0
		{1, 1, 1},
		{2, 1, 0},
		{-20, -21, 0},
		{-21, -20, 1},
		{-21, -21, 1},
	}

	for i, test := range tests {
		testByteCode := makeTestByteCode(
			uint8(opcode.Push), int32ToBytes(test.x),
			uint8(opcode.Push), int32ToBytes(test.y),
			uint8(opcode.LTE),
		)
		testExpected := item(test.answer)

		stack, err := Execute(testByteCode, nil, nil)
		if err != nil {
			t.Error(err)
		}
		result := stack.pop()
		if testExpected != result {
			t.Errorf("test[%d]:stack.pop() result wrong - expected=%d, got=%d", i, testExpected, result)
		}
	}
}

func TestGT(t *testing.T) {
	tests := []struct {
		x      int32
		y      int32
		answer int
	}{
		{1, 2, 0}, // true = 1 , false = 0
		{2, 1, 1},
		{-20, -21, 1},
		{-21, 20, 0},
	}

	for i, test := range tests {
		testByteCode := makeTestByteCode(
			uint8(opcode.Push), int32ToBytes(test.x),
			uint8(opcode.Push), int32ToBytes(test.y),
			uint8(opcode.GT),
		)
		testExpected := item(test.answer)

		stack, err := Execute(testByteCode, nil, nil)
		if err != nil {
			t.Error(err)
		}
		result := stack.pop()
		if testExpected != result {
			t.Errorf("test[%d]:stack.pop() result wrong - expected=%d, got=%d", i, testExpected, result)
		}
	}
}

func TestGTE(t *testing.T) {
	tests := []struct {
		x      int32
		y      int32
		answer int
	}{
		{1, 2, 0}, // true = 1 , false = 0
		{2, 1, 1},
		{2, 2, 1},
		{-20, -21, 1},
		{-21, -21, 1},
		{-21, 20, 0},
	}

	for i, test := range tests {
		testByteCode := makeTestByteCode(
			uint8(opcode.Push), int32ToBytes(test.x),
			uint8(opcode.Push), int32ToBytes(test.y),
			uint8(opcode.GTE),
		)
		testExpected := item(test.answer)

		stack, err := Execute(testByteCode, nil, nil)
		if err != nil {
			t.Error(err)
		}
		result := stack.pop()
		if testExpected != result {
			t.Errorf("test[%d]:stack.pop() result wrong - expected=%d, got=%d", i, testExpected, result)
		}
	}
}

func TestEQ(t *testing.T) {
	tests := []struct {
		x      int32
		y      int32
		answer int
	}{
		{1, 1, 1}, // true = 1 , false = 0
		{2, 1, 0},
		{-20, -20, 1},
		{-20, -21, 0},
	}

	for i, test := range tests {
		testByteCode := makeTestByteCode(
			uint8(opcode.Push), int32ToBytes(test.x),
			uint8(opcode.Push), int32ToBytes(test.y),
			uint8(opcode.EQ),
		)
		testExpected := item(test.answer)

		stack, err := Execute(testByteCode, nil, nil)
		if err != nil {
			t.Error(err)
		}
		result := stack.pop()
		if testExpected != result {
			t.Errorf("test[%d]:stack.pop() result wrong - expected=%d, got=%d", i, testExpected, result)
		}
	}
}

func TestNOT(t *testing.T) {
	tests := []struct {
		x      int32
		answer int32
	}{
		{0, -1}, // ^x = -x-1  0x00000000 -> 0xFFFFFFFF
		{3, -4},
	}

	for i, test := range tests {
		testByteCode := makeTestByteCode(
			uint8(opcode.Push), int32ToBytes(test.x),
			uint8(opcode.NOT),
		)
		testExpected := item(test.answer)

		stack, err := Execute(testByteCode, nil, nil)
		if err != nil {
			t.Error(err)
		}
		result := stack.pop()
		if testExpected != result {
			t.Errorf("test[%d]:stack.pop() result wrong - expected=%d, got=%d", i, testExpected, result)
		}
	}
}

func TestPop(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), int32ToBytes(1),
		uint8(opcode.Push), int32ToBytes(2),
		uint8(opcode.Pop),
	)
	testExpected := []item{1}

	stack, err := Execute(testByteCode, nil, nil)
	if err != nil {
		t.Error(err)
	}

	if len(stack.items) != 1 {
		t.Errorf("Invalid stack size - expected=%d, got =%d", len(testExpected), stack.len())
	}

	for i, item := range stack.items {
		if testExpected[i] != item {
			t.Errorf("Stack item is incorrect - expected=%d, got=%d", testExpected[i], item)
		}
	}
}

func TestPush(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), int32ToBytes(1),
		uint8(opcode.Push), int32ToBytes(2),
	)

	testExpected := []item{1, 2}

	stack, err := Execute(testByteCode, nil, nil)
	if err != nil {
		t.Error(err)
	}

	if len(stack.items) != 2 {
		t.Errorf("Invalid stack size - expected=%d, got =%d", len(testExpected), stack.len())
	}

	for i, item := range stack.items {
		if testExpected[i] != item {
			t.Errorf("Stack item is incorrect - expected=%d, got=%d", testExpected[i], item)
		}
	}
}

func TestPush_invalid(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), uint8(opcode.Push),
		uint8(opcode.Push), int32ToBytes(3),
	)

	_, err := Execute(testByteCode, nil, nil)

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

// TODO: implement test cases :-)
func TestLoadFunc(t *testing.T) {

}

// TODO: implement test cases :-)
func TestLoadArgs(t *testing.T) {

}

// TODO: implement test cases :-)
func TestReturning(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCallFunc_function(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCallFunc_arguments(t *testing.T) {

}

func TestDUP(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), int32ToBytes(1),
		uint8(opcode.DUP),
	)

	testExpected := []item{1, 1}

	stack, err := Execute(testByteCode, nil, nil)
	if err != nil {
		t.Error(err)
	}

	if len(stack.items) != 2 {
		t.Errorf("Invalid stack size - expected=%d, got =%d", len(testExpected), stack.len())
	}

	for i, item := range stack.items {
		if testExpected[i] != item {
			t.Errorf("Stack item is incorrect - expected=%d, got=%d", testExpected[i], item)
		}
	}
}

func TestSWAP(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), int32ToBytes(1),
		uint8(opcode.Push), int32ToBytes(2),
		uint8(opcode.SWAP),
	)

	testExpected := []item{2, 1}

	stack, err := Execute(testByteCode, nil, nil)
	if err != nil {
		t.Error(err)
	}

	if len(stack.items) != 2 {
		t.Errorf("Invalid stack size - expected=%d, got =%d", len(testExpected), stack.len())
	}

	for i, item := range stack.items {
		if testExpected[i] != item {
			t.Errorf("Stack item is incorrect - expected=%d, got=%d", testExpected[i], item)
		}
	}
}
