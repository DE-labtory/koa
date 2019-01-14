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

package translate

import (
	"fmt"

	"github.com/DE-labtory/koa/opcode"
)

type Bytecode struct {
	// Raw byte code like 0x10211234...
	RawByte []byte
	// Assemble code like Push 1 Push 2 Add 2 3 Pop ...
	AsmCode []string
}

// Emit() adds raw byte code and assemble code to Bytecode.
// Then, returns the position of this instruction.
// TODO: return pc position w/ test cases :-)
func (b *Bytecode) Emit(operator opcode.Type, operands ...[]byte) (uint64, error) {
	b.RawByte = append(b.RawByte, byte(operator))

	opStr, err := operator.ToString()
	if err != nil {
		return 0, err
	}
	b.AsmCode = append(b.AsmCode, opStr)

	for _, o := range operands {
		b.RawByte = append(b.RawByte, o...)

		operand := fmt.Sprintf("%x", o)
		b.AsmCode = append(b.AsmCode, operand)
	}

	return 0, nil
}
