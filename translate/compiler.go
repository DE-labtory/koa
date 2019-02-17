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

// Declare() saves the start point of function.
func (m FuncMap) Declare(signature string, asm Asm) {
	funcSig := abi.Selector(signature)
	m[string(funcSig)] = len(asm.AsmCodes)
}

// TODO: implement me w/ test cases :-)
// CompileContract() compiles a smart contract.
// returns bytecode and error.
func CompileContract(c ast.Contract) (Asm, error) {
	asm := &Asm{
		AsmCodes: make([]AsmCode, 0),
	}

	memTracer := NewMemEntryTable()

	funcMap := FuncMap{}
	for _, f := range c.Functions {
		funcMap.Declare(f.Signature(), *asm)

		if err := compileFunction(*f, asm, memTracer); err != nil {
			return *asm, err
		}
	}

	if err := generateFuncJumper(asm); err != nil {
		return *asm, err
	}

	return *asm, nil
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
func generateFuncJumper(bytecode *Asm) error {
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
func compileFunction(f ast.FunctionLiteral, bytecode *Asm, tracer MemTracer) error {
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
func compileStatement(s ast.Statement, bytecode *Asm, tracer MemTracer) error {
	switch statement := s.(type) {
	case *ast.AssignStatement:
		return compileAssignStatement(statement, bytecode, tracer)

	case *ast.ReturnStatement:
		return compileReturnStatement(statement, bytecode)

	case *ast.IfStatement:
		return compileIfStatement(statement, bytecode, tracer)

	case *ast.BlockStatement:
		return compileBlockStatement(statement, bytecode, tracer)

	case *ast.ExpressionStatement:
		return compileExpressionStatement(statement, bytecode, tracer)

	case *ast.FunctionLiteral:
		return compileFunctionLiteral(statement, bytecode)

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
func compileAssignStatement(s *ast.AssignStatement, asm *Asm, tracer MemTracer) error {
	if err := compileExpression(s.Value, asm, tracer); err != nil {
		return err
	}

	memEntry := tracer.Define(s.Variable.Value)
	size, err := encoding.EncodeOperand(memEntry.Size)
	if err != nil {
		return err
	}

	offset, err := encoding.EncodeOperand(memEntry.Offset)
	if err != nil {
		return err
	}

	asm.Emerge(opcode.Push, size)
	asm.Emerge(opcode.Push, offset)
	asm.Emerge(opcode.Mstore)
	return nil
}

// TODO: implement me w/ test cases :-)
func compileReturnStatement(s *ast.ReturnStatement, bytecode *Asm) error {
	return nil
}

// compileIfStatement() compiles a 'if statement'.
//
// Ex)
//
// translate
// 	'if (expression){
// 		// Consequence...
//  }else {
// 		// Alternative...
//  }'
// to
//  'push <expression> push <pc-to-jumpdst-1> jumpi <Consequence...> push <pc-to-end-of-jumpdst-2> jump jumpdst-1 <Alternative...> jumpdst-2'
//
func compileIfStatement(s *ast.IfStatement, asm *Asm, tracer MemTracer) error {

	if err := compileExpression(s.Condition, asm, tracer); err != nil {
		return err
	}

	if s.Alternative != nil {
		return compileIfElse(s, asm, tracer)
	}

	return compileIf(s, asm, tracer)
}

func compileIfElse(s *ast.IfStatement, asm *Asm, tracer MemTracer) error {

	l1 := len(asm.AsmCodes)
	if err := compileBlockStatement(s.Consequence, asm, tracer); err != nil {
		return err
	}
	// 'push <expression> <Consequence...>'

	l2 := len(asm.AsmCodes)
	if err := compileBlockStatement(s.Alternative, asm, tracer); err != nil {
		return err
	}
	// 'push <Consequence...> <Alternative...>'

	l3 := len(asm.AsmCodes)
	pc2al, err := encoding.EncodeOperand(l2 + 6)
	if err != nil {
		return err
	}

	asm.EmergeAt(l1, opcode.Jumpi)
	asm.EmergeAt(l1, opcode.Push, pc2al)
	// 'push <expression> push <pc-to-Alternative> jumpi <Consequence...> <Alternative...>'

	pc2EndOfAlter, err := encoding.EncodeOperand(l3 + 6)
	if err != nil {
		return err
	}

	asm.EmergeAt(l2+3, opcode.Jump)
	asm.EmergeAt(l2+3, opcode.Push, pc2EndOfAlter)
	// 'push <expression> push <pc-to-Alternative> jumpi <Consequence...> push <pc-to-end-of-Alternative> jump <Alternative...>'

	return nil
}

func compileIf(s *ast.IfStatement, asm *Asm, tracer MemTracer) error {

	l1 := len(asm.AsmCodes)
	if err := compileBlockStatement(s.Consequence, asm, tracer); err != nil {
		return err
	}
	l2 := len(asm.AsmCodes)
	// 'push <expression> <Consequence...>'

	pc2al, err := encoding.EncodeOperand(l2 + 1)
	if err != nil {
		return err
	}

	asm.EmergeAt(l1, opcode.Jump)
	asm.EmergeAt(l1, opcode.Push, pc2al)
	// 'push <expression> push <pc-to-end-of-Consequence> jump <Consequence...>'

	return nil
}

func compileBlockStatement(s *ast.BlockStatement, bytecode *Asm, tracer MemTracer) error {
	for _, statement := range s.Statements {
		if err := compileStatement(statement, bytecode, tracer); err != nil {
			return err
		}
	}

	return nil
}

func compileExpressionStatement(s *ast.ExpressionStatement, bytecode *Asm, tracer MemTracer) error {
	return compileExpression(s.Expr, bytecode, tracer)
}

// TODO: implement me w/ test cases :-)
func compileFunctionLiteral(s *ast.FunctionLiteral, bytecode *Asm) error {
	return nil
}

// TODO: implement me w/ test cases :-)
// compileExpression() compiles a expression in statement.
// Generates and adds ouput to bytecode.
func compileExpression(e ast.Expression, asm *Asm, tracer MemTracer) error {
	switch expr := e.(type) {
	case *ast.CallExpression:
		return compileCallExpression(expr, asm)

	case *ast.InfixExpression:
		return compileInfixExpression(expr, asm, tracer)

	case *ast.PrefixExpression:
		return compilePrefixExpression(expr, asm, tracer)

	case *ast.IntegerLiteral:
		return compilePrimitive(expr.Value, asm)

	case *ast.StringLiteral:
		return compilePrimitive(expr.Value, asm)

	case *ast.BooleanLiteral:
		return compilePrimitive(expr.Value, asm)

	case *ast.Identifier:
		return compileIdentifier(expr, asm, tracer)

	case *ast.ParameterLiteral:
		return compileParameterLiteral(expr, asm)

	default:
		return errors.New("compileExpression() error")
	}
}

// TODO: implement me w/ test cases :-)
func compileCallExpression(e *ast.CallExpression, asm *Asm) error {
	return nil
}

func compileInfixExpression(e *ast.InfixExpression, asm *Asm, tracer MemTracer) error {
	if err := compileExpression(e.Left, asm, tracer); err != nil {
		return err
	}

	if err := compileExpression(e.Right, asm, tracer); err != nil {
		return err
	}

	switch e.Operator {
	case ast.Plus:
		asm.Emerge(opcode.Add)
	case ast.Minus:
		asm.Emerge(opcode.Sub)
	case ast.Asterisk:
		asm.Emerge(opcode.Mul)
	case ast.Slash:
		asm.Emerge(opcode.Div)
	case ast.Mod:
		asm.Emerge(opcode.Mod)

		//comparison
	case ast.LT:
		asm.Emerge(opcode.LT)
	case ast.GT:
		asm.Emerge(opcode.GT)
	case ast.LTE:
		asm.Emerge(opcode.LTE)
	case ast.GTE:
		asm.Emerge(opcode.GTE)
	case ast.EQ:
		asm.Emerge(opcode.EQ)
	case ast.NOT_EQ:
		asm.Emerge(opcode.EQ)
		asm.Emerge(opcode.NOT)
	case ast.LAND:
		asm.Emerge(opcode.And)
	case ast.LOR:
		asm.Emerge(opcode.Or)
	}

	return nil
}

func compilePrefixExpression(e *ast.PrefixExpression, asm *Asm, tracer MemTracer) error {
	if err := compileExpression(e.Right, asm, tracer); err != nil {
		return err
	}

	switch e.Operator {
	case ast.Bang:
		asm.Emerge(opcode.NOT)
	case ast.Minus:
		asm.Emerge(opcode.Minus)
	default:
		return fmt.Errorf("unknown operator %s", e.Operator.String())
	}

	return nil
}

func compilePrimitive(value interface{}, asm *Asm) error {
	operand, err := encoding.EncodeOperand(value)
	if err != nil {
		return err
	}

	asm.Emerge(opcode.Push, operand)
	return nil
}

func compileIdentifier(e *ast.Identifier, asm *Asm, tracer MemTracer) error {
	memEntry, err := tracer.GetEntry(e.Value)
	if err != nil {
		return err
	}

	size, err := encoding.EncodeOperand(memEntry.Size)
	if err != nil {
		return err
	}

	offset, err := encoding.EncodeOperand(memEntry.Offset)
	if err != nil {
		return err
	}

	asm.Emerge(opcode.Push, size)
	asm.Emerge(opcode.Push, offset)
	asm.Emerge(opcode.Mload)

	return nil
}

// TODO: implement me w/ test cases :-)
func compileParameterLiteral(e *ast.ParameterLiteral, bytecode *Asm) error {
	return nil
}
