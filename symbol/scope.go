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
	"sort"
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
func (s *Scope) Get(name string) Symbol {
	var scope = s
	for scope != nil {
		obj, ok := scope.store[name]
		if ok {
			return obj
		}
		scope = scope.outer
	}
	return nil
}

// Setter set a variable to target scope's map
func (s *Scope) Set(name string, val Symbol) Symbol {
	s.store[name] = val
	return val
}

func (s *Scope) String() string {
	var out bytes.Buffer
	scope := s

	for scope != nil {
		// sort scope.store's keys alphabetically
		keys := make([]string, 0)
		for k := range scope.store {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		out.WriteString("[ Scope ]\n")
		for _, k := range keys {
			out.WriteString(fmt.Sprintf("key:%s, symbol:%s\n",
				k, scope.store[k].String()))
		}

		scope = scope.outer
	}
	return out.String()
}
