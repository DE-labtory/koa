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

package ast

import (
	"github.com/DE-labtory/koa/parse"
)

// Node represent ast node
type Node interface {
	String() string
}

// Represent Statement
type Statement interface {
	Node
	do()
}

// Represent Expression
type Expression interface {
	Node
	produce()
}

// Represent Program.
// Program consists of multiple statements.
type Program struct {
	Statements []Statement
}

// Represent identifier
type Identifier struct {
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) produce() {}

// Represent assign statement
type AssignStatement struct {
	Type  parse.Token
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) do() {}

// TODO: implement me w/ test cases :-)
func (as *AssignStatement) String() string {
	return ""
}

// Represent return statement
type ReturnStatement struct {
	ReturnValue Expression
}

func (rs *ReturnStatement) do() {}

// TODO: implement me w/ test cases :-)
func (rs *ReturnStatement) String() string {
	return ""
}

// Represent string literal
type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) produce() {}

// TODO: implement me w/ test cases :-)
func (sl *StringLiteral) String() string {
	return ""
}

// Represent integer literal
type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) produce() {}
func (il *IntegerLiteral) String() string {
	return string(il.Value)
}

// Represent boolean literal
type Boolean struct {
	Value bool
}

func (b *Boolean) produce() {}

// TODO: implement me w/ test cases :-)
func (b *Boolean) String() string {
	return ""
}

// Represent prefix expression
type PrefixExpression struct {
	Operator parse.Token
	Right    Expression
}

func (pe *PrefixExpression) produce() {}

// TODO: implement me w/ test cases :-)
func (pe *PrefixExpression) String() string {
	return ""
}

// Represent infix expression
type InfixExpression struct {
	Left     Expression
	Operator parse.Token
	Right    Expression
}

func (ie *InfixExpression) produce() {}

// TODO: implement me w/ test cases :-)
func (ie *InfixExpression) String() string {
	return ""
}

// Represent if, else expression
type IfExpression struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) do() {}

// TODO: implement me w/ test cases :-)
func (ie *IfExpression) String() string {
	return ""
}

// Represent block statements, statements contained between
// right brace, left brace
type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) do() {}

// TODO: implement me w/ test cases :-)
func (bs *BlockStatement) String() string {
	return ""
}
