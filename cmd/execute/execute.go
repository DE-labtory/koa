package execute

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/DE-labtory/koa"
	"github.com/DE-labtory/koa/abi"
	"github.com/urfave/cli"
)

var executeCmd = cli.Command{
	Name:    "execute",
	Aliases: []string{"e"},
	Usage:   "koa execute [raw byte code] [function name] [args...]",
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 2 {
			return errors.New("you must input at least byte code and function name")
		}
		if len(c.Args()) == 2 {
			return execute(c.Args().Get(0), c.Args().Get(1), nil)
		}
		return execute(c.Args().Get(0), c.Args().Get(1), c.Args()[2:])
	},
}

func Cmd() cli.Command {
	return executeCmd
}

func execute(rawByteCode string, functionName string, args []string) error {
	fnSel := abi.Selector(functionName)
	params, err := encodeParams(args)
	if err != nil {
		return err
	}

	contractDecoding, err := hex.DecodeString(rawByteCode)
	if err != nil {
		return err
	}

	result, err := koa.Execute(contractDecoding, fnSel, params)
	if err != nil {
		return err
	}

	printExecuteResult(result)

	return nil
}

func encodeParams(params []string) ([]byte, error) {
	ps := make([]interface{}, len(params))
	for idx, oneParam := range params {
		// check param is integer
		if iVal, err := strconv.ParseInt(oneParam, 10, 64); err == nil {
			ps[idx] = iVal
			continue
		}

		// check param is bool
		if oneParam == "true" {
			ps[idx] = true
		}
		if oneParam == "false" {
			ps[idx] = false
		}

		// otherwise string
		ps[idx] = oneParam
	}
	result, err := abi.Encode(ps...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func printExecuteResult(result []byte) {
	fmt.Printf("execute Result: %s\n", string(result))
}
