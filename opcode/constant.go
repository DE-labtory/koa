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
	// Calculate bit-and popped two items and push to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a&b]
	// [x]      [x]
	//
	And Type = 0x06

	// Pop the first two items in the stack.
	// Calculate bit-or popped two items and push to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a|b]
	// [x]      [x]
	//
	Or Type = 0x07

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
	// minus the value and push it to the stack.
	//
	// Ex)
	// [a]       [-a]
	// [b]  ==>  [b]
	// [x]       [x]
	//
	Minus Type = 0x16

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
	// [offset]
	// [size]         [memory[offset:offset+size]]
	// [x]       ==>  [x]
	// [y]            [y]
	//
	Mload Type = 0x22

	// Pop the first two items in the stack.
	// Store the value with the first item(offset) as the offset of the memory.
	//
	// Ex)
	// [offset]
	// [size]
	// [value]   ==>
	// [y]            [y]
	//
	// memory[offset:offset+size] = value
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

	// Jump to position at which function was called
	// Ex)
	// [value]
	// [funcSel]
	// [position]  ==>  [value]
	// [y]              [y]
	Returning Type = 0x26

	// Jump to last position (Terminate the contract)
	Revert Type = 0x27

	// pop the data which present specific pc to jump
	// The opcode pointed to by pc must be JumpDst.
	Jump Type = 0x28

	// jumpDst should be where the jump will be.
	JumpDst Type = 0x29

	// Pop the first two items in the stack.
	// First item pointed to where to jump.
	// Second item should be bool data that decide to jump.
	// If second item is false, jump to first item(pc) pointed to
	Jumpi Type = 0x30

	// Duplicate data that exists at the top of the stack.
	//
	// Ex)
	//           [a]
	// [a]  ==>  [a]
	// [b]       [b]
	DUP Type = 0x31

	// Swap the first two items in the stack
	//
	// Ex)
	//
	// [a]       [b]
	// [b]  ==>  [a]
	// [c]       [c]
	SWAP Type = 0x32
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
	case 0x06:
		return "And", nil
	case 0x07:
		return "Or", nil
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
	case 0x16:
		return "Minus", nil
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
	case 0x27:
		return "Revert", nil
	case 0x28:
		return "Jump", nil
	case 0x29:
		return "JumpDst", nil
	case 0x30:
		return "Jumpi", nil
	case 0x31:
		return "DUP", nil
	case 0x32:
		return "SWAP", nil

	default:
		return "", errors.New("String() error - Not defined opcode")
	}
}
