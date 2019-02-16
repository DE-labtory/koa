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

	"bytes"

	"github.com/DE-labtory/koa/abi"
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
		uint8(opcode.Push), int64ToBytes(1),
		uint8(opcode.Push), int64ToBytes(2),
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
		uint8(opcode.Push), int64ToBytes(-20),
		uint8(opcode.Push), int64ToBytes(30),
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
		uint8(opcode.Push), int64ToBytes(3),
		uint8(opcode.Push), int64ToBytes(5),
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
		uint8(opcode.Push), int64ToBytes(-3),
		uint8(opcode.Push), int64ToBytes(5),
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
		uint8(opcode.Push), int64ToBytes(50),
		uint8(opcode.Push), int64ToBytes(20),
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
		uint8(opcode.Push), int64ToBytes(-20),
		uint8(opcode.Push), int64ToBytes(50),
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
		uint8(opcode.Push), int64ToBytes(14),
		uint8(opcode.Push), int64ToBytes(5),
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
		uint8(opcode.Push), int64ToBytes(-20),
		uint8(opcode.Push), int64ToBytes(6),
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
		uint8(opcode.Push), int64ToBytes(14),
		uint8(opcode.Push), int64ToBytes(5),
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
		uint8(opcode.Push), int64ToBytes(-20),
		uint8(opcode.Push), int64ToBytes(6),
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
		uint8(opcode.Push), int64ToBytes(0xAC), // 000...10101100
		uint8(opcode.Push), int64ToBytes(0xF0), // 000...11110000
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
		uint8(opcode.Push), int64ToBytes(0xAC), //  000...10101100
		uint8(opcode.Push), int64ToBytes(0xF0), //  000...11110000
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
		x      int64
		y      int64
		answer int
	}{
		{1, 2, 1}, // true = 1 , false = 0
		{2, 1, 0},
		{-20, -21, 0},
		{-21, -20, 1},
	}

	for i, test := range tests {
		testByteCode := makeTestByteCode(
			uint8(opcode.Push), int64ToBytes(test.x),
			uint8(opcode.Push), int64ToBytes(test.y),
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
		x      int64
		y      int64
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
			uint8(opcode.Push), int64ToBytes(test.x),
			uint8(opcode.Push), int64ToBytes(test.y),
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
		x      int64
		y      int64
		answer int
	}{
		{1, 2, 0}, // true = 1 , false = 0
		{2, 1, 1},
		{-20, -21, 1},
		{-21, 20, 0},
	}

	for i, test := range tests {
		testByteCode := makeTestByteCode(
			uint8(opcode.Push), int64ToBytes(test.x),
			uint8(opcode.Push), int64ToBytes(test.y),
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
		x      int64
		y      int64
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
			uint8(opcode.Push), int64ToBytes(test.x),
			uint8(opcode.Push), int64ToBytes(test.y),
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
		x      int64
		y      int64
		answer int
	}{
		{1, 1, 1}, // true = 1 , false = 0
		{2, 1, 0},
		{-20, -20, 1},
		{-20, -21, 0},
	}

	for i, test := range tests {
		testByteCode := makeTestByteCode(
			uint8(opcode.Push), int64ToBytes(test.x),
			uint8(opcode.Push), int64ToBytes(test.y),
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
		x      int64
		answer int64
	}{
		{0, -1}, // ^x = -x-1  0x00000000 -> 0xFFFFFFFF
		{3, -4},
	}

	for i, test := range tests {
		testByteCode := makeTestByteCode(
			uint8(opcode.Push), int64ToBytes(test.x),
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
		uint8(opcode.Push), int64ToBytes(1),
		uint8(opcode.Push), int64ToBytes(2),
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
		uint8(opcode.Push), int64ToBytes(1),
		uint8(opcode.Push), int64ToBytes(2),
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
		uint8(opcode.Push), int64ToBytes(3),
	)

	_, err := Execute(testByteCode, nil, nil)

	if err != ErrInvalidOpcode {
		t.Error("The desired error was not found")
	}
}

func TestMload(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), int64ToBytes(8), // size
		uint8(opcode.Push), int64ToBytes(0), // offset
		uint8(opcode.Mload),
		uint8(opcode.Push), int64ToBytes(8), // size
		uint8(opcode.Push), int64ToBytes(8), // offset
		uint8(opcode.Mload),
	)

	testMemory := NewMemory()
	testMemory.Resize(16)
	testMemory.Sets(0, 8, int64ToBytes(40))
	testMemory.Sets(8, 8, int64ToBytes(20))

	testExpected := []struct {
		value  []byte
		size   uint64
		offset uint64
	}{
		{
			value:  []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14},
			size:   8,
			offset: 0,
		},
		{
			value:  []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x28},
			size:   8,
			offset: 8,
		},
	}

	stack, err := Execute(testByteCode, testMemory, nil)
	if err != nil {
		t.Error(err)
	}

	if stack.len() != 2 {
		t.Errorf("Invalid stack size - expected=%d, got =%d", 0, stack.len())
	}

	for _, test := range testExpected {
		value := stack.pop()
		if !bytes.Equal(int64ToBytes(int64(value)), test.value) {
			t.Errorf("Invalid memory value - expected=%x, got=%x", test.value, value)
		}
	}
}

