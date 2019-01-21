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
	if m.sp+1 < len(m.buf) {
		m.sp++
	}
	return ret
}

func (m *mockTokenBuffer) Peek(n peekNumber) Token {
	return m.buf[m.sp+int(n)]
}

var mockError = errors.New("error occurred for some reason")

// TestParserOnly tests three things
//
// 1. "contract" keyword with its open-brace & close-brace
// 2. When there's single & multiple function inside contract
// 3. When there's statements other than function literal
//
func TestParserOnly(t *testing.T) {
	tests := []struct {
		tokens      []Token
		expected    string
		expectedErr error
	}{
		{
			tokens: []Token{
				{Type: Contract, Val: "contract"},
				{Type: Lbrace, Val: "{"},
				{Type: Rbrace, Val: "}"},
				{Type: Eof},
			},
			expected: `
contract {
}`,
		},
		{
			tokens: []Token{
				{Type: Contract, Val: "contract"},
				{Type: Lbrace, Val: "{"},
				{Type: Eof},
			},
			expected:    "",
			expectedErr: parseError{Eof, "expectNext() : expected [RBRACE], but got [EOF]"},
		},
		{
			tokens: []Token{
				{Type: Contract, Val: "contract"},
				{Type: Lbrace, Val: "{"},
				{Type: Function, Val: "func"},
				{Type: Ident, Val: "foo"},
				{Type: Lparen, Val: "("},
				{Type: Rparen, Val: ")"},
				{Type: Lbrace, Val: "{"},
				{Type: Rbrace, Val: "}"},
				{Type: Rbrace, Val: "}"},
				{Type: Eof},
			},
			expected: `
contract {
func foo() void {

}
}`,
		},
		{
			tokens: []Token{
				{Type: Contract, Val: "contract"},
				{Type: Lbrace, Val: "{"},
				{Type: Function, Val: "func"},
				{Type: Ident, Val: "foo"},
				{Type: Lparen, Val: "("},
				{Type: Rparen, Val: ")"},
				{Type: Lbrace, Val: "{"},
				{Type: Rbrace, Val: "}"},
				{Type: Function, Val: "func"},
				{Type: Ident, Val: "bar"},
				{Type: Lparen, Val: "("},
				{Type: Rparen, Val: ")"},
				{Type: Lbrace, Val: "{"},
				{Type: Rbrace, Val: "}"},
				{Type: Rbrace, Val: "}"},
				{Type: Eof},
			},
			expected: `
contract {
func foo() void {

}
func bar() void {

}
}`,
		},
		{
			tokens: []Token{
				{Type: Contract, Val: "contract"},
				{Type: Lbrace, Val: "{"},
				{Type: IntType, Val: "int"},
				{Type: Ident, Val: "a"},
				{Type: Assign, Val: "="},
				{Type: Int, Val: "1"},
				{Type: Rbrace, Val: "}"},
				{Type: Eof},
			},
			expected:    ``,
			expectedErr: parseError{IntType, "expectNext() : expected [RBRACE], but got [INT_TYPE]"},
		},
	}

	for i, tt := range tests {
		buf := &mockTokenBuffer{tt.tokens, 0}
		stmt, err := Parse(buf)

		if err != nil && err != tt.expectedErr {
			t.Errorf(`test[%d] - Wrong error returned expected="%v", got="%v"`,
				i, tt.expectedErr, err)
			continue
		}

		if err == nil && stmt.String() != tt.expected {
			t.Errorf(`test[%d] - Wrong result returned expected="%s", got="%s"`,
				i, tt.expected, stmt.String())
		}
	}
}

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
			expected: CALL,
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
			expected: CALL,
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
			token:         Int,
			expectedBool:  true,
			expectedError: nil,
		},
		{
			token:         Minus,
			expectedBool:  false,
			expectedError: errors.New("IDENT, expectNext() : expected [MINUS], but got [IDENT]"),
		},
		{
			token:         Plus,
			expectedBool:  true,
			expectedError: nil,
		},
		{
			token:         Rbrace,
			expectedBool:  false,
			expectedError: errors.New("ASTERISK, expectNext() : expected [RBRACE], but got [ASTERISK]"),
		},
	}

	for i, test := range tests {
		retError := expectNext(&tokenBuf, test.token)
		if retError != nil && retError.Error() != test.expectedError.Error() {
			t.Fatalf("test[%d] - expectNext() result wrong.\n"+
				"expected error: %s\n"+
				"got error: %s", i, test.expectedError.Error(), retError.Error())
		}

		if retError != nil {
			tokenBuf.Read()
		}
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
		{expected: nil, expectedErr: errors.New(`strconv.ParseInt: parsing "a": invalid syntax`)},
		{expected: nil, expectedErr: errors.New("STRING, parseIntegerLiteral() error - abcdefg is not integer")},
		{expected: &ast.IntegerLiteral{Value: -13}, expectedErr: nil},
	}

	for i, test := range tests {
		// For debugging
		tokenBuf.sp = i
		exp, err := parseIntegerLiteral(&tokenBuf)
		if err != nil && err.Error() != test.expectedErr.Error() {
			t.Fatalf("test[%d] - TestParseIntegerLiteral() wrong error. expected=%s, got=%s",
				i, test.expectedErr, err.Error())
		}

		if exp != nil && exp.String() != test.expected.String() {
			t.Fatalf("test[%d] - TestParseIntegerLiteral() wrong error. expected=%s, got=%s",
				i, test.expectedErr, err.Error())
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
		{expected: &ast.BooleanLiteral{true}, expectedErr: nil},
		{expected: &ast.BooleanLiteral{false}, expectedErr: nil},
		{expected: nil, expectedErr: errors.New(`strconv.ParseBool: parsing "azzx": invalid syntax`)},
		{expected: nil, expectedErr: errors.New("STRING, parseBooleanLiteral() error - abcdefg is not bool")},
	}

	for i, test := range tests {
		exp, err := parseBooleanLiteral(&tokenBuf)

		if err != nil && err.Error() != test.expectedErr.Error() {
			t.Fatalf(`test[%d] - TestParseBooleanLiteral() wrong error. expected="%s", got="%s"`,
				i, test.expectedErr.Error(), err.Error())
		}

		lit, ok := exp.(*ast.BooleanLiteral)
		if err == nil && !ok {
			t.Fatalf("test[%d] - TestParseBooleanLiteral() returned expression is not *ast.BooleanLiteral", i)
		}

		if err == nil && lit.String() != test.expected.String() {
			t.Fatalf(`test[%d] - TestParseBooleanLiteral() wrong result. expected="%s", got="%s"`,
				i, test.expected, lit.String())
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
		// For debbuging
		tokenBuf.sp = i
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

func TestParseFunctionLiteral(t *testing.T) {
	bufs := [][]Token{
		{
			// func example (a int, b string) {}
			{Type: Function, Val: "func"},
			{Type: Ident, Val: "example"},
			{Type: Lparen, Val: "("},
			{Type: Ident, Val: "a"},
			{Type: IntType, Val: "int"},
			{Type: Comma, Val: ","},
			{Type: Ident, Val: "b"},
			{Type: StringType, Val: "string"},
			{Type: Rparen, Val: ")"},
			{Type: Lbrace, Val: "{"},
			{Type: Rbrace, Val: "}"},
		},
		{
			// func name (a int, b string) {
			//	int c = 5
			// }
			{Type: Function, Val: "func"},
			{Type: Ident, Val: "name"},
			{Type: Lparen, Val: "("},
			{Type: Ident, Val: "a"},
			{Type: IntType, Val: "int"},
			{Type: Comma, Val: ","},
			{Type: Ident, Val: "b"},
			{Type: StringType, Val: "string"},
			{Type: Rparen, Val: ")"},
			{Type: Lbrace, Val: "{"},
			{Type: IntType, Val: "int"},
			{Type: Ident, Val: "c"},
			{Type: Assign, Val: "="},
			{Type: Int, Val: "5"},
			{Type: Rbrace, Val: "}"},
		},
		{
			// error case
			{Type: Lbrace, Val: "{"},
		},
		{
			// func example () string {}
			{Type: Function, Val: "func"},
			{Type: Ident, Val: "example"},
			{Type: Lparen, Val: "("},
			{Type: Rparen, Val: ")"},
			{Type: StringType, Val: "string"},
			{Type: Lbrace, Val: "{"},
			{Type: Rbrace, Val: "}"},
		},
		{
			// func example () invalid {}
			{Type: Function, Val: "func"},
			{Type: Ident, Val: "example"},
			{Type: Lparen, Val: "("},
			{Type: Rparen, Val: ")"},
			{Type: Illegal, Val: "invalid"},
			{Type: Lbrace, Val: "{"},
			{Type: Rbrace, Val: "}"},
		},
		{
			// func example () invalid {}
			{Type: Function, Val: "func"},
			{Type: Ident, Val: "example"},
			{Type: Lparen, Val: "("},
			{Type: Rparen, Val: ")"},
			{Type: Int, Val: "1"},
			{Type: Lbrace, Val: "{"},
			{Type: Rbrace, Val: "}"},
		},
	}
	tests := []string{
		"func example(Parameter : (Identifier: a, Type: int), Parameter : (Identifier: b, Type: string)) void {\n\n}",
		"func name(Parameter : (Identifier: a, Type: int), Parameter : (Identifier: b, Type: string)) void {\nint [IDENT, c] = 5\n}",
		"LBRACE, expectNext() : expected [FUNCTION], but got [LBRACE]",
		"func example() string {\n\n}",
		"ILLEGAL, invalid function return type",
		"INT, invalid function return type",
	}
	prefixParseFnMap[Int] = parseIntegerLiteral
	infixParseFnMap[Plus] = parseInfixExpression

	for i, test := range tests {
		buf := mockTokenBuffer{bufs[i], 0}
		exp, err := parseFunctionLiteral(&buf)
		if err != nil && err.Error() != test {
			t.Fatalf("test[%d] - TestParseFunctionLiteral() wrong error\n"+
				"expected: %s\n"+
				"got: %s", i, test, err.Error())
		}

		if exp != nil && exp.String() != test {
			t.Fatalf("test[%d] - TestParseFunctionLiteral wrong result\n"+
				"expected: %s\n"+
				"got: %s", i, test, exp.String())
		}
	}
}

func TestParseFunctionParameter(t *testing.T) {
	bufs := [][]Token{
		{
			{Type: Ident, Val: "a"},
			{Type: IntType, Val: "int"},
			{Type: Comma, Val: ","},
			{Type: Ident, Val: "b"},
			{Type: StringType, Val: "string"},
			{Type: Rparen, Val: ")"},
		},
		{
			{Type: Ident, Val: "arg"},
			{Type: BoolType, Val: "bool"},
			{Type: Rparen, Val: ")"},
		},
		{
			{Type: Rparen, Val: ")"},
		},
		{
			{Type: Ident, Val: "arg"},
			{Type: IntType, Val: "int"},
			{Type: Rbrace, Val: "}"},
		},
	}
	tests := [][]string{
		{
			"Parameter : (Identifier: a, Type: int)",
			"Parameter : (Identifier: b, Type: string)",
		},
		{
			"Parameter : (Identifier: arg, Type: bool)",
		},
		{
			"Parameter : ()",
		},
		{
			"RBRACE, expectNext() : expected [RPAREN], but got [RBRACE]",
		},
	}

	for i, test := range tests {
		buf := mockTokenBuffer{bufs[i], 0}
		idents, err := parseFunctionParameters(&buf)
		if err != nil && err.Error() != test[0] {
			t.Fatalf("test[%d] - TestParseFunctionParameter() wrong error.\n"+
				"expected: %s\n"+
				"got: %s", i, test[0], err.Error())
		} else {
			for j, ident := range idents {
				if ident.String() != tests[i][j] {
					t.Fatalf("test[%d-%d] - TestParseFunctionParameter() failed.\n"+
						"expected: %s\n"+
						"got: %s", i, j, tests[i][j], ident)
				}
			}
		}
	}
}

func TestMakePrefixExpression(t *testing.T) {
	initParseFnMap()

	tests := []struct {
		tokens      []Token
		expected    string
		expectedErr error
	}{
		{
			tokens: []Token{
				{Type: Minus, Val: "-"},
				{Type: Int, Val: "1"},
			},
			expected: "(-1)",
		},
		{
			tokens: []Token{
				{Type: Minus, Val: "-"},
				{Type: Ident, Val: "a"},
			},
			expected: "(-a)",
		},
		{
			tokens: []Token{
				{Type: Bang, Val: "!"},
				{Type: True, Val: "true"},
			},
			expected: "(!true)",
		},
		{
			tokens: []Token{
				{Type: Bang, Val: "!"},
				{Type: Bang, Val: "!"},
				{Type: True, Val: "false"},
			},
			expected: "(!(!false))",
		},
		{
			tokens: []Token{
				{Type: Bang, Val: "!"},
				{Type: Minus, Val: "-"},
				{Type: Ident, Val: "foo"},
			},
			expected: "(!(-foo))",
		},
		{
			tokens: []Token{
				{Type: Minus, Val: "-"},
				{Type: Bang, Val: "!"},
				{Type: Ident, Val: "foo"},
			},
			expected: "(-(!foo))",
		},
		{
			tokens: []Token{
				{Type: Minus, Val: "-"},
				{Type: True, Val: "true"},
			},
			expected: "(-true)",
		},
		{
			tokens: []Token{
				{Type: Bang, Val: "!"},
				{Type: String, Val: "hello"},
			},
			expected: `(!"hello")`,
		},
	}

	for i, tt := range tests {
		buf := &mockTokenBuffer{tt.tokens, 0}
		exp, err := makePrefixExpression(buf)

		if err != nil && err != tt.expectedErr {
			t.Errorf(`test[%d] - Wrong error returned expected="%v", got="%v"`,
				i, tt.expectedErr, err)
			continue
		}

		if err == nil && exp.String() != tt.expected {
			t.Errorf(`test[%d] - Wrong result returned expected="%s", got="%s"`,
				i, tt.expected, exp.String())
		}
	}
}

func TestMakeInfixExpression(t *testing.T) {
	initParseFnMap()

	bufs := [][]Token{
		{{Type: Plus, Val: "+"}, {Type: Int, Val: "2"}, {Type: Asterisk, Val: "*"},
			{Type: Int, Val: "3"}, {Type: Eof, Val: ""}},
		{{Type: Asterisk, Val: "*"}, {Type: Int, Val: "242"}, {Type: Plus, Val: "+"},
			{Type: Int, Val: "312"}, {Type: Eof, Val: ""}},
		{{Type: Minus, Val: "-"}, {Type: Int, Val: "15"}, {Type: Minus, Val: "-"},
			{Type: Int, Val: "55"}, {Type: Eof, Val: ""}},
		{{Type: Minus, Val: "-"}, {Type: Int, Val: "2"}, {Type: Asterisk, Val: "*"},
			{Type: Int, Val: "3"}, {Type: Plus, Val: "+"}, {Type: Int, Val: "4"}, {Type: Eof, Val: ""}},
		{{Type: Plus, Val: "+"}, {Type: Plus, Val: "+"}, {Type: Eof, Val: ""}},
	}

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
		{expected: "PLUS, prefix parse function not defined"},
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
		{{Type: Int, Val: "-10"}, {Type: Asterisk, Val: "*"}, {Type: Int, Val: "15"}, {Type: Plus, Val: "+"},
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
		{expected: "(-10 * 15)"},
		{expected: "PLUS, prefix parse function not defined"},
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
		"RBRACE, expectNext() : expected [RPAREN], but got [RBRACE]",
	}

	for i, test := range tests {
		mockBuf := mockTokenBuffer{bufs[i], 0}
		exp, err := parseGroupedExpression(&mockBuf)

		if err != nil && err.Error() != test {
			t.Fatalf("test[%d] - TestParseGroupedExpression() wrong error.\n"+
				"expected=%s,\n"+
				"got=%s",
				i, test, err.Error())
		}

		if exp != nil && exp.String() != test {
			t.Fatalf("test[%d] - TestParseGroupedExpression() wrong answer.\n"+
				"expected=%s,\n"+
				"got=%s",
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
		"INT_TYPE, expectNext() : expected [RETURN], but got [INT_TYPE]",
	}

	for i, test := range tests {
		mockBufs := mockTokenBuffer{bufs[i], 0}
		exp, err := parseReturnStatement(&mockBufs)
		if err != nil && err.Error() != test {
			t.Fatalf("test[%d] - TestParseReturnStatement() wrong error.\n"+
				"expected=%s,\n"+
				"got=%s",
				i, test, err.Error())
		}

		if exp != nil && exp.String() != test {
			t.Fatalf("test[%d] - TestParseReturnStatement() wrong result.\n"+
				"expected=%s,\n"+
				"got=%s",
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
					{Type: Bang, Val: "!"},
					{Type: True, Val: "true"},
					{Type: Eof}},
				sp: 0,
			},
			"!", "true",
		},
		{
			&mockTokenBuffer{
				buf: []Token{
					{Type: Bang, Val: "!"},
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
			t.Errorf(`tests[%d] - Returned error is "%s"`,
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
			function: &ast.Identifier{Value: "add"},
			expected: `function add( (1 + 2) )`,
		},
		{
			function: &ast.Identifier{Value: "testFunc"},
			expected: `function testFunc( "a", "b", 5 )`,
		},
		{
			function: &ast.Identifier{Value: "errorFunc"},
			expected: "ASTERISK, prefix parse function not defined",
		},
		{
			function: &ast.Identifier{Value: "complexFunc"},
			expected: `function complexFunc( ("a" + "b"), (5 * 3) )`,
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
		`function testFunction(  )`,
		`function testFunction( 1 )`,
		`function testFunction( "a", "b", 5 )`,
		`function testFunction( ("a" + "b"), (5 * 3) )`,
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
			Function: &ast.Identifier{Value: "testFunction"},
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
		expectedErr           error
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
			"string", "[IDENT, a]", `"hello"`,
			nil,
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
			"int", "[IDENT, myInt]", "1",
			nil,
		},
		{
			&mockTokenBuffer{
				buf: []Token{
					{Type: BoolType, Val: "bool"},
					{Type: Ident, Val: "ddd"},
					{Type: Assign, Val: "="},
					{Type: True, Val: "true"},
					{Type: Eof}},
				sp: 0,
			},
			"bool", "[IDENT, ddd]", "true",
			nil,
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
			"int", "[IDENT, ddd2]", `"iam_string"`,
			nil,
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
			"bool", "[IDENT, foo]", `"iam_string"`,
			nil,
		},
		{
			&mockTokenBuffer{
				buf: []Token{
					{Type: BoolType, Val: "bool"},
					{Type: String, Val: "foo"},
					{Type: Assign, Val: "="},
					{Type: String, Val: "iam_string"},
					{Type: Eof}},
				sp: 0,
			},
			"bool", "[IDENT, foo]", `"iam_string"`,
			parseError{String, "token is not identifier"},
		},
		{
			&mockTokenBuffer{
				buf: []Token{
					{Type: BoolType, Val: "bool"},
					{Type: Ident, Val: "foo"},
					{Type: String, Val: "iam_string"},
					{Type: Eof}},
				sp: 0,
			},
			"bool", "[IDENT, foo]", `"iam_string"`,
			parseError{String, "token is not assign"},
		},
	}

	prefixParseFnMap[Int] = parseIntegerLiteral
	prefixParseFnMap[String] = parseStringLiteral
	prefixParseFnMap[True] = parseBooleanLiteral
	prefixParseFnMap[False] = parseBooleanLiteral

	for i, tt := range tests {
		exp, err := parseAssignStatement(tt.tokenBuffer)

		if err != nil && err != tt.expectedErr {
			t.Errorf(`tests[%d] - Returned err is not "%s", but got "%s"`,
				i, tt.expectedErr.Error(), err.Error())
		}

		if err == nil && exp.Type.String() != tt.expectedDataStructure {
			t.Errorf("tests[%d] - Type is not %s but got %s",
				i, tt.expectedDataStructure, exp.Type.String())
		}

		if err == nil && exp.Variable.String() != tt.expectedIdent {
			t.Errorf("tests[%d] - Variable is not %s but got %s",
				i, tt.expectedIdent, exp.Variable.String())
		}

		if err == nil && exp.Value.String() != tt.expectedVal {
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
	initParseFnMap()

	bufs := [][]Token{
		{
			{Type: If, Val: "if"},
			{Type: Lparen, Val: "("},
			{Type: True, Val: "true"},
			{Type: Rparen, Val: ")"},
			{Type: Lbrace, Val: "{"},
			{Type: IntType, Val: "int"},
			{Type: Ident, Val: "a"},
			{Type: Assign, Val: "="},
			{Type: Int, Val: "0"},
			{Type: Eol, Val: "\n"},
			{Type: Rbrace, Val: "}"},
			{Type: Eol, Val: "\n"},
		},
		{
			{Type: If, Val: "if"},
			{Type: Lparen, Val: "("},
			{Type: Ident, Val: "a"},
			{Type: EQ, Val: "=="},
			{Type: Int, Val: "5"},
			{Type: Rparen, Val: ")"},
			{Type: Lbrace, Val: "{"},
			{Type: IntType, Val: "int"},
			{Type: Ident, Val: "a"},
			{Type: Assign, Val: "="},
			{Type: Int, Val: "1"},
			{Type: Eol, Val: "\n"},
			{Type: Rbrace, Val: "}"},
			{Type: Else, Val: "else"},
			{Type: Lbrace, Val: "{"},
			{Type: StringType, Val: "string"},
			{Type: Ident, Val: "b"},
			{Type: Assign, Val: "="},
			{Type: String, Val: "example"},
			{Type: Eol, Val: "\n"},
			{Type: Rbrace, Val: "}"},
			{Type: Eol, Val: "\n"},
		},
		{
			{Type: IntType, Val: "int"},
			{Type: Ident, Val: "a"},
			{Type: Assign, Val: "="},
			{Type: Int, Val: "5"},
		},
		{
			{Type: If, Val: "if"},
			{Type: Lparen, Val: "("},
			{Type: True, Val: "true"},
			{Type: Rparen, Val: ")"},
			{Type: IntType, Val: "int"},
			{Type: Ident, Val: "a"},
			{Type: Assign, Val: "="},
			{Type: Int, Val: "0"},
			{Type: Eol, Val: "\n"},
			{Type: Rbrace, Val: "}"},
			{Type: Eol, Val: "\n"},
		},
		{
			{Type: If, Val: "if"},
			{Type: Lparen, Val: "("},
			{Type: True, Val: "true"},
			{Type: Rbrace, Val: "}"},
			{Type: Lbrace, Val: "{"},
			{Type: IntType, Val: "int"},
			{Type: Ident, Val: "a"},
			{Type: Assign, Val: "="},
			{Type: Int, Val: "0"},
			{Type: Eol, Val: "\n"},
			{Type: Rbrace, Val: "}"},
			{Type: Eol, Val: "\n"},
		},
	}

	tests := []string{
		"if ( true ) { int [IDENT, a] = 0 }",
		`if ( (a == 5) ) { int [IDENT, a] = 1 } else { string [IDENT, b] = "example" }`,
		"INT_TYPE, is not a If",
		"INT_TYPE, is not a Left brace",
		"RBRACE, is not a Right paren",
	}

	for i, test := range tests {
		mockBuf := mockTokenBuffer{bufs[i], 0}
		stmt, err := parseIfStatement(&mockBuf)

		if err != nil && err.Error() != test {
			t.Fatalf("test[%d] - TestParseIfStatement() wrong error. expected=%s got=%s",
				i, test, err.Error())
		}

		if stmt != nil && stmt.String() != test {
			t.Fatalf("test[%d] - TestParseIfStatement() wrong result. expected=%s, got=%s",
				i, test, stmt.String())
		}

	}
}

func TestParseBlockStatement(t *testing.T) {
	prefixParseFnMap[True] = parseBooleanLiteral
	prefixParseFnMap[False] = parseBooleanLiteral
	prefixParseFnMap[Int] = parseIntegerLiteral
	prefixParseFnMap[Ident] = parseIdentifier
	prefixParseFnMap[String] = parseStringLiteral
	infixParseFnMap[EQ] = parseInfixExpression
	bufs := [][]Token{
		{
			{Type: IntType, Val: "int"},
			{Type: Ident, Val: "a"},
			{Type: Assign, Val: "="},
			{Type: Int, Val: "0"},
			{Type: Eol, Val: "\n"},
			{Type: Rbrace, Val: "}"},
		},
		{
			{Type: IntType, Val: "int"},
			{Type: Ident, Val: "a"},
			{Type: Assign, Val: "="},
			{Type: Int, Val: "0"},
			{Type: Eol, Val: "\n"},
			{Type: StringType, Val: "string"},
			{Type: Ident, Val: "b"},
			{Type: Assign, Val: "="},
			{Type: String, Val: "abc"},
			{Type: Eol, Val: "\n"},
			{Type: Rbrace, Val: "}"},
		},
		{
			{Type: IntType, Val: "int"},
			{Type: Ident, Val: "a"},
			{Type: Assign, Val: "="},
			{Type: Int, Val: "0"},
			{Type: Eol, Val: "\n"},
			{Type: StringType, Val: "string"},
			{Type: Ident, Val: "b"},
			{Type: Assign, Val: "="},
			{Type: String, Val: "abc"},
			{Type: Eol, Val: "\n"},
			{Type: BoolType, Val: "bool"},
			{Type: Ident, Val: "c"},
			{Type: Assign, Val: "="},
			{Type: True, Val: "true"},
			{Type: Eol, Val: "\n"},
			{Type: Rbrace, Val: "}"},
		},
	}
	tests := []string{
		"int [IDENT, a] = 0",
		`int [IDENT, a] = 0/string [IDENT, b] = "abc"`,
		`int [IDENT, a] = 0/string [IDENT, b] = "abc"/bool [IDENT, c] = true`,
	}

	for i, test := range tests {
		mockBuf := mockTokenBuffer{bufs[i], 0}
		exp, err := parseBlockStatement(&mockBuf)

		if err != nil && err.Error() != test {
			t.Fatalf("test[%d] - TestParseBlockStatement() wrong error. expected=%s, got=%s",
				i, test, err.Error())
		}

		if exp != nil && exp.String() != test {
			t.Fatalf("test[%d] - TestParseBlockStatement() wrong result. expected=%s, got=%s",
				i, test, exp.String())
		}
	}
}

func TestParseStatement(t *testing.T) {
	initParseFnMap()

	tests := []struct {
		tokens       []Token
		expectedErr  error
		expectedStmt string
	}{
		// tests for IntType
		{
			tokens: []Token{
				{Type: IntType, Val: "int"},
				{Type: Ident, Val: "a"},
				{Type: Assign, Val: "="},
				{Type: Int, Val: "1"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: "int [IDENT, a] = 1",
		},
		{
			tokens: []Token{
				{Type: IntType, Val: "int"},
				{Type: Ident, Val: "a"},
				{Type: Assign, Val: "="},
				{Type: Int, Val: "1"},
				{Type: Plus, Val: "+"},
				{Type: Int, Val: "2"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: "int [IDENT, a] = (1 + 2)",
		},
		{
			tokens: []Token{
				{Type: IntType, Val: "int"},
				{Type: Ident, Val: "a"},
				{Type: Assign, Val: "="},
				{Type: Int, Val: "1"},
				{Type: Plus, Val: "+"},
				{Type: Int, Val: "2"},
				{Type: Asterisk, Val: "*"},
				{Type: Int, Val: "3"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: "int [IDENT, a] = (1 + (2 * 3))",
		},
		{
			tokens: []Token{
				{Type: IntType, Val: "int"},
				{Type: Ident, Val: "a"},
				{Type: Assign, Val: "="},
				{Type: String, Val: "1"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `int [IDENT, a] = "1"`,
		},

		// tests for StringType
		{
			tokens: []Token{
				{Type: StringType, Val: "string"},
				{Type: Ident, Val: "abb"},
				{Type: Assign, Val: "="},
				{Type: String, Val: "do not merge, rebase!"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `string [IDENT, abb] = "do not merge, rebase!"`,
		},
		{
			tokens: []Token{
				{Type: StringType, Val: "string"},
				{Type: Ident, Val: "abb"},
				{Type: Assign, Val: "="},
				{Type: String, Val: "hello,*+"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `string [IDENT, abb] = "hello,*+"`,
		},
		{
			tokens: []Token{
				{Type: StringType, Val: "string"},
				{Type: Ident, Val: "abb"},
				{Type: Assign, Val: "="},
				{Type: Int, Val: "1"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `string [IDENT, abb] = 1`,
		},

		// tests for BoolType
		{
			tokens: []Token{
				{Type: BoolType, Val: "bool"},
				{Type: Ident, Val: "asdf"},
				{Type: Assign, Val: "="},
				{Type: True, Val: "true"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `bool [IDENT, asdf] = true`,
		},
		{
			tokens: []Token{
				{Type: BoolType, Val: "bool"},
				{Type: Ident, Val: "asdf"},
				{Type: Assign, Val: "="},
				{Type: False, Val: "false"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `bool [IDENT, asdf] = false`,
		},
		{
			tokens: []Token{
				{Type: BoolType, Val: "bool"},
				{Type: Ident, Val: "asdf"},
				{Type: Assign, Val: "="},
				{Type: Int, Val: "1"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `bool [IDENT, asdf] = 1`,
		},

		// tests for If statement
		{
			tokens: []Token{
				{Type: If, Val: "if"},
				{Type: Lparen, Val: "("},
				{Type: True, Val: "true"},
				{Type: Rparen, Val: ")"},
				{Type: Lbrace, Val: "{"},
				{Type: Rbrace, Val: "}"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `if ( true ) {  }`,
		},
		{
			tokens: []Token{
				{Type: If, Val: "if"},
				{Type: Lparen, Val: "("},
				{Type: Int, Val: "1"},
				{Type: Plus, Val: "+"},
				{Type: Int, Val: "2"},
				{Type: EQ, Val: "=="},
				{Type: Int, Val: "3"},
				{Type: Rparen, Val: ")"},
				{Type: Lbrace, Val: "{"},
				{Type: Rbrace, Val: "}"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `if ( ((1 + 2) == 3) ) {  }`,
		},
		{
			tokens: []Token{
				{Type: If, Val: "if"},
				{Type: Lparen, Val: "("},
				{Type: True, Val: "true"},
				{Type: Rparen, Val: ")"},
				{Type: Lbrace, Val: "{"},
				{Type: Int, Val: "1"},
				{Type: Plus, Val: "+"},
				{Type: Int, Val: "2"},
				{Type: EQ, Val: "=="},
				{Type: Int, Val: "3"},
				{Type: Rbrace, Val: "}"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  parseError{Int, "token is not identifier"},
			expectedStmt: ``,
		},
		{
			tokens: []Token{
				{Type: If, Val: "if"},
				{Type: Lparen, Val: "("},
				{Type: True, Val: "true"},
				{Type: Rparen, Val: ")"},
				{Type: Lbrace, Val: "{"},
				{Type: IntType, Val: "int"},
				{Type: Ident, Val: "a"},
				{Type: Assign, Val: "="},
				{Type: Int, Val: "2"},
				{Type: Eol, Val: "\n"},
				{Type: Rbrace, Val: "}"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `if ( true ) { int [IDENT, a] = 2 }`,
		},
		{
			tokens: []Token{
				{Type: If, Val: "if"},
				{Type: Lparen, Val: "("},
				{Type: True, Val: "true"},
				{Type: Rparen, Val: ")"},
				{Type: Lbrace, Val: "{"},
				{Type: IntType, Val: "int"},
				{Type: Ident, Val: "a"},
				{Type: Assign, Val: "="},
				{Type: Int, Val: "2"},
				{Type: Eol, Val: "\n"},
				{Type: Rbrace, Val: "}"},
				{Type: Else, Val: "else"},
				{Type: Lbrace, Val: "{"},
				{Type: Rbrace, Val: "}"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `if ( true ) { int [IDENT, a] = 2 } else {  }`,
		},
		{
			tokens: []Token{
				{Type: If, Val: "if"},
				{Type: Lparen, Val: "("},
				{Type: True, Val: "true"},
				{Type: Rparen, Val: ")"},
				{Type: Lbrace, Val: "{"},
				{Type: IntType, Val: "int"},
				{Type: Ident, Val: "a"},
				{Type: Assign, Val: "="},
				{Type: Int, Val: "2"},
				{Type: Eol, Val: "\n"},
				{Type: Rbrace, Val: "}"},
				{Type: Else, Val: "else"},
				{Type: Lbrace, Val: "{"},
				{Type: StringType, Val: "string"},
				{Type: Ident, Val: "b"},
				{Type: Assign, Val: "="},
				{Type: String, Val: "hello"},
				{Type: Eol, Val: "\n"},
				{Type: Rbrace, Val: "}"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `if ( true ) { int [IDENT, a] = 2 } else { string [IDENT, b] = "hello" }`,
		},

		// tests for Return statement
		{
			tokens: []Token{
				{Type: Return, Val: "return"},
				{Type: Ident, Val: "asdf"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `return asdf`,
		},
		{
			tokens: []Token{
				{Type: Return, Val: "return"},
				{Type: Lparen, Val: "("},
				{Type: Int, Val: "1"},
				{Type: Plus, Val: "+"},
				{Type: Int, Val: "2"},
				{Type: Asterisk, Val: "*"},
				{Type: Int, Val: "3"},
				{Type: Rparen, Val: ")"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `return (1 + (2 * 3))`,
		},
		{
			tokens: []Token{
				{Type: Return, Val: "return"},
				{Type: Lparen, Val: "("},
				{Type: Ident, Val: "add"},
				{Type: Lparen, Val: "("},
				{Type: Int, Val: "1"},
				{Type: Comma, Val: ","},
				{Type: Int, Val: "2"},
				{Type: Rparen, Val: ")"},
				{Type: Plus, Val: "+"},
				{Type: Int, Val: "2"},
				{Type: Asterisk, Val: "*"},
				{Type: Int, Val: "3"},
				{Type: Rparen, Val: ")"},
				{Type: Eol, Val: "\n"},
			},
			expectedErr:  nil,
			expectedStmt: `return (function add( 1, 2 ) + (2 * 3))`,
		},

		// tests for Default
		{
			tokens: []Token{
				{Type: Int, Val: "1"},
			},
			expectedErr:  parseError{Int, "token is not identifier"},
			expectedStmt: ``,
		},
	}

	for i, tt := range tests {
		buf := &mockTokenBuffer{tt.tokens, 0}
		stmt, err := parseStatement(buf)

		if err != nil && err != tt.expectedErr {
			t.Errorf(`test[%d] - parseStatement wrong error. expected="%v", got="%v"`,
				i, tt.expectedErr, err)
			continue
		}

		if err == nil && stmt.String() != tt.expectedStmt {
			t.Errorf(`test[%d] - parseStatement wrong result. expected="%s", got="%s"`,
				i, tt.expectedStmt, stmt.String())
		}
	}
}

func TestParseExpressionStatement(t *testing.T) {
	initParseFnMap()

	tests := []struct {
		tokens       []Token
		expectedStmt string
		expectedErr  string
	}{
		{
			// add()
			[]Token{
				{Type: Ident, Val: "add"},
				{Type: Lparen, Val: "("},
				{Type: Rparen, Val: ")"},
			},
			"function add(  )",
			"",
		},
		{
			// read(x int)
			[]Token{
				{Type: Ident, Val: "read"},
				{Type: Lparen, Val: "("},
				{Type: Ident, Val: "x"},
				{Type: Rparen, Val: ")"},
			},
			"function read( x )",
			"",
		},
		{
			// testFunction(a int, b string)
			[]Token{
				{Type: Ident, Val: "testFunction"},
				{Type: Lparen, Val: "("},
				{Type: Ident, Val: "a"},
				{Type: Comma, Val: ","},
				{Type: Ident, Val: "b"},
				{Type: Rparen, Val: ")"},
			},
			"function testFunction( a, b )",
			"",
		},
		{
			// testFunction(a int b string) <= error case
			[]Token{
				{Type: Ident, Val: "testFunction"},
				{Type: Lparen, Val: "("},
				{Type: Ident, Val: "a"},
				{Type: IntType, Val: "int"},
				{Type: Ident, Val: "b"},
				{Type: IntType, Val: "string"},
				{Type: Rparen, Val: ")"},
			},
			"",
			"INT_TYPE, expectNext() : expected [RPAREN], but got [INT_TYPE]",
		},
		{
			// 1() <= error case
			[]Token{
				{Type: Int, Val: "1"},
				{Type: Lparen, Val: "("},
				{Type: Rparen, Val: ")"},
			},
			"",
			"INT, token is not identifier",
		},
		{
			// add) <= error case
			[]Token{
				{Type: Ident, Val: "add"},
				{Type: Rparen, Val: ")"},
			},
			"",
			"RPAREN, expectNext() : expected [LPAREN], but got [RPAREN]",
		},
	}

	for i, test := range tests {
		stmt, err := parseExpressionStatement(&mockTokenBuffer{test.tokens, 0})
		if stmt != nil && stmt.String() != test.expectedStmt {
			t.Fatalf("test[%d] - TestParseFunctionStatement wrong answer.\n"+
				"expected= %s\n"+
				"got= %s", i, test.expectedStmt, stmt.String())
		}

		if err != nil && err.Error() != test.expectedErr {
			t.Fatalf("test[%d] - TestParseFunctionStatement wrong error.\n"+
				"expected= %s\n"+
				"got= %s", i, test.expectedErr, err.Error())
		}
	}
}
