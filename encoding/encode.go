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
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type convertToBytes func(int64) string
type EncodingBytes int

const (
	NO_PADDING    EncodingBytes = 0
	FOUR_PADDING  EncodingBytes = 4
	EIGHT_PADDING EncodingBytes = 8
)

// In koa, we use hexadecimal encoding

type EncodeError struct {
	Operand interface{}
}

func (e EncodeError) Error() string {
	return fmt.Sprintf("EncodeOperand() error - operand %v could not encoded", e.Operand)
}

// EncodeOperand() encodes operand to bytes.
func EncodeOperand(operand interface{}, convertBytes EncodingBytes) ([]byte, error) {
	convertFunc, err := selectConvertFunc(convertBytes)
	if err != nil {
		return nil, err
	}

	switch op := operand.(type) {
	case int:
		return encodeInt(int64(op), convertFunc)

	case int64:
		return encodeInt(op, convertFunc)

	case string:
		return encodeString(op)

	case bool:
		return encodeBool(op, convertFunc)

	default:
		return nil, EncodeError{op}
	}
}

func selectConvertFunc(bytes EncodingBytes) (convertToBytes, error) {
	switch bytes {
	case NO_PADDING:
		return convertToByte, nil
	case FOUR_PADDING:
		return convertTo4Bytes, nil
	case EIGHT_PADDING:
		return convertTo8Bytes, nil
	default:
		return nil, errors.New("EncodingBytes does not match")
	}
}

// Encode integer to hexadecimal bytes
// ex) When EncodingBytes is 8 : int 123 => 0x000000000000007b
// ex) When EncodingBytes is 4 : int 123 => 0x0000007b
// ex) When EncodingBytes is 0 : int 123 => 0x7b
func encodeInt(operand int64, convertFunc convertToBytes) ([]byte, error) {
	s := convertFunc(operand)

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
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
// ex) When EncodingBytes is 8 : bool true => 0x0000000000000001
// ex) When EncodingBytes is 4 : bool true => 0x00000001
// ex) When EncodingBytes is 0 : bool true => 0x01

// ex) When EncodingBytes is 8 : bool false => 0x0000000000000000
// ex) When EncodingBytes is 4 : bool false => 0x00000000
// ex) When EncodingBytes is 0 : bool false => 0x00
func encodeBool(operand bool, convertFunc convertToBytes) ([]byte, error) {
	var src string

	if operand {
		src = convertFunc(1)
	} else {
		src = convertFunc(0)
	}

	dst, err := hex.DecodeString(src)
	if err != nil {
		return nil, err
	}

	return dst, nil

}

func convertToByte(operand int64) string {
	var zeroSet string

	src := strconv.FormatUint(uint64(operand), 16)

	if len(src)%2 == 1 {
		zeroSet += "0"
	}

	return zeroSet + src
}

// convert to 4 byte
func convertTo4Bytes(operand int64) string {
	var zeroSet string

	src := strconv.FormatUint(uint64(operand), 16)
	diff := 8 - len(src)

	for ; diff > 0; diff-- {
		zeroSet += "0"
	}
	return zeroSet + src
}

// convert to 8 byte
func convertTo8Bytes(operand int64) string {
	var zeroSet string

	src := strconv.FormatUint(uint64(operand), 16)
	diff := 16 - len(src)

	for ; diff > 0; diff-- {
		zeroSet += "0"
	}
	return zeroSet + src
}
