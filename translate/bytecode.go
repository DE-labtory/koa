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

// Bytecode is generated by compiling.
type Bytecode struct {
	RawByte []byte
	AsmCode []string
}

// Emerge() translates instruction to bytecode
// An operand of operands should be 4 bytes.
func (b *Bytecode) Emerge(operator opcode.Type, operands ...[]byte) int {
	// Translate operator to byte
	b.RawByte = append(b.RawByte, byte(operator))

	// Translate operator to assembly
	opStr, err := operator.String()
	if err != nil {
		return 0
	}
	b.AsmCode = append(b.AsmCode, opStr)

	for _, o := range operands {
		// Translate operands to byte
		b.RawByte = append(b.RawByte, o...)

		// Translate operands to assembly
		operand := fmt.Sprintf("%x", o)
		b.AsmCode = append(b.AsmCode, operand)
	}

	// Returns next bytecode position
	return len(b.AsmCode)
}
