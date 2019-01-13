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

	"github.com/DE-labtory/koa/ast"
	"github.com/DE-labtory/koa/encoding"
	"github.com/DE-labtory/koa/opcode"
)

// Compile() compiles statements in ast.program.
// Statements would be compiled to byte code.
// TODO: implement w/ test cases :-)
func Compile(program ast.Program) ([]byte, error) {
	bin := make([]byte, 0)

	for _, s := range program.Statements {
		b, err := compileNode(s)
		if err != nil {
			return nil, err
		}

		bin = append(bin, b...)
	}

	return bin, nil
}

// emit() generates a byte code with operator and operands.
// Then, returns byte code.
func emit(operator opcode.Type, operands ...[]byte) []byte {
	b := make([]byte, 0)

	b = append(b, byte(operator))

	for _, o := range operands {
		b = append(b, o...)
	}

	return b
}

// compileNode() compiles a node in statement.
// This function will be executed recursively.
func compileNode(node ast.Node) ([]byte, error) {
	// Nodes are many kinds.
	switch node := node.(type) {
	case *ast.Identifier:
		return compileIdentifier(*node)

	case *ast.AssignStatement:
		return compileAssignStatement(*node)

	case *ast.StringLiteral:
		return compileString(*node)

	case *ast.IntegerLiteral:
		return compileInteger(*node)

	case *ast.BooleanLiteral:
		return compileBoolean(*node)

	case *ast.PrefixExpression:
		return compilePrefixExpression(*node)

	default:
		err := errors.New("compileNode() error - " + node.String() + " could not compiled")
		return nil, err
	}
}

// TODO: implement w/ test cases :-)
func compileIdentifier(node ast.Identifier) ([]byte, error) {
	return nil, nil
}

// TODO: implement w/ test cases :-)
func compileAssignStatement(node ast.AssignStatement) ([]byte, error) {
	return nil, nil
}

// TODO: implement w/ test cases :-)
func compileString(node ast.StringLiteral) ([]byte, error) {
	return nil, nil
}

// TODO: implement w/ test cases :-)
func compileInteger(node ast.IntegerLiteral) ([]byte, error) {
	return nil, nil
}

func compileBoolean(node ast.BooleanLiteral) ([]byte, error) {
	operand, err := encoding.EncodeOperand(node.Value)
	if err != nil {
		return nil, err
	}

	b := emit(opcode.Push, operand)
	return b, nil
}

// TODO: implement w/ test cases :-)
func compilePrefixExpression(node ast.PrefixExpression) ([]byte, error) {
	return nil, nil
}
