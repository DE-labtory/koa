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

	"github.com/DE-labtory/koa/abi"
	"github.com/DE-labtory/koa/ast"
	"github.com/DE-labtory/koa/encoding"
	"github.com/DE-labtory/koa/opcode"
)

type FuncMap map[string]int

func (m FuncMap) Declare(f ast.FunctionLiteral, b Bytecode) {
	funcSel := abi.Selector(f.Signature())
	m[string(funcSel)] = len(b.RawByte)
}

// TODO: implement me w/ test cases :-)
// CompileContract() compiles a smart contract.
// returns bytecode and error.
func CompileContract(c ast.Contract) (Bytecode, error) {
	bytecode := &Bytecode{
		RawByte: make([]byte, 0),
		AsmCode: make([]string, 0),
	}

	memTracer := NewMemEntryTable()

	funcMap := FuncMap{}

	for _, f := range c.Functions {
		funcMap.Declare(*f, *bytecode)

		if err := compileFunction(*f, bytecode, memTracer); err != nil {
			return *bytecode, err
		}
	}

	if err := expectFuncJumper(c, funcMap); err != nil {
		return *bytecode, err
	}

	if err := generateFuncJumper(c, bytecode, funcMap); err != nil {
		return *bytecode, err
	}

	return *bytecode, nil
}

// Expects a size of the function jumper and updates the function map.
func expectFuncJumper(c ast.Contract, funcMap FuncMap) error {
	expectedJpr := NewBytecode()
	if err := generateFuncJumper(c, expectedJpr, funcMap); err != nil {
		return err
	}

	// Updates the function map with a size of function jumper.
	for funcSel, pc := range funcMap {
		funcMap[funcSel] = pc + len(expectedJpr.RawByte)
	}

	return nil
}

// Generates the function jumper in front of the bytecode.
func generateFuncJumper(c ast.Contract, bytecode *Bytecode, funcMap FuncMap) error {
	funcJpr := NewBytecode()

	// Loads the function selector of call function.
	funcJpr.Emerge(opcode.LoadFunc)

	// Compares and finds the corresponding function selector.
	for _, f := range c.Functions {
		if err := compileFunctionLiteral(f, funcJpr, funcMap); err != nil {
			return err
		}
	}

	// No match to any function selector, Revert!
	funcJpr.Emerge(opcode.JumpDst)
	operand, err := encoding.EncodeOperand(0)
	if err != nil {
		return err
	}
	funcJpr.Emerge(opcode.Push, operand)
	funcJpr.Emerge(opcode.DUP)
	funcJpr.Emerge(opcode.Returning)

	// Add function jumper to bytecode
	// [bytecode] => [function jumper][bytecode]
	funcJpr.Append(bytecode)
	bytecode = funcJpr

	return nil
}

// TODO: implement me w/ test cases :-)
func compileFunctionLiteral(s *ast.FunctionLiteral, bytecode *Bytecode, funcMap FuncMap) error {
	// Duplicates the function selector to find.
	bytecode.Emerge(opcode.DUP)
	// Pushes the function selector of this function literal.
	selector := string(abi.Selector(s.Signature()))
	funcSel, err := encoding.EncodeOperand(selector)
	if err != nil {
		return err
	}
	bytecode.Emerge(opcode.Push, funcSel)
	bytecode.Emerge(opcode.EQ)
	// If the result is equal, Pushed the pc to jump.
	pc, err := encoding.EncodeOperand(funcMap[selector])
	if err != nil {
		return err
	}
	bytecode.Emerge(opcode.Push, pc)
	bytecode.Emerge(opcode.Jumpi)

	return nil
}

