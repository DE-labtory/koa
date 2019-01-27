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
)

type SymbolType string

const (
	IntegerSymbol  = "INTEGER"
	BooleanSymbol  = "BOOLEAN"
	StringSymbol   = "STRING"
	FunctionSymbol = "FUNCTION"
)

type Symbol interface {
	Type() SymbolType
	String() string
}

// Represent Integer symbol
type Integer struct {
	Value int64
}

func (i *Integer) Type() SymbolType {
	return IntegerSymbol
}

func (i *Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

// Represent Boolean Object
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() SymbolType {
	return BooleanSymbol
}

func (b *Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

// Represent String Object
type String struct {
	Value string
}

func (s *String) Type() SymbolType {
	return StringSymbol
}

func (s *String) String() string {
	return fmt.Sprintf("%s", s.Value)
}

// Represent Function symbol
// Name represents function's name.
// Scope represents function value's scope.
type Function struct {
	name  string
	Scope *Scope
}

func (f *Function) Type() SymbolType {
	return FunctionSymbol
}

func (f *Function) String() string {
	return fmt.Sprintf("func %s", f.name)
}
