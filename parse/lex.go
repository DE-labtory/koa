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

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// emitter is the interface to emit the token to the client(parser).
type emitter interface {
	emit(t Token)
}

type Lexer struct {
	tokench chan Token
}

func NewLexer(input string) *Lexer {

	l := &Lexer{
		tokench: make(chan Token, 2),
	}

	go l.run(input)
	return l
}

// run runs the state machine for the lexer.
func (l *Lexer) run(input string) {

	state := &state{
		input: input,
	}

	for stateFn := defaultStateFn; stateFn != nil; {
		stateFn = stateFn(state, l)
	}

	close(l.tokench)
}

// emit passes an token back to the client.
func (l *Lexer) emit(t Token) {
	l.tokench <- t
}

// The process of generating a token from an input string(codes) is generally implemented
// by defining a state and determining how to process the state.
// After the state is processed, it goes to the next state and it is repeated to determine
// how to process again through the switch statement.
//
// Example)
//
// // One interation:
// switch state {
// case state1:
//    state = action1()
// case state2:
//	  state = action2()
// case state3:
//    state = action3()
// }
//
// In the above code, if a new state is returned through action2() and the state is checked
// again with a switch. But we already know what state comes after action2(), and it would be better
// if we could execute the corresponding action without switch.
//
//
// The above code can be changed to execute an action, returns the next state as a state function.
// Recursive definition but simple and clear
//
// func run(){
//     for stateFn := startState; state != nil{
//         stateFn = stateFn(lexer)
//     }
// }
//
// stateFn determines how to scan the current state.
// stateFn also returns the stateFn to be scanned next after scanning the current state.
type stateFn func(*state, emitter) stateFn

// NextToken returns the next token from the input.
// Called by the parser, not in the lexing goroutine.
func (l *Lexer) NextToken() Token {
	return <-l.tokench
}

// State has the input(codes) as a string and has the current position and the line.
type state struct {
	input string
	start Pos
	end   Pos
	line  int
	width Pos
}

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

// Cut return a token and set start position to pos
func (s *state) cut(t TokenType) Token {
	token := Token{t, s.input[s.start:s.end], s.end, s.line}
	s.start = s.end

	return token
}

// Next returns the next rune in the input.
func (s *state) next() rune {
	if int(s.end) >= len(s.input) {
		s.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(s.input[s.end:])
	s.width = Pos(w)
	s.end += s.width

	if r == '\n' {
		s.line++
	}
	return r
}

// TODO: implement me w/ test cases :-)
// Peek returns but does not consume the next byte in the input.
func (s *state) peek() rune {
	ch := s.next()
	s.backup()

	return ch
}

// TODO: implement me w/ test cases :-)
// Backup steps back one rune. Can only be called once per call of next.
func (s *state) backup() {
	s.end -= s.width
	// Correct newline count.
	if s.width == 1 && s.input[s.end] == '\n' {
		s.line--
	}
}

// Accept consumes the next byte if it's from the valid set.
func (s *state) accept(valid string) bool {
	if strings.IndexRune(valid, s.next()) >= 0 {
		return true
	}
	s.backup()
	return false
}

// AcceptRun consumes a run of byte from the valid set.
func (s *state) acceptRun(valid string) {
	for strings.IndexRune(valid, s.next()) >= 0 {
	}
	s.backup()
}

//
//	Bang     // !
//  Assign   // =
//
//	Plus     // +
//	Minus    // -
//	Asterisk // *
//	Slash    // /
//	Mod      // %
//
//	LT     // <
//	GT     // >
//	LTE    // <=
//	GTE    // >=
//	EQ     // ==
//	NOT_EQ // !=
//
//	Comma // ,
//
//	Lparen // (
//	Rparen // )
//	Lbrace // {
//	Rbrace // }
//
// TODO: LTE, GTE, spaceStateFn, numberStateFn, identifierStateFn case
func defaultStateFn(s *state, e emitter) stateFn {

	switch ch := s.next(); {

	case ch == '!':
		if s.peek() == '=' {
			s.next()
			e.emit(s.cut(NOT_EQ))
		} else {
			e.emit(s.cut(Bang))
		}
	case ch == '=':
		if s.peek() == '=' {
			s.next()
			// ==
			e.emit(s.cut(EQ))
		} else {
			// =
			e.emit(s.cut(Assign))
		}

	case ch == '+':
		e.emit(s.cut(Plus))
	case ch == '-':
		e.emit(s.cut(Minus))
	case ch == '/':
		e.emit(s.cut(Slash))
	case ch == '*':
		e.emit(s.cut(Asterisk))
	case ch == '%':
		e.emit(s.cut(Mod))

	case ch == '<':
		e.emit(s.cut(LT))
	case ch == '>':
		e.emit(s.cut(GT))

	case ch == ')':
		e.emit(s.cut(Rparen))
	case ch == '(':
		e.emit(s.cut(Lparen))
	case ch == '}':
		e.emit(s.cut(Rbrace))
	case ch == '{':
		e.emit(s.cut(Lbrace))
	case ch == ',':
		e.emit(s.cut(Comma))
	case ch == eof:
		e.emit(s.cut(Eof))
	default:
		e.emit(s.cut(Illegal))
	}

	return defaultStateFn
}

// NumberStateFn scans an alphanumeric. ex) 123, 4001, 232
// After reading Number, it returns DefaultStateFn.
// number = { decimal_digit }
func numberStateFn(s *state, e emitter) stateFn {
	s.accept("+-")

	const digits = "0123456789"

	if !s.accept(digits) {
		e.emit(Token{Illegal, "Invalid function call: numberStateFn", s.end, s.line})
		return defaultStateFn
	}

	//accept first digit
	s.next()

	for s.accept(digits) {
		s.next()
	}

	e.emit(s.cut(Int))
	return defaultStateFn
}

// IdentifierStateFn scans an identifiers. ex) a, b, add
// After reading a identifier, it returns DefaultStateFn.
//
// when a input of a state is "int abc = 5" and start of state is 4,
// IdentifierStateFn should emit "abc" and return DefaultStateFn
//
// identifier = letter { letter | unicode_digit }.
func identifierStateFn(s *state, e emitter) stateFn {
	if !(unicode.IsLetter(s.peek()) || s.peek() == '_') {
		errToken := Token{Illegal, "Invalid function call: identifierStateFn", s.end, s.line}
		e.emit(errToken)
		return defaultStateFn
	}

	s.next() // accept first letter
	for isAlphaNumeric(s.peek()) {
		s.next()
	}

	//lookup keywords map and return tokenType
	tok := LookupIdent(s.input[s.start:s.end])
	e.emit(s.cut(tok))
	return defaultStateFn
}

// SpaceStateFn scans an space. ex) `\t`, `" "`
// After ignoring all spaces, it returns DefaultStateFn.
//
func spaceStateFn(s *state, e emitter) stateFn {
	const spaceChars = " \t"

	if !s.accept(spaceChars) {
		errToken := Token{Illegal, "Invalid function call: spaceStateFn", s.end, s.line}
		e.emit(errToken)
		return defaultStateFn
	}

	//absorb " ", "\t"
	s.acceptRun(spaceChars)
	s.cut(Illegal)
	return defaultStateFn
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
