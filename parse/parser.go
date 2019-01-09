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
	"fmt"
	"strconv"

	"github.com/DE-labtory/koa/ast"
)

// OperatorTypeMap maps TokenType with OperatorType by doing this
// we can remove dependency for token's string value
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
	nextType := buf.Peek(NEXT).Type
	if nextType != t {
		return false, parseError{
			nextType,
			fmt.Sprintf("expectNext() : expected [%s], but got [%s]",
				TokenTypeMap[t], TokenTypeMap[nextType]),
		}
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
	prefixParseFn func(TokenBuffer) (ast.Expression, error)
	infixParseFn  func(TokenBuffer, ast.Expression) (ast.Expression, error)
)

var prefixParseFnMap = map[TokenType]prefixParseFn{}
var infixParseFnMap = map[TokenType]infixParseFn{}

// Parse function create an abstract syntax tree
func Parse(buf TokenBuffer) (*ast.Program, error) {
	initParseFnMap()

	prog := &ast.Program{}
	prog.Statements = []ast.Statement{}

	for buf.Peek(CURRENT).Type != Eof {
		stmt, err := parseStatement(buf)
		if err != nil {
			return nil, err
		}

		prog.Statements = append(prog.Statements, stmt)
	}

	return prog, nil
}

// TODO: Currently there's no parsing function for statement
// TODO: so when complete parsing statement function, add that function
// TODO: to the map
func initParseFnMap() {
	prefixParseFnMap[Ident] = parseIdentifier
	prefixParseFnMap[Int] = parseIntegerLiteral
	prefixParseFnMap[Bang] = parsePrefixExpression
	prefixParseFnMap[Minus] = parsePrefixExpression
	prefixParseFnMap[True] = parseBooleanLiteral
	prefixParseFnMap[False] = parseBooleanLiteral
	prefixParseFnMap[Lparen] = parseGroupedExpression

	infixParseFnMap[Plus] = parseInfixExpression
	infixParseFnMap[Minus] = parseInfixExpression
	infixParseFnMap[Asterisk] = parseInfixExpression
	infixParseFnMap[Slash] = parseInfixExpression
	infixParseFnMap[Mod] = parseInfixExpression
	infixParseFnMap[EQ] = parseInfixExpression
	infixParseFnMap[NOT_EQ] = parseInfixExpression
	infixParseFnMap[LT] = parseInfixExpression
	infixParseFnMap[GT] = parseInfixExpression
	infixParseFnMap[LTE] = parseInfixExpression
	infixParseFnMap[GTE] = parseInfixExpression
}

// TODO: implement me w/ test cases :-)
func parseStatement(buf TokenBuffer) (ast.Statement, error) {
	switch buf.Peek(CURRENT).Type {
	case Return:
		return parseReturnStatement(buf)
	default:
		return nil, nil
	}
}

// ParseExpression parse expression in two ways, first
// by considering expression as prefix, next as infix
func parseExpression(buf TokenBuffer, pre precedence) (ast.Expression, error) {
	exp, err := makePrefixExpression(buf)
	if err != nil {
		return exp, err
	}

	exp, err = makeInfixExpression(buf, exp, pre)
	if err != nil {
		return exp, err
	}

	return exp, nil
}

// ParseExpAsPrefix retrieves prefix parse function from
// map, then parse expression with that function if exist.
func makePrefixExpression(buf TokenBuffer) (ast.Expression, error) {
	curTok := buf.Peek(CURRENT)

	fn := prefixParseFnMap[curTok.Type]

	if fn == nil {
		return nil, parseError{
			curTok.Type,
			"prefix parse function not defined",
		}
	}
	exp, err := fn(buf)
	if err != nil {
		return nil, err
	}

	return exp, nil
}

