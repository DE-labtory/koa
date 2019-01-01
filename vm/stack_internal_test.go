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

import (
	"testing"
)

func TestStack_new(t *testing.T) {
	stack := newStack()
	if cap(stack.items) != stackMaxSize {
		t.Fatalf("Stack initializing failed")
	}
}

func TestStack_pushN(t *testing.T) {
	stack := newStack()

	stack.pushN(item(1), item(2), item(3))

	if len(stack.items) != 3 {
		t.Fatalf("PushN did not work properly")
	}
}

func TestStack_push(t *testing.T) {
	stack := newStack()

	stack.push(item(1))
	stack.push(item(2))
	stack.push(item(3))
	if len(stack.items) != 3 {
		t.Fatalf("Push did not work properly")
	}
}

func TestStack_pop(t *testing.T) {
	stack := newStack()

	stack.push(item(1))
	stack.push(item(2))
	stack.push(item(3))

	for i := len(stack.items); i > 0; i-- {
		if stack.pop() != item(i) {
			t.Fatalf("Pop error at stack point %d", i)
		}
	}

	if len(stack.items) != 0 {
		t.Fatal("Stack has some datas")
	}
}

func TestStack_len(t *testing.T) {
	stack := newStack()

	stack.push(item(1))
	stack.push(item(2))
	stack.push(item(3))

	if len(stack.items) != stack.len() {
		t.Fatal("function len is invalid")
	}
}

func TestStack_print(t *testing.T) {
	stack := newStack()

	stack.push(item(1))
	stack.push(item(2))
	stack.push(item(3))

	//	### stack ###
	//	0    1
	//	1    2
	//	2    3
	//	#############

	stack.print()
}
