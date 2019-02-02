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
	"bytes"
	"fmt"
	"strconv"
	"strings"
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

// Represent Contract.
// Contract consists of multiple functions.
type Contract struct {
	Functions []*FunctionLiteral
}

func (c *Contract) do() {}
func (c *Contract) String() string {
	var buf bytes.Buffer

	// start by change line for readability
	buf.WriteString("\ncontract {\n")

	for _, fn := range c.Functions {
		buf.WriteString(fn.String() + "\n")
	}
	buf.WriteString("}")

	return buf.String()
}

// Represent identifier
type Identifier struct {
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) produce() {}

// Operator represent operator between expression
type Operator int

const (
	_        Operator = iota
	Plus              // +
	Minus             // -
	Bang              // !
	Asterisk          // *
	Slash             // /
	Mod               // %
	LT                // <
	GT                // >
	LTE               // <=
	GTE               // >=
	EQ                // ==
	NOT_EQ            // !=
	LAND              // &&
	LOR               // ||
)

var OperatorMap = map[Operator]string{
	Plus:     "+",
	Minus:    "-",
	Bang:     "!",
	Asterisk: "*",
	Slash:    "/",
	Mod:      "%",
	LT:       "<",
	GT:       ">",
	LTE:      "<=",
	GTE:      ">=",
	EQ:       "==",
	NOT_EQ:   "!=",
	LAND:     "&&",
	LOR:      "||",
}

func (o Operator) String() string {
	return OperatorMap[o]
}

// DataStructure represent identifier's data structure
// e.g. string, int, bool
type DataStructure int

const (
	_ DataStructure = iota
	IntType
	StringType
	BoolType
	VoidType
)

var DataStructureMap = map[DataStructure]string{
	IntType:    "int",
	StringType: "string",
	BoolType:   "bool",
	VoidType:   "void",
}

func (ds DataStructure) String() string {
	return DataStructureMap[ds]
}

// Represent assign statement
type AssignStatement struct {
	Type     DataStructure
	Variable Identifier
	Value    Expression
}

func (as *AssignStatement) do() {}

func (as *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Type.String() + " ")
	out.WriteString(as.Variable.Value + " = ")
	out.WriteString(as.Value.String())
	return out.String()
}

// Represent return statement
type ReturnStatement struct {
	ReturnValue Expression
}

func (rs *ReturnStatement) do() {}

func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", rs.ReturnValue.String())
}

// Represent if statement
type IfStatement struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) do() {}

func (is *IfStatement) String() string {
	if is.Alternative == nil {
		return fmt.Sprintf("if ( %s ) { %s }", is.Condition.String(), is.Consequence.String())
	}
	return fmt.Sprintf("if ( %s ) { %s } else { %s }", is.Condition.String(), is.Consequence.String(),
		is.Alternative.String())
}

// FunctionLiteral represents function definition
// e.g. func foo(int a) { ... }
type FunctionLiteral struct {
	Name       *Identifier
	Parameters []*ParameterLiteral
	Body       *BlockStatement
	ReturnType DataStructure
}

func (fl *FunctionLiteral) do() {}

// TODO: Add test cases when Type field is added to Identifier
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("func " + fl.Name.String() + "(")

	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.ReturnType.String() + " {\n")
	out.WriteString(fl.Body.String() + "\n")
	out.WriteString("}")

	return out.String()
}

func (f *FunctionLiteral) Signature() string {

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	return "func " + f.Name.String() + "(" + strings.Join(params, ", ") + ")"
}

// Represent block statement
type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) do() {}

func (bs *BlockStatement) String() string {
	if bs.Statements == nil {
		return ""
	}

	str := make([]string, 0)
	for _, s := range bs.Statements {
		str = append(str, s.String())
	}
	return strings.Join(str, "\n")
}

// Represent function statement
type ExpressionStatement struct {
	Expr Expression
}

func (es *ExpressionStatement) do() {}

func (es *ExpressionStatement) String() string {
	return es.Expr.String()
}

// Represent string literal
type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) produce() {}

func (sl *StringLiteral) String() string {
	return sl.Value
	//return fmt.Sprintf("\"%s\"", sl.Value)
}

// Represent integer literal
type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) produce() {}

func (il *IntegerLiteral) String() string {
	return strconv.FormatInt(il.Value, 10)
}

// Represent Boolean expression
type BooleanLiteral struct {
	Value bool
}

func (bl *BooleanLiteral) produce() {}

func (bl *BooleanLiteral) String() string {
	return strconv.FormatBool(bl.Value)
}

// Represent Function Parameter expression
type ParameterLiteral struct {
	Identifier *Identifier
	Type       DataStructure
}

func (pl *ParameterLiteral) produce() {}

func (pl *ParameterLiteral) String() string {
	return fmt.Sprintf("Parameter : (Identifier: %s, Type: %s)", pl.Identifier.String(), pl.Type.String())
}

// Represent prefix expression
type PrefixExpression struct {
	Operator
	Right Expression
}

func (pe *PrefixExpression) produce() {}

func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator.String(), pe.Right.String())
}

// Repersent Infix expression
type InfixExpression struct {
	Left Expression
	Operator
	Right Expression
}

func (ie *InfixExpression) produce() {}

func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator.String(), ie.Right.String())
}

// Represent Call expression
type CallExpression struct {
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) produce() {}

func (ce *CallExpression) String() string {
	strs := make([]string, 0)
	for _, exps := range ce.Arguments {
		strs = append(strs, exps.String())
	}
	return fmt.Sprintf("function %s( %s )", ce.Function.String(), strings.Join(strs, ", "))
}
