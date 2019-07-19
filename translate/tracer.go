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

package translate

import (
	"fmt"
)

const EntrySize = 8

type EntryError struct {
	Id string
}

func (e EntryError) Error() string {
	return fmt.Sprintf("[%s] definition doesn't exist", e.Id)
}

type MemTracer interface {
	MemDefiner
	MemGetter
}

// Define() saves an variable to EntryMap and increase the MemoryCounter.
// This should be used when compiles the assign statement.
// ex)
// a = 5 -> Define("a", 5)
// b = "abc" -> Define("b", "abc")
type MemDefiner interface {
	Define(id string) MemEntry
}

// MemGetter gets the data of the memory entry.
// GetOffsetOfEntry() returns the offset of the memory entry corresponding the Id.
// GetSizeOfEntry() returns the size of the memory entry corresponding the Id.
type MemGetter interface {
	Entry(id string) (MemEntry, error)
	Counter() int
	MemSize() int
}

// MemEntry saves size and offset of the value which the variable has.
type MemEntry struct {
	Offset int
	Size   int
}

// MemEntryTable is used to know the location of the memory
type MemEntryTable struct {
	EntryMap      map[string]MemEntry
	MemoryCounter int
	Outer         *MemEntryTable
}

func NewMemEntryTable() *MemEntryTable {
	return &MemEntryTable{
		EntryMap:      make(map[string]MemEntry),
		MemoryCounter: 0,
		Outer:         nil,
	}
}

func NewEnclosedMemEntryTable(memEntryTable *MemEntryTable) *MemEntryTable {
	m := NewMemEntryTable()
	m.Outer = memEntryTable
	m.MemoryCounter = memEntryTable.MemoryCounter
	return m
}

func (m *MemEntryTable) Out() *MemEntryTable {
	m.Outer.MemoryCounter = m.MemoryCounter
	return m.Outer
}

func (m *MemEntryTable) Define(id string) MemEntry {
	entry := MemEntry{
		Offset: m.MemoryCounter,
	}

	entry.Size = EntrySize
	m.MemoryCounter += EntrySize
	m.EntryMap[id] = entry

	return entry
}

func (m MemEntryTable) Entry(id string) (MemEntry, error) {
	entry, ok := m.EntryMap[id]
	if !ok {
		return MemEntry{}, EntryError{
			Id: id,
		}
	}

	return entry, nil
}

// TODO: implement test cases :-)
func (m MemEntryTable) Counter() int {
	return m.MemoryCounter
}

// TODO: implement test cases :-)
func (m MemEntryTable) MemSize() int {
	return m.MemoryCounter
}
