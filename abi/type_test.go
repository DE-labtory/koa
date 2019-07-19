/*
 * Copyright 2018-2019 De-labtory
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

func TestNewType(t *testing.T) {
	tests := []struct {
		Type         string
		expectedType abi.ParamType
	}{
		{
			Type:         "int64",
			expectedType: abi.Integer64,
		},
		{
			Type:         "bool",
			expectedType: abi.Boolean,
		},
		{
			Type:         "string",
			expectedType: abi.String,
		},
	}

	for _, test := range tests {
		Type, err := abi.NewType(test.Type)
		if err != nil {
			t.Error(err)
		}

		if Type.Type != test.expectedType {
			t.Errorf("Invalid Type : %s", Type.Type)
		}
	}
}
