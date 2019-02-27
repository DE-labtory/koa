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

var objTypeMap = map[ast.DataStructure]ObjectType{
	ast.IntType:    IntegerObject,
	ast.StringType: StringObject,
	ast.BoolType:   BooleanObject,
	ast.VoidType:   VoidObject,
}

// Resolver traverse AST and resolve symbols, (1) check symbol's scope
// (2) check type of symbol
type Resolver struct {
	scope *Scope

	// types manage expression its own type
	types map[ast.Expression]ObjectType

	// defs manage identifier its own object
	defs map[*ast.Identifier]Object

	// fns manager function its own type
	fns map[*ast.FunctionLiteral]ObjectType
}

// Define scope error
type DupError struct {
	object Object
}

func (e DupError) Error() string {
	return fmt.Sprintf("%s is already defined.", e.object.String())
}

// Define data type error
type TypeError struct {
	target ObjectType
	object ObjectType
}

func (e TypeError) Error() string {
	return fmt.Sprintf("Type mismatched : %s, %s", e.target, e.object)
}

// Define not exists error
type NoExistError struct {
	identifier *ast.Identifier
}

func (e NoExistError) Error() string {
	return fmt.Sprintf("Object not exists : %s", e.identifier.Name)
}

// Define scope structure
type Scope struct {
	store map[string]Object

	child  []*Scope
	parent *Scope
}

func (s *Scope) Get(name string) Object {
	scope := s
	for scope.parent != nil {
		if obj, ok := scope.store[name]; ok {
			return obj
		}
		scope = scope.parent
	}

	if obj, ok := scope.store[name]; ok {
		return obj
	}
	scope = scope.parent

	return nil
}

func NewResolver() *Resolver {
	s := NewScope()
	return &Resolver{
		scope: s,
		types: make(map[ast.Expression]ObjectType),
		defs:  make(map[*ast.Identifier]Object),
		fns:   make(map[*ast.FunctionLiteral]ObjectType),
	}
}

// typeof returns the type of expression exp,
// or return InvalidSymbol if not found
func (r *Resolver) typeOf(exp ast.Expression) ObjectType {
	if st, ok := r.types[exp]; ok {
		return st
	}
	return InvalidSymbol
}

// objectOf returns the object denoted by the specified identifier
// or return nil if not found
func (r *Resolver) objectOf(id *ast.Identifier) Object {
	if sym, ok := r.defs[id]; ok {
		return sym
	}
	return nil
}

func NewScope() *Scope {
	return &Scope{
		store:  make(map[string]Object),
		parent: nil,
	}
}

func NewEnclosedScope(p *Scope) *Scope {
	s := &Scope{
		store:  make(map[string]Object),
		child:  make([]*Scope, 0),
		parent: p,
	}
	p.child = append(p.child, s)
	return s
}

func ResolveContract(c *ast.Contract) error {
	r := NewResolver()

	for _, f := range c.Functions {
		name := f.Name.Name
		obj := &Function{Name: name}

		r.scope.store[name] = obj
		r.fns[f] = FunctionObject
		r.defs[f.Name] = obj
	}

	for _, f := range c.Functions {
		if err := ResolveFunction(f, r); err != nil {
			return err
		}
	}
	return nil
}

func ResolveFunction(f *ast.FunctionLiteral, r *Resolver) error {
	r.scope = NewEnclosedScope(r.scope)

	if err := ResolveFunctionParameter(f.Parameters, r); err != nil {
		return nil
	}

	if err := ResolveFunctionBody(f.Body, r); err != nil {
		return err
	}

	r.scope = r.scope.parent
	return nil
}

func ResolveFunctionParameter(p []*ast.ParameterLiteral, r *Resolver) error {
	for _, p := range p {
		if obj, ok := r.scope.store[p.Identifier.Name]; ok {
			return DupError{
				object: obj,
			}
		}

		var obj Object
		var t ObjectType
		switch p.Type {
		case ast.IntType:
			obj = &Integer{Name: p.Identifier}
			t = IntegerObject
		case ast.StringType:
			obj = &String{Name: p.Identifier}
			t = StringObject
		case ast.BoolType:
			obj = &Boolean{Name: p.Identifier}
			t = BooleanObject
		}
		r.scope.store[p.Identifier.Name] = obj
		r.types[p] = t
		r.defs[p.Identifier] = obj
	}
	return nil
}

func ResolveFunctionBody(b *ast.BlockStatement, r *Resolver) error {
	if err := ResolveBlockStatement(b, r); err != nil {
		return err
	}
	return nil
}