func TestMstore(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), int64ToBytes(20), // value
		uint8(opcode.Push), int64ToBytes(8), // size
		uint8(opcode.Push), int64ToBytes(0), // offset
		uint8(opcode.Mstore),
		uint8(opcode.Push), int64ToBytes(40), // value
		uint8(opcode.Push), int64ToBytes(8), // size
		uint8(opcode.Push), int64ToBytes(8), // offset
		uint8(opcode.Mstore),
	)

	testExpected := []struct {
		value  []byte
		size   uint64
		offset uint64
	}{
		{
			value:  []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14},
			size:   8,
			offset: 0,
		},
		{
			value:  []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x28},
			size:   8,
			offset: 8,
		},
	}

	memory := NewMemory()

	stack, err := Execute(testByteCode, memory, nil)
	if err != nil {
		t.Error(err)
	}

	if stack.len() != 0 {
		t.Errorf("Invalid stack size - expected=%d, got =%d", 0, stack.len())
	}

	for _, test := range testExpected {
		value := memory.GetVal(test.offset, test.size)
		if !bytes.Equal(value, test.value) {
			t.Errorf("Invalid memory value - expected=%x, got=%x", test.value, value)
		}
	}
}

func TestLoadFunc(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.LoadFunc),
	)

	callFunc := &CallFunc{
		Func: abi.Selector("foo(int64,bool,string)"),
	}

	testExpected := []byte{0x00, 0x00, 0x00, 0x00, 0xf2, 0x61, 0xd0, 0x09}

	stack, err := Execute(testByteCode, nil, callFunc)
	if err != nil {
		t.Error(err)
	}

	if stack.len() != 1 {
		t.Errorf("Invalid stack size - expected=%d, got =%d", 0, stack.len())
	}

	if !bytes.Equal(int64ToBytes(int64(stack.pop())), testExpected) {
		t.Error("Invalid")
	}
}

func TestLoadArgs(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), int64ToBytes(0), // value
		uint8(opcode.LoadArgs),
		uint8(opcode.Push), int64ToBytes(1), // value
		uint8(opcode.LoadArgs),
		uint8(opcode.Push), int64ToBytes(2), // value
		uint8(opcode.LoadArgs),
	)

	encodedParams, err := abi.Encode(50, "HelloKOA", 256)
	if err != nil {
		t.Error(err)
	}
	callFunc := &CallFunc{
		Args: encodedParams,
	}

	testExpected := []struct {
		value []byte
	}{
		{
			value: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00},
		},
		{
			value: []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x4b, 0x4f, 0x41},
		},
		{
			value: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x32},
		},
	}

	stack, err := Execute(testByteCode, nil, callFunc)
	if err != nil {
		t.Error(err)
	}

	if stack.len() != 3 {
		t.Errorf("Invalid stack size - expected=%d, got =%d", 0, stack.len())
	}

	for _, expected := range testExpected {
		value := stack.pop()
		if !bytes.Equal(int64ToBytes(int64(value)), expected.value) {
			t.Errorf("Invalid bytes - expected = %x, got = %x", expected.value, value)
		}
	}

}

