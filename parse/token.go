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

type TokenType int

type Token struct {
	Type   TokenType
	Val    string
	Column Pos
	Line   int
}

// End of file
const Eof = -1

const (
	// ILLEGAL Token
	Illegal TokenType = iota

	// Identifiers + literals
	Ident  // add, foobar, x, y, ...
	Int    // 1343456
	String // "hello world"
	Bool   // true, false

	Assign   // =
	Plus     // +
	Minus    // -
	Bang     // !
	Asterisk // *
	Slash    // /
	Mod      // %

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
)
