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

// The parser holds a lexer.
type Parser struct {
	l         *Lexer
	errors    []string
	curToken  Token
	nextToken Token
}

func NewParser(l *Lexer) *Parser {
	return &Parser{
		l:      l,
		errors: []string{},
	}
}

// Parse function create an abstract syntax tree
func (p *Parser) Parse() ast.Program {
	return ast.Program{}
}