func TestReturning(t *testing.T) {
	testByteCode := makeTestByteCode( //  op code index
		uint8(opcode.Push), int64ToBytes(13), // 0 , 1
		uint8(opcode.Push), int64ToBytes(1), // 2 , 3 -> function selector
		uint8(opcode.Push), int64ToBytes(2), // 4 , 5
		uint8(opcode.Returning),             // 6
		uint8(opcode.Push), int64ToBytes(2), // 7 , 8
		uint8(opcode.Push), int64ToBytes(3), // 9 , 10
		uint8(opcode.Push), int64ToBytes(4), // 11 , 12
		uint8(opcode.JumpDst), // 13 ( jump to here! )
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

func TestRevert(t *testing.T) {
	testByteCode := makeTestByteCode( //  op code index
		uint8(opcode.Push), int64ToBytes(1), // 0 , 1
		uint8(opcode.Push), int64ToBytes(2), // 2 , 3
		uint8(opcode.Revert),                // 4
		uint8(opcode.Push), int64ToBytes(3), // 5 , 6
		uint8(opcode.Push), int64ToBytes(4), // 7 , 8
		uint8(opcode.Push), int64ToBytes(5), // 9 , 10
		uint8(opcode.JumpDst), // 11 ( jump to here! )
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

// TODO: implement test cases :-)
func TestCallFunc_function(t *testing.T) {

}

// TODO: implement test cases :-)
func TestCallFunc_arguments(t *testing.T) {

}

func TestJumpOp(t *testing.T) {
	testByteCode := makeTestByteCode( //  op code index
		uint8(opcode.Push), int64ToBytes(1), // 0 , 1
		uint8(opcode.Push), int64ToBytes(7), // 2 , 3
		uint8(opcode.Jump),                  // 4
		uint8(opcode.Push), int64ToBytes(2), // 5 , 6
		uint8(opcode.JumpDst),               // 7 ( jump to here! )
		uint8(opcode.Push), int64ToBytes(3), // 8 , 9
	)

	testExpected := []item{1, 3}

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

func TestJumpiJump(t *testing.T) {
	testByteCode := makeTestByteCode( //  op code index
		uint8(opcode.Push), int64ToBytes(1), // 0 , 1
		uint8(opcode.Push), int64ToBytes(1), // 2 , 3(false)
		uint8(opcode.Push), int64ToBytes(9), // 4 , 5
		uint8(opcode.Jumpi),                 // 6
		uint8(opcode.Push), int64ToBytes(2), // 7 , 8
		uint8(opcode.JumpDst),               // 9 ( jump to here! )
		uint8(opcode.Push), int64ToBytes(3), // 10 , 11
	)

	testExpected := []item{1, 3}

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

func TestJumpiNotJump(t *testing.T) {
	testByteCode := makeTestByteCode( //  op code index
		uint8(opcode.Push), int64ToBytes(1), // 0 , 1
		uint8(opcode.Push), int64ToBytes(1), // 2 , 3(true)
		uint8(opcode.Push), int64ToBytes(9), // 4 , 5
		uint8(opcode.Jumpi),                 // 6
		uint8(opcode.Push), int64ToBytes(2), // 7 , 8
		uint8(opcode.JumpDst),               // 9 ( jump to here! )
		uint8(opcode.Push), int64ToBytes(3), // 10 , 11
	)

	testExpected := []item{1, 2, 3}

	stack, err := Execute(testByteCode, nil, nil)
	if err != nil {
		t.Error(err)
	}

	if len(stack.items) != 3 {
		t.Errorf("Invalid stack size - expected=%d, got =%d", len(testExpected), stack.len())
	}

	for i, item := range stack.items {
		if testExpected[i] != item {
			t.Errorf("Stack item is incorrect - expected=%d, got=%d", testExpected[i], item)
		}
	}
}

func TestDUP(t *testing.T) {
	testByteCode := makeTestByteCode(
		uint8(opcode.Push), int64ToBytes(1),
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
		uint8(opcode.Push), int64ToBytes(1),
		uint8(opcode.Push), int64ToBytes(2),
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
