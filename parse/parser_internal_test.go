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

import (
	"errors"
	"testing"

	"github.com/DE-labtory/koa/ast"
)

type mockTokenBuffer struct {
	buf []Token
	sp  int
}

func (m *mockTokenBuffer) Read() Token {
	ret := m.buf[m.sp]
	m.sp++
	return ret
}

func (m *mockTokenBuffer) Peek(n peekNumber) Token {
	return m.buf[m.sp+int(n)]
}

func TestParse_curTokenIs(t *testing.T) {
	tokens := []Token{
		{Type: Int, Val: "1"},
		{Type: Ident, Val: "ADD"},
		{Type: Plus, Val: "+"},
		{Type: Asterisk, Val: "*"},
		{Type: Lparen, Val: "("},
	}
	tokenBuf := mockTokenBuffer{tokens, 0}
	tests := []struct {
		tokenType TokenType
		expected  bool
	}{
		{
			tokenType: Int,
			expected:  true,
		},
		{
			tokenType: Ident,
			expected:  true,
		},
		{
			tokenType: Mod,
			expected:  false,
		},
		{
			tokenType: Rbrace,
			expected:  false,
		},
		{
			tokenType: Lparen,
			expected:  true,
		},
	}

	for i, test := range tests {
		ret := curTokenIs(&tokenBuf, test.tokenType)
		if ret != test.expected {
			t.Fatalf("test[%d] - curTokenIs() result wrong. expected=%t, got=%t", i, test.expected, ret)
		}
		tokenBuf.Read()
	}
}

func TestParse_nextTokenIs(t *testing.T) {
	tokens := []Token{
		{Type: Int, Val: "1"},
		{Type: Ident, Val: "ADD"},
		{Type: Plus, Val: "+"},
		{Type: Asterisk, Val: "*"},
		{Type: Lparen, Val: "("},
	}
	tokenBuf := mockTokenBuffer{tokens, 0}
	tests := []struct {
		tokenType TokenType
		expected  bool
	}{
		{
			tokenType: Ident,
			expected:  true,
		},
		{
			tokenType: Plus,
			expected:  true,
		},
		{
			tokenType: Minus,
			expected:  false,
		},
		{
			tokenType: Rbrace,
			expected:  false,
		},
	}

	for i, test := range tests {
		ret := nextTokenIs(&tokenBuf, test.tokenType)
		if ret != test.expected {
			t.Fatalf("test[%d] - nextTokenIs() result wrong. expected=%t, got=%t", i, test.expected, ret)
		}
		tokenBuf.Read()
	}
}

func TestParse_curPrecedence(t *testing.T) {
	tokens := []Token{
		{Type: Int, Val: "1"},
		{Type: Ident, Val: "ADD"},
		{Type: Plus, Val: "+"},
		{Type: Asterisk, Val: "*"},
		{Type: Lparen, Val: "("},
	}
	tokenBuf := mockTokenBuffer{tokens, 0}
	tests := []struct {
		expected precedence
	}{
		{
			expected: 0,
		},
		{
			expected: 0,
		},
		{
			expected: SUM,
		},
		{
			expected: PRODUCT,
		},
		{
			expected: LOWEST,
		},
	}

	for i, test := range tests {
		ret := curPrecedence(&tokenBuf)
		if ret != test.expected {
			t.Fatalf("test[%d] - curPrecedence() result wrong. expected=%d, got=%d", i, test.expected, ret)
		}
		tokenBuf.Read()
	}
}

func TestParse_nextPrecedence(t *testing.T) {
	tokens := []Token{
		{Type: Int, Val: "1"},
		{Type: Ident, Val: "ADD"},
		{Type: Plus, Val: "+"},
		{Type: Asterisk, Val: "*"},
		{Type: Lparen, Val: "("},
	}
	tokenBuf := mockTokenBuffer{tokens, 0}
	tests := []struct {
		expected precedence
	}{
		{
			expected: 0,
		},
		{
			expected: SUM,
		},
		{
			expected: PRODUCT,
		},
		{
			expected: LOWEST,
		},
	}

	for i, test := range tests {
		ret := nextPrecedence(&tokenBuf)
		if ret != test.expected {
			t.Fatalf("test[%d] - curPrecedence() result wrong. expected=%d, got=%d", i, test.expected, ret)
		}
		tokenBuf.Read()
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

func TestExpectNext(t *testing.T) {
	err := errors.New("expectNext() : expecting token and next token are different")
	tokens := []Token{
		{Type: Int, Val: "1"},
		{Type: Ident, Val: "ADD"},
		{Type: Plus, Val: "+"},
		{Type: Asterisk, Val: "*"},
		{Type: Lparen, Val: "("},
	}
	tokenBuf := mockTokenBuffer{tokens, 0}
	tests := []struct {
		token         TokenType
		expectedBool  bool
		expectedError error
	}{
		{
			token:         Ident,
			expectedBool:  true,
			expectedError: nil,
		},
		{
			token:         Minus,
			expectedBool:  false,
			expectedError: err,
		},
		{
			token:         Asterisk,
			expectedBool:  true,
			expectedError: nil,
		},
		{
			token:         Rbrace,
			expectedBool:  false,
			expectedError: err,
		},
	}

	for i, test := range tests {
		retBool, retError := expectNext(&tokenBuf, test.token)
		if retBool != test.expectedBool {
			t.Fatalf("test[%d] - expectNext() result wrong.\n"+
				"expected bool: %t, error: %d\n"+
				"got bool: %t, error: %d", i, test.expectedBool, test.expectedError, retBool, retError)
		}
		tokenBuf.Read()
	}
}

func TestParseIdentifier(t *testing.T) {
	tokens := []Token{
		{Type: Int, Val: "1"},
		{Type: Ident, Val: "ADD"},
		{Type: Plus, Val: "+"},
		{Type: Asterisk, Val: "*"},
		{Type: Lparen, Val: "("},
	}
	tokenBuf := mockTokenBuffer{tokens, 0}
	tests := []struct {
		expected     ast.Expression
		expectedErrs string
	}{
		{
			expected:     nil,
			expectedErrs: "parseIdentifier() - 1 is not a identifier",
		},
		{
			expected: &ast.Identifier{
				Value: "ADD",
			},
		},
		{
			expected:     nil,
			expectedErrs: "parseIdentifier() - + is not a identifier",
		},
		{
			expected:     nil,
			expectedErrs: "parseIdentifier() - * is not a identifier",
		},
		{
			expected:     nil,
			expectedErrs: "parseIdentifier() - ( is not a identifier",
		},
	}

	for i, test := range tests {

		exp, errs := parseIdentifier(&tokenBuf)

		if errs != nil && errs[0].Error() != test.expectedErrs {
			t.Fatalf("test[%d] - wrong error. expected=%s, got=%s", i, test.expectedErrs, errs[0])
		}

		switch exp {
		case nil:
			if test.expected != nil {
				t.Fatalf("test[%d] - wrong result. expected=%s, got=%s", i, test.expected.String(), exp.String())
			}
			tokenBuf.Read()
		case &ast.Identifier{Value: exp.String()}:
			if exp.String() != exp.String() {
				t.Fatalf("test[%d] - wrong result. expected=%s, got=%s", i, test.expected.String(), exp.String())
			}
		}
	}
}
