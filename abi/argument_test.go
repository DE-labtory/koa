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

package abi_test

import (
	"testing"

	"github.com/DE-labtory/koa/abi"
)

func TestArguments_Pack(t *testing.T) {
	tests := []struct {
		arguments    abi.Arguments
		expectedPack string
	}{
		{
			arguments: abi.Arguments{
				abi.Argument{
					Type: abi.Type{abi.String},
				},
				abi.Argument{
					Type: abi.Type{abi.String},
				},
				abi.Argument{
					Type: abi.Type{abi.String},
				},
			},
			expectedPack: "string,string,string",
		},
		{
			arguments: abi.Arguments{
				abi.Argument{
					Type: abi.Type{abi.String},
				},
				abi.Argument{
					Type: abi.Type{abi.Boolean},
				},
				abi.Argument{
					Type: abi.Type{abi.Integer64},
				},
			},
			expectedPack: "string,bool,int64",
		},
		{
			arguments: abi.Arguments{
				abi.Argument{
					Type: abi.Type{abi.Integer64},
				},
				abi.Argument{
					Type: abi.Type{abi.Integer64},
				},
				abi.Argument{
					Type: abi.Type{abi.Integer64},
				},
			},
			expectedPack: "int64,int64,int64",
		},
	}

	for _, test := range tests {
		if test.arguments.Pack() != test.expectedPack {
			t.Errorf("Invalid result off Arguments.Pack(). expected : %s, but got %s", test.arguments.Pack(), test.expectedPack)
		}
	}
}
