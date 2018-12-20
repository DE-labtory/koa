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

package parse

type Lexer struct {
	tokench chan Token
}

func New(input string) *Lexer {

	l := &Lexer{
		tokench: make(chan Token, 2),
	}

	go l.run(input)
	return l
}

func (l *Lexer) run(input string) {

	state := newState(input)
	var token Token

	for stateFn := DefaultStateFn; stateFn != nil; {
		stateFn, state, token = stateFn(state)
		l.tokench <- token
	}

	close(l.tokench)
}

type stateFn func(state) (stateFn, state, Token)

// nextToken returns the next token from the input.
// Called by the parser, not in the lexing goroutine.
func (l *Lexer) Next() Token {
	return <-l.tokench
}

type state struct {
	input string
	start int
	pos   int
	line  int
}

func (s *state) emit(t TokenType) Token {
	return Token{t, s.input[s.start:s.pos], s.start, s.pos, s.line}
}

func newState(input string) state {
	return state{
		input: input,
	}
}

func DefaultStateFn(s state) (stateFn, state, Token) {

	return DefaultStateFn, state{}, Token{}
}

func NumberStateFN(s state) (stateFn, state, Token) {

	return DefaultStateFn, state{}, Token{}
}
