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
	"bytes"
	"fmt"
)

// NewScope() makes a new scope.
func NewScope() *Scope {
	s := make(map[string]Symbol)
	return &Scope{s, nil}
}

// NewEnclosedScope makes a new scope which has outer scope.
func NewEnclosedScope(outer *Scope) *Scope {
	scope := NewScope()
	scope.outer = outer
	return scope
}

// Scope represent a variable's scope.
type Scope struct {
	store map[string]Symbol
	outer *Scope
}

// Getter return's variable's scope.
// If a scope doesn't have a variable,
// Program will search in outer scope.
func (s *Scope) Get(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	if !ok && s.outer != nil {
		obj, ok = s.outer.Get(name)
	}
	return obj, ok
}

// Setter set a variable to target scope's map
func (s *Scope) Set(name string, val Symbol) Symbol {
	s.store[name] = val
	return val
}

func (s *Scope) String() string {
	var out bytes.Buffer
	scope := s
	div := ""
	for scope != nil {
		out.WriteString(div)
		out.WriteString("[ Scope ] ")
		out.WriteString(fmt.Sprintln(scope.store))
		div += "  "
		scope = scope.outer
	}
	return out.String()
}
