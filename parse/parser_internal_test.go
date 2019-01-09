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

var mockError = errors.New("error occurred for some reason")

func TestCurTokenIs(t *testing.T) {
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

func TestNextTokenIs(t *testing.T) {
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

func TestCurPrecedence(t *testing.T) {
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
			expected: LOWEST,
		},
		{
			expected: LOWEST,
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

func TestNextPrecedence(t *testing.T) {
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
			expected: LOWEST,
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
			expectedErrs: "INT, parseIdentifier() - 1 is not a identifier",
		},
		{
			expected: &ast.Identifier{
				Value: "ADD",
			},
		},
		{
			expected:     nil,
			expectedErrs: "PLUS, parseIdentifier() - + is not a identifier",
		},
		{
			expected:     nil,
			expectedErrs: "ASTERISK, parseIdentifier() - * is not a identifier",
		},
		{
			expected:     nil,
			expectedErrs: "LPAREN, parseIdentifier() - ( is not a identifier",
		},
	}

	for i, test := range tests {

		exp, err := parseIdentifier(&tokenBuf)

		if err != nil && err.Error() != test.expectedErrs {
			t.Fatalf("test[%d] - wrong error. expected=%s, got=%s", i, test.expectedErrs, err)
		}

		switch exp {
		case nil:
			if test.expected != nil {
				t.Fatalf("test[%d] - wrong result. expected=%s, got=%s", i, test.expected.String(), exp.String())
			}
		case &ast.Identifier{Value: exp.String()}:
			if exp.String() != exp.String() {
				t.Fatalf("test[%d] - wrong result. expected=%s, got=%s", i, test.expected.String(), exp.String())
			}
		}
	}
}

func TestParseIntegerLiteral(t *testing.T) {
	tokens := []Token{
		{Type: Int, Val: "12"},
		{Type: Int, Val: "55"},
		{Type: Int, Val: "a"},
		{Type: String, Val: "abcdefg"},
		{Type: Int, Val: "-13"},
	}
	tokenBuf := mockTokenBuffer{tokens, 0}
	tests := []struct {
		expected    *ast.IntegerLiteral
		expectedErr error
	}{
		{expected: &ast.IntegerLiteral{Value: 12}, expectedErr: nil},
		{expected: &ast.IntegerLiteral{Value: 55}, expectedErr: nil},
		{expected: nil, expectedErr: errors.New("strconv.ParseInt: parsing \"a\": invalid syntax")},
		{expected: nil, expectedErr: errors.New("STRING, parseIntegerLiteral() error - abcdefg is not integer")},
		{expected: &ast.IntegerLiteral{Value: -13}, expectedErr: nil},
	}

	for i, test := range tests {
		exp, err := parseIntegerLiteral(&tokenBuf)

		switch err != nil {
		case true:
			if err.Error() != test.expectedErr.Error() {
				t.Fatalf("test[%d] - TestParseIntegerLiteral() wrong error. expected=%s, got=%s",
					i, test.expectedErr, err.Error())
			}
		}

		switch exp != nil {
		case true:
			if exp.String() != test.expected.String() {
				t.Fatalf("test[%d] - TestParseIntegerLiteral() wrong error. expected=%s, got=%s",
					i, test.expectedErr, err.Error())
			}
		}
	}
}

func TestParseBooleanLiteral(t *testing.T) {
	tokens := []Token{
		{Type: True, Val: "true"},
		{Type: False, Val: "false"},
		{Type: True, Val: "azzx"},
		{Type: String, Val: "abcdefg"},
	}
	tokenBuf := mockTokenBuffer{tokens, 0}
	tests := []struct {
		expected    *ast.BooleanLiteral
		expectedErr error
	}{
		{expected: &ast.BooleanLiteral{Value: true}, expectedErr: nil},
		{expected: &ast.BooleanLiteral{Value: false}, expectedErr: nil},
		{expected: nil, expectedErr: errors.New("strconv.ParseBool: parsing \"azzx\": invalid syntax")},
		{expected: nil, expectedErr: errors.New("STRING, parseBooleanLiteral() error - abcdefg is not bool")},
	}

	for i, test := range tests {
		exp, err := parseBooleanLiteral(&tokenBuf)

		switch err != nil {
		case true:
			if err.Error() != test.expectedErr.Error() {
				t.Fatalf("test[%d] - TestParseBooleanLiteral() wrong error. expected=%s, got=%s",
					i, test.expectedErr, err.Error())
			}
		}

		switch exp != nil {
		case true:
			if exp.String() != test.expected.String() {
				t.Fatalf("test[%d] - TestParseBooleanLiteral() wrong result. expected=%s, got=%s",
					i, test.expected, exp.String())
			}
		}
	}
}

func TestParseStringLiteral(t *testing.T) {
	tokens := []Token{
		{Type: String, Val: "hello"},
		{Type: String, Val: "hihi"},
		{Type: Int, Val: "3"},
		{Type: String, Val: "koa zzang"},
	}
	tokenBuf := mockTokenBuffer{tokens, 0}
	tests := []struct {
		expected    *ast.StringLiteral
		expectedErr error
	}{
		{expected: &ast.StringLiteral{Value: "hello"}, expectedErr: nil},
		{expected: &ast.StringLiteral{Value: "hihi"}, expectedErr: nil},
		{expected: nil, expectedErr: errors.New("INT, parseStringLiteral() error - 3 is not string")},
		{expected: &ast.StringLiteral{Value: "koa zzang"}, expectedErr: nil},
	}

	for i, test := range tests {
		exp, err := parseStringLiteral(&tokenBuf)

		switch err != nil {
		case true:
			if err.Error() != test.expectedErr.Error() {
				t.Fatalf("test[%d] - TestParseStringLiteral() wrong error. expected=%s, got=%s",
					i, test.expectedErr, err.Error())
			}
		}

		switch exp != nil {
		case true:
			if exp.String() != test.expected.String() {
				t.Fatalf("test[%d] - TestParseStringLiteral() wrong result. expected=%s, got=%s",
					i, test.expected, exp.String())
			}
		}
	}
}

func TestParseExpAsPrefix(t *testing.T) {
	// setup mockTokenBuffer
	tokens := []Token{
		{Type: Ident, Val: "a"},
		{Type: String, Val: "hello"},
		{Type: Plus, Val: "+"},
	}
	buf := mockTokenBuffer{tokens, 0}

	// mock prefixParseFnMap
	// In the case of Identifier, return normal expression
	prefixParseFnMap[String] = func(buf TokenBuffer) (ast.Expression, error) {
		return &ast.StringLiteral{"hello"}, nil
	}
	// In the case of Asterisk, return with errors
	prefixParseFnMap[Plus] = func(buf TokenBuffer) (ast.Expression, error) {
		return &ast.PrefixExpression{}, mockError
	}

	tests := []struct {
		expectedExpression ast.Expression
		expectedError      error
	}{
		{
			nil,
			parseError{Ident, "prefix parse function not defined"},
		},
		{
			&ast.StringLiteral{"hello"},
			nil,
		},
		{
			nil,
			mockError,
		},
	}

	for i, test := range tests {
		exp, err := makePrefixExpression(&buf)
		buf.Read()

		if exp != nil && exp.String() != test.expectedExpression.String() {
			t.Errorf("tests[%d] - Returned statements is not %s but got %s",
				i, test.expectedExpression.String(), exp.String())
		}

		if err != nil && err != test.expectedError {
			t.Errorf("tests[%d] - Returend error is not %s but got %s",
				i, test.expectedError, err)
		}
	}
}

func TestMakeInfixExpression(t *testing.T) {
	bufs := [][]Token{
		{{Type: Plus, Val: "+"}, {Type: Int, Val: "2"}, {Type: Asterisk, Val: "*"},
			{Type: Int, Val: "3"}, {Type: Eof, Val: ""}},
		{{Type: Asterisk, Val: "*"}, {Type: Int, Val: "242"}, {Type: Plus, Val: "+"},
			{Type: Int, Val: "312"}, {Type: Eof, Val: ""}},
		{{Type: Asterisk, Val: "-"}, {Type: Int, Val: "15"}, {Type: Plus, Val: "-"},
			{Type: Int, Val: "55"}, {Type: Eof, Val: ""}},
		{{Type: Minus, Val: "-"}, {Type: Int, Val: "2"}, {Type: Asterisk, Val: "*"},
			{Type: Int, Val: "3"}, {Type: Plus, Val: "+"}, {Type: Int, Val: "4"}, {Type: Eof, Val: ""}},
		{{Type: Plus, Val: "+"}, {Type: Plus, Val: "+"}, {Type: Eof, Val: ""}},
	}

	infixParseFnMap[Plus] = parseInfixExpression
	infixParseFnMap[Asterisk] = parseInfixExpression
	infixParseFnMap[Minus] = parseInfixExpression

	prefixParseFnMap[Int] = parseIntegerLiteral

	prefixes := []ast.IntegerLiteral{
		{Value: 1},
		{Value: 121},
		{Value: -10},
		{Value: 1},
		{Value: 1}}

	tests := []struct {
		expected string
	}{
		{expected: "(1 + (2 * 3))"},
		{expected: "((121 * 242) + 312)"},
		{expected: "((-10 - 15) - 55)"},
		{expected: "((1 - (2 * 3)) + 4)"},
		{expected: mockError.Error()},
	}

	// Expected value is
	//      +
	//     / \
	//    1   *
	//       / \
	//      2  3
	// result String() : 1+(2*3)

	for i, test := range tests {
		buf := mockTokenBuffer{bufs[i], 0}

		//exp, _ := makePrefixExpression(&buf)
		prefix := prefixes[i]
		exp, err := makeInfixExpression(&buf, &prefix, LOWEST)
		if err != nil && test.expected != err.Error() {
			t.Fatalf("test[%d] - TestMakeInfixExpression() wrong error. expected=%s, got=%s",
				i, test.expected, err.Error())
		}
		if err == nil && test.expected != exp.String() {
			t.Fatalf("test[%d] - TestMakeInfixExpression() wrong result. expected=%s, got=%s",
				i, test.expected, exp.String())
		}
	}
}

func TestParseInfixExpression(t *testing.T) {
	bufs := [][]Token{
		{{Type: Int, Val: "1"}, {Type: Plus, Val: "+"}, {Type: Int, Val: "2"}, {Type: Asterisk, Val: "*"},
			{Type: Int, Val: "3"}, {Type: Eof, Val: ""}},
		{{Type: Int, Val: "121"}, {Type: Asterisk, Val: "*"}, {Type: Int, Val: "242"}, {Type: Plus, Val: "+"},
			{Type: Int, Val: "312"}, {Type: Eof, Val: ""}},
		{{Type: Int, Val: "-10"}, {Type: Asterisk, Val: "-"}, {Type: Int, Val: "15"}, {Type: Plus, Val: "-"},
			{Type: Int, Val: "55"}, {Type: Eof, Val: ""}},
		{{Type: Int, Val: "1"}, {Type: Plus, Val: "+"}, {Type: Plus, Val: "+"}, {Type: Eof, Val: ""}},
	}

	infixParseFnMap[Plus] = parseInfixExpression
	infixParseFnMap[Asterisk] = parseInfixExpression
	prefixParseFnMap[Int] = parseIntegerLiteral

	tests := []struct {
		expected string
	}{
		{expected: "(1 + (2 * 3))"},
		{expected: "(121 * 242)"},
		{expected: "(-10 - 15)"},
		{expected: mockError.Error()},
	}

	lefts := []ast.IntegerLiteral{{Value: 1}, {Value: 121}, {Value: -10}, {Value: 1}}

	for i, test := range tests {
		buf := mockTokenBuffer{bufs[i], 1}

		left := lefts[i]
		exp, err := parseInfixExpression(&buf, &left)

		if err != nil && test.expected != err.Error() {
			t.Fatalf("test[%d] - TestMakeInfixExpression() wrong error. expected=%s, got=%s",
				i, test.expected, err.Error())
		}

		if err == nil && test.expected != exp.String() {
			t.Fatalf("test[%d] - TestMakeInfixExpression() wrong result. expected=%s, got=%s",
				i, test.expected, exp.String())
		}
	}
}

func TestParseGroupedExpression(t *testing.T) {
	prefixParseFnMap[Lparen] = parseGroupedExpression
	prefixParseFnMap[Int] = parseIntegerLiteral
	prefixParseFnMap[Ident] = parseIdentifier
	infixParseFnMap[Plus] = parseInfixExpression
	infixParseFnMap[Minus] = parseInfixExpression
	bufs := [][]Token{
		{
			// ( 2 + 1 )
			{Type: Lparen, Val: "("}, {Type: Int, Val: "2"}, {Type: Plus, Val: "+"}, {Type: Int, Val: "1"}, {Type: Rparen, Val: ")"},
			{Type: Eol, Val: "\n"},
		},
		{
			// ( a + ( 1 - 2 ) )
			{Type: Lparen, Val: "("}, {Type: Ident, Val: "a"}, {Type: Plus, Val: "+"}, {Type: Lparen, Val: "("},
			{Type: Int, Val: "1"}, {Type: Minus, Val: "-"}, {Type: Int, Val: "2"}, {Type: Rparen, Val: ")"}, {Type: Rparen, Val: ")"},
			{Type: Eol, Val: "\n"},
		},
		{
			// ( a + ( 1 - 2 ) + 3 )
			{Type: Lparen, Val: "("}, {Type: Ident, Val: "a"}, {Type: Plus, Val: "+"}, {Type: Lparen, Val: "("},
			{Type: Int, Val: "1"}, {Type: Minus, Val: "-"}, {Type: Int, Val: "2"}, {Type: Rparen, Val: ")"},
			{Type: Plus, Val: "+"}, {Type: Int, Val: "3"}, {Type: Rparen, Val: ")"}, {Type: Eol, Val: "\n"},
		},
		{
			{Type: Lparen, Val: "("}, {Type: Int, Val: "2"}, {Type: Plus, Val: "+"}, {Type: Int, Val: "1"},
			{Type: Rbrace, Val: "}"}, {Type: Eol, Val: "\n"},
		},
	}
	tests := []string{
		"(2 + 1)",
		"(a + (1 - 2))",
		"((a + (1 - 2)) + 3)",
		"RBRACE, is not Right paren",
	}

	for i, test := range tests {
		mockBuf := mockTokenBuffer{bufs[i], 0}
		exp, err := parseGroupedExpression(&mockBuf)

		if err != nil && err.Error() != test {
			t.Fatalf("test[%d] - TestParseGroupedExpression() wrong error. expected=%s, got=%s",
				i, test, err.Error())
		}

		if exp != nil && exp.String() != test {
			t.Fatalf("test[%d] - TestParseGroupedExpression() wrong answer. expected=%s, got=%s",
				i, test, exp.String())
		}
	}
}

func TestParseReturnStatement(t *testing.T) {
	prefixParseFnMap[True] = parseBooleanLiteral
	prefixParseFnMap[False] = parseBooleanLiteral
	prefixParseFnMap[Int] = parseIntegerLiteral
	infixParseFnMap[Plus] = parseInfixExpression
	infixParseFnMap[Asterisk] = parseInfixExpression

	bufs := [][]Token{
		{
			{Type: Return, Val: "return"},
			{Type: True, Val: "true"},
			{Type: Eol, Val: "\n"},
		},
		{
			{Type: Return, Val: "return"},
			{Type: Int, Val: "1"},
			{Type: Plus, Val: "+"},
			{Type: Int, Val: "2"},
			{Type: Eol, Val: "\n"},
		},
		{
			{Type: Return, Val: "return"},
			{Type: Int, Val: "1"},
			{Type: Plus, Val: "+"},
			{Type: Int, Val: "2"},
			{Type: Asterisk, Val: "*"},
			{Type: Int, Val: "3"},
			{Type: Eol, Val: "\n"},
		},
		{
			{Type: IntType, Val: "int"},
			{Type: Ident, Val: "a"},
			{Type: Assign, Val: "="},
			{Type: Int, Val: "1"},
			{Type: Eol, Val: "\n"},
		},
	}
	tests := []string{
		"return true",
		"return (1 + 2)",
		"return (1 + (2 * 3))",
		"INT_TYPE, parseReturnStatement() error - Statement must be started with return",
	}

	for i, test := range tests {
		mockBufs := mockTokenBuffer{bufs[i], 0}
		exp, err := parseReturnStatement(&mockBufs)
		if err != nil && err.Error() != test {
			t.Fatalf("test[%d] - TestParseReturnStatement() wrong error. expected=%s, got=%s",
				i, test, err.Error())
		}

		if exp != nil && exp.String() != test {
			t.Fatalf("test[%d] - TestParseReturnStatement() wrong result. expected=%s, got=%s",
				i, test, exp.String())
		}
	}
}

func TestParsePrefixExpression(t *testing.T) {
	tests := []struct {
		tokenBuffer      TokenBuffer
		expectedOperator string
		expectedRight    string
	}{
		{
			&mockTokenBuffer{
				buf: []Token{
					{Type: Minus, Val: "-"},
					{Type: Int, Val: "1"},
					{Type: Eof}},
				sp: 0,
			},
			"-", "1",
		},
		{
			&mockTokenBuffer{
				buf: []Token{
					{Type: Minus, Val: "-"},
					{Type: Int, Val: "3333"},
					{Type: Eof}},
				sp: 0,
			},
			"-", "3333",
		},
		{
			&mockTokenBuffer{
				buf: []Token{
					{Type: Minus, Val: "!"},
					{Type: True, Val: "true"},
					{Type: Eof}},
				sp: 0,
			},
			"!", "true",
		},
		{
			&mockTokenBuffer{
				buf: []Token{
					{Type: Minus, Val: "!"},
					{Type: False, Val: "false"},
					{Type: Eof}},
				sp: 0,
			},
			"!", "false",
		},
	}

	prefixParseFnMap[Int] = parseIntegerLiteral
	prefixParseFnMap[True] = parseBooleanLiteral
	prefixParseFnMap[False] = parseBooleanLiteral

	for i, tt := range tests {
		exp, err := parsePrefixExpression(tt.tokenBuffer)
		if err != nil {
			t.Errorf("tests[%d] - Returned error is \"%s\"",
				i, err)
		}

		expression, ok := exp.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp is not *ast.PrefixExpression. got=%T", exp)
		}

		if expression.Operator.String() != tt.expectedOperator {
			t.Errorf("tests[%d] - Type is not %s but got %s",
				i, tt.expectedOperator, expression.Operator.String())
		}

		if expression.Right.String() != tt.expectedRight {
			t.Errorf("tests[%d] - Value is not %s but got %s",
				i, tt.expectedRight, expression.Right.String())
		}
	}
}

