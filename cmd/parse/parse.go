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
	"io/ioutil"

	"fmt"

	"github.com/DE-labtory/koa/ast"
	parser "github.com/DE-labtory/koa/parse"
	"github.com/urfave/cli"
)

const (
	newLine      = "\n"
	emptySpace   = "    "
	middleItem   = "├── "
	continueItem = "│   "
	lastItem     = "└── "
)

var parseCmd = cli.Command{
	Name:    "parse",
	Aliases: []string{"p"},
	Usage:   "koa parse [filePath]",
	Action: func(c *cli.Context) error {
		return parse(c.Args().Get(0))
	},
}

func Cmd() cli.Command {
	return parseCmd
}

func parse(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	l := parser.NewLexer(string(file))
	buf := parser.NewTokenBuffer(l)
	contract, err := parser.Parse(buf)
	if err != nil {
		return err
	}

	fmt.Println(PrintContract(contract))
	return nil
}

func PrintContract(contract *ast.Contract) string {
	return "contract" + newLine + printFunctions(contract.Functions)
}

func printFunctions(funcs []*ast.FunctionLiteral) string {
	var result string
	for i, f := range funcs {
		result += printFunction(f, len(funcs)-1 == i)
	}

	return result
}

func printFunction(f *ast.FunctionLiteral, isLast bool) string {
	var result string

	if isLast {
		result += lastItem
	} else {
		result += middleItem
	}

	result += f.Signature() + newLine
	result += printStatements(f.Body.Statements, []bool{}, isLast)
	return result
}

func printStatements(statements []ast.Statement, spaces []bool, isLastf bool) string {
	var result string

	for i, s := range statements {
		isLast := i == len(statements)-1

		if isLastf {
			result += emptySpace
		} else {
			result += continueItem
		}

		switch statement := s.(type) {
		case *ast.IfStatement:
			result += printText("if "+statement.Condition.String(), spaces, isLast)

			if statement.Consequence != nil {
				result += printStatements(statement.Consequence.Statements, append(spaces, isLast), isLastf)
			}

			if statement.Alternative != nil {
				result += printText("else", spaces, isLast)
				result += printStatements(statement.Alternative.Statements, append(spaces, isLast), isLastf)
			}

		default:
			result += printText(statement.String(), spaces, isLast)
		}
	}

	return result
}

func printText(text string, spaces []bool, last bool) string {
	var result string
	for _, space := range spaces {
		if space {
			result += emptySpace
		} else {
			result += continueItem
		}
	}

	indicator := middleItem
	if last {
		indicator = lastItem
	}

	return result + indicator + text + newLine
}
