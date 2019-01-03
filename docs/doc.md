# Koa

The project is inspired by [the simplicity](https://blockstream.com/simplicity.pdf) and the [ivy-bitcoin](https://github.com/ivy-lang/ivy-bitcoin).

The koa project is to create a high-level language that has `more expressions` than the bitcoin script.

### Lexical analysis

<p align="center"><img src="../image/lexer-diagram.png" width="600px" height="600px"></p>
The first step in the compiler is `lexical analysis` or `scanning`. Lexical analysis reads the stream of characters that make up the source code and groups these letters into a "meaningful permutation" form called lexemes. The lexical analyzer takes each lexeme as a `token` and passes it to the next step, syntax analysis(parser). (For reference, the` lexical analyzer` is simply abbreviated as `lexer`)

### Virtual Machine

#### Basic Architecture

<p align="center"><img src="../image/vm-architecture.png" width="570px" height="350px"></p>

#### Execution Model

<p align="center"><img src="../image/vm-execution.png" width="600px" height="350px"></p>