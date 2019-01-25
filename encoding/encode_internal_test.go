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
	"testing"
)

func TestEncodeInt(t *testing.T) {
	tests := []struct {
		operand  int
		expected []byte
	}{
		{
			operand:  1,
			expected: []byte{0x00, 0x00, 0x00, 0x01},
		},
		{
			operand:  23,
			expected: []byte{0x00, 0x00, 0x00, 0x17},
		},
		{
			operand:  456,
			expected: []byte{0x00, 0x00, 0x01, 0xc8},
		},
	}

	for i, test := range tests {
		op := test.operand
		bytecode, err := encodeInt(op)

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
		operand  string
		expected []byte
	}{
		{
			operand:  "abc",
			expected: []byte{0x61, 0x62, 0x63},
		},
		{
			operand:  "1234567890",
			expected: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30},
		},
		{
			operand:  "~!@#$%^&*()_+",
			expected: []byte{0x7e, 0x21, 0x40, 0x23, 0x24, 0x25, 0x5e, 0x26, 0x2a, 0x28, 0x29, 0x5f, 0x2b},
		},
		{
			operand:  "12!@qw",
			expected: []byte{0x31, 0x32, 0x21, 0x40, 0x71, 0x77},
		},
	}

	for i, test := range tests {
		op := test.operand
		bytecode, err := encodeString(op)

		if err != nil {
			t.Fatalf("test[%d] - encodeString() had error. err=%v", i, err)

		}

		if !bytes.Equal(bytecode, test.expected) {
			t.Fatalf("test[%d] - encodeString() result wrong. expected=%x, got=%x", i, test.expected, bytecode)
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
			expected: []byte{0x00, 0x00, 0x00, 0x01},
		},
		{
			operand:  false,
			expected: []byte{0x00, 0x00, 0x00, 0x00},
		},
	}

	for i, test := range tests {
		op := test.operand
		bytecode, err := encodeBool(op)

		if err != nil {
			t.Fatalf("test[%d] - encodeBool() had error. err=%v", i, err)
		}

		if !bytes.Equal(bytecode, test.expected) {
			t.Fatalf("test[%d] - encodeBool() result wrong. expected=%x, got=%x", i, test.expected, bytecode)
		}

		if len(bytecode) != 4 {
			t.Fatalf("test[%d] - encodeBool() result wrong. expected=4, got=%x", i, bytecode)
		}
	}
}
