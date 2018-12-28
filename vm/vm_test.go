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
	"testing"

	"github.com/DE-labtory/koa/opcode"
	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	testInstruction := []Item{Item(opcode.Push), 1, Item(opcode.Push), 2}
	testExpected := []Item{1, 2}

	stack, err := Execute(testInstruction)

	assert.NoError(t, err)
	for i, item := range stack.item {
		assert.Equal(t, testExpected[i], item)
	}
}

func TestAdd(t *testing.T) {
	testInstruction := []Item{
		Item(opcode.Push), 1,
		Item(opcode.Push), 2,
		Item(opcode.Add),
	}
	testExpected := []Item{3}

	stack, err := Execute(testInstruction)

	assert.NoError(t, err)
	assert.Equal(t, testExpected[0], stack.Pop())
}
