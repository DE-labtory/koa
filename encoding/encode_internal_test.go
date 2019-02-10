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
	"bytes"
	"errors"
	"testing"
)

func TestEncodeInt(t *testing.T) {
	tests := []struct {
		operand  int64
		expected []byte
	}{
		{
			operand:  1,
			expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			operand:  23,
			expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x17},
		},
		{
			operand:  456,
			expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xc8},
		},
	}

	for i, test := range tests {
		op := test.operand
		bytecode, err := encodeInt(op, convertTo8Bytes)

		if err != nil {
			t.Fatalf("test[%d] - encodeInt() had error. err=%v", i, err)
		}

		if !bytes.Equal(bytecode, test.expected) {
			t.Fatalf("test[%d] - encodeInt() result wrong. expected=%x, got=%x", i, test.expected, bytecode)
		}
	}
}

func TestEncodeString(t *testing.T) {
	tests := []struct {
		operand     string
		expected    []byte
		expectedErr error
	}{
		{
			operand:     "abc",
			expected:    []byte{0x61, 0x62, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00},
			expectedErr: nil,
		},
		{
			operand:     "12345678",
			expected:    []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38},
			expectedErr: nil,
		},
		{
			operand:     "~!@#$%^&*()_+",
			expected:    nil,
			expectedErr: errors.New("Length of string must shorter than 8"),
		},
		{
			operand:  "12!@qw",
			expected: []byte{0x31, 0x32, 0x21, 0x40, 0x71, 0x77, 0x00, 0x00},
		},
	}

	for i, test := range tests {
		op := test.operand
		bytecode, err := encodeString(op)

		if bytecode != nil && !bytes.Equal(bytecode, test.expected) {
			t.Fatalf("test[%d] - encodeString() result wrong. expected=%x, got=%x", i, test.expected, bytecode)
		}

		if err != nil && err.Error() != test.expectedErr.Error() {
			t.Fatalf("test[%d] - encodeString() result error. expected=%x, got=%x",
				i, test.expectedErr.Error(), err.Error())
		}
	}
}
func TestEncodeBool(t *testing.T) {
	tests := []struct {
		operand  bool
		expected []byte
	}{
		{
			operand:  true,
			expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			operand:  false,
			expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}

	for i, test := range tests {
		op := test.operand
		bytecode, err := encodeBool(op, convertTo8Bytes)

		if err != nil {
			t.Fatalf("test[%d] - encodeBool() had error. err=%v", i, err)
		}

		if !bytes.Equal(bytecode, test.expected) {
			t.Fatalf("test[%d] - encodeBool() result wrong. expected=%x, got=%x", i, test.expected, bytecode)
		}

		if len(bytecode) != 8 {
			t.Fatalf("test[%d] - encodeBool() result wrong. expected=8, got=%x", i, bytecode)
		}
	}
}

func Test_convertToByte(t *testing.T) {

	tests := []struct {
		operand  int64
		expected string
	}{
		{
			operand:  0,
			expected: "00",
		},
		{
			operand:  1,
			expected: "01",
		},
		{
			operand:  214748,
			expected: "0346dc",
		},
		{
			operand:  429496,
			expected: "068db8",
		},
	}

	for i, test := range tests {
		op := test.operand
		convertedCode := convertToByte(op)

		if len(convertedCode) != len(test.expected) {
			t.Fatalf("test[%d] - convertToByte() result wrong. \n expected= %d, \n got=%d", i, len(test.expected), len(convertedCode))
		}

		if convertedCode != test.expected {
			t.Fatalf("test[%d] - convertToBytes() result wrong. \n expected=%s, \n got=%s", i, test.expected, convertedCode)
		}
	}
}

func Test_convertTo4Bytes(t *testing.T) {

	tests := []struct {
		operand  int64
		expected string
	}{
		{
			operand:  0,
			expected: "00000000",
		},
		{
			operand:  1,
			expected: "00000001",
		},
		{
			operand:  214748,
			expected: "000346dc",
		},
		{
			operand:  429496,
			expected: "00068db8",
		},
	}

	for i, test := range tests {
		op := test.operand
		convertedCode := convertTo4Bytes(op)

		if len(convertedCode) != 8 {
			t.Fatalf("test[%d] - convertTo4Byte() result wrong. \n expected=16, \n got=%d", i, len(convertedCode))
		}

		if convertedCode != test.expected {
			t.Fatalf("test[%d] - convertTo4Bytes() result wrong. \n expected=%s, \n got=%s", i, test.expected, convertedCode)
		}
	}
}

func Test_convertTo8Bytes(t *testing.T) {

	tests := []struct {
		operand  int64
		expected string
	}{
		{
			operand:  0,
			expected: "0000000000000000",
		},
		{
			operand:  1,
			expected: "0000000000000001",
		},
		{
			operand:  214748,
			expected: "00000000000346dc",
		},
		{
			operand:  429496,
			expected: "0000000000068db8",
		},
	}

	for i, test := range tests {
		op := test.operand
		convertedCode := convertTo8Bytes(op)

		if len(convertedCode) != 16 {
			t.Fatalf("test[%d] - convertTo8Byte() result wrong. \n expected=16, \n got=%d", i, len(convertedCode))
		}

		if convertedCode != test.expected {
			t.Fatalf("test[%d] - convertTo8Bytes() result wrong. \n expected=%s, \n got=%s", i, test.expected, convertedCode)
		}
	}
}
