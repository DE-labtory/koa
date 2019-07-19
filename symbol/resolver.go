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

package symbol

import (
	"github.com/DE-labtory/koa/ast"
)

// Resolver traverse AST and resolve symbols, (1) check symbol's scope
// (2) check type of symbol
type Resolver struct {
	scope *Scope

	// types manage expression its own type
	types map[ast.Expression]SymbolType

	// defs manage identifier its own object
	defs map[*ast.Identifier]Symbol
}

func NewResolver() *Resolver {
	return &Resolver{
		scope: NewScope(),
		types: make(map[ast.Expression]SymbolType),
		defs:  make(map[*ast.Identifier]Symbol),
	}
}

// typeof returns the type of expression exp,
// or return InvalidSymbol if not found
func (r *Resolver) typeOf(exp ast.Expression) SymbolType {
	return ""
}

// objectOf returns the object denoted by the specified identifier
// or return nil if not found
func (r *Resolver) objectOf(id *ast.Identifier) Symbol {
	return nil
}

func (r *Resolver) ResolveContract(c *ast.Contract) error {
	return nil
}