func TestParseCallExpression(t *testing.T) {
	prefixParseFnMap[Int] = parseIntegerLiteral
	prefixParseFnMap[String] = parseStringLiteral
	infixParseFnMap[Plus] = parseInfixExpression
	infixParseFnMap[Asterisk] = parseInfixExpression
	bufs := [][]Token{
		{
			{Type: Lparen, Val: "("},
			{Type: Int, Val: "1"},
			{Type: Plus, Val: "+"},
			{Type: Int, Val: "2"},
			{Type: Rparen, Val: ")"},
		},
		{
			{Type: Lparen, Val: "("},
			{Type: String, Val: "a"},
			{Type: Comma, Val: ","},
			{Type: String, Val: "b"},
			{Type: Comma, Val: ","},
			{Type: Int, Val: "5"},
			{Type: Rparen, Val: ")"},
		},
		{
			{Type: Lparen, Val: "("},
			{Type: String, Val: "a"},
			{Type: Comma, Val: ","},
			{Type: Asterisk, Val: "*"},
			{Type: Comma, Val: ","},
			{Type: Int, Val: "5"},
			{Type: Rparen, Val: ")"},
		},
		{
			{Type: Lparen, Val: "("},
			{Type: String, Val: "a"},
			{Type: Plus, Val: "+"},
			{Type: String, Val: "b"},
			{Type: Comma, Val: ","},
			{Type: Int, Val: "5"},
			{Type: Asterisk, Val: "*"},
			{Type: Int, Val: "3"},
			{Type: Rparen, Val: ")"},
		},
	}
	tests := []struct {
		function ast.Expression
		expected string
	}{
		{
			function: &ast.StringLiteral{Value: "add"},
			expected: "function \"add\"( (1 + 2) )",
		},
		{
			function: &ast.StringLiteral{Value: "testFunc"},
			expected: "function \"testFunc\"( \"a\", \"b\", 5 )",
		},
		{
			function: &ast.StringLiteral{Value: "errorFunc"},
			expected: "ASTERISK, prefix parse function not defined",
		},
		{
			function: &ast.StringLiteral{Value: "complexFunc"},
			expected: "function \"complexFunc\"( (\"a\" + \"b\"), (5 * 3) )",
		},
	}

	for i, test := range tests {
		mockBuf := mockTokenBuffer{bufs[i], 0}
		exp, err := parseCallExpression(&mockBuf, test.function)
		if err != nil && err.Error() != test.expected {
			t.Fatalf("test[%d] - parseCallExpression() wrong error. expected=%s, got=%s",
				i, test.expected, err.Error())
		}
		if exp != nil && exp.String() != test.expected {
			t.Fatalf("test[%d] - parseCallExpression() wrong answer. expected=%s, got=%s",
				i, test.expected, exp.String())
		}
	}
}

