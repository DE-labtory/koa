# Koa

There are two well known blockchain these days, bitcoin and ethereum. and bitcoin has bitcoin script and ethereum has solidity for programming its own smart contract. Both have pros and cons:

**In the case of bitcoin**, it has no state concept and bitcoin script is basically low-level language and has little operation so the capability what it can do is restricted. On the other hand, because of its simplicity of how it works and for bitcoin has no state, we can easily do static analysis — how fast this script will run.

**In the case of ethereum**, it has state concept and solidity designed as high-level language, the solidity developer can program more intuitively, and ethereum smart contract can do a lot of things. (and yes this is also because ethereum has state) On the other hand, as it is designed as high-level language, developer can put infinite-loop by mistake on their smart contract which won’t finish forever and this can make bad effect on network. plus as ethereum has states it is difficulty to do static analysis.


This project is inspired by [the simplicity](https://blockstream.com/simplicity.pdf) and the [ivy-bitcoin](https://github.com/ivy-lang/ivy-bitcoin). Both are aim to high-level crypto-currency language. And “Simplicity” is focuses on functional language without states, loops which enables static analysis to calculate upper bound for computational resources needed easily.

The koa project is to create a high-level language that has `more expressions` than the bitcoin script.

## Contents

* [Lexical analysis](#lexical-analysis)
* [Syntax analysis](#syntax-analysis)
* [Virtual machine](#virtual-machine)



### <a name="lexical-analysis">Lexical analysis</a>

<p align="center"><img src="../image/lexer-diagram.png" width="600px" height="600px"></p>
The first step in the compiler is `lexical analysis` or `scanning`. Lexical analysis reads the stream of characters that make up the source code and groups these letters into a "meaningful permutation" form called lexemes. The lexical analyzer takes each lexeme as a `token` and passes it to the next step, syntax analysis(parser). (For reference, the` lexical analyzer` is simply abbreviated as `lexer`)

For example, in the diagram lexer take the source code; ‘func main() { return 0 }’, then lexer reads code character by character; ‘f’, ‘u’, ‘n’, ‘c’. At the time lexer read ‘c’ lexer knows ‘fun’ + ‘c’ is meaningful word, keyword for function, then lexer cut ‘func’ characters from text(code) and make **token** for that word. Lexer keep this work until we meet `eof` . **In a nutshell lexer group characters and make tokens**.

#### Token

What is token? We can see ‘func’ word as raw data, but without processing that data, those data cannot be easily used from other components. And token is doing that job for us, token is data structure which helps data to be expressed structurally.

```go
type TokenType int

type Token struct {
   Type   TokenType
   Val    string
   Column Pos
   Line   int
}
```

That’s our token defined in project. `Type` is type for word, `Val` for word value. With this `Token` structure the other components like parser can do its job more efficiently and code will be much maintainable and scalable.

#### State and Action

Our lexer design is greatly inspired by [golang template package](https://github.com/golang/go/tree/master/src/text/template/parse) which use **state an action** concept. Actually [go-ethereum](https://github.com/ethereum/go-ethereum/blob/master/core/asm/lexer.go) also use this concept.

- **State** represents where the lexer is from the given input text and what we expect to see next.
- **Action** represents what we are going to do in current state with a piece of input

We can see lexer jobs — read character, generate token, move on to next character— as take the action with current state and move on to next state. After each action, you know where you want to be, the new state is the result of the action.

#### State function

```go
// stateFn determines how to scan the current state.
// stateFn also returns the stateFn to be scanned next after scanning the current state.
type stateFn func(*state, emitter) stateFn
```

This is our state function declaration. State function take current state and emitter, return another state function. Returned state function is based on the current state and knows what to do next. I know that state function definition is quite recursive but this helps keep things simple and clear

```go
// emitter is the interface to emit the token to the client(parser).
type emitter interface {
   emit(t Token)
}
```

And you may have curiosity, what does `emitter` do for us. You may have noticed that we know how to lexing the given inputs, but don’t know how to pass the generated tokens to the client which is probably something like parser. This is why we need `emitter` , `emitter` simply pass the token to the client using one of go features, channel. We are going to see how `emitter`works in a few seconds.

#### Run our state machine

```go
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
```

This is our lexer `run` method which takes the input string — source code — and make `state` with our input. And in the for-loop state function call with the state as argument then return value of state function and this is the new state function. We can see lexer is passed to the state function as `emitter` , don’t be nervous we see this later how lexer implements `emitte` interface. From now, we just need to keep it mind how our state machine works:

**take the current state, do action, walk over to next state.**

What is the advantage of doing this? Well, first of all, we don’t have to check everytime what state we are in. That’s not our concern. We are always in the right place. The only thing to do in our machine is just run state function until we meet nil state function.

#### Run our machine concurrently

We don’t talk much about how to emit the token we generate to the client and I think this is the right time. The idea is we are going to run the lexer as a go routine with the client probably like parser so the two independent machines do their jobs, whenever the lexer has a new thing the client will lob it and do their own work. This mechanism can be done by go channel.

Channel is one of the greatest features in go language and yes complex, but in our lexer it is just a way to deliver data over to another program which may be running completely independent.

```go
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
// emit passes an token back to the client.
func (l *Lexer) emit(t Token) {
   l.tokench <- t
}
```

That’s our lexer definition it has just token channel which are going to be used when emitting token to the client. And we can see in `NewLexer` start to run machine using go-routine.



### <a name="syntax-analysis">Syntax analysis</a>

#### Concept

<p align="center"><img src="../image/parser-diagram.png" width="600px" height="600px"></p>

Syntax analysis reads the stream of `token` which is generated by `lexer`, then make AST(Abstract Synstax Tree) which is passed to compiler.  Koa parser is 'Pratt parser' which is easy to make, modulable and scalable. **The main idea of 'Pratt parser' is each `token` has its own parsing functions.** (infix parsing function, prefix parsing function). In the above diagram `Program` denotes root node of AST and AST consists of slice of `Statements`.



### <a name="virtual-machine">Virtual Machine</a>

#### Basic Architecture

<p align="center"><img src="../image/vm-architecture.png" width="570px" height="350px"></p>

#### Execution Model

<p align="center"><img src="../image/vm-execution.png" width="600px" height="350px"></p>
