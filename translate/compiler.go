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

func (m FuncMap) Declare(signature string, b Bytecode) {
	funcSel := abi.Selector(signature)
	m[string(funcSel)] = len(b.AsmCode)
}

// TODO: implement me w/ test cases :-)
// CompileContract() compiles a smart contract.
// returns bytecode and error.
func CompileContract(c ast.Contract) (Bytecode, error) {
	bytecode := &Bytecode{
		RawByte: make([]byte, 0),
		AsmCode: make([]string, 0),
	}

	funcMap := FuncMap{}
	if err := expectFuncJmpr(c, bytecode, funcMap); err != nil {
		return *bytecode, err
	}

	memTracer := NewMemEntryTable()
	for _, f := range c.Functions {
		funcMap.Declare(f.Signature(), *bytecode)

		if err := compileFunction(*f, bytecode, memTracer); err != nil {
			return *bytecode, err
		}
	}

	funcJmpr, err := generateFuncJmpr(c, funcMap)
	if err != nil {
		return *bytecode, err
	}

	if err := relocateFuncJmpr(bytecode, funcJmpr); err != nil {
		return *bytecode, err
	}

	return *bytecode, nil
}

// Expects a size of the function jumper and emerges with the unmeaningful value.
func expectFuncJmpr(c ast.Contract, bytecode *Bytecode, funcMap FuncMap) error {
	// Pushes the location of revert with the unmeaningful value.
	operand, err := encoding.EncodeOperand(0)
	if err != nil {
		return err
	}
	bytecode.Emerge(opcode.Push, operand)

	// Loads the function selector of call function.
	bytecode.Emerge(opcode.LoadFunc)

	// Adds the logic to compare and find the corresponding function selector with the unmeaningful value.
	funcMap.Declare("FuncJmpr", *bytecode)
	for range c.Functions {
		if err := compileFuncSelector(bytecode, string(abi.Selector("")), 0); err != nil {
			return err
		}
	}

	// No match to any function selector, Revert!
	funcMap.Declare("Revert", *bytecode)
	bytecode.Emerge(opcode.Returning)

	return nil
}

// Generates the function jumper bytecode.
func generateFuncJmpr(c ast.Contract, funcMap FuncMap) (*Bytecode, error) {
	funcJmpr := NewBytecode()

	// Pushes the location of revert.
	operand, err := encoding.EncodeOperand(funcMap[string(abi.Selector("Revert"))])
	if err != nil {
		return funcJmpr, err
	}
	funcJmpr.Emerge(opcode.Push, operand)

	// Loads the function selector of call function.
	funcJmpr.Emerge(opcode.LoadFunc)

	// Adds the logic to compare and find the corresponding function selector.
	for _, f := range c.Functions {
		selector := string(abi.Selector(f.Signature()))
		funcDst := funcMap[selector]
		if err := compileFuncSelector(funcJmpr, selector, funcDst); err != nil {
			return funcJmpr, err
		}
	}

	// No match to any function selector, Revert!
	funcJmpr.Emerge(opcode.Returning)

	return funcJmpr, nil
}

// Relocate the function jumper of bytecode with new function jumper.
func relocateFuncJmpr(bytecode *Bytecode, funcJmpr *Bytecode) error {
	if len(bytecode.RawByte) < len(funcJmpr.RawByte) {
		return fmt.Errorf("Can't relocate the function jumper. Bytecode=%x, FuncJmpr=%x", bytecode.RawByte, funcJmpr.RawByte)
	}

	if len(bytecode.AsmCode) < len(funcJmpr.AsmCode) {
		return fmt.Errorf("Can't relocate the function jumper. Bytecode=%x, FuncJmpr=%x", bytecode.AsmCode, funcJmpr.AsmCode)
	}

	for i := range funcJmpr.RawByte {
		bytecode.RawByte[i] = funcJmpr.RawByte[i]
	}

	for i := range funcJmpr.AsmCode {
		bytecode.AsmCode[i] = funcJmpr.AsmCode[i]
	}

	return nil
}

// Compiles a function of function jumper with its function selector
func compileFuncSelector(bytecode *Bytecode, funcSelector string, funcDst int) error {
	// Duplicates the function selector to find.
	bytecode.Emerge(opcode.DUP)
	// Pushes the function selector of this function literal.
	selector, err := encoding.EncodeOperand(funcSelector)
	if err != nil {
		return err
	}
	bytecode.Emerge(opcode.Push, selector)
	bytecode.Emerge(opcode.EQ)
	// If the result is equal, pushed the destination to jump.
	dst, err := encoding.EncodeOperand(funcDst)
	if err != nil {
		return err
	}
	bytecode.Emerge(opcode.Push, dst)
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