func TestParseCallArguments(t *testing.T) {
	prefixParseFnMap[Int] = parseIntegerLiteral
	prefixParseFnMap[String] = parseStringLiteral
	infixParseFnMap[Plus] = parseInfixExpression
	infixParseFnMap[Asterisk] = parseInfixExpression
	bufs := [][]Token{
		{
			{Type: Lparen, Val: "("},
			{Type: Rparen, Val: ")"},
		},
		{
			{Type: Lparen, Val: "("},
			{Type: Int, Val: "1"},
			{Type: Rparen, Val: ")"},
		},
		{
			{Type: Lparen, Val: "("},
			{Type: String, Val: "a"},
			{Type: Comma, Val: ","},
			{Type: String, Val: "b"},
			{Type: Comma, Val: ","},
			{Type: Int, Val: "5"},
			{Type: Rparen, Val: ")"},
		},
		{
			{Type: Lparen, Val: "("},
			{Type: String, Val: "a"},
			{Type: Plus, Val: "+"},
			{Type: String, Val: "b"},
			{Type: Comma, Val: ","},
			{Type: Int, Val: "5"},
			{Type: Asterisk, Val: "*"},
			{Type: Int, Val: "3"},
			{Type: Rparen, Val: ")"},
		},
		{
			{Type: Lparen, Val: "("},
			{Type: Asterisk, Val: "*"},
			{Type: Rparen, Val: ")"},
		},
		{
			{Type: Lparen, Val: "("},
			{Type: String, Val: "a"},
			{Type: Comma, Val: ","},
			{Type: Asterisk, Val: "*"},
			{Type: Comma, Val: ","},
			{Type: Int, Val: "5"},
			{Type: Rparen, Val: ")"},
		},
	}
	tests := []string{
		"function \"testFunction\"(  )",
		"function \"testFunction\"( 1 )",
		"function \"testFunction\"( \"a\", \"b\", 5 )",
		"function \"testFunction\"( (\"a\" + \"b\"), (5 * 3) )",
		"ASTERISK, prefix parse function not defined",
		"ASTERISK, prefix parse function not defined",
	}

	for i, test := range tests {
		mockBufs := mockTokenBuffer{bufs[i], 0}
		exp, err := parseCallArguments(&mockBufs)

		if err != nil && err.Error() != test {
			t.Fatalf("test[%d] - TestParseCallArguments() wrong error. expected=%s, got=%s",
				i, test, err.Error())
		}

		mockFn := &ast.CallExpression{
			Function: &ast.StringLiteral{Value: "testFunction"},
		}
		mockFn.Arguments = exp
		if exp != nil && mockFn.String() != test {
			t.Fatalf("test[%d] - TestParseCallArguments() wrong error. expected=%s, got=%s",
				i, test, mockFn.String())
		}
	}
}

