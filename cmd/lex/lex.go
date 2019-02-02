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
	"fmt"
	"io/ioutil"

	"github.com/DE-labtory/koa/parse"
	"github.com/urfave/cli"
)

var lexCmd = cli.Command{
	Name:    "lex",
	Aliases: []string{"l"},
	Usage:   "koa lex [filePath]",
	Action: func(c *cli.Context) error {
		return lex(c.Args().Get(0))
	},
}

func Cmd() cli.Command {
	return lexCmd
}

func lex(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	l := parse.NewLexer(string(file))
	printTokens(l)
	return nil
}

func printTokens(l *parse.Lexer) {
	for token := l.NextToken(); token.Type != parse.Eof; {
		fmt.Println(token)
		token = l.NextToken()
	}
}
