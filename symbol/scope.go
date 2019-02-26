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

import "github.com/DE-labtory/koa/ast"

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

func GenerateScope(c *ast.Contract) *Scope {
	s := NewScope()
	for _, f := range c.Functions {
		ScopingFunctionParameter(f, s)
		ScopingFunctionBody(f, s)
	}
	return nil
}

func ScopingFunctionParameter(f *ast.FunctionLiteral, s *Scope) {
	for _, p := range f.Parameters {
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
}

func ScopingFunctionBody(b *ast.FunctionLiteral, s *Scope) {
	for _, stmt := range b.Body.Statements {
		switch implStmt := stmt.(type) {
		case *ast.AssignStatement:
			{

			}
		case *ast.ReturnStatement:
			{

			}
		case *ast.BlockStatement:
			{

			}
		case *ast.ExpressionStatement:
			{

			}
		case *ast.ReassignStatement:
			{

			}
		case *ast.IfStatement:
			{

			}
		}
	}
}

func ScopingExpression(exp *ast.Expression, scope *Scope) {

}

func ScopingIdentifier(ident *ast.Identifier, scope *Scope) {

}
