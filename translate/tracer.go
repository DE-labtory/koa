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
	"github.com/DE-labtory/koa/encoding"
)

// MemoryTableEntry saves size and offset of the value which the variable has.
type MemoryTableEntry struct {
	Offset uint
	Size   uint
}

// MemoryTable is used to know the location of the memory
type MemoryTable struct {
	EntryMap      map[string]MemoryTableEntry
	MemoryCounter uint
}

func NewMemoryTable() *MemoryTable {
	return &MemoryTable{
		EntryMap:      make(map[string]MemoryTableEntry),
		MemoryCounter: 0,
	}
}

// Define() saves an variable to EntryMap and increase the MemoryCounter.
// ex)
// a = 5 -> Define("a", 5)
// b = "abc" -> Define("b", "abc")
func (m *MemoryTable) Define(id string, value interface{}) (MemoryTableEntry, error) {
	entry := MemoryTableEntry{
		Offset: m.MemoryCounter,
	}

	encodedValue, err := encoding.EncodeOperand(value)
	if err != nil {
		return entry, err
	}

	size := uint(len(encodedValue))
	entry.Size = size
	m.MemoryCounter += size

	m.EntryMap[id] = entry

	return entry, nil
}
