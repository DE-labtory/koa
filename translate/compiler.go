/*
 * Copyright 2018-2019 De-labtory
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

	// Keep the size of the memory with createMemSizePlaceholder.
	if err := createMemSizePlaceholder(asm); err != nil {
		return *asm, err
	}

	// Keep the size of the jumper with createFuncJmprPlaceholder.
	funcMap := FuncMap{}
	if err := createFuncJmprPlaceholder(c, asm, funcMap); err != nil {
		return *asm, err
	}

	// Compile the functions in contract.
	memTracer := NewMemEntryTable()

	for _, f := range c.Functions {
		funcMap.Declare(f.Signature(), *asm)

		if err := compileFunction(*f, asm, memTracer); err != nil {
			return *asm, err
		}
	}

	// Compile Memory size with updated memory table.
	// And replace expected memory size with new memory size of the memory table.
	if err := compileMemSize(asm, memTracer); err != nil {
		return *asm, err
	}

	// Compile Function jumper with updated FuncMap.
	// And replace expected function jumper with new function jumper
	if err := compileFuncJmpr(c, asm, funcMap); err != nil {
		return *asm, err
	}

	return *asm, nil
}

// TODO: implement test cases :-)
// Create a placeholder to calculate a size of the memory.
// It emerges with the unmeaningful value.
func createMemSizePlaceholder(asm *Asm) error {
	operand, err := encoding.EncodeOperand(0)
	if err != nil {
		return err
	}
	asm.Emerge(opcode.Push, operand)
	asm.Emerge(opcode.Msize)

	return nil
}

// Create a placeholder to calculate a size of the function jumper.
// It emerges with the unmeaningful value.
func createFuncJmprPlaceholder(c ast.Contract, asm *Asm, funcMap FuncMap) error {
	// Pushes the location of revert with the unmeaningful value.
	if err := compileProgramEndPoint(asm, 0); err != nil {
		return err
	}

	// Loads the function selector of call function.
	asm.Emerge(opcode.LoadFunc)

	// Adds the logic to compare and find the corresponding function selector with the unmeaningful value.
	funcMap.Declare("FuncJmpr", *asm)
	for range c.Functions {
		if err := compileFuncSel(asm, abi.Selector(""), 0); err != nil {
			return err
		}
	}

	// No match to any function selector, Revert!
	funcMap.Declare("Revert", *asm)
	compileExit(asm)

	return nil
}

// Pushed the location of revert to exit the program.
func compileProgramEndPoint(asm *Asm, revertDst int) error {
	operand, err := encoding.EncodeOperand(revertDst)
	if err != nil {
		return err
	}
	asm.Emerge(opcode.Push, operand)

	return nil
}

// compileExit compiles exiting the program.
// If jumps to here, exit the program.
func compileExit(asm *Asm) {
	asm.Emerge(opcode.Exit)
}

// TODO: implement test cases :-)
// Generates a bytecode of memory size.
func compileMemSize(asm *Asm, tracer MemTracer) error {
	operand, err := encoding.EncodeOperand(tracer.MemSize())
	if err != nil {
		return err
	}

	if err := asm.ReplaceOperandAt(1, operand); err != nil {
		return err
	}

	return nil
}

// Generates a bytecode of function jumper.
func compileFuncJmpr(c ast.Contract, asm *Asm, funcMap FuncMap) error {
	funcJmpr := &Asm{
		AsmCodes: make([]AsmCode, 0),
	}

	// Pushes the location of revert.
	revertDst := funcMap[string(abi.Selector("Revert"))]
	if err := compileProgramEndPoint(funcJmpr, revertDst); err != nil {
		return err
	}

	// Loads the function selector of call function.
	funcJmpr.Emerge(opcode.LoadFunc)

	// Adds the logic to compare and find the corresponding function selector.
	for _, f := range c.Functions {
		selector := abi.Selector(f.Signature())
		funcDst := funcMap[string(selector)]

		if err := compileFuncSel(funcJmpr, selector, funcDst); err != nil {
			return err
		}
	}

	// No match to any function selector, Revert!
	compileExit(funcJmpr)

	// Replace expected function jumper with new function jumper.
	if err := fillFuncJmpr(asm, *funcJmpr); err != nil {
		return err
	}

	return nil
}

// Fill the function jumper in the location of function jumper placeholder.
func fillFuncJmpr(asm *Asm, funcJmpr Asm) error {
	if len(asm.AsmCodes) < len(funcJmpr.AsmCodes) {
		return fmt.Errorf("Can't fill the function jumper. Bytecode=%x, FuncJmpr=%x", asm.AsmCodes, funcJmpr.AsmCodes)
	}

	for i, asmCode := range funcJmpr.AsmCodes {
		asm.AsmCodes[i+3] = asmCode
	}

	return nil
}

// Compiles function jumper logic to find a function with its function selector
func compileFuncSel(asm *Asm, funcSel []byte, funcDst int) error {
	// Duplicates the function selector to find.
	asm.Emerge(opcode.DUP)
	// Pushes the function selector of this function literal.
	selector, err := encoding.EncodeOperand(funcSel)
	if err != nil {
		return err
	}
	asm.Emerge(opcode.Push, selector)
	asm.Emerge(opcode.EQ)
	asm.Emerge(opcode.NOT)
	// If the result is false (the condition is true), jump to the destination of function.
	dst, err := encoding.EncodeOperand(funcDst)
	if err != nil {
		return err
	}
	asm.Emerge(opcode.Push, dst)
	asm.Emerge(opcode.Jumpi)

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

// compileFunction() compiles a function in contract.
// Generates and adds output to bytecode.
func compileFunction(f ast.FunctionLiteral, bytecode *Asm, tracer *MemEntryTable) error {
	closedTracer := NewEnclosedMemEntryTable(tracer)
	for i, param := range f.Parameters {
		if err := compileParameter(*param, i, bytecode, closedTracer); err != nil {
			return err
		}
	}

	statements := f.Body.Statements
	for _, s := range statements {
		if err := compileStatement(s, bytecode, closedTracer); err != nil {
			return err
		}
	}

	tracer = closedTracer.Out()
	return nil
}

// compileParameter() compiles parameters in a function.
func compileParameter(p ast.ParameterLiteral, argNum int, bytecode *Asm, tracer MemTracer) error {
	entry := tracer.Define(p.Identifier.String())
	// Load an argument
	operand, err := encoding.EncodeOperand(argNum)
	if err != nil {
		return err
	}
	bytecode.Emerge(opcode.Push, operand)
	bytecode.Emerge(opcode.LoadArgs)
	// Push size of the argument
	size, err := encoding.EncodeOperand(entry.Size)
	if err != nil {
		return err
	}
	bytecode.Emerge(opcode.Push, size)
	// Push offset of the argument
	offset, err := encoding.EncodeOperand(entry.Offset)
	if err != nil {
		return err
	}
	bytecode.Emerge(opcode.Push, offset)
	// Save the argument in the memory
	bytecode.Emerge(opcode.Mstore)

	return nil
}

// compileStatement() compiles a statement in function.
// Generates and adds output to bytecode.
func compileStatement(s ast.Statement, bytecode *Asm, tracer MemTracer) error {
	switch statement := s.(type) {
	case *ast.AssignStatement:
		return compileAssignStatement(statement, bytecode, tracer)

	case *ast.ReturnStatement:
		return compileReturnStatement(statement, bytecode, tracer)

	case *ast.IfStatement:
		return compileIfStatement(statement, bytecode, tracer)

	case *ast.BlockStatement:
		return compileBlockStatement(statement, bytecode, tracer)

	case *ast.ExpressionStatement:
		return compileExpressionStatement(statement, bytecode, tracer)

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

	memEntry := tracer.Define(s.Variable.Name)
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

// compileReturnStatement compiles 'return' keyword
//
// PROTOCOL:
//   if return value of return statement is nil, then
//   return value zero
func compileReturnStatement(s *ast.ReturnStatement, asm *Asm, tracer MemTracer) error {

	var retVal ast.Expression

	if s.ReturnValue == nil {
		retVal = &ast.IntegerLiteral{Value: 0}
	} else {
		retVal = s.ReturnValue
	}

	if err := compileExpression(retVal, asm, tracer); err != nil {
		return err
	}

	asm.Emerge(opcode.Returning)

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
	// 'push <expression>

	asm.Emerge(opcode.Push, []byte(fmt.Sprintf("%d", -1)))
	// 'push <expression> push <-1(will be replaced)>'

	l1 := len(asm.AsmCodes)
	asm.Emerge(opcode.Jumpi)
	// 'push <expression> push <-1(will be replaced)> jumpi'
	if err := compileBlockStatement(s.Consequence, asm, tracer); err != nil {
		return err
	}
	// 'push <expression> push <-1(will be replaced)> jumpi <Consequence...>'
	asm.Emerge(opcode.Push, []byte(fmt.Sprintf("%d", -1)))
	l2 := len(asm.AsmCodes)
	// 'push <expression> push <-1(will be replaced)> jumpi <Consequence...> push <pc-to-end-of-Alternative>'
	asm.Emerge(opcode.Jump)
	// 'push <expression> push <-1(will be replaced)> jumpi <Consequence...> push <pc-to-end-of-Alternative> jump'

	if err := compileBlockStatement(s.Alternative, asm, tracer); err != nil {
		return err
	}

	// 'push <expression> push <pc-to-Alternative> jumpi <Consequence...> push <pc-to-end-of-Alternative> jump <Alternative...>'
	l3 := len(asm.AsmCodes)
	pc2al, err := encoding.EncodeOperand(l2 + 1)
	if err != nil {
		return err
	}
	asm.ReplaceOperandAt(l1-1, pc2al)

	pc2EndOfAlter, err := encoding.EncodeOperand(l3)
	if err != nil {
		return err
	}

	asm.ReplaceOperandAt(l2-1, pc2EndOfAlter)

	return nil
}

func compileIf(s *ast.IfStatement, asm *Asm, tracer MemTracer) error {
	// 'push <expression>

	asm.Emerge(opcode.Push, []byte(fmt.Sprintf("%d", -1)))
	// 'push <expression> push <-1(will be replaced)>'

	l1 := len(asm.AsmCodes)
	asm.Emerge(opcode.Jumpi)
	// 'push <expression> push <-1(will be replaced)> jumpi'
	if err := compileBlockStatement(s.Consequence, asm, tracer); err != nil {
		return err
	}
	// 'push <expression> push <-1(will be replaced)> jumpi <Consequence...>'

	l2 := len(asm.AsmCodes)
	pc2al, err := encoding.EncodeOperand(l2)
	if err != nil {
		return err
	}

	asm.ReplaceOperandAt(l1-1, pc2al)
	// 'push <expression> push <pc-to-end-of-Consequence> jumpi <Consequence...>'

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
	if err := compileExpression(s.Expr, bytecode, tracer); err != nil {
		return err
	}

	// Clear the stack.
	bytecode.Emerge(opcode.Pop)

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

	default:
		return fmt.Errorf("Undefined operator %s", e.Operator.String())
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
	memEntry, err := tracer.Entry(e.Name)
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
