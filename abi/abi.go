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

package abi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DE-labtory/koa/ast"
)

type ABI struct {
	Methods []Method
}

func New(abiJSON string) (ABI, error) {
	reader := strings.NewReader(abiJSON)
	dec := json.NewDecoder(reader)

	var abi ABI
	if err := dec.Decode(&abi); err != nil {
		return ABI{}, err
	}

	return abi, nil
}

// UnmarshalJSON is implementation of json.Decoder's UnmarshalJSON
func (abi *ABI) UnmarshalJSON(data []byte) error {
	var methods []Method

	if err := json.Unmarshal(data, &methods); err != nil {
		return err
	}

	for _, method := range methods {
		abi.Methods = append(abi.Methods, method)
	}

	return nil
}

// Extract ABI from function of ast.
func ExtractAbiFromFunction(f ast.FunctionLiteral) (Method, error) {
	method := Method{
		Name:      f.Name.String(),
		Arguments: make([]Argument, 0),
		Output:    Argument{},
	}

	args := make([]Argument, 0)

	for _, param := range f.Parameters {
		t, err := convertAstTypeToAbi(param.Type)
		if err != nil {
			return Method{}, err
		}

		arg := Argument{
			Name: param.Identifier.String(),
			Type: t,
		}

		args = append(args, arg)
	}

	method.Arguments = args

	t, err := convertAstTypeToAbi(f.ReturnType)
	if err != nil {
		return Method{}, err
	}

	method.Output = Argument{
		Name: "",
		Type: t,
	}

	return method, nil
}

func convertAstTypeToAbi(p ast.DataStructure) (Type, error) {
	switch p {
	case ast.IntType:
		return NewType("int")
	case ast.StringType:
		return NewType("string")
	case ast.BoolType:
		return NewType("bool")
	default:
		return Type{}, fmt.Errorf("Unknown paramter type. got=%v", p)
	}
}
