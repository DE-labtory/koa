package koa

import (
	"github.com/DE-labtory/koa/parse"
	"github.com/DE-labtory/koa/translate"
	"github.com/DE-labtory/koa/vm"
)

func Compile(input string) (translate.Asm, error) {
	ast, err := parse.Parse(
		parse.NewTokenBuffer(
			parse.NewLexer(input)))

	if err != nil {
		return translate.Asm{}, err
	}

	asm, err := translate.CompileContract(*ast)
	if err != nil {
		return asm, err
	}

	return asm, nil
}

func Execute(rawByteCode []byte, callFunc *vm.CallFunc) (int64, error) {
	_, err := vm.Execute(rawByteCode, vm.NewMemory(), callFunc)
	if err != nil {
		return 0, err
	}

	return 0, nil
}
