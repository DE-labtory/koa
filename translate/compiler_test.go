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

package translate_test

import (
	"testing"

	"github.com/DE-labtory/koa/abi"
	"github.com/DE-labtory/koa/ast"
	"github.com/DE-labtory/koa/translate"
)

func TestFuncMap_Declare(t *testing.T) {
	tests := []struct {
		funcLiteral ast.FunctionLiteral
		b           translate.Bytecode
		pc          int
		len         int
	}{
		{
			funcLiteral: makeFuncLiteral("foo", ast.VoidType),
			b: translate.Bytecode{
				RawByte: []byte{01, 02, 03, 04},
				AsmCode: []string{"01020304"},
			},
			pc:  1,
			len: 1,
		},
		{
			funcLiteral: makeFuncLiteral("bar", ast.IntType),
			b: translate.Bytecode{
				RawByte: []byte{05, 06},
				AsmCode: []string{"05", "06"},
			},
			pc:  2,
			len: 2,
		},
		{
			funcLiteral: makeFuncLiteral("sam", ast.StringType),
			b: translate.Bytecode{
				RawByte: []byte{11, 12, 13, 14, 15, 16},
				AsmCode: []string{"11", "12", "13", "14", "15", "16"},
			},
			pc:  6,
			len: 3,
		},
	}

	fMap := translate.FuncMap{}

	for i, test := range tests {
		fMap.Declare(test.funcLiteral.Signature(), test.b)

		result := fMap[string(abi.Selector(test.funcLiteral.Signature()))]
		if test.pc != result {
			t.Fatalf("test[%d] - Declare() pc result wrong. expected=%d, got=%d",
				i, test.pc, result)
		}

		if test.len != len(fMap) {
			t.Fatalf("test[%d] - Declare() FuncMap result wrong. expected=%d, got=%d", i, test.len, len(fMap))
		}
	}
}

// TODO: implement test cases :-)
func TestCompileContract(t *testing.T) {

}

func makeFuncLiteral(funcName string, retType ast.DataStructure, params ...*ast.ParameterLiteral) ast.FunctionLiteral {
	funcLiteral := ast.FunctionLiteral{
		Name: &ast.Identifier{
			Value: funcName,
		},
		Parameters: []*ast.ParameterLiteral{},
		Body:       &ast.BlockStatement{},
		ReturnType: retType,
	}

	for _, param := range params {
		funcLiteral.Parameters = append(funcLiteral.Parameters, param)
	}

	return funcLiteral
}
