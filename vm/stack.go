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

package vm

import "fmt"

const (
	stackMaxSize = 1024
)

type item int64

// Stack is an object for basic Stack operations. Items popped to the Stack are
// expected not to be changed and modified.
type Stack struct {
	items []item
}

func newStack() *Stack {
	return &Stack{items: make([]item, 0, stackMaxSize)}
}

func (s *Stack) Push(d item) {
	s.items = append(s.items, d)
}

// Push n number of data([]data) to Stack
func (s *Stack) PushN(ds ...item) {
	s.items = append(s.items, ds...)
}

func (s *Stack) Pop() item {
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *Stack) Len() int {
	return len(s.items)
}

func (s *Stack) Dup() {
	s.Push(s.items[s.Len()-1])
}

func (s *Stack) Swap() {
	s.items[s.Len()-2], s.items[s.Len()-1] = s.items[s.Len()-1], s.items[s.Len()-2]
}

// Print dumps the content of the Stack
func (s *Stack) Print() {
	fmt.Println("### Stack ###")
	if len(s.items) > 0 {
		for i, val := range s.items {
			fmt.Printf("%-3d  %v\n", i, val)
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("#############")
}
