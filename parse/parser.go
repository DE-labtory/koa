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
	"github.com/DE-labtory/koa/ast"
)

// TokenBuffer provide tokenized token, we can read from this buffer
// or just peek the token
type TokenBuffer interface {
	// Read retrieve token from buffer
	Read() Token

	// Peek just take token from buffer but not change the
	// buffer states
	Peek() Token
}

// TODO: implement me w/ test cases :-)
func peekTokenTypeIs(buf TokenBuffer, t TokenType) bool {
	return buf.Peek().Type == t
}

// nextTokenIs helps to check whether peek token type
// if peek token type is what we want return it
// TODO: implement me w/ test cases :-)
func nextTokenIs(buf TokenBuffer, t TokenType) (bool, Token, error) {
	return false, Token{}, nil
}

// parseError contains error which happened during
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

	for buf.Peek().Type != Eof {
		stmt, e := parseStatement(buf)
		if len(errs) != 0 {
			errs = append(errs, e...)
			break
		}

		prog.Statements = append(prog.Statements, stmt)
	}

	return prog, errs
}

func parseStatement(buf TokenBuffer) (ast.Statement, []error) {
	switch buf.Peek().Type {
	default:
		return parseExpressionStatement(buf)
	}
}

// TODO: implement me w/ test cases :-)
func parseExpressionStatement(buf TokenBuffer) (*ast.Statement, []error) {
	return nil, nil
}

// TODO: implement me w/ test cases :-)
func parseExpression(buf TokenBuffer, pre int) (ast.Expression, []error) {
	return nil, nil
}
