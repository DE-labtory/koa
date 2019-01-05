// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

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
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
)

var ErrInvalidMemory = errors.New("Invalid memory reference")

type Memory struct {
	data []byte
	cost uint32
}

func NewMemory() *Memory {
	return &Memory{
		data: make([]byte, 0),
		cost: 0,
	}
}

// Set sets offset to value
func (m *Memory) Set(offset uint32, value byte) {
	if offset > uint32(m.Len()) {
		panic(ErrInvalidMemory)
	}
	m.data[offset] = value
}

// Set4 sets the 4 bytes starting at offset to the value of val, left-padded with zeroes to
// 4 bytes.
func (m *Memory) Set4(offset uint32, value []byte) {
	tmp := new(big.Int)
	tmp.SetBytes(value)

	if offset+4 > uint32(m.Len()) {
		panic(ErrInvalidMemory)
	}
	copy(m.data[offset:offset+4], []byte{0, 0, 0, 0})
	math.ReadBits(tmp, m.data[offset:offset+4])
}

// Sets sets offset + size to value
func (m *Memory) Sets(offset, size uint32, value []byte) {
	if size > 0 {
		if offset+size > uint32(m.Len()) {
			panic(ErrInvalidMemory)
		}
		copy(m.data[offset:offset+size], value)
	}
}

// Get returns offset + size as a new slice
func (m *Memory) GetVal(offset, size uint32) []byte {
	if size == 0 {
		return nil
	}

	if uint32(m.Len()) > offset {
		cpy := make([]byte, size)
		copy(cpy, m.data[offset:offset+size])

		return cpy
	}

	return nil
}

// GetPtr returns the offset + size
func (m *Memory) GetPtr(offset, size uint32) []byte {
	if size == 0 {
		return nil
	}

	if len(m.data) > int(offset) {
		return m.data[offset : offset+size]
	}

	return nil
}

// Resize resizes the memory to size
func (m *Memory) Resize(size uint32) {
	if uint32(m.Len()) < size {
		m.data = append(m.data, make([]byte, size-uint32(m.Len()))...)
	}
}

func (m *Memory) Len() int {
	return len(m.data)
}

func (m *Memory) Data() []byte {
	return m.data
}

func (m *Memory) Cost() uint32 {
	return m.cost
}

// Print dumps the content of the memory.
func (m *Memory) Print() {
	fmt.Printf("### mem %d bytes ###\n", len(m.data))
	if len(m.data) > 0 {
		addr := 0
		for i := 0; i+4 <= len(m.data); i += 4 {
			fmt.Printf("%03d: % x\n", addr, m.data[i:i+4])
			addr++
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("####################")
}
