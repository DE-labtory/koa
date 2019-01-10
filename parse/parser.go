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
	"fmt"
	"strconv"

	"github.com/DE-labtory/koa/ast"
)

type precedence int

const (
	_ precedence = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // function(X)
)

var precedenceMap = map[TokenType]precedence{
	Assign:   LOWEST,
	Plus:     SUM,
	Minus:    SUM,
	Bang:     PREFIX,
	Asterisk: PRODUCT,
	Slash:    PRODUCT,
	Mod:      PRODUCT,

	LT:     LESSGREATER,
	GT:     LESSGREATER,
	LTE:    LESSGREATER,
	GTE:    LESSGREATER,
	EQ:     EQUALS,
	NOT_EQ: EQUALS,

	Comma: LOWEST,

	Lparen: LOWEST,
	Rparen: LOWEST,
	Lbrace: LOWEST,
	Rbrace: LOWEST,

	Eol: LOWEST,
}

// Operator map
var operatorTypeMap = map[TokenType]ast.OperatorType{
	Plus:     ast.Plus,
	Minus:    ast.Minus,
	Bang:     ast.Bang,
	Asterisk: ast.Asterisk,
	Slash:    ast.Slash,
	Mod:      ast.Mod,
	LT:       ast.LT,
	GT:       ast.GT,
	LTE:      ast.LTE,
	GTE:      ast.GTE,
	EQ:       ast.EQ,
	NOT_EQ:   ast.NOT_EQ,
}

// PeekNumber restrict peek count from the TokenBuffer
type peekNumber int

const (
	CURRENT peekNumber = iota
	NEXT
)

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

func curTokenIs(buf TokenBuffer, t TokenType) bool {
	return buf.Peek(CURRENT).Type == t
}

func nextTokenIs(buf TokenBuffer, t TokenType) bool {
	return buf.Peek(NEXT).Type == t
}

// ExpectNext helps to check whether next token is
// type of token we want, and if true then return it
// otherwise return with error
func expectNext(buf TokenBuffer, t TokenType) (bool, error) {
	if buf.Peek(NEXT).Type != t {
		return false, errors.New("expectNext() : expecting token and next token are different")
	}
	return true, nil
}

func curPrecedence(buf TokenBuffer) precedence {
	if p, ok := precedenceMap[buf.Peek(CURRENT).Type]; ok {
		return p
	}
	return LOWEST
}
func nextPrecedence(buf TokenBuffer) precedence {
	if p, ok := precedenceMap[buf.Peek(NEXT).Type]; ok {
		return p
	}
	return LOWEST
}

// ParseError contains error which happened during
// parsing tokens
type parseError struct {
	tokenType TokenType
	reason    string
}

func (e parseError) Error() string {
	return fmt.Sprintf("%s, %s", TokenTypeMap[e.tokenType], e.reason)
}

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
	case Return:
		return parseReturnStatement(buf)
	default:
		return nil, nil
	}
}

// ParseExpression parse expression in two ways, first
// by considering expression as prefix, next as infix
func parseExpression(buf TokenBuffer, pre precedence) (ast.Expression, []error) {
	errs := make([]error, 0)
	exp, es := makePrefixExpression(buf)
	if len(es) != 0 {
		errs = append(errs, es...)
		return exp, errs
	}

	exp, es = makeInfixExpression(buf, exp, pre)
	if len(es) != 0 {
		errs = append(errs, es...)
		return exp, errs
	}

	return exp, errs
}

// ParseExpAsPrefix retrieves prefix parse function from
// map, then parse expression with that function if exist.
func makePrefixExpression(buf TokenBuffer) (ast.Expression, []error) {
	errs := make([]error, 0)
	curTok := buf.Peek(CURRENT)

	fn := prefixParseFnMap[curTok.Type]
	if fn == nil {
		e := parseError{
			curTok.Type,
			"prefix parse function not defined",
		}
		return nil, append(errs, e)
	}

	exp, e := fn(buf)
	if len(e) != 0 {
		return nil, append(errs, e...)
	}

	return exp, errs
}

