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

package parse

import "testing"

type MockTokenBuffer struct {
	buf []Token
	sp  int
}

func (m *MockTokenBuffer) Read() Token {
	ret := m.buf[m.sp]
	m.sp++
	return ret
}

func (m *MockTokenBuffer) Peek(n peekNumber) Token {
	return m.buf[m.sp+int(n)]
}

func TestParse_curTokenIs_nextTokenIs(t *testing.T) {
	tokenBuf := MockTokenBuffer{
		buf: []Token{
			{
				Type:   Int,
				Val:    "1",
				Column: 1,
				Line:   5,
			},
			{
				Type:   Ident,
				Val:    "ADD",
				Column: 2,
				Line:   8,
			},
		},
	}

	ret := curTokenIs(&tokenBuf, Int)
	if !ret {
		t.Errorf("Expected value is true")
	}

	ret = nextTokenIs(&tokenBuf, Ident)
	if !ret {
		t.Errorf("Expected value is true")
	}
}

func TestPeekNumber_IsValid(t *testing.T) {
	tests := []struct {
		n        peekNumber
		expected bool
	}{
		{
			n:        peekNumber(0),
			expected: true,
		},
		{
			n:        peekNumber(1),
			expected: true,
		},
		{
			n:        peekNumber(2),
			expected: false,
		},
		{
			n:        peekNumber(-1),
			expected: false,
		},
	}

	for i, test := range tests {
		n := test.n
		if n.isValid() != test.expected {
			t.Fatalf("test[%d] - isValid() result wrong. expected=%t, got=%t", i, test.expected, n.isValid())
		}
	}
}
