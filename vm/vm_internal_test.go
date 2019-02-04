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

package vm

import (
	"bytes"
	"testing"
)

func TestCallFuncArguments(t *testing.T) {
	tests := []struct {
		n        int
		args     []byte
		expected []byte
	}{
		{
			n: 0,
			args: []byte{
				0x00, // 0x0 pointer
				0x00, // 0x1 pointer
				0x00, // 0x2 pointer
				0x04, // 0x3 pointer
				0x00, // 0x4 size
				0x00, // 0x5 size
				0x00, // 0x6 size
				0x01, // 0x7 size
				0x0a, // 0x8 value
			},
			expected: []byte{0x0a},
		},
		{
			n: 1,
			args: []byte{
				0x00, // 0x00 pointer
				0x00, // 0x01 pointer
				0x00, // 0x02 pointer
				0x08, // 0x03 pointer

				0x00, // 0x04 pointer
				0x00, // 0x05 pointer
				0x00, // 0x06 pointer
				0x0d, // 0x07 pointer

				0x00, // 0x08 size
				0x00, // 0x09 size
				0x00, // 0x0a size
				0x01, // 0x0b size
				0x0a, // 0x0c value

				0x00, // 0x0d size
				0x00, // 0x0e size
				0x00, // 0x0f size
				0x01, // 0x10 size
				0x0b, // 0x11 value
			},
			expected: []byte{0x0b},
		},
		{
			n: 2,
			args: []byte{
				0x00, // 0x00 pointer
				0x00, // 0x01 pointer
				0x00, // 0x02 pointer
				0x0c, // 0x03 pointer

				0x00, // 0x04 pointer
				0x00, // 0x05 pointer
				0x00, // 0x06 pointer
				0x11, // 0x07 pointer

				0x00, // 0x08 pointer
				0x00, // 0x09 pointer
				0x00, // 0x0a pointer
				0x16, // 0x0b pointer

				0x00, // 0x0c size
				0x00, // 0x0d size
				0x00, // 0x0e size
				0x01, // 0x0f size
				0x0a, // 0x10 value

				0x00, // 0x11 size
				0x00, // 0x12 size
				0x00, // 0x13 size
				0x01, // 0x14 size
				0x0b, // 0x15 value

				0x00, // 0x16 size
				0x00, // 0x17 size
				0x00, // 0x18 size
				0x01, // 0x19 size
				0x0c, // 0x1a value
			},
			expected: []byte{0x0c},
		},
		// test when size is higher than 1
		{
			n: 0,
			args: []byte{
				0x00, // 0x0 pointer
				0x00, // 0x1 pointer
				0x00, // 0x2 pointer
				0x04, // 0x3 pointer

				0x00, // 0x4 size
				0x00, // 0x5 size
				0x00, // 0x6 size
				0x02, // 0x7 size
				0x0a, // 0x8 value
				0x0b, // 0x9 value
			},
			expected: []byte{0x0a, 0x0b},
		},
		{
			n: 1,
			args: []byte{
				0x00, // 0x00 pointer
				0x00, // 0x01 pointer
				0x00, // 0x02 pointer
				0x08, // 0x03 pointer

				0x00, // 0x04 pointer
				0x00, // 0x05 pointer
				0x00, // 0x06 pointer
				0x0e, // 0x07 pointer

				0x00, // 0x08 size
				0x00, // 0x09 size
				0x00, // 0x0a size
				0x02, // 0x0b size
				0x0a, // 0x0c value
				0x0b, // 0x0d value

				0x00, // 0x0e size
				0x00, // 0x0f size
				0x00, // 0x10 size
				0x03, // 0x11 size
				0x01, // 0x12 value
				0x02, // 0x13 value
				0x03, // 0x14 value
			},
			expected: []byte{0x01, 0x02, 0x03},
		},
	}

	for i, tt := range tests {
		cf := CallFunc{Args: tt.args}
		result := cf.arguments(tt.n)

		if !bytes.Equal(result, tt.expected) {
			t.Errorf("test[%d] - Wrong arguments returned expected=%v, got=%v",
				i, tt.expected, result)
		}
	}
}

func TestCallFuncArgumentsPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Test case did not panic")
		}
	}()

	CallFunc{}.arguments(-1)
}
