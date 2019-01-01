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

import "fmt"

const (
	stackMaxSize = 1024
)

type item uint32

// Stack is an object for basic stack operations. Items popped to the stack are
// expected not to be changed and modified.
type stack struct {
	items []item
}

func newStack() *stack {
	return &stack{items: make([]item, 0, stackMaxSize)}
}

// TODO: implement me w/ test cases :-)
func (s *stack) push(d item) {
	s.items = append(s.items, d)
}

// TODO: implement me w/ test cases :-)
// Push n number of data([]data) to stack
func (s *stack) pushN(ds ...item) {

}

// TODO: implement me w/ test cases :-)
func (s *stack) pop() item {
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

// TODO: implement me w/ test cases :-)
func (s *stack) len() int {
	return len(s.items)
}

// TODO: implement me w/ test cases :-)
// Print dumps the content of the stack
func (s *stack) print() {
	fmt.Println("### stack ###")
	if len(s.items) > 0 {
		for i, val := range s.items {
			fmt.Printf("%-3d  %v\n", i, val)
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("#############")
}
