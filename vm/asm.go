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

// Converts rawByteCode to assembly code.
func assemble(rawByteCode []byte) (*asm, error) {
	analysis()
	return &asm{}, nil
}

// Do some analysis step (calculating the cost of running the code)
func analysis() {

}

// Assemble Reader read assembly codes and can jump to certain assembly code
type asmReader interface {
	next() hexer
	jump(i uint64)
}

type hexer interface {
	hex() []uint8
}

type asm struct {
	code  []hexer
	cost  uint64
	index uint64
}

func (a *asm) next() hexer {
	return nil
}

func (a *asm) jump(pc uint64) {

}

func (a *asm) print() {

}
