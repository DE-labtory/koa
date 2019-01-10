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

package parse_test

import (
	"testing"

	"github.com/DE-labtory/koa/parse"
)

func TestLexer_NextToken(t *testing.T) {
	input := `
	3 / 10
	int a = 5
	int b = 315 + (5 * 7) / 3 - 10
	<= >= == != = { } , 
	`

	tests := []struct {
		expectedType  parse.TokenType
		expectedValue string
	}{
		{parse.Eol, "\n"},
		{parse.Int, "3"},
		{parse.Slash, "/"},
		{parse.Int, "10"},
		{parse.Eol, "\n"},

		{parse.IntType, "int"},
		{parse.Ident, "a"},
		{parse.Assign, "="},
		{parse.Int, "5"},
		{parse.Eol, "\n"},

		{parse.IntType, "int"},
		{parse.Ident, "b"},
		{parse.Assign, "="},
		{parse.Int, "315"},
		{parse.Plus, "+"},
		{parse.Lparen, "("},
		{parse.Int, "5"},
		{parse.Asterisk, "*"},
		{parse.Int, "7"},
		{parse.Rparen, ")"},
		{parse.Slash, "/"},
		{parse.Int, "3"},
		{parse.Minus, "-"},
		{parse.Int, "10"},
		{parse.Eol, "\n"},

		{parse.LTE, "<="},
		{parse.GTE, ">="},
		{parse.EQ, "=="},
		{parse.NOT_EQ, "!="},
		{parse.Assign, "="},
		{parse.Lbrace, "{"},
		{parse.Rbrace, "}"},
		{parse.Comma, ","},
		{parse.Eol, "\n"},
	}

	l := parse.NewLexer(input)
	for i, test := range tests {
		token := l.NextToken()

		if token.Type != test.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, parse.TokenTypeMap[test.expectedType], parse.TokenTypeMap[token.Type])
		}

		if token.Val != test.expectedValue {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, test.expectedValue, token.Val)
		}
	}
}