func TestParseAssignStatement(t *testing.T) {
	tests := []struct {
		tokenBuffer           TokenBuffer
		expectedDataStructure string
		expectedIdent         string
		expectedVal           string
	}{
		{
			&mockTokenBuffer{
				buf: []Token{
					{Type: StringType, Val: "string"},
					{Type: Ident, Val: "a"},
					{Type: Assign, Val: "="},
					{Type: String, Val: "hello"},
					{Type: Eof}},
				sp: 0,
			},
			"string", "val: a, type: IDENT", "\"hello\"",
		},
		{
			&mockTokenBuffer{
				buf: []Token{
					{Type: IntType, Val: "int"},
					{Type: Ident, Val: "myInt"},
					{Type: Assign, Val: "="},
					{Type: Int, Val: "1"},
					{Type: Eof}},
				sp: 0,
			},
			"int", "val: myInt, type: IDENT", "1",
		},
		{
			&mockTokenBuffer{
				buf: []Token{
					{Type: BoolType, Val: "bool"},
					{Type: Ident, Val: "ddd"},
					{Type: Assign, Val: "="},
					{Type: Bool, Val: "true"},
					{Type: Eof}},
				sp: 0,
			},
			"bool", "val: ddd, type: IDENT", "true",
		},
		{
			// type mismatch tc - int ddd2 = "iam_string"
			&mockTokenBuffer{
				buf: []Token{
					{Type: IntType, Val: "int"},
					{Type: Ident, Val: "ddd2"},
					{Type: Assign, Val: "="},
					{Type: String, Val: "iam_string"},
					{Type: Eof}},
				sp: 0,
			},
			"int", "val: ddd2, type: IDENT", "\"iam_string\"",
		},
		{
			// type mismatch tc - bool foo = "iam_string"
			&mockTokenBuffer{
				buf: []Token{
					{Type: BoolType, Val: "bool"},
					{Type: Ident, Val: "foo"},
					{Type: Assign, Val: "="},
					{Type: String, Val: "iam_string"},
					{Type: Eof}},
				sp: 0,
			},
			"bool", "val: foo, type: IDENT", "\"iam_string\"",
		},
	}

	prefixParseFnMap[Int] = parseIntegerLiteral
	prefixParseFnMap[String] = parseStringLiteral
	prefixParseFnMap[Bool] = parseBooleanLiteral

	for i, tt := range tests {
		exp, err := parseAssignStatement(tt.tokenBuffer)

		if err != nil {
			t.Errorf("tests[%d] - Returned error is \"%s\"",
				i, err)
		}

		if exp.Type.String() != tt.expectedDataStructure {
			t.Errorf("tests[%d] - Type is not %s but got %s",
				i, tt.expectedDataStructure, exp.Type.String())
		}

		if exp.Variable.String() != tt.expectedIdent {
			t.Errorf("tests[%d] - Variable is not %s but got %s",
				i, tt.expectedIdent, exp.Variable.String())
		}

		if exp.Value.String() != tt.expectedVal {
			t.Errorf("tests[%d] - Value is not %s but got %s",
				i, tt.expectedVal, exp.Value.String())
		}
	}
}

