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

package main

import (
	"bufio"
	"fmt"
	"io"

	"os"

	"github.com/DE-labtory/koa/parse"
	"github.com/fatih/color"
)

const PROMPT = ">> "
const EXIT = "exit"

const koa = `


	#    #  ####    ##   
	#   #  #    #  #  #
	####   #    # #    #
	#  #   #    # ######
	#   #  #    # #    #       
	#    #  ####  #    #       @DE-labtory/koa v0.0.1


`

func printLogo() {
	color.Yellow(koa)
	bold := color.New(color.Bold)
	bold.Printf("The project is inspired by the simplicity and the ivy-bitcoin. The koa project is to create \na high-level language that has more expressions than the bitcoin script and is simpler and easy to analyze than soldity(ethereum).\n\n")
}

func main() {
	printLogo()
	run(os.Stdin, os.Stdout)
}

func run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == EXIT {
			fmt.Print("bye")
			return
		}

		l := parse.NewLexer(line)
		printTokens(l)
	}
}

func printTokens(l *parse.Lexer) {
	for token := l.NextToken(); token.Type != parse.Eof; {
		fmt.Println(token)
		token = l.NextToken()
	}
}
