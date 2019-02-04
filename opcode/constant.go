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

package opcode

import "errors"

type Type uint8

const (

	// Pop the first two items in the stack.
	// Add popped two items and push to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a+b]
	// [x]      [x]
	//
	Add Type = 0x01

	// Pop the first two items in the stack.
	// Multiply popped two items and push to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a*b]
	// [x]      [x]
	//
	Mul Type = 0x02

	// Pop the first two items in the stack.
	// Subtract popped two items and push to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a-b]
	// [x]      [x]
	//
	Sub Type = 0x03

	// Pop the first two items in the stack.
	// Divide popped two items and push to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a/b]
	// [x]      [x]
	//
	Div Type = 0x04

	// Pop the first two items in the stack.
	// Mod popped two items and push to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a%b]
	// [x]      [x]
	//
	Mod Type = 0x05

	// Pop the first two items in the stack.
	// Check if the left operand(first popped item) is less than the right operand(second popped item).
	// If it is true, push true to the stack. If not push false to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a<b]
	// [x]      [x]
	//
	LT Type = 0x10

	// Pop the first two items in the stack.
	// Check if the left operand(first popped item) is equal to or less than the right operand(second popped item).
	// If it is true, push true to the stack. If not push false to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a<=b]
	// [x]      [x]
	//
	LTE Type = 0x11

	// Pop the first two items in the stack.
	// Check if the left operand(first popped item) is greater than the right operand(second popped item).
	// If it is true, push true to the stack. If not push false to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a>b]
	// [x]      [x]
	//
	GT Type = 0x12

	// Pop the first two items in the stack.
	// Check if the left operand(first popped item) is equal to or greater than the right operand(second popped item).
	// If it is true, push true to the stack. If not push false to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a>=b]
	// [x]      [x]
	//
	GTE Type = 0x13

	// Pop the first two items in the stack.
	// Check if the left operand(first popped item) is equal to the right operand(second popped item).
	// If it is true, push true to the stack. If not push false to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a == b]
	// [x]      [x]
	//
	EQ Type = 0x14

	// Pop the first item in the stack.
	// Reverse the sign and push it to the stack.
	//
	// Ex)
	// [a]       [~a]
	// [b]  ==>  [b]
	// [x]       [x]
	//
	NOT Type = 0x15

	// Pop the first item in the stack.
	//
	// Ex)
	// [a]
	// [b]  ==>  [b]
	// [x]       [x]
	//
	Pop Type = 0x20

	// Push the 32bits(uint32) item .
	//
	// Ex)
	//           [a]
	// [b]  ==>  [b]
	// [x]       [x]
	//
	Push Type = 0x21

	// Pop the first item in the stack.
	// Load a value from memory and push it to the stack
	//
	// Ex)
	// [offset]       [memory[offset:offset+32]]
	// [x]       ==>  [b]
	// [y]            [x]
	//
	Mload Type = 0x22

	// Pop the first two items in the stack.
	// Store the value with the first item(offset) as the offset of the memory.
	//
	// Ex)
	// [offset]
	// [value]   ==>
	// [y]            [y]
	//
	// memory[offset:offset+32] = value
	Mstore Type = 0x23

	// Get the function selector information in the CallFunc.
	// `Func` data is 4bytes of Keccak(function(params))
	//
	// Ex)
	//           [CallFunc.Func]
	// [x]  ==>  [x]
	// [y]       [y]
	// [x] is 4byte which encoded as function selector
	LoadFunc Type = 0x24

	// Get the function arguments in the CallFunc.
	// 'Args' is arguments encoded as a bytes value according to abi.
	// Args's first 4 bytes is pointing to encoded size and value.
	// Ex)
	//
	// [index]  ==>  [Callfunc.Args[index]]
	// [y]           [y]
	LoadArgs Type = 0x25

	// Push the data which stored in 'memory' to stack or remain in memory when return type is string
	//
	// Ex)
	//           [offset]
	//      ==>  [size]
	// [y]       [y]
	Returning Type = 0x26
)

// Change the bytecode of an opcode to string.
func (p Type) String() (string, error) {
	switch p {
	case 0x01:
		return "Add", nil
	case 0x02:
		return "Mul", nil
	case 0x03:
		return "Sub", nil
	case 0x04:
		return "Div", nil
	case 0x05:
		return "Mod", nil
	case 0x10:
		return "LT", nil
	case 0x11:
		return "LTE", nil
	case 0x12:
		return "GT", nil
	case 0x13:
		return "GTE", nil
	case 0x14:
		return "EQ", nil
	case 0x15:
		return "NOT", nil
	case 0x20:
		return "Pop", nil
	case 0x21:
		return "Push", nil
	case 0x22:
		return "Mload", nil
	case 0x23:
		return "Mstore", nil
	case 0x24:
		return "LoadFunc", nil
	case 0x25:
		return "LoadArgs", nil
	case 0x26:
		return "Returning", nil

	default:
		return "", errors.New("String() error - Not defined opcode")
	}
}
