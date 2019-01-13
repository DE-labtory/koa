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

package translate

import (
	"errors"
	"fmt"

	"github.com/DE-labtory/koa/ast"
	"github.com/DE-labtory/koa/encoding"
	"github.com/DE-labtory/koa/opcode"
)

// Compile() compiles statements in ast.program.
// Statements would be compiled to byte code.
// TODO: implement w/ test cases :-)
func Compile(program ast.Program) ([]byte, error) {
	o := &Bytecode{
		RawByte: make([]byte, 0),
		AsmCode: make([]string, 0),
		PC:      0,
	}

	for _, s := range program.Statements {
		err := compileNode(s, o)
		if err != nil {
			return nil, err
		}
	}

	return o.RawByte, nil
}

// emit() generates a byte code with operator and operands.
// Then, saves the byte code and assemble code to output.
func emit(bytecode *Bytecode, operator opcode.Type, operands ...[]byte) {
	b := make([]byte, 0)
	s := make([]string, 0)

	b = append(b, byte(operator))
	s = append(s, operator.ToString())

	for _, o := range operands {
		b = append(b, o...)

		operand := fmt.Sprintf("%x", o)
		s = append(s, operand)
	}

	bytecode.RawByte = append(bytecode.RawByte, b...)
	bytecode.AsmCode = append(bytecode.AsmCode, s...)
}

// compileNode() compiles a node in statement.
// This function will be executed recursively.
// TODO: implement w/ test cases :-)
func compileNode(node ast.Node, bytecode *Bytecode) error {
	// Nodes are many kinds.
	switch node := node.(type) {
	case *ast.Identifier:
		return compileIdentifier(*node)

	case *ast.AssignStatement:
		return compileAssignStatement(*node)

	case *ast.ReturnStatement:
		return nil

	case *ast.StringLiteral:
		return compileString(*node)

	case *ast.IntegerLiteral:
		return compileInteger(*node)

	case *ast.BooleanLiteral:
		return compileBoolean(*node, bytecode)

	case *ast.PrefixExpression:
		return compilePrefixExpression(*node)

	case *ast.InfixExpression:
		return nil

	case *ast.CallExpression:
		return nil

	default:
		return errors.New("compileNode() error - " + node.String() + " could not compiled")
	}
}

// TODO: implement w/ test cases :-)
func compileIdentifier(node ast.Identifier) error {
	return nil
}

// TODO: implement w/ test cases :-)
func compileAssignStatement(node ast.AssignStatement) error {
	return nil
}

// TODO: implement w/ test cases :-)
func compileString(node ast.StringLiteral) error {
	return nil
}

// TODO: implement w/ test cases :-)
func compileInteger(node ast.IntegerLiteral) error {
	return nil
}

func compileBoolean(node ast.BooleanLiteral, bytecode *Bytecode) error {
	operand, err := encoding.EncodeOperand(node.Value)
	if err != nil {
		return err
	}

	emit(bytecode, opcode.Push, operand)
	return nil
}

// TODO: implement w/ test cases :-)
func compilePrefixExpression(node ast.PrefixExpression) error {
	return nil
}
