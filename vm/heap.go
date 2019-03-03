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

type Hid int
type Hvalue []byte

type DataMap map[Hid]Hvalue

type Heap struct {
	data DataMap
}

func NewHeap() *Heap {
	return &Heap{
		data: make(map[Hid]Hvalue),
	}
}

func (h *Heap) Set(id Hid, value Hvalue) {
	h.data[id] = value
}

func (h *Heap) Value(id Hid) Hvalue {
	return h.data[id]
}

func (h *Heap) Size() int {
	return len(h.data)
}
