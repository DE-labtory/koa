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

package lex_test

import (
	"testing"

	"github.com/DE-labtory/koa/cmd/lex"
)

func TestCmd(t *testing.T) {
	cmd := lex.Cmd()

	if cmd.Name != "lex" {
		t.Errorf("wrong name - expected=lex, got=%s", cmd.Name)
	}

	if len(cmd.Aliases) != 1 {
		t.Errorf("length of aliases should be 1 - got=%d", len(cmd.Aliases))
	}

	if cmd.Aliases[0] != "l" {
		t.Errorf("value of aliases[0] should be 'l' - got=%s", cmd.Aliases[0])
	}

	if cmd.Usage != "koa lex [filePath]" {
		t.Errorf("usage should be 'koa lex [filePath]' - got=%s", cmd.Usage)
	}
}
