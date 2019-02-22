package compile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/DE-labtory/koa/abi"

	parser "github.com/DE-labtory/koa/parse"
	"github.com/DE-labtory/koa/translate"
	"github.com/urfave/cli"
)

type Result struct {
	Abi     *abi.ABI
	Asm     string
	RawByte string
}

var compileCmd = cli.Command{
	Name:    "compile",
	Aliases: []string{"c"},
	Usage:   "koa compile [filepath]",
	Action: func(c *cli.Context) error {
		return compile(c.Args().Get(0))
	},
}

func Cmd() cli.Command {
	return compileCmd
}

func compile(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	contract, err := parser.Parse(
		parser.NewTokenBuffer(
			parser.NewLexer(string(file))))

	if err != nil {
		return err
	}

	asm, err := translate.CompileContract(*contract)
	if err != nil {
		return err
	}

	ab, err := translate.ExtractAbi(*contract)
	if err != nil {
		return err
	}

	if err := PrintCompileResult(asm, ab); err != nil {
		return err
	}

	return nil
}

func PrintCompileResult(asm translate.Asm, ab *abi.ABI) error {
	result := Result{
		Abi:     ab,
		Asm:     asm.String(),
		RawByte: fmt.Sprintf("%x", asm.ToRawByteCode()),
	}

	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}
