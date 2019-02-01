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

package lex

import (
	"testing"

	"github.com/kami-zh/go-capturer"
)

func Test_lex(t *testing.T) {
	expectedOutput :=
		`[CONTRACT, contract]
[LBRACE, {]
[FUNCTION, func]
[IDENT, hello]
[LPAREN, (]
[RPAREN, )]
[STRING_TYPE, string]
[LBRACE, {]
[RETURN, return]
[STRING, "hello world!"]
[SEMICOLON, \n]
[RBRACE, }]
[SEMICOLON, \n]
[RBRACE, }]
[SEMICOLON, \n]
`

	out := capturer.CaptureStdout(func() {
		err := lex("../../test/hello.koa")
		if err != nil {
			t.Errorf("failed to lex test.koa, %s", err.Error())
		}
	})

	if out != expectedOutput {
		t.Errorf("wrong outout - expected=%s, got=%s",
			out, expectedOutput)
	}
}
