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

type TokenBuffer interface {
	Read() Token
	Peek() Token
}

// The parser holds a lexer.
type Parser struct {
	errors []string
}

func NewParser() *Parser {
	return &Parser{
		errors: []string{},
	}
}

// Parse function create an abstract syntax tree
func (p *Parser) Parse(buf TokenBuffer) *ast.Program {
	prog := &ast.Program{}
	prog.Statements = []ast.Statement{}

	// TODO: need to break out for-loop when EOF
	for {
		stmt := parseStatement(buf)
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
	}

	return prog
}

func parseStatement(buf TokenBuffer) ast.Statement {
	switch buf.Peek().Type {
	case LET:
		return parseLetStatement(buf)
	default:
		return nil
	}
}

func parseLetStatement(buf TokenBuffer) ast.Statement {
	return nil
}