// TODO: implement me w/ test cases :-)
func compileParameterLiteral(e *ast.ParameterLiteral, bytecode *Bytecode) error {
	return nil
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
func compileFunction(f ast.FunctionLiteral, bytecode *Bytecode, tracer MemTracer) error {
	// TODO: generate function identifier with Keccak256()

	statements := f.Body.Statements
	for _, s := range statements {
		if err := compileStatement(s, bytecode, tracer); err != nil {
			return err
		}
	}

	return nil
}

// TODO: implement me w/ test cases :-)
// compileStatement() compiles a statement in function.
// Generates and adds output to bytecode.
func compileStatement(s ast.Statement, bytecode *Bytecode, tracer MemTracer) error {
	switch statement := s.(type) {
	case *ast.AssignStatement:
		return compileAssignStatement(statement, bytecode, tracer)

	case *ast.ReturnStatement:
		return compileReturnStatement(statement, bytecode)

	case *ast.IfStatement:
		return compileIfStatement(statement, bytecode)

	case *ast.BlockStatement:
		return compileBlockStatement(statement, bytecode, tracer)

	case *ast.ExpressionStatement:
		return compileExpressionStatement(statement, bytecode)

	default:
		return nil
	}
}

// compileAssignStatement() compiles a assign statement.
//
// Ex)
//
// translate
// 	'int a = 5'
// to
// 	'Push 5 Push <size of a> Push <offset of a> Mstore'
//
// stack will be
//
// 	[offset]
// 	[size]
// 	[value]
//
func compileAssignStatement(s *ast.AssignStatement, bytecode *Bytecode, memDefiner MemDefiner) error {
	if err := compileExpression(s.Value, bytecode); err != nil {
		return err
	}

	memEntry := memDefiner.Define(s.Variable.Value)
	size, err := encoding.EncodeOperand(memEntry.Size)
	if err != nil {
		return err
	}

	offset, err := encoding.EncodeOperand(memEntry.Offset)
	if err != nil {
		return err
	}

	bytecode.Emerge(opcode.Push, size)
	bytecode.Emerge(opcode.Push, offset)
	bytecode.Emerge(opcode.Mstore)
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

func compileBlockStatement(s *ast.BlockStatement, bytecode *Bytecode, tracer MemTracer) error {
	for _, statement := range s.Statements {
		if err := compileStatement(statement, bytecode, tracer); err != nil {
			return err
		}
	}

	return nil
}

func compileExpressionStatement(s *ast.ExpressionStatement, bytecode *Bytecode) error {
	return compileExpression(s.Expr, bytecode)
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

	default:
		return errors.New("compileExpression() error")
	}
}

// TODO: implement me w/ test cases :-)
func compileCallExpression(e *ast.CallExpression, bytecode *Bytecode) error {
	return nil
}

func compileInfixExpression(e *ast.InfixExpression, bytecode *Bytecode) error {
	if err := compileExpression(e.Left, bytecode); err != nil {
		return err
	}

	if err := compileExpression(e.Right, bytecode); err != nil {
		return err
	}

	switch e.Operator {
	case ast.Plus:
		bytecode.Emerge(opcode.Add)
	case ast.Minus:
		bytecode.Emerge(opcode.Sub)
	case ast.Asterisk:
		bytecode.Emerge(opcode.Mul)
	case ast.Slash:
		bytecode.Emerge(opcode.Div)
	case ast.Mod:
		bytecode.Emerge(opcode.Mod)

		//comparison
	case ast.LT:
		bytecode.Emerge(opcode.LT)
	case ast.GT:
		bytecode.Emerge(opcode.GT)
	case ast.LTE:
		bytecode.Emerge(opcode.LTE)
	case ast.GTE:
		bytecode.Emerge(opcode.GTE)
	case ast.EQ:
		bytecode.Emerge(opcode.EQ)
	case ast.NOT_EQ:
		bytecode.Emerge(opcode.EQ)
		bytecode.Emerge(opcode.NOT)
	case ast.LAND:
		bytecode.Emerge(opcode.And)
	case ast.LOR:
		bytecode.Emerge(opcode.Or)
	}

	return nil
}

func compilePrefixExpression(e *ast.PrefixExpression, bytecode *Bytecode) error {
	if err := compileExpression(e.Right, bytecode); err != nil {
		return err
	}

	switch e.Operator {
	case ast.Bang:
		bytecode.Emerge(opcode.NOT)
	case ast.Minus:
		bytecode.Emerge(opcode.Minus)
	default:
		return fmt.Errorf("unknown operator %s", e.Operator.String())
	}

	return nil
}

func compileIntegerLiteral(e *ast.IntegerLiteral, bytecode *Bytecode) error {
	operand, err := encoding.EncodeOperand(e.Value)
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
	operand, err := encoding.EncodeOperand(e.Value)
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
