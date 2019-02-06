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

	"github.com/DE-labtory/koa/symbol"

	"github.com/DE-labtory/koa/ast"
)

// OperatorTypeMap maps TokenType with OperatorType. By doing this
// we can remove dependency for token's string value
var operatorMap = map[TokenType]ast.Operator{
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
	Land:     ast.LAND,
	Lor:      ast.LOR,
}

// datastructureMap maps TokenType with Datastructure. By doing this
// we can remove dependency for token's string value
var datastructureMap = map[TokenType]ast.DataStructure{
	IntType:    ast.IntType,
	StringType: ast.StringType,
	BoolType:   ast.BoolType,
	VoidType:   ast.VoidType,
}

// precedence determine which token is going to be grouped first when
// parsing expression with Pratt Parsing.
//
// For example, 1 + 2 * 3, * has higher precedence value then +. So 2 * 3
// is grouped first and result would be (1 + (2 * 3))
type precedence int

const (
	_ precedence = iota
	LOWEST
	LOR         // ||
	LAND        // &&
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

	Lparen: CALL,

	Eol:  LOWEST,
	Land: LAND,
	Lor:  LOR,
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
func expectNext(buf TokenBuffer, t TokenType) error {
	tok := buf.Peek(CURRENT)
	if tok.Type != t {
		return ExpectError{
			tok,
			t,
		}
	}
	buf.Read()
	return nil
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

func consumeSemi(buf TokenBuffer) {
	for curTokenIs(buf, Semicolon) {
		buf.Read()
	}
}

// Error contains error which happened during
// parsing tokens
type Error struct {
	source Token
	reason string
}

func (e Error) Error() string {
	return fmt.Sprintf("[line %d, column %d] [%s] %s",
		e.source.Line, e.source.Column, TokenTypeMap[e.source.Type], e.reason)
}

// ExpectError happens during parsing expectNext
type ExpectError struct {
	source   Token
	expected TokenType
}

func (e ExpectError) Error() string {
	return fmt.Sprintf("[line %d, column %d] expected [%s], but got [%s]",
		e.source.Line, e.source.Column, TokenTypeMap[e.expected], TokenTypeMap[e.source.Type])
}

// dupSymError occur when there is duplicated symbol
type DupSymError struct {
	source Token
}

func (e DupSymError) Error() string {
	return fmt.Sprintf("[line %d, column %d] symbol [%s] already exist",
		e.source.Line, e.source.Column, e.source.Val)
}

type (
	prefixParseFn func(TokenBuffer) (ast.Expression, error)
	infixParseFn  func(TokenBuffer, ast.Expression) (ast.Expression, error)
)

var prefixParseFnMap = map[TokenType]prefixParseFn{}
var infixParseFnMap = map[TokenType]infixParseFn{}

// scope keeps symbols that shows on tokens, every time scope meet symbol,
// trying to check whether symbol with same name already exist, if true
// then throw error, if not add that symbol to scope.
var scope *symbol.Scope

// updateScopeSymbol checks whether token value is exist in scope first,
// if exist, then throw error, if not make symbol with token value then add
// to scope
func updateScopeSymbol(token Token, symType symbol.SymbolType) error {
	if sym := scope.Get(token.Val); sym != nil {
		return DupSymError{token}
	}

	switch symType {
	case symbol.IntegerSymbol:
		// TODO: add symbol to scope
	case symbol.BooleanSymbol:
		// TODO: add symbol to scope
	case symbol.StringSymbol:
		// TODO: add symbol to scope
	case symbol.FunctionSymbol:
		scope.Set(token.Val, &symbol.Function{Name: token.Val})
	default:
		return Error{
			token,
			fmt.Sprintf("unexpected symbol type [%s]", symType),
		}
	}

	return nil
}

// enterScope creates new scope than converts it to existing scope
func enterScope() {
	innerScope := symbol.NewScope()
	innerScope.SetOuter(scope)
	scope = innerScope
}

// leaveScope converts current scope's outer to existing scope
func leaveScope() {
	outerScope := scope.GetOuter()
	scope = outerScope
}

// Parse creates an abstract syntax tree
func Parse(buf TokenBuffer) (*ast.Contract, error) {
	initParseFnMap()

	scope = symbol.NewScope()

	contract := &ast.Contract{}
	contract.Functions = []*ast.FunctionLiteral{}

	if err := parseContractStart(buf); err != nil {
		return nil, err
	}

	for buf.Peek(CURRENT).Type == Function {
		fn, err := parseFunctionLiteral(buf)
		if err != nil {
			return nil, err
		}

		contract.Functions = append(contract.Functions, fn)
	}

	if err := parseContractEnd(buf); err != nil {
		return nil, err
	}

	return contract, nil
}

// parseContractStart validates whether given token stream is
// starts with "contract" keyword with left-brace, otherwise throw error
func parseContractStart(buf TokenBuffer) error {
	if err := expectNext(buf, Contract); err != nil {
		return err
	}
	if err := expectNext(buf, Lbrace); err != nil {
		return err
	}
	return nil
}

// parseContractEnd validates whether contracts finish with
// right-brace, otherwise throw error
func parseContractEnd(buf TokenBuffer) error {
	if err := expectNext(buf, Rbrace); err != nil {
		return err
	}
	if err := expectNext(buf, Semicolon); err != nil {
		return err
	}
	return nil
}

// initParseFnMap initialize parsing function for each tokens.
// Each token can have at most two parsing function:
//
//   - infix-parsing function
//   - prefix-parsing function
//
func initParseFnMap() {
	prefixParseFnMap[Ident] = parseIdentifier
	prefixParseFnMap[Int] = parseIntegerLiteral
	prefixParseFnMap[String] = parseStringLiteral
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
	infixParseFnMap[Land] = parseInfixExpression
	infixParseFnMap[Lor] = parseInfixExpression
	infixParseFnMap[Lparen] = parseCallExpression
}

// parseStatement parse statement which don't produce value
func parseStatement(buf TokenBuffer) (ast.Statement, error) {
	switch tt := buf.Peek(CURRENT).Type; tt {
	case IntType:
		return parseAssignStatement(buf)
	case BoolType:
		return parseAssignStatement(buf)
	case StringType:
		return parseAssignStatement(buf)
	case If:
		return parseIfStatement(buf)
	case Return:
		return parseReturnStatement(buf)
	default:
		return parseExpressionStatement(buf)
	}
}

// ParseExpression parse expression in two ways, first
// by considering expression as prefix, next as infix
//
// Parsing expression is done in Pratt Parsing way. So each
// token has its own parsing function. And each token has its
// parsing precedence.
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
		return nil, Error{
			curTok,
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
	for !curTokenIs(buf, Semicolon) && pre < curPrecedence(buf) {
		token := buf.Peek(CURRENT)
		fn := infixParseFnMap[token.Type]
		if fn == nil {
			return nil, Error{
				token,
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

// parseInfixExprsesion parse expression when current token in TokenBuffer
// works as infix of expression.
//
// Infix parsing is based on a precedence of given token which is defined
// in precedenceMap
func parseInfixExpression(buf TokenBuffer, left ast.Expression) (ast.Expression, error) {
	var err error
	curTok := buf.Read()

	expression := &ast.InfixExpression{
		Left:     left,
		Operator: operatorMap[curTok.Type],
	}

	precedence := precedenceMap[curTok.Type]
	expression.Right, err = parseExpression(buf, precedence)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

// parsePrefixExprsesion parse expression when current token in TokenBuffer
// works as prefix of expression.
//
// Prefix parsing is based on a precedence of given token which is defined
// in precedenceMap.
func parsePrefixExpression(buf TokenBuffer) (ast.Expression, error) {
	token := buf.Read()
	op := operatorMap[token.Type]

	right, err := parseExpression(buf, PREFIX)
	if err != nil {
		return nil, err
	}

	switch op {
	case ast.Bang:
		switch right.(type) {
		case *ast.StringLiteral:
			return nil, Error{
				token,
				fmt.Sprintf("parsePrefixExpression() - Invalid prefix of %s", right.String()),
			}
		}
	case ast.Minus:
		switch right.(type) {
		case *ast.BooleanLiteral:
			return nil, Error{
				token,
				fmt.Sprintf("parsePrefixExpression() - Invalid prefix of %s", right.String()),
			}
		}
	}

	exp := &ast.PrefixExpression{
		Operator: op,
		Right:    right,
	}
	return exp, nil
}

// parseIdentifier parse identifier.
func parseIdentifier(buf TokenBuffer) (ast.Expression, error) {
	token := buf.Read()
	if token.Type != Ident {
		return nil, ExpectError{token, Ident}
	}

	// TODO: check whether variable name exist,
	// TODO: but do not update symbol table in this case
	// TODO: if not, return error

	return &ast.Identifier{Value: token.Val}, nil
}

// parseIntegerLiteral parse integer literal.
func parseIntegerLiteral(buf TokenBuffer) (ast.Expression, error) {
	token := buf.Read()
	if token.Type != Int {
		return nil, ExpectError{token, Int}
	}

	value, err := strconv.ParseInt(token.Val, 0, 64)
	if err != nil {
		return nil, err
	}

	lit := &ast.IntegerLiteral{Value: value}
	return lit, nil
}

// parseBooleanLiteral parse boolean literal.
func parseBooleanLiteral(buf TokenBuffer) (ast.Expression, error) {
	token := buf.Read()
	if token.Type != True && token.Type != False {
		return nil, ExpectError{token, BoolType}
	}

	val, err := strconv.ParseBool(token.Val)
	if err != nil {
		return nil, err
	}

	lit := &ast.BooleanLiteral{Value: val}
	return lit, nil
}

// parseStringLiteral parse string value which is
// going to be assigned to variable
func parseStringLiteral(buf TokenBuffer) (ast.Expression, error) {
	token := buf.Read()
	if token.Type != String {
		return nil, ExpectError{token, String}
	}

	return &ast.StringLiteral{Value: token.Val}, nil
}

// parseFunctionLiteral parse functional expression
// first parse name, and parse parameter, body
func parseFunctionLiteral(buf TokenBuffer) (*ast.FunctionLiteral, error) {
	lit := &ast.FunctionLiteral{}
	var err error

	if err = expectNext(buf, Function); err != nil {
		return nil, err
	}

	token := buf.Read()
	if token.Type != Ident {
		return nil, ExpectError{token, Ident}
	}

	if err := updateScopeSymbol(token, symbol.FunctionSymbol); err != nil {
		return nil, err
	}

	lit.Name = &ast.Identifier{Value: token.Val}

	if err = expectNext(buf, Lparen); err != nil {
		return nil, err
	}

	if lit.Parameters, err = parseFunctionParameterList(buf); err != nil {
		return nil, err
	}

	if lit.ReturnType, err = parseFunctionReturnType(buf); err != nil {
		return nil, err
	}

	if lit.Body, err = parseBlockStatement(buf); err != nil {
		return nil, err
	}

	consumeSemi(buf)

	return lit, nil
}

// parseFunctionReturnType parse function's return data structure type
func parseFunctionReturnType(buf TokenBuffer) (ast.DataStructure, error) {
	peekTok := buf.Peek(CURRENT)

	ds, ok := datastructureMap[peekTok.Type]
	if !ok && peekTok.Type != Lbrace {
		return 0, Error{
			peekTok,
			"invalid function return type",
		}
	}

	if !ok {
		ds = ast.VoidType
	} else {
		buf.Read()
	}

	return ds, nil
}

// parseFunctionParameters parse function's parameters which
// separated by comma
func parseFunctionParameterList(buf TokenBuffer) ([]*ast.ParameterLiteral, error) {
	identifiers := []*ast.ParameterLiteral{}
	if err := expectNext(buf, Rparen); err == nil {
		return identifiers, nil
	}

	ident, err := parseFunctionParameter(buf)
	if err != nil {
		return nil, err
	}
	identifiers = append(identifiers, ident)

	for curTokenIs(buf, Comma) {
		buf.Read()

		ident, err := parseFunctionParameter(buf)
		if err != nil {
			return nil, err
		}
		identifiers = append(identifiers, ident)
	}

	if err = expectNext(buf, Rparen); err != nil {
		return nil, err
	}

	return identifiers, nil
}

func parseFunctionParameter(buf TokenBuffer) (*ast.ParameterLiteral, error) {
	token := buf.Read()
	if token.Type != Ident {
		return nil, ExpectError{
			token,
			Ident,
		}
	}

	// TODO: check whether variable name exist,
	// TODO: but do not update symbol table in this case
	// TODO: if not, return error

	ident := &ast.ParameterLiteral{
		Identifier: &ast.Identifier{Value: token.Val},
	}

	token = buf.Read()
	ds, ok := datastructureMap[token.Type]
	if !ok {
		return nil, Error{
			token,
			"Function parameter type missed",
		}
	}
	ident.Type = ds

	return ident, nil
}

// parseReturnStatement parse "return" keyword with its expression
func parseReturnStatement(buf TokenBuffer) (ast.Statement, error) {
	if err := expectNext(buf, Return); err != nil {
		return nil, err
	}

	stmt := &ast.ReturnStatement{}

	exp, err := parseExpression(buf, LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.ReturnValue = exp

	consumeSemi(buf)

	return stmt, nil
}

// parseGroupedExpression parse grouped expression which
// grouped using parenthesis
func parseGroupedExpression(buf TokenBuffer) (ast.Expression, error) {
	buf.Read()
	exp, err := parseExpression(buf, LOWEST)
	if err != nil {
		return nil, err
	}
	if err = expectNext(buf, Rparen); err != nil {
		return nil, err
	}

	return exp, nil
}

// parseAssignStatement parse assign statements which assign values
// to its identifier. e.g. int a = 1
func parseAssignStatement(buf TokenBuffer) (*ast.AssignStatement, error) {
	stmt := &ast.AssignStatement{}

	token := buf.Read()

	ds := datastructureMap[token.Type]
	stmt.Type = ds

	token = buf.Read()
	if token.Type != Ident {
		return nil, ExpectError{
			token,
			Ident,
		}
	}

	// TODO: check whether variable name exist
	// TODO: if not, add symbol to the scope

	stmt.Variable = ast.Identifier{
		Value: token.Val,
	}

	if err := expectNext(buf, Assign); err != nil {
		return nil, err
	}

	exp, err := parseExpression(buf, LOWEST)
	if err != nil {
		return nil, err
	}

	stmt.Value = exp

	consumeSemi(buf)

	return stmt, nil
}

// parseCallExpression parse function call
func parseCallExpression(buf TokenBuffer, fn ast.Expression) (ast.Expression, error) {
	exp := &ast.CallExpression{Function: fn}

	var err error
	exp.Arguments, err = parseCallArguments(buf)
	if err != nil {
		return nil, err
	}

	consumeSemi(buf)

	return exp, nil
}

// parseCallArguments parse arguments of function call
func parseCallArguments(buf TokenBuffer) ([]ast.Expression, error) {
	args := []ast.Expression{}
	if err := expectNext(buf, Lparen); err != nil {
		return nil, err
	}

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

	consumeSemi(buf)

	if err := expectNext(buf, Rparen); err != nil {
		return nil, err
	}

	return args, nil
}

// parseIfStatement parse if-else statement. Else statement is optional
func parseIfStatement(buf TokenBuffer) (*ast.IfStatement, error) {
	if err := expectNext(buf, If); err != nil {
		return nil, err
	}

	if err := expectNext(buf, Lparen); err != nil {
		return nil, err
	}

	expression := &ast.IfStatement{}
	var err error
	expression.Condition, err = parseExpression(buf, LOWEST)
	if err != nil {
		return nil, err
	}

	if err := expectNext(buf, Rparen); err != nil {
		return nil, err
	}

	expression.Consequence, err = parseBlockStatement(buf)
	if err != nil {
		return nil, err
	}

	if curTokenIs(buf, Else) {
		buf.Read()

		expression.Alternative, err = parseBlockStatement(buf)
		if err != nil {
			return nil, err
		}
	}

	consumeSemi(buf)

	return expression, nil
}

// parseBlockStatement parse block statement.
// PROTOCOL:
//   reading token from TokenBuffer **only and must** be done in
//   parsing statement/expression function. And in parseBlockStatement,
//   parser reads from left-brace until meeting right-brace.
//   and consume right-brace.
//
//  parseBlockStatement parse: { ... } <-- left-brace + statements + right-brace
//
func parseBlockStatement(buf TokenBuffer) (*ast.BlockStatement, error) {
	if err := expectNext(buf, Lbrace); err != nil {
		return nil, err
	}

	// TODO: enter scope

	block := &ast.BlockStatement{}
	curToken := buf.Peek(CURRENT)

	for curToken.Type != Rbrace && curToken.Type != Eof {
		stmt, err := parseStatement(buf)
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		curToken = buf.Peek(CURRENT)
	}

	if curTokenIs(buf, Rbrace) {
		buf.Read()
	}

	// TODO: leave scope

	return block, nil
}

func parseExpressionStatement(buf TokenBuffer) (*ast.ExpressionStatement, error) {
	stmt := &ast.ExpressionStatement{}
	token := buf.Read()
	if token.Type != Ident {
		return nil, ExpectError{
			token,
			Ident,
		}
	}

	ident := &ast.Identifier{Value: token.Val}
	exp, err := parseCallExpression(buf, ident)
	if err != nil {
		return nil, err
	}

	stmt.Expr = exp
	return stmt, nil
}
