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
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
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
// ex) int 123 => 0x7b
func encodeInt(operand int) ([]byte, error) {
	s := convertTo4Bytes(operand)

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	if len(b) != 4 {
		return nil, errors.New("Integer size is not 32 bits")
	}

	return b, nil
}

// Encode string to hexadecimal bytes
// ex) string "abc" => 0x616263
func encodeString(operand string) ([]byte, error) {
	src := hex.EncodeToString([]byte(operand))

	if len(src)&1 == 1 {
		src = "0" + src
	}

	dst, err := hex.DecodeString(src)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// Encode boolean to hexadecimal bytes
// ex) bool true => 0x01
// ex) bool false => 0x00
func encodeBool(operand bool) ([]byte, error) {
	var src string

	if operand {
		src = convertTo4Bytes(1)
	} else {
		src = convertTo4Bytes(0)
	}

	dst, err := hex.DecodeString(src)
	if err != nil {
		return nil, err
	}

	return dst, nil

}

// convert to 4 byte
func convertTo4Bytes(operand int) string {
	var zeroSet string

	src := strconv.FormatUint(uint64(operand), 16)
	diff := 8 - len(src)

	for ; diff > 0; diff-- {
		zeroSet += "0"
	}
	return zeroSet + src

}
