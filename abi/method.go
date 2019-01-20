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

package abi

type Method struct {
	Name      string
	Arguments Arguments
	Output    Argument
}

// Signature returns function's signature according to the ABI spec.
//
// Example
// function foo(uint32 a, int b) = "foo(uint32,int256)"
func (method Method) Signature() string {
	return ""
}

// ID return function's id using function selector
func (method Method) ID() []byte {
	return nil
}
