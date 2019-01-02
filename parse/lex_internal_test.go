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
	"testing"
)

func TestState_cut(t *testing.T) {
	tests := []struct {
		inputState     state
		inputTokenType TokenType
		expectedToken  Token
	}{
		{
			inputState:     state{input: "hello", start: 0, end: 5, line: 0},
			inputTokenType: Ident,
			expectedToken:  Token{Line: 0, Val: "hello", Type: Ident, Column: 5},
		},
		{
			inputState:     state{input: "5", start: 0, end: 1, line: 0},
			inputTokenType: Int,
			expectedToken:  Token{Line: 0, Val: "5", Type: Int, Column: 1},
		},
	}

	for i, test := range tests {

		state := test.inputState
		token := state.cut(test.inputTokenType)
		if token != test.expectedToken {
			t.Fatalf("tests[%d] - Token wrong. expected=%q, got=%q",
				i, test.expectedToken, token)
		}
	}
}

func TestState_next(t *testing.T) {
	input := "hello world! \n hello world2"

	tests := []struct {
		expectedCh rune
	}{
		{rune('h')},
		{rune('e')},
		{rune('l')},
		{rune('l')},
		{rune('o')},
		{rune(' ')},
		{rune('w')},
		{rune('o')},
		{rune('r')},
		{rune('l')},
		{rune('d')},
		{rune('!')},
		{rune(' ')},
		{rune('\n')},
		{rune(' ')},
		{rune('h')},
		{rune('e')},
		{rune('l')},
		{rune('l')},
		{rune('o')},
		{rune(' ')},
		{rune('w')},
		{rune('o')},
		{rune('r')},
		{rune('l')},
		{rune('d')},
		{rune('2')},
		{rune(-1)},
	}

	s := state{
		input: input,
	}

	for i, test := range tests {
		ch := s.next()
		if ch != test.expectedCh {
			t.Errorf("tests[%d] - rune wrong. expected=%q, got=%q",
				i, test.expectedCh, ch)
		}
	}

	if s.line != 1 {
		t.Errorf("line is not 1. got=%d",
			s.line)
	}
}

func TestState_backup(t *testing.T) {
	input := "hello\n"

	tests := []struct {
		expectedPeekCh rune
	}{
		{rune('h')},
		{rune('e')},
		{rune('l')},
		{rune('l')},
		{rune('o')},
		{rune('\n')},
		{rune(-1)},
	}

	s := state{
		input: input,
	}

	s.next()
	for i, test := range tests {
		s.backup()
		ch := s.next()

		if ch != test.expectedPeekCh {
			t.Errorf("tests[%d] - rune wrong. expected=%q, got=%q",
				i, test.expectedPeekCh, ch)
		}

		if ch == rune('\n') {
			if s.line != 1 {
				t.Errorf("line is not 1. got=%d",
					s.line)
			}
		}

		s.next()
	}
}

func TestState_peek(t *testing.T) {
	input := "hello"

	tests := []struct {
		expectedPeekCh rune
	}{
		{rune('h')},
		{rune('e')},
		{rune('l')},
		{rune('l')},
		{rune('o')},
		{rune(eof)},
	}

	s := state{
		input: input,
	}

	for i, test := range tests {
		ch := s.peek()
		if ch != test.expectedPeekCh {
			t.Errorf("tests[%d] - rune wrong. expected=%q, got=%q",
				i, test.expectedPeekCh, ch)
		}
		s.next()
	}
}

type MockEmitter struct {
	emitFunc func(t Token)
}

func (m MockEmitter) emit(t Token) {
	m.emitFunc(t)
}

// TODO: LTE, GTE, spaceStateFn, numberStateFn, identifierStateFn case
func TestLex_defaultStateFn(t *testing.T) {

	//	Bang     // !
	//  Assign   // =
	//
	//	Plus     // +
	//	Minus    // -
	//	Asterisk // *
	//	Slash    // /
	//	Mod      // %
	//
	//	LT     // <
	//	GT     // >
	//	LTE    // <=
	//	GTE    // >=
	//	EQ     // ==
	//	NOT_EQ // !=
	//
	//	Comma // ,
	//
	//	Lparen // (
	//	Rparen // )
	//	Lbrace // {
	//	Rbrace // }
	tests := []struct {
		input             string
		expectedTokenType TokenType
	}{
		{"!", Bang},
		{"=", Assign},
		{"+", Plus},
		{"-", Minus},
		{"/", Slash},
		{"%", Mod},
		{"<", LT},
		{">", GT},
		{"==", EQ},
		{"!=", NOT_EQ},
		{",", Comma},
		{"(", Lparen},
		{")", Rparen},
		{"{", Lbrace},
		{"}", Rbrace},
		{"", Eof},
	}

	for i, test := range tests {

		s := &state{
			input: test.input,
		}

		e := MockEmitter{}
		e.emitFunc = func(tok Token) {
			if tok.Type != test.expectedTokenType {
				t.Errorf("tests[%d] - wrong token type. expected=%s, got=%s",
					i, TokenTypeMap[test.expectedTokenType], TokenTypeMap[tok.Type])
			}
		}

		defaultStateFn(s, e)
	}
}