// MakeInfixExpression retrieves infix parse function from map
// then parse expression with that function if exist.
func makeInfixExpression(buf TokenBuffer, exp ast.Expression, pre precedence) (ast.Expression, error) {
	var err error
	expression := exp

	for pre < curPrecedence(buf) {
		fn := infixParseFnMap[buf.Peek(CURRENT).Type]
		if fn == nil {
			return nil, parseError{
				buf.Peek(CURRENT).Type,
				"infix parse function not defined",
			}
		}

		expression, err = fn(buf, expression)
		if err != nil {
			return nil, err
		}
	}
	return expression, nil
}

func parseInfixExpression(buf TokenBuffer, left ast.Expression) (ast.Expression, error) {
	var err error
	curTok := buf.Read()

	expression := &ast.InfixExpression{
		Left: left,
		Operator: ast.Operator{
			Type: operatorTypeMap[curTok.Type],
			Val:  ast.OperatorVal(curTok.Val),
		},
	}
	precedence := precedenceMap[curTok.Type]

	expression.Right, err = parseExpression(buf, precedence)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func parsePrefixExpression(buf TokenBuffer) (ast.Expression, error) {
	var err error
	tok := buf.Read()

	exp := &ast.PrefixExpression{
		Operator: ast.Operator{
			Type: operatorTypeMap[tok.Type],
			Val:  ast.OperatorVal(tok.Val),
		},
	}

	exp.Right, err = parseExpression(buf, PREFIX)
	if err != nil {
		return nil, err
	}

	return exp, nil
}

func parseIdentifier(buf TokenBuffer) (ast.Expression, error) {
	token := buf.Read()

	if token.Type != Ident {
		return nil, parseError{
			token.Type,
			fmt.Sprintf("parseIdentifier() - %s is not a identifier", token.Val),
		}
	}
	return &ast.Identifier{Value: token.Val}, nil
}

func parseIntegerLiteral(buf TokenBuffer) (ast.Expression, error) {
	token := buf.Read()

	if token.Type != Int {
		return nil, parseError{
			token.Type,
			fmt.Sprintf("parseIntegerLiteral() error - %s is not integer", token.Val),
		}
	}

	value, err := strconv.ParseInt(token.Val, 0, 64)
	if err != nil {
		return nil, err
	}

	lit := &ast.IntegerLiteral{Value: value}
	return lit, nil
}

func parseBooleanLiteral(buf TokenBuffer) (ast.Expression, error) {
	token := buf.Read()

	if token.Type != True && token.Type != False {
		return nil, parseError{
			token.Type,
			fmt.Sprintf("parseBooleanLiteral() error - %s is not bool", token.Val),
		}
	}

	value, err := strconv.ParseBool(token.Val)
	if err != nil {
		return nil, err
	}

	lit := &ast.BooleanLiteral{Value: value}
	return lit, nil
}

func parseStringLiteral(buf TokenBuffer) (ast.Expression, error) {
	token := buf.Read()

	if token.Type != String {
		return nil, parseError{
			token.Type,
			fmt.Sprintf("parseStringLiteral() error - %s is not string", token.Val),
		}
	}

	return &ast.StringLiteral{Value: token.Val}, nil
}

func parseReturnStatement(buf TokenBuffer) (ast.Statement, error) {
	token := buf.Read()
	if token.Type != Return {
		return nil, parseError{
			token.Type,
			fmt.Sprintf("parseReturnStatement() error - Statement must be started with return"),
		}
	}

	stmt := &ast.ReturnStatement{}
	for !curTokenIs(buf, Eol) {
		exp, err := parseExpression(buf, LOWEST)
		if err != nil {
			return nil, err
		}
		stmt.ReturnValue = exp
	}

	return stmt, nil
}

func parseGroupedExpression(buf TokenBuffer) (ast.Expression, error) {
	buf.Read()
	exp, err := parseExpression(buf, LOWEST)
	if err != nil {
		return nil, err
	}

	tok := buf.Read()
	if tok.Type != Rparen {
		err := parseError{
			tok.Type,
			fmt.Sprintf("is not Right paren"),
		}
		return nil, err
	}
	return exp, nil
}

func parseAssignStatement(buf TokenBuffer) (*ast.AssignStatement, error) {
	stmt := &ast.AssignStatement{}

	tok := buf.Read()
	if !isDataStructure(tok) {
		return nil, parseError{
			tok.Type,
			"token is not data structure type",
		}
	}

	stmt.Type = ast.DataStructure{
		Type: ast.DataStructureType(tok.Type),
		Val:  ast.DataStructureVal(tok.Val),
	}

	tok = buf.Read()
	if tok.Type != Ident {
		return nil, parseError{
			tok.Type,
			"token is not identifier",
		}
	}

	stmt.Variable = ast.Identifier{
		Value: tok.String(),
	}

	if buf.Read().Type != Assign {
		return nil, parseError{
			tok.Type,
			"token is not assign",
		}
	}

	exp, err := parseExpression(buf, LOWEST)
	if err != nil {
		return nil, err
	}

	stmt.Value = exp

	return stmt, nil
}

func parseCallExpression(buf TokenBuffer, fn ast.Expression) (ast.Expression, error) {
	exp := &ast.CallExpression{Function: fn}

	var err error
	exp.Arguments, err = parseCallArguments(buf)
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func parseCallArguments(buf TokenBuffer) ([]ast.Expression, error) {
	args := []ast.Expression{}
	buf.Read()
	if curTokenIs(buf, Rparen) {
		return args, nil
	}

	exp, err := parseExpression(buf, LOWEST)
	if err != nil {
		return nil, err
	}
	args = append(args, exp)

	for curTokenIs(buf, Comma) {
		buf.Read()
		exp, err := parseExpression(buf, LOWEST)
		if err != nil {
			return nil, err
		}
		args = append(args, exp)
	}
	return args, nil
}

func isDataStructure(tok Token) bool {
	return tok.Type == StringType ||
		tok.Type == IntType ||
		tok.Type == BoolType
}

func parseIfStatement(buf TokenBuffer) (*ast.IfStatement, []error) {
	errs := make([]error, 0)
	var err []error
	tok := buf.Read()
	if tok.Type != If {
		errs = append(errs, errors.New("parseIfStatement() error - No IF statement"))
		return nil, errs
	}

	tok = buf.Read()
	if tok.Type != Lparen {
		errs = append(errs, errors.New("parseIfStatement() error - Condition must be stared with left paren"))
		return nil, errs
	}

	expression := &ast.IfStatement{}
	expression.Condition, err = parseExpression(buf, LOWEST)
	if len(err) > 0 {
		errs = append(errs, errors.New("parseIfStatement() error - Condition parsing error"))
		return nil, errs
	}

	tok = buf.Read()
	if tok.Type != Rparen {
		errs = append(errs, errors.New("parseIfStatement() error - Condition must be ended with right paren"))
		return nil, errs
	}

	tok = buf.Read()
	if tok.Type != Lbrace {
		errs = append(errs, errors.New("parseIfStatement() error - Block must be started with left brace"))
		return nil, errs
	}

	expression.Consequence, err = parseBlockStatement(buf)
	if len(err) > 0 {
		errs = append(errs, errors.New("parseIfStatement() error - block statement parsing error"))
		return nil, errs
	}

	return expression, nil
}

func parseBlockStatement(buf TokenBuffer) (*ast.BlockStatement, []error) {
	errs := make([]error, 0)
	block := &ast.BlockStatement{
		Statements: []ast.Statement{},
	}

	for !curTokenIs(buf, Rbrace) && !curTokenIs(buf, Eof) {
		stmt, err := parseStatement(buf)
		if len(err) > 0 {
			errs = append(errs, errors.New("parseBlockStatement() error - statement parsing error"))
			return nil, err
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		buf.Read()
	}

	return block, nil
}
