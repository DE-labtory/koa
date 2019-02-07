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

	"github.com/DE-labtory/koa/abi"
	"github.com/DE-labtory/koa/ast"
	"github.com/DE-labtory/koa/encoding"
	"github.com/DE-labtory/koa/opcode"
)

// TODO: implement me w/ test cases :-)
// CompileContract() compiles a smart contract.
// returns bytecode and error.
func CompileContract(c ast.Contract) (Bytecode, error) {
	bytecode := &Bytecode{
		RawByte: make([]byte, 0),
		AsmCode: make([]string, 0),
	}

	for _, f := range c.Functions {
		if err := compileFunction(*f, bytecode); err != nil {
			return *bytecode, err
		}
	}

	if err := generateFuncJumper(bytecode); err != nil {
		return *bytecode, err
	}

	return *bytecode, nil
}

func ExtractAbi(c ast.Contract) (*abi.ABI, error) {
	abiMethods, err := toAbiMethods(c.Functions)
	if err != nil {
		return nil, err
	}

	return &abi.ABI{
		Methods: abiMethods,
	}, nil
}

// TODO: implement me w/ test cases :-)
// Generates a bytecode of function jumper.
func generateFuncJumper(bytecode *Bytecode) error {
	return nil
}

// Generates the ABI of functions in contract.
// Then, adds the ABI to bytecode.
func toAbiMethods(functions []*ast.FunctionLiteral) ([]abi.Method, error) {
	methods := make([]abi.Method, 0)

	for _, f := range functions {
		m, err := abi.ExtractAbiFromFunction(*f)
		if err != nil {
			return nil, err
		}
		methods = append(methods, m)
	}

	return methods, nil
}

// TODO: implement me w/ test cases :-)
// compileFunction() compiles a function in contract.
// Generates and adds output to bytecode.
func compileFunction(f ast.FunctionLiteral, bytecode *Bytecode) error {
	// TODO: generate function identifier with Keccak256()

	statements := f.Body.Statements
	for _, s := range statements {
		if err := compileStatement(s, bytecode); err != nil {
			return err
		}
	}

	return nil
}

// TODO: implement me w/ test cases :-)
// compileStatement() compiles a statement in function.
// Generates and adds output to bytecode.
func compileStatement(s ast.Statement, bytecode *Bytecode) error {
	switch statement := s.(type) {
	case *ast.AssignStatement:
		return compileAssignStatement(statement, bytecode)

	case *ast.ReturnStatement:
		return compileReturnStatement(statement, bytecode)

	case *ast.IfStatement:
		return compileIfStatement(statement, bytecode)

	case *ast.BlockStatement:
		return compileBlockStatement(statement, bytecode)

	case *ast.ExpressionStatement:
		return compileExpressionStatement(statement, bytecode)

	case *ast.FunctionLiteral:
		return compileFunctionLiteral(statement, bytecode)

	default:
		return nil
	}
}

// TODO: implement me w/ test cases :-)
func compileAssignStatement(s *ast.AssignStatement, bytecode *Bytecode) error {
	return nil
}

// TODO: implement me w/ test cases :-)
func compileReturnStatement(s *ast.ReturnStatement, bytecode *Bytecode) error {
	return nil
}

// TODO: implement me w/ test cases :-)
func compileIfStatement(s *ast.IfStatement, bytecode *Bytecode) error {
	return nil
}

func compileBlockStatement(s *ast.BlockStatement, bytecode *Bytecode) error {
	for _, statement := range s.Statements {
		if err := compileStatement(statement, bytecode); err != nil {
			return err
		}
	}

	return nil
}

// TODO: implement me w/ test cases :-)
func compileExpressionStatement(s *ast.ExpressionStatement, bytecode *Bytecode) error {
	return nil
}

// TODO: implement me w/ test cases :-)
func compileFunctionLiteral(s *ast.FunctionLiteral, bytecode *Bytecode) error {
	return nil
}

// TODO: implement me w/ test cases :-)
// compileExpression() compiles a expression in statement.
// Generates and adds ouput to bytecode.
func compileExpression(e ast.Expression, bytecode *Bytecode) error {
	switch expr := e.(type) {
	case *ast.CallExpression:
		return compileCallExpression(expr, bytecode)

	case *ast.InfixExpression:
		return compileInfixExpression(expr, bytecode)

	case *ast.PrefixExpression:
		return compilePrefixExpression(expr, bytecode)

	case *ast.IntegerLiteral:
		return compileIntegerLiteral(expr, bytecode)

	case *ast.StringLiteral:
		return compileStringLiteral(expr, bytecode)

	case *ast.BooleanLiteral:
		return compileBooleanLiteral(expr, bytecode)

	case *ast.Identifier:
		return compileIdentifier(expr, bytecode)

	case *ast.ParameterLiteral:
		return compileParameterLiteral(expr, bytecode)

	default:
		return errors.New("compileExpression() error")
	}
}

// TODO: implement me w/ test cases :-)
func compileCallExpression(e *ast.CallExpression, bytecode *Bytecode) error {
	return nil
}

// TODO: implement me w/ test cases :-)
func compileInfixExpression(e *ast.InfixExpression, bytecode *Bytecode) error {
	return nil
}

// TODO: implement me w/ test cases :-)
func compilePrefixExpression(e *ast.PrefixExpression, bytecode *Bytecode) error {
	return nil
}

func compileIntegerLiteral(e *ast.IntegerLiteral, bytecode *Bytecode) error {
	operand, err := encoding.EncodeOperand(e.Value, encoding.EIGHT_PADDING)
	if err != nil {
		return err
	}

	bytecode.Emerge(opcode.Push, operand)
	return nil
}

// TODO: implement me w/ test cases :-)
func compileStringLiteral(e *ast.StringLiteral, bytecode *Bytecode) error {
	return nil

}

func compileBooleanLiteral(e *ast.BooleanLiteral, bytecode *Bytecode) error {
	operand, err := encoding.EncodeOperand(e.Value, encoding.EIGHT_PADDING)
	if err != nil {
		return err
	}

	bytecode.Emerge(opcode.Push, operand)
	return nil
}

// TODO: implement me w/ test cases :-)
func compileIdentifier(e *ast.Identifier, bytecode *Bytecode) error {
	return nil
}

// TODO: implement me w/ test cases :-)
func compileParameterLiteral(e *ast.ParameterLiteral, bytecode *Bytecode) error {
	return nil
}
