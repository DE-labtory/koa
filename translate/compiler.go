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

import "github.com/DE-labtory/koa/ast"

// TODO: parameter should be ast.Contract
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

		// TODO: use abi.ExtractAbi(f)
	}

	if err := generateFuncJumper(bytecode); err != nil {
		return *bytecode, err
	}

	return *bytecode, nil
}

// TODO: implement me w/ test cases :-)
// Generates a bytecode of function jumper.
func generateFuncJumper(bytecode *Bytecode) error {
	return nil
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

// TODO: implement me w/ test cases :-)
func compileBlockStatement(s *ast.BlockStatement, bytecode *Bytecode) error {
	return nil
}
