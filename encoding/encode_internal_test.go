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

func Test_encodeInt(t *testing.T) {
	tests := []struct {
		operand  int
		expected []byte
	}{
		{
			operand:  1,
			expected: []byte{0x01},
		},
		{
			operand:  23,
			expected: []byte{0x17},
		},
		{
			operand:  456,
			expected: []byte{0x01, 0xc8},
		},
	}

	for i, test := range tests {
		op := test.operand
		bytecode, err := encodeInt(op)

		if err != nil {
			t.Fatalf("test[%d] - encodeInt() had error. err=%v", i, err)
		}

		if !bytes.Equal(bytecode, test.expected) {
			t.Errorf("test[%d] - encodeInt() result wrong. expected=%x, got=%x", i, test.expected, bytecode)
		}
	}
}
