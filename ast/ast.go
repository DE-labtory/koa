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

type OperatorType int

const (
	_        OperatorType = iota
	Plus                  // +
	Minus                 // -
	Bang                  // !
	Asterisk              // *
	Slash                 // /
	Mod                   // %
	LT                    // <
	GT                    // >
	LTE                   // <=
	GTE                   // >=
	EQ                    // ==
	NOT_EQ                // !=
)

type OperatorVal string

// Operator represent operator between expression
type Operator struct {
	Type OperatorType
	Val  OperatorVal
}

type DataStructureType int

const (
	_ DataStructureType = iota
	Int
	String
	Bool
)

type DataStructureVal string

// DataStructure represent identifier's data structure
// e.g. string, int, bool
type DataStructure struct {
	Type DataStructureType
	Val  DataStructureVal
}

// Represent assign statement
type AssignStatement struct {
	Type     DataStructure
	Variable *Identifier
	Value    Expression
}

func (as *AssignStatement) do() {}

// TODO: implement me w/ test cases :-)
func (as *AssignStatement) String() string {
	return ""
}

// Represent string literal
type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) produce() {}

func (sl *StringLiteral) String() string {
	return sl.Value
}

// Represent integer literal
type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) produce() {}

func (il *IntegerLiteral) String() string {
	return string(il.Value)
}

// Represent Boolean expression
type BooleanLiteral struct {
	Value bool
}

func (bl *BooleanLiteral) produce() {}

func (bl *BooleanLiteral) String() string {
	if bl.Value {
		return "true"
	}
	return "false"
}

// Represent prefix expression
type PrefixExpression struct {
	Operator
	Right Expression
}

func (pe *PrefixExpression) produce() {}

// TODO: implement me w/ test cases :-)
func (pe *PrefixExpression) String() string {
	return ""
}
