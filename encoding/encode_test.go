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

package encoding_test

import (
	"bytes"
	"testing"

	"github.com/DE-labtory/koa/encoding"
)

func TestEncodeOperand(t *testing.T) {
	tests := []struct {
		operand      interface{}
		expectedByte []byte
		expectedErr  error
	}{
		{
			operand:      true,
			expectedByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			expectedErr:  nil,
		},
		{
			operand:      false,
			expectedByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expectedErr:  nil,
		},
		{
			operand:      "abc",
			expectedByte: []byte{0x61, 0x62, 0x63},
			expectedErr:  nil,
		},
		{
			operand:      "1234567890",
			expectedByte: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30},
			expectedErr:  nil,
		},
		{
			operand:      "~!@#$%^&*()_+",
			expectedByte: []byte{0x7e, 0x21, 0x40, 0x23, 0x24, 0x25, 0x5e, 0x26, 0x2a, 0x28, 0x29, 0x5f, 0x2b},
			expectedErr:  nil,
		},
		{
			operand:      "12!@qw",
			expectedByte: []byte{0x31, 0x32, 0x21, 0x40, 0x71, 0x77},
			expectedErr:  nil,
		},
		{
			operand:      int64(1),
			expectedByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			expectedErr:  nil,
		},
		{
			operand:      int64(23),
			expectedByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x17},
			expectedErr:  nil,
		},
		{
			operand:      int64(456),
			expectedByte: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xc8},
			expectedErr:  nil,
		},
		{
			operand:      'c',
			expectedByte: nil,
			expectedErr:  encoding.EncodeError{'c'},
		},
	}

	for i, test := range tests {
		op := test.operand
		byteCode, err := encoding.EncodeOperand(op)

		if !bytes.Equal(byteCode, test.expectedByte) {
			t.Fatalf("test[%d] - EncodeOperand() result wrong. expectedByte=%x, got=%x", i, test.expectedByte, byteCode)
		}

		if err != test.expectedErr {
			t.Fatalf("test[%d] - EncodeOperand() result wrong. expectedByte=%x, got=%x", i, test.expectedByte, byteCode)
		}
	}
}