func ResolveStatement(stmt ast.Statement, r *Resolver) error {
	var err error
	switch implStmt := stmt.(type) {
	case *ast.AssignStatement:
		err = ResolveAssignStatement(implStmt, r)
	case *ast.BlockStatement:
		err = ResolveBlockStatement(implStmt, r)
	case *ast.ExpressionStatement:
		// TODO
	case *ast.IfStatement:
		err = ResolveIfStatement(implStmt, r)
	}

	if err != nil {
		return err
	}
	return nil
}

func ResolveAssignStatement(aStmt *ast.AssignStatement, r *Resolver) error {
	// check if variable which has same name exists
	if obj, ok := r.scope.store[aStmt.Variable.Name]; ok {
		return DupError{
			object: obj,
		}
	}

	name := aStmt.Variable.Name
	switch aStmt.Type {
	case ast.IntType:
		r.scope.store[name] = &Integer{Name: &ast.Identifier{Name: name}}
	case ast.BoolType:
		r.scope.store[name] = &Boolean{Name: &ast.Identifier{Name: name}}
	case ast.StringType:
		r.scope.store[name] = &String{Name: &ast.Identifier{Name: name}}
	}

	ot, err := ResolveExpression(aStmt.Value, r)
	if err != nil {
		return err
	}
	if ot != objTypeMap[aStmt.Type] {
		return TypeError{
			target: objTypeMap[aStmt.Type],
			object: ot,
		}
	}

	switch aStmt.Type {
	case ast.IntType:
		r.types[aStmt.Value] = IntegerObject
		r.defs[&aStmt.Variable] = &Integer{Name: &aStmt.Variable}
	case ast.BoolType:
		r.types[aStmt.Value] = BooleanObject
		r.defs[&aStmt.Variable] = &Boolean{Name: &aStmt.Variable}
	case ast.StringType:
		r.types[aStmt.Value] = StringObject
		r.defs[&aStmt.Variable] = &String{Name: &aStmt.Variable}
	}

	return nil
}

func ResolveBlockStatement(bStmt *ast.BlockStatement, r *Resolver) error {
	r.scope = NewEnclosedScope(r.scope)
	for _, stmt := range bStmt.Statements {
		if err := ResolveStatement(stmt, r); err != nil {
			return err
		}
	}
	r.scope = r.scope.parent
	return nil
}

func ResolveIfStatement(ifStmt *ast.IfStatement, r *Resolver) error {
	if ifStmt.Consequence != nil {
		if err := ResolveBlockStatement(ifStmt.Consequence, r); err != nil {
			return err
		}
	}
	if ifStmt.Alternative != nil {
		if err := ResolveBlockStatement(ifStmt.Alternative, r); err != nil {
			return err
		}
	}
	return nil
}

// ResolveIdentifier checks whether target identifier is in the scope.
func ResolveIdentifier(identifier *ast.Identifier, r *Resolver) (ObjectType, error) {
	if obj := r.scope.Get(identifier.Name); obj != nil {
		return obj.Type(), nil
	}

	return InvalidSymbol, NoExistError{
		identifier: identifier,
	}
}

func ResolveExpression(exp ast.Expression, r *Resolver) (ObjectType, error) {
	switch implExp := exp.(type) {
	case *ast.PrefixExpression:
		return ResolvePrefixExpression(implExp, r)
	case *ast.InfixExpression:
		return ResolveInfixExpression(implExp, r)
	case *ast.Identifier:
		return ResolveIdentifier(implExp, r)
	case *ast.CallExpression:
		return ResolveCallExpression(implExp, r)
	case *ast.IntegerLiteral:
		return IntegerObject, nil
	case *ast.BooleanLiteral:
		return BooleanObject, nil
	case *ast.StringLiteral:
		return StringObject, nil
	}
	return InvalidSymbol, nil
}

func ResolveCallExpression(exp *ast.CallExpression, r *Resolver) (ObjectType, error) {
	return InvalidSymbol, nil
}

func ResolvePrefixExpression(exp ast.Expression, r *Resolver) (ObjectType, error) {
	return InvalidSymbol, nil
}

func ResolveInfixExpression(exp *ast.InfixExpression, r *Resolver) (ObjectType, error) {
	if ot, err := ResolveExpression(exp.Left, r); err != nil {
		return ot, err
	}

	if ot, err := ResolveExpression(exp.Right, r); err != nil {
		return ot, err
	}

	return InvalidSymbol, nil
}