// MakeInfixExpression retrieves infix parse function from map
// then parse expression with that function if exist.
func makeInfixExpression(buf TokenBuffer, exp ast.Expression, pre precedence) (ast.Expression, []error) {
	errs := make([]error, 0)
	expression := exp

	for pre < curPrecedence(buf) {
		fn := infixParseFnMap[buf.Peek(CURRENT).Type]
		if fn == nil {
			errs = append(errs, parseError{buf.Peek(CURRENT).Type, "infix parse function not defined"})
			return nil, errs
		}

		var err []error
		expression, err = fn(buf, expression)
		if len(err) > 0 {
			return nil, append(errs, err...)
		}
	}
	return expression, nil
}

func parseInfixExpression(buf TokenBuffer, left ast.Expression) (ast.Expression, []error) {
	curTok := buf.Read()
	expression := &ast.InfixExpression{
		Left: left,
		Operator: ast.Operator{
			Type: operatorTypeMap[curTok.Type],
			Val:  ast.OperatorVal(curTok.Val),
		},
	}
	precedence := precedenceMap[curTok.Type]

	var err []error
	expression.Right, err = parseExpression(buf, precedence)
	if len(err) > 0 {
		return nil, err
	}

	return expression, nil
}

func parseIdentifier(buf TokenBuffer) (ast.Expression, []error) {
	errs := make([]error, 0)
	token := buf.Read()

	if token.Type != Ident {
		errs = append(errs, errors.New("parseIdentifier() - "+token.Val+" is not a identifier"))
		return nil, errs
	}
	return &ast.Identifier{Value: token.Val}, nil
}

func parseIntegerLiteral(buf TokenBuffer) (ast.Expression, []error) {
	token := buf.Read()
	errs := make([]error, 0)

	if token.Type != Int {
		errs = append(errs, errors.New("parseIntegerLiteral() error - "+token.Val+" is not integer"))
		return nil, errs
	}

	value, err := strconv.ParseInt(token.Val, 0, 64)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	lit := &ast.IntegerLiteral{Value: value}
	return lit, nil
}

func parseBooleanLiteral(buf TokenBuffer) (ast.Expression, []error) {
	token := buf.Read()
	errs := make([]error, 0)

	if token.Type != Bool {
		errs = append(errs, errors.New("parseBooleanLiteral() error - "+token.Val+" is not bool"))
		return nil, errs
	}

	value, err := strconv.ParseBool(token.Val)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	lit := &ast.BooleanLiteral{Value: value}
	return lit, nil
}

func parseStringLiteral(buf TokenBuffer) (ast.Expression, []error) {
	token := buf.Read()
	errs := make([]error, 0)

	if token.Type != String {
		errs = append(errs, errors.New("parseStringLiteral() error - "+token.Val+" is not string"))
		return nil, errs
	}

	return &ast.StringLiteral{Value: token.Val}, nil
}

func parseReturnStatement(buf TokenBuffer) (ast.Statement, []error) {
	errs := make([]error, 0)
	token := buf.Read()
	if token.Type != Return {
		errs = append(errs, errors.New("parseReturnStatement() error - Statement must be started with return"))
		return nil, errs
	}

	stmt := &ast.ReturnStatement{}
	for !curTokenIs(buf, Eol) {
		exp, err := parseExpression(buf, LOWEST)
		if len(err) > 0 {
			errs = append(err, errors.New("parseReturnStatement() error - pareseExpression error"))
			return nil, errs
		}
		stmt.ReturnValue = exp
	}

	return stmt, nil
}

// TODO: implement me w/ test cases :-)
func parseAssignStatement(buf TokenBuffer) (*ast.AssignStatement, []error) {
	return nil, nil
}

// TODO: implement me w/ test cases :-) (This used for calling built-in function)
func parseCallExpression(buf TokenBuffer, fn ast.Expression) (ast.Expression, []error) {
	return nil, nil
}
