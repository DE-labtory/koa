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

// TODO: implement me w/ test cases :-)
// Define() saves an variable to EntryMap and increase the MemoryCounter.
func (m *MemoryTable) Define(id string, value interface{}) {}

// TODO: implement me w/ test cases :-)
// Use() increase the MemoryCounter.
// When saves any data to memory, you should use Use() function.
// Then, it will return the position of that data.
func (m *MemoryTable) Use(size uint) uint {
	return 0
}
