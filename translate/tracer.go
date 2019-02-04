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

import "github.com/pkg/errors"

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

func New() *MemoryTable {
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

	switch v := value.(type) {
	case int:
		entry.Size = 8
		m.MemoryCounter += 8

	case bool:
		entry.Size = 8
		m.MemoryCounter += 8

	case string:
		size := uint(len(v))
		entry.Size = size
		m.MemoryCounter += size

	default:
		return entry, errors.New("Not defined type definition")
	}

	m.EntryMap[id] = entry

	return entry, nil
}
