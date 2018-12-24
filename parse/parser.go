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
	// Read retrieve token from buffer and we can choose how many
	// token we're going to read from this buffer
	Read(amount int) []Token

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

// The parser holds a lexer.
type Parser struct {
	errors []error
}

func NewParser() *Parser {
	return &Parser{
		errors: make([]error, 0),
	}
}

// Parse function create an abstract syntax tree
func (p *Parser) Parse(buf TokenBuffer) *ast.Program {
	prog := &ast.Program{}
	prog.Statements = []ast.Statement{}

	// TODO: need to break out for-loop when EOF
	for {
		stmt := p.parseStatement(buf)
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
	}

	return prog
}

func (p *Parser) parseStatement(buf TokenBuffer) ast.Statement {
	switch buf.Peek().Type {
	case LET:
		return p.parseLetStatement(buf)
	default:
		return p.parseExpressionStatement(buf)
	}
}

// TODO: implement me w/ test cases :-)
func (p *Parser) parseLetStatement(buf TokenBuffer) ast.Statement {
	return nil
}

// TODO: implement me w/ test cases :-)
func (p *Parser) parseExpressionStatement(buf TokenBuffer) *ast.Statement {
	return nil
}
