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

import "github.com/DE-labtory/koa/ast"

type precedence int

const (
	_ precedence = iota
)

var precedenceMap = map[TokenType]precedence{}

// PeekNumber restrict peek count from the TokenBuffer
type peekNumber int

const (
	CURRENT peekNumber = iota
	NEXT
)

// TODO: implement me w/ test cases :-)
func (p peekNumber) isValid() bool {
	return p == CURRENT || p == NEXT
}

// TokenBuffer provide tokenized token, we can read from this buffer
// or just peek the token
type TokenBuffer interface {
	// Read retrieve token from buffer. and change the
	// buffer states
	Read() Token

	// Peek take token as many as n from buffer but not change the
	// buffer states
	Peek(n peekNumber) Token
}

// TODO: implement me w/ test cases :-)
func curTokenIs(buf TokenBuffer, t TokenType) bool {
	return false
}

// TODO: implement me w/ test cases :-)
func nextTokenIs(buf TokenBuffer, t TokenType) bool {
	return false
}

// ExpectNext helps to check whether next token is
// type of token we want, and if true then return it
// otherwise return with error
// TODO: implement me w/ test cases :-)
func expectNext(buf TokenBuffer, t TokenType) (bool, Token, error) {
	return false, Token{}, nil
}

// TODO: implement me w/ test cases :-)
func curPrecedence(buf TokenBuffer) precedence {
	return 0
}

// TODO: implement me w/ test cases :-)
func nextPrecedence(buf TokenBuffer) precedence {
	return 0
}

// ParseError contains error which happened during
// parsing tokens
type parseError struct{}

func (e *parseError) Error() string { return "" }

type (
	prefixParseFn func(TokenBuffer) (ast.Expression, []error)
	infixParseFn  func(TokenBuffer, ast.Expression) (ast.Expression, []error)
)

var prefixParseFnMap = map[TokenType]prefixParseFn{}
var infixParseFnMap = map[TokenType]infixParseFn{}

// Parse function create an abstract syntax tree
func Parse(buf TokenBuffer) (*ast.Program, []error) {
	errs := []error{}
	prog := &ast.Program{}
	prog.Statements = []ast.Statement{}

	for buf.Peek(CURRENT).Type != Eof {
		stmt, e := parseStatement(buf)
		if len(errs) != 0 {
			errs = append(errs, e...)
			break
		}

		prog.Statements = append(prog.Statements, stmt)
	}

	return prog, errs
}

// TODO: implement me w/ test cases :-)
func parseStatement(buf TokenBuffer) (ast.Statement, []error) {
	switch buf.Peek(CURRENT).Type {
	default:
		return nil, nil
	}
}

// TODO: implement me w/ test cases :-)
func parseExpression(buf TokenBuffer, pre int) (ast.Expression, []error) {
	return nil, nil
}
