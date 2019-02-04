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

package translate_test

import (
	"testing"

	"github.com/DE-labtory/koa/encoding"
	"github.com/DE-labtory/koa/translate"
)

func TestMemoryTable_Define(t *testing.T) {
	tests := []struct {
		id           string
		value        interface{}
		expectedSize uint
		err          error
	}{
		{
			id:           "aInteger",
			value:        int64(1),
			expectedSize: 8,
			err:          nil,
		},
		{
			id:           "aBoolean",
			value:        true,
			expectedSize: 8,
			err:          nil,
		},
		{
			id:           "aString",
			value:        "abc",
			expectedSize: 3,
			err:          nil,
		},
		{
			id:           "aNotDefined",
			value:        []byte{01, 02, 03},
			expectedSize: 0,
			err: encoding.EncodeError{
				Operand: []byte{01, 02, 03},
			},
		},
	}

	mTable := translate.NewMemoryTable()

	for i, test := range tests {
		prevOffset := mTable.MemoryCounter

		entry, err := mTable.Define(test.id, test.value)

		if err != nil && err.Error() != test.err.Error() {
			t.Fatalf("test[%d] - Define() error wrong. expected=%v, err=%v", i, test.err, err)
		}

		if entry.Size != test.expectedSize {
			t.Fatalf("test[%d] - Define() result wrong. expected=%d, got=%d", i, test.expectedSize, entry.Size)
		}

		if entry.Offset != prevOffset {
			t.Fatalf("test[%d] - Define() result wrong. expected=%d, got=%d", i, prevOffset, entry.Offset)
		}

		if err == nil && mTable.EntryMap[test.id] != entry {
			t.Fatalf("test[%d] - Define() result wrong. expected=%v, got=%v", i, mTable.EntryMap[test.id], entry)
		}

		expectedMemoryCounter := prevOffset + test.expectedSize
		if mTable.MemoryCounter != expectedMemoryCounter {
			t.Fatalf("test]%d] - Define() result wrong. expected=%d, got=%d", i, expectedMemoryCounter, mTable.MemoryCounter)
		}
	}
}
