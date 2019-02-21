package koa

import (
	"encoding/binary"

	"github.com/DE-labtory/koa/abi"
	"github.com/DE-labtory/koa/parse"
	"github.com/DE-labtory/koa/translate"
	"github.com/DE-labtory/koa/vm"
)

func Compile(input string) (translate.Asm, abi.ABI, error) {
	ast, err := parse.Parse(
		parse.NewTokenBuffer(
			parse.NewLexer(input)))

	if err != nil {
		return translate.Asm{}, abi.ABI{}, err
	}

	asm, err := translate.CompileContract(*ast)
	if err != nil {
		return asm, abi.ABI{}, err
	}

	a, err := translate.ExtractAbi(*ast)
	if err != nil {
		return asm, abi.ABI{}, err
	}

	return asm, *a, nil
}

func Execute(rawByteCode []byte, function []byte, args []byte) ([]byte, error) {
	callFunc := &vm.CallFunc{
		Func: function,
		Args: args,
	}

	stack, err := vm.Execute(rawByteCode, vm.NewMemory(), callFunc)
	if err != nil {
		return nil, err
	}

	output := Bytes(int64(stack.Pop()))

	return output, nil
}

func Bytes(item int64) []byte {
	byteSlice := make([]byte, 8)
	binary.BigEndian.PutUint64(byteSlice, uint64(item))
	return byteSlice
}
