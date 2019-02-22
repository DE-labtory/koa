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

package repl

import (
	"fmt"
	"io"

	"github.com/DE-labtory/koa/translate"

	"bufio"
	"os"

	compile_cmd "github.com/DE-labtory/koa/cmd/compile"
	lex_cmd "github.com/DE-labtory/koa/cmd/lex"
	parse_cmd "github.com/DE-labtory/koa/cmd/parse"
	"github.com/DE-labtory/koa/parse"
	"github.com/fatih/color"
)

const PROMPT = ">> "
const EXIT = "exit()"

const koa = `


	#    #  ####    ##   
	#   #  #    #  #  #
	####   #    # #    #
	#  #   #    # ######
	#   #  #    # #    #       
	#    #  ####  #    #       @DE-labtory/koa v0.1.0


`

func PrintLogo() {
	color.Yellow(koa)
	bold := color.New(color.Bold)
	bold.Printf("github: https://github.com/DE-labtory/koa \n\n")
	fmt.Printf("The project is inspired by the simplicity and the ivy-bitcoin. The koa project is to create \na high-level language that has more expressions than the bitcoin script and is simpler and easy to analyze than soldity(ethereum).\n\n")
	bold.Print("Use exit() or Ctrl-c to exit \n")
}

func Run() {
	PrintLogo()
	run(os.Stdin, os.Stdout)
}

func run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	bold := color.New(color.Bold)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == EXIT {
			fmt.Println("bye")
			return
		}

		l := parse.NewLexer(line)
		l2 := parse.NewLexer(line)

		buf := parse.NewTokenBuffer(l)
		contract, err := parse.Parse(buf)

		asm, err := translate.CompileContract(*contract)
		if err != nil {
			color.Red(err.Error())
			continue
		}

		ab, err := translate.ExtractAbi(*contract)
		if err != nil {
			color.Red(err.Error())
		}

		if err != nil {
			color.Red(err.Error())
			continue
		}

		bold.Println("-->>   LEX RESULT   <<-----------------------------------------------")
		lex_cmd.PrintTokens(l2)
		fmt.Println()

		bold.Println("-->>  PARSE RESULT  <<-----------------------------------------------")
		fmt.Println(parse_cmd.PrintContract(contract))
		fmt.Println()

		bold.Println("-->> COMPILE RESULT <<-----------------------------------------------")
		if err := compile_cmd.PrintCompileResult(asm, ab); err != nil {
			color.Red(err.Error())
			continue
		}
		fmt.Println()
	}
}