// TestParseExpression tests strings which combine prefix and
// infix expression
func TestParseExpression(t *testing.T) {
	initParseFnMap()
	tokens := [][]Token{
		{
			{Type: Minus, Val: "-"},
			{Type: Ident, Val: "a"},
			{Type: Asterisk, Val: "*"},
			{Type: Ident, Val: "b"},
			{Type: Eof},
		},

		{
			{Type: Bang, Val: "!"},
			{Type: Minus, Val: "-"},
			{Type: Ident, Val: "b"},
			{Type: Eof},
		},
		{
			{Type: Minus, Val: "-"},
			{Type: Int, Val: "33"},
			{Type: Slash, Val: "/"},
			{Type: Int, Val: "67"},
			{Type: Plus, Val: "+"},
			{Type: Ident, Val: "a"},
			{Type: Eof},
		},
		{
			{Type: Int, Val: "33"},
			{Type: Mod, Val: "%"},
			{Type: Minus, Val: "-"},
			{Type: Int, Val: "67"},
			{Type: Plus, Val: "+"},
			{Type: Ident, Val: "a"},
			{Type: Asterisk, Val: "*"},
			{Type: Ident, Val: "c"},
			{Type: Eof},
		},
		{
			{Type: Int, Val: "33"},
			{Type: Mod, Val: "%"},
			{Type: Lparen, Val: "("},
			{Type: Minus, Val: "-"},
			{Type: Int, Val: "67"},
			{Type: Plus, Val: "+"},
			{Type: Ident, Val: "a"},
			{Type: Rparen, Val: ")"},
			{Type: Asterisk, Val: "*"},
			{Type: Ident, Val: "c"},
			{Type: Eof},
		},
		{
			{Type: Minus, Val: "-"},
			{Type: Int, Val: "33"},
			{Type: Slash, Val: "/"},
			{Type: Int, Val: "67"},
			{Type: LT, Val: "<"},
			{Type: Ident, Val: "a"},
			{Type: Asterisk, Val: "*"},
			{Type: Int, Val: "67"},
			{Type: Eof},
		},
		{
			{Type: Minus, Val: "-"},
			{Type: Int, Val: "33"},
			{Type: Slash, Val: "/"},
			{Type: Int, Val: "67"},
			{Type: GTE, Val: ">="},
			{Type: Ident, Val: "a"},
			{Type: Plus, Val: "+"},
			{Type: Int, Val: "67"},
			{Type: Mod, Val: "%"},
			{Type: Ident, Val: "z"},
			{Type: Eof},
		},
	}

	tests := []struct {
		expected string
		err      error
	}{
		{"((-a) * b)", nil},
		{"(!(-b))", nil},
		{"(((-33) / 67) + a)", nil},
		{"((33 % (-67)) + (a * c))", nil},
		{"((33 % ((-67) + a)) * c)", nil},
		{"(((-33) / 67) < (a * 67))", nil},
		{"(((-33) / 67) >= (a + (67 % z)))", nil},
	}

	for i, tt := range tests {
		buf := &mockTokenBuffer{tokens[i], 0}
		exp, err := parseExpression(buf, LOWEST)

		if err != nil {
			t.Fatalf("test[%d] - parseExpression() with wrong error. expected=%s, got=%s",
				i, tt.err, err)
		}

		if err == nil && exp.String() != tt.expected {
			t.Fatalf("test[%d] - parseExpression() with wrong expression. expected=%s, got=%s",
				i, tt.expected, exp.String())
		}
	}
}

func TestParseIfStatement(t *testing.T) {
	prefixParseFnMap[Bool] = parseBooleanLiteral

	bufs := [][]Token{
		{
			{Type: If, Val: "if"},
			{Type: Lparen, Val: "("},
			{Type: Bool, Val: "true"},
			{Type: Rparen, Val: ")"},
			{Type: Lbrace, Val: "{"},
			{Type: IntType, Val: "int"},
			{Type: Ident, Val: "a"},
			{Type: Rbrace, Val: "}"},
		},
	}

	tests := []string{
		// This will be "if ( true ) { int a }"
		"if ( true ) {  }",
	}

	for i, test := range tests {
		mockBuf := mockTokenBuffer{bufs[i], 0}
		result, err := parseIfStatement(&mockBuf)

		if len(err) > 0 && err[0].Error() != test {
			t.Fatalf("test[%d] - TestParseIfStatement() wrong error. expected=%s got=%s",
				i, test, err[0].Error())
		}

		if result.String() != test {
			t.Fatalf("test[%d] - TestParseIfStatement() wrong result. expected=%s, got=%s",
				i, test, result.String())
		}

	}
}

func TestParseBlockStatement(t *testing.T) {

}
