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
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DE-labtory/koa/cmd/lex"
	"github.com/DE-labtory/koa/cmd/parse"
	"github.com/DE-labtory/koa/cmd/repl"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

const PROMPT = ">> "
const EXIT = "exit()"

const koa = `


	#    #  ####    ##   
	#   #  #    #  #  #
	####   #    # #    #
	#  #   #    # ######
	#   #  #    # #    #       
	#    #  ####  #    #       @DE-labtory/koa v0.0.1


`

func PrintLogo() {
	color.Yellow(koa)
	bold := color.New(color.Bold)
	fmt.Printf("The project is inspired by the simplicity and the ivy-bitcoin. The koa project is to create \na high-level language that has more expressions than the bitcoin script and is simpler and easy to analyze than soldity(ethereum).\n\n")
	bold.Print("Use exit() or Ctrl-c to exit \n")
}

func main() {

	app := cli.NewApp()
	app.Name = "koa"
	app.Version = "0.0.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "koa",
			Email: "koa@Delabtory.io",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "",
			Usage: "set config",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "set debug mode",
		},
	}
	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, lex.Cmd())
	app.Commands = append(app.Commands, parse.Cmd())
	app.Action = func(c *cli.Context) error {
		repl.Run()
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
