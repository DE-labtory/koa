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

package encoding

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"errors"
)

// In koa, we use hexadecimal encoding
type EncodeError struct {
	Operand interface{}
}

func (e EncodeError) Error() string {
	return fmt.Sprintf("EncodeOperand() error - operand %v could not encoded", e.Operand)
}

// EncodeOperand() encodes operand to bytes.
func EncodeOperand(operand interface{}) ([]byte, error) {
	switch op := operand.(type) {
	case int:
		return encodeInt(int64(op))

	case int64:
		return encodeInt(op)

	case string:
		return encodeString(op)

	case bool:
		return encodeBool(op)

	default:
		return nil, EncodeError{op}
	}
}

// Encode integer to hexadecimal bytes
// ex) int 123 => 0x000000000000007b
func encodeInt(operand int64) ([]byte, error) {
	byteSlice := make([]byte, 8)
	binary.BigEndian.PutUint64(byteSlice, uint64(operand))

	return byteSlice, nil
}

// Encode string to hexadecimal bytes
// ex) string "abc" => 0x6162630000000000
func encodeString(operand string) ([]byte, error) {
	if len(operand) > 8 {
		return nil, errors.New("Length of string must shorter than 8")
	}

	src := hex.EncodeToString([]byte(operand))
	if len(src)&1 == 1 {
		src = "0" + src
	}

	dst, err := hex.DecodeString(src)
	if err != nil {
		return nil, err
	}

	for len(dst) < 8 {
		dst = append(dst, 0)
	}

	return dst, nil
}

// Encode boolean to hexadecimal bytes
// ex) bool true  => 0x0000000000000001
// ex) bool false => 0x0000000000000000
func encodeBool(operand bool) ([]byte, error) {
	byteSlice := make([]byte, 8)

	if operand {
		binary.BigEndian.PutUint64(byteSlice, 1)
	} else {
		binary.BigEndian.PutUint64(byteSlice, 0)
	}

	return byteSlice, nil
}
