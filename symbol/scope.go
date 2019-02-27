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

package symbol

import (
	"fmt"

	"github.com/DE-labtory/koa/ast"
)

// Define scope errors
type DupError struct {
	Identifier ast.Identifier
}

func (e DupError) Error() string {
	return fmt.Sprintf("%s is already defined.", e.Identifier.Name)
}

// Define scope structure
type Scope struct {
	store map[string]Symbol

	child  []*Scope
	parent *Scope
}

func NewScope() *Scope {
	return &Scope{
		store:  make(map[string]Symbol),
		parent: nil,
	}
}

func NewEnclosedScope(p *Scope) *Scope {
	s := &Scope{
		store:  make(map[string]Symbol),
		child:  make([]*Scope, 0),
		parent: p,
	}
	p.child = append(p.child, s)
	return s
}

func ScopingContract(c *ast.Contract, s *Scope) error {
	for _, f := range c.Functions {
		if err := ScopingFunction(f, s); err != nil {
			return err
		}
	}
	return nil
}

func ScopingFunction(f *ast.FunctionLiteral, s *Scope) error {
	s = NewEnclosedScope(s)
	if err := ScopingFunctionParameter(f.Parameters, s); err != nil {
		return nil
	}

	if err := ScopingFunctionBody(f.Body, s); err != nil {
		return err
	}

	return nil
}

func ScopingFunctionParameter(f []*ast.ParameterLiteral, s *Scope) error {
	for _, p := range f {
		switch p.Type {
		case ast.IntType:
			{
				s.store[p.Identifier.Name] = &Integer{Name: p.Identifier}
			}
		case ast.StringType:
			{
				s.store[p.Identifier.Name] = &String{Name: p.Identifier}
			}
		case ast.BoolType:
			{
				s.store[p.Identifier.Name] = &Boolean{Name: p.Identifier}
			}
		}
	}
	return nil
}

func ScopingFunctionBody(f *ast.BlockStatement, scope *Scope) error {
	if err := ScopingBlockStatement(f, scope); err != nil {
		return err
	}
	return nil
}

func ScopingStatement(stmt ast.Statement, scope *Scope) error {
	var err error
	switch implStmt := stmt.(type) {
	case *ast.AssignStatement:
		err = ScopingAssignStatement(implStmt, scope)
	case *ast.BlockStatement:
		err = ScopingBlockStatement(implStmt, scope)
	case *ast.ExpressionStatement:
		{

		}
	case *ast.IfStatement:
		{

		}
	}

	if err != nil {
		return err
	}
	return nil
}

func ScopingAssignStatement(stmt *ast.AssignStatement, scope *Scope) error {
	if ok := scope.store[stmt.Variable.Name]; ok != nil {
		return DupError{
			Identifier: stmt.Variable,
		}
	}

	if err := ScopingIdentifier(stmt.Variable, stmt.Type, scope); err != nil {
		return err
	}
	return nil
}

func ScopingBlockStatement(bStmt *ast.BlockStatement, scope *Scope) error {
	scope = NewEnclosedScope(scope)
	for _, stmt := range bStmt.Statements {
		if err := ScopingStatement(stmt, scope); err != nil {
			return err
		}
	}
	return nil
}

func ScopingIfStatement(iStmt *ast.IfStatement, scope *Scope) {

}

func ScopingIdentifier(identifier ast.Identifier, ds ast.DataStructure, scope *Scope) error {
	if ok := scope.store[identifier.Name]; ok != nil {
		return DupError{Identifier: identifier}
	}

	name := identifier.Name
	switch ds {
	case ast.IntType:
		scope.store[name] = &Integer{Name: &ast.Identifier{Name: name}}
	case ast.StringType:
		scope.store[name] = &String{Name: &ast.Identifier{Name: name}}
	case ast.BoolType:
		scope.store[name] = &Boolean{Name: &ast.Identifier{Name: name}}
	}

	return nil
}
