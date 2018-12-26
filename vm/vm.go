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

package vm

type VM struct {
	// stack []object.Object
	sp int // Always points to the next value. Top of stack is stack[sp-1]
}

func NewVM() *VM {
	return nil
}

func (vm *VM) Run() error {
	return nil
}

// push() will interact with stack
func (vm *VM) push() error {
	return nil
}

// pop() will interact with stack
func (vm *VM) pop() error {
	return nil
}
