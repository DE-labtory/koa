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

import "fmt"

type TokenType int

type Token struct {
	Type   TokenType
	Val    string
	Column Pos
	Line   int
}

func (t Token) String() string {
	if t.Type == Eol {
		return fmt.Sprintf("[EOL, End of line]")
	}

	if t.Type == Eof {
		return fmt.Sprintf("[EOF, End of file]")
	}
	return fmt.Sprintf("[%s, %s]", TokenTypeMap[t.Type], t.Val)
}

// End of file
const eof = -1

const (
	// ILLEGAL Token
	Illegal TokenType = iota

	// Identifiers + literals
	Ident    // add, foobar, x, y, ...
	Int      // 1343456
	String   // "hello world"
	Function // func
	Contract // contract

	IntType
	StringType
	BoolType

	Assign   // =
	Plus     // +
	Minus    // -
	Bang     // !
	Asterisk // *
	Slash    // /
	Mod      // %

	Plus_assign     // +=
	Minus_assign    // -=
	Asterisk_assign // *=
	Slash_assign    // /=
	Mod_assign      // %=

	Land // &&
	Lor  // ||
	Inc  // ++
	Dec  //--

	LT     // <
	GT     // >
	LTE    // <=
	GTE    // >=
	EQ     // ==
	NOT_EQ // !=

	Comma // ,

	Lparen // (
	Rparen // )
	Lbrace // {
	Rbrace // }

	True   // true
	False  // false
	If     // if
	Else   // else
	Return // return
	Eof    // end of file
	Eol    // end of line
)

// TokenTypeMap mapping TokenType with its
// string, this helps debugging
var TokenTypeMap = map[TokenType]string{
	Illegal: "ILLEGAL",

	Ident:    "IDENT",
	Int:      "INT",
	String:   "STRING",
	Function: "FUNCTION",
	Contract: "CONTRACT",

	IntType:    "INT_TYPE",
	StringType: "STRING_TYPE",
	BoolType:   "BOOL_TYPE",

	Assign:   "ASSIGN",
	Plus:     "PLUS",
	Minus:    "MINUS",
	Bang:     "BANG",
	Asterisk: "ASTERISK",
	Slash:    "SLASH",
	Mod:      "MOD",

	Plus_assign:     "PLUS_ASSIGN",
	Minus_assign:    "MINUS_ASSIGN",
	Asterisk_assign: "ASTERISK_ASSIGN",
	Slash_assign:    "SLASH_ASSIGN",
	Mod_assign:      "MOD_ASSIGN",

	Land: "LAND",
	Lor:  "LOR",
	Inc:  "INC",
	Dec:  "DEC",

	LT:     "LT",
	GT:     "GT",
	LTE:    "LTE",
	GTE:    "GTE",
	EQ:     "EQ",
	NOT_EQ: "NOT_EQ",

	Comma: "COMMA",

	Lparen: "LPAREN",
	Rparen: "RPAREN",
	Lbrace: "LBRACE",
	Rbrace: "RBRACE",

	True:   "TRUE",
	False:  "FALSE",
	If:     "IF",
	Else:   "ELSE",
	Return: "RETURN",

	Eof: "EOF",
	Eol: "EOL",
}

var keywords = map[string]TokenType{
	"contract": Contract,
	"func":     Function,
	"if":       If,
	"else":     Else,
	"int":      IntType,
	"string":   StringType,
	"return":   Return,
	"true":     True,
	"false":    False,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}
