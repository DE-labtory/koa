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
	// Check if the left operand(first popped item) is less than the right operand(second popped item).
	// If it is true, push true to the stack. If not push false to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a>b]
	// [x]      [x]
	//
	GT Type = 0x11

	// Pop the first two items in the stack.
	// Check if the left operand(first popped item) is equal to the right operand(second popped item).
	// If it is true, push true to the stack. If not push false to the stack.
	//
	// Ex)
	// [a]
	// [b]  ==> [a == b]
	// [x]      [x]
	//
	EQ Type = 0x12

	// Pop the first item in the stack.
	// Reverse the sign and push it to the stack.
	//
	// Ex)
	// [a]       [-a]
	// [b]  ==>  [b]
	// [x]       [x]
	//
	NOT Type = 0x13

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
)
