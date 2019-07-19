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

package abi

import (
	"encoding/json"
	"strings"
)

type Argument struct {
	Name string
	Type Type
}

type Arguments []Argument

type ArgumentMarshaling struct {
	Name string
	Type string
}

// UnmarshalJSON implements json.Unmarshaler interface
func (argument *Argument) UnmarshalJSON(data []byte) error {
	var arg ArgumentMarshaling
	err := json.Unmarshal(data, &arg)
	if err != nil {
		return err
	}
	argument.Type, err = NewType(arg.Type)
	if err != nil {
		return err
	}
	argument.Name = arg.Name

	return nil
}

// Pack returns series of Arguments type
// Example
// function foo(int64 a, bool b)
// Pack() extract "int64,bool"
func (arguments Arguments) Pack() string {
	var packedTypes []string

	for _, argument := range arguments {
		packedTypes = append(packedTypes, string(argument.Type.Type))
	}

	return strings.Join(packedTypes, ",")
}
