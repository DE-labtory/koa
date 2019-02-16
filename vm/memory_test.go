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
	"bytes"
	"testing"
)

func TestMemory_New(t *testing.T) {

}

func TestMemory_Set(t *testing.T) {
	memory := NewMemory()

	tests := []struct {
		offset uint64
		data   []byte
	}{
		{0x00, []byte{0x00}},
		{0x01, []byte{0x01}},
		{0x02, []byte{0x02}},
		{0x03, []byte{0x03}},
	}

	testsExpected := []struct {
		offset uint64
		data   []byte
		size   uint64
	}{
		{0x00, []byte{0x00}, 1},
		{0x01, []byte{0x01}, 1},
		{0x02, []byte{0x02}, 1},
		{0x03, []byte{0x03}, 1},
	}

	memory.Resize(8)

	for _, test := range tests {
		memory.Set(test.offset, test.data[0])
	}

	for _, expected := range testsExpected {
		data := memory.GetVal(expected.offset, expected.size)

		if !bytes.Equal(expected.data, data) {
			t.Error("Invalid memory value")
		}
	}
}

func TestMemory_Set8(t *testing.T) {
	memory := NewMemory()

	tests := []struct {
		offset uint64
		data   []byte
	}{
		{0x00, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}},
		{0x08, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}},
		{0x10, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}},
		{0x18, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}},
	}

	testsExpected := []struct {
		offset uint64
		data   []byte
		size   uint64
	}{
		{0x00, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}, 8},
		{0x08, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}, 8},
		{0x10, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}, 8},
		{0x18, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}, 8},
	}

	memory.Resize(32)

	for _, test := range tests {
		memory.Set8(test.offset, test.data)
	}

	for _, expected := range testsExpected {
		data := memory.GetVal(expected.offset, expected.size)

		if !bytes.Equal(expected.data, data) {
			t.Error("Invalid memory value")
		}
	}
}

func TestMemory_Set8_dynamic(t *testing.T) {
	memory := NewMemory()

	tests := []struct {
		offset uint64
		data   []byte
	}{
		{0x00, []byte{0x00, 0x01, 0x02, 0x03}},
		{0x08, []byte{0x00, 0x01, 0x02}},
		{0x10, []byte{0x00, 0x01}},
		{0x18, []byte{0x00}},
	}

	testsExpected := []struct {
		offset uint64
		data   []byte
		size   uint64
	}{
		{0x00, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}, 8},
		{0x08, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02}, 8},
		{0x10, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, 8},
		{0x18, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 8},
	}

	memory.Resize(32)

	for _, test := range tests {
		memory.Set8(test.offset, test.data)
	}

	for _, expected := range testsExpected {
		data := memory.GetVal(expected.offset, expected.size)

		if !bytes.Equal(expected.data, data) {
			t.Error("Invalid memory value")
		}
	}
}

func TestMemory_Sets(t *testing.T) {
	memory := NewMemory()

	tests := []struct {
		offset uint64
		data   []byte
		size   uint64
	}{
		{0x00, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 8},
		{0x08, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}, 7},
		{0x0F, []byte{0x00, 0x01, 0x02, 0x03, 0x04}, 5},
		{0x14, []byte{0x00, 0x01}, 2},
	}

	testsExpected := []struct {
		offset uint64
		data   []byte
		size   uint64
	}{
		{0x00, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 8},
		{0x08, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}, 7},
		{0x0F, []byte{0x00, 0x01, 0x02, 0x03, 0x04}, 5},
		{0x14, []byte{0x00, 0x01}, 2},
	}

	memory.Resize(32)

	for _, test := range tests {
		memory.Sets(test.offset, test.size, test.data)
	}

	for _, expected := range testsExpected {
		data := memory.GetVal(expected.offset, expected.size)

		if !bytes.Equal(expected.data, data) {
			t.Error("Invalid memory value")
		}
	}
}

func TestMemory_GetPtr(t *testing.T) {
	memory := NewMemory()

	test := struct {
		offset uint64
		data   []byte
		size   uint64
	}{
		0x00,
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03},
		8,
	}

	testModified := struct {
		offset uint64
		data   []byte
		size   uint64
	}{
		0x00,
		[]byte{0x00, 0x00, 0x00, 0x00, 0x04, 0x05, 0x06, 0x07},
		8,
	}

	testExpected := struct {
		offset uint64
		data   []byte
		size   uint64
	}{
		0x00,
		[]byte{0x00, 0x00, 0x00, 0x00, 0x04, 0x05, 0x06, 0x07},
		8,
	}

	memory.Resize(32)

	memory.Set8(test.offset, test.data)

	data := memory.GetPtr(test.offset, test.size)
	if !bytes.Equal(test.data, data) {
		t.Error("Invalid memory value")
	}

	memory.Sets(testModified.offset, testModified.size, testModified.data)
	if !bytes.Equal(testExpected.data, data) {
		t.Error("Invalid memory value")
	}

}

func TestMemory_Resize(t *testing.T) {
	memory := NewMemory()

	memory.Resize(32)
	if cap(memory.data) != 32 {
		t.Error("Invalid memory size")
	}

	memory.Resize(64)
	if cap(memory.data) != 64 {
		t.Error("Invalid memory size")
	}

	// Can't resize less than before
	memory.Resize(32)
	if cap(memory.data) != 64 {
		t.Error("Invalid memory size")
	}
}

func TestMemory_Len(t *testing.T) {
	memory := NewMemory()

	test := struct {
		offset uint64
		data   []byte
		size   uint64
	}{
		0x00,
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03},
		8,
	}

	memory.Resize(8)

	memory.Set8(test.offset, test.data)

	if memory.Len() != len(memory.data) {
		t.Error("Invalid memory size")
	}
}

func TestMemory_Data(t *testing.T) {
	memory := NewMemory()

	testBytes := []struct {
		offset uint64
		data   []byte
		size   uint64
	}{
		{0x00, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}, 8},
		{0x08, []byte{0x00, 0x00, 0x00, 0x00, 0x04, 0x05, 0x06, 0x07}, 8},
		{0x10, []byte{0x00, 0x00, 0x00, 0x00, 0x08, 0x09, 0x0A, 0x0B}, 8},
		{0x18, []byte{0x00, 0x00, 0x00, 0x00, 0x0C, 0x0D, 0x0E, 0x0F}, 8},
	}

	testExpected := struct {
		offset uint64
		data   []byte
		size   uint64
	}{
		0x00,
		[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
		16,
	}

	memory.Resize(32)

	for _, test := range testBytes {
		memory.Set8(test.offset, test.data)
	}

	if bytes.Equal(memory.Data(), testExpected.data) {
		t.Error("Invalid memory data")
	}
}