func TestNumberStateFn(t *testing.T) {
	tests := []struct {
		input        string
		expectedType TokenType
		expectedVal  string
	}{
		{"1", Int, "1"},
		{"0", Int, "0"},
		{"123", Int, "123"},
		{"9990909", Int, "9990909"},
		{"+909", Int, "+909"},
		{"-909", Int, "-909"},
		{"-012", Int, "-012"}, //accept 0122
		{"_121", Illegal, "Invalid function call: numberStateFn"},
		{"+-121", Illegal, "Invalid function call: numberStateFn"},
		{"+_11", Illegal, "Invalid function call: numberStateFn"},
	}

	for i, test := range tests {
		s := &state{input: test.input}
		e := MockEmitter{}
		e.emitFunc = func(tok Token) {
			if tok.Type != test.expectedType {
				t.Errorf("tests[%d] - Wrong token type", i)
			}
			if tok.Val != test.expectedVal {
				t.Errorf("tests[%d] - rune wrong. expected=%s, got=%s", i, test.expectedVal, tok.Val)
			}
		}
		numberStateFn(s, e)
	}
}

func TestIdentifierStateFn(t *testing.T) {
	tests := []struct {
		input        string
		expectedType TokenType
		expectedVal  string
	}{
		{"abc", Ident, "abc"},
		{"_abc", Ident, "_abc"},
		{"a1", Ident, "a1"},
		{"a1123", Ident, "a1123"},
		{"_var123", Ident, "_var123"},
		{"_var_123", Ident, "_var_123"},
	}

	for i, test := range tests {
		s := &state{input: test.input}
		e := MockEmitter{}
		e.emitFunc = func(tok Token) {
			if tok.Type != Ident {
				t.Errorf("tests[%d] - Wrong token type", i)
			}
			if tok.Val != test.expectedVal {
				t.Errorf("tests[%d] - rune wrong. expected=%s, got=%s", i, test.expectedVal, tok.Val)
			}
		}
		identifierStateFn(s, e)
	}
}

func TestSpaceStateFn(t *testing.T) {
	tests := []struct {
		input string
	}{
		{" "},
		{"         "},
		{"\t"},
		{"\t\t\t\t\t"},
		{"\t \t    \t    \t        \t\t\t"},
	}

	for i, test := range tests {
		s := &state{input: test.input}
		e := MockEmitter{}

		spaceStateFn(s, e)
		if ch := s.next(); ch != rune(eof) {
			t.Errorf("tests[%d] - not fully absorbed", i)
		}
	}
}

func TestIsSpace(t *testing.T) {
	tests := []struct {
		input rune
	}{
		{' '},
		{'\t'},
	}

	for i, test := range tests {
		if !isSpace(test.input) {
			t.Errorf("tests[%d] - error occur", i)
		}
	}
}

func TestIsAlphaNumeric(t *testing.T) {
	tests := []struct {
		input rune
	}{
		{'a'},
		{'A'},
		{'_'},
		{'1'},
		{'2'},
		{'0'},
		{'z'},
		{'Z'},
	}

	errortests := []struct {
		input rune
	}{
		{'&'},
		{'*'},
		{'^'},
		{'%'},
		{'$'},
		{'#'},
		{'@'},
		{'!'},
	}

	for i, test := range tests {
		if !isAlphaNumeric(test.input) {
			t.Errorf("tests[%d] - %q is not an alphanumeric", i, test.input)
		}
	}
	for i, test := range errortests {
		if isAlphaNumeric(test.input) {
			t.Errorf("error tests[%d] - %q is alphanumeric", i, test.input)
		}
	}
}

func TestIsDigit(t *testing.T) {
	tests := []struct {
		input rune
	}{
		{'1'},
		{'2'},
		{'3'},
		{'4'},
		{'5'},
		{'6'},
		{'7'},
		{'8'},
		{'9'},
		{'0'},
	}

	for i, test := range tests {
		if !isDigit(test.input) {
			t.Errorf("tests[%d] - %q is not alphanumeric", i, test.input)
		}
	}
}
