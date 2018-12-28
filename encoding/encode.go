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
	"log"
	"strconv"
)

// In koa, we use hexadecimal encoding

// TODO: implement w/ test cases :-)
func Encode(operand interface{}) []byte {
	switch op := operand.(type) {
	case int:
		return nil
	case string:
		return nil
	case bool:
		b, err := encodeBool(op)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		return b
	default:
		return nil
	}
}

// Encode integer to hexadecimal bytes
// ex) int 123 => 0x7b
func encodeInt(operand int) ([]byte, error) {
	s := strconv.FormatInt(int64(operand), 16)

	// Encoded byte length should be even number
	if len(s)%2 == 1 {
		s = "0" + s
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Encode string to hexadecimal bytes
// ex) string "abc" => 0x616263
// TODO: implement w/ test cases :-)
func encodeString(operand string) ([]byte, error) {
	return nil, nil
}

// Encode boolean to hexadecimal bytes
// ex) bool true => 0x01
// ex) bool false => 0x00
// TODO: implement w/ test cases :-)
func encodeBool(operand bool) ([]byte, error) {
	return nil, nil
}
