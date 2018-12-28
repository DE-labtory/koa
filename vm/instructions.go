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

type InstructionBuffer interface {
	GetItems(size int) []Item
	JumpTo(i int)
}

type Instructions struct {
	Buffer []Item
	ip     int
}

func NewInstructions(instructions []Item) Instructions {
	return Instructions{
		Buffer: instructions,
		ip:     0,
	}
}

func (inst *Instructions) GetItems(size int) []Item {
	item := make([]Item, 0)
	item = append(item, inst.Buffer[inst.ip+1:inst.ip+1+size]...)

	inst.ip += size
	return item
}

func (inst *Instructions) JumpTo(i int) {
	inst.ip += i
}
