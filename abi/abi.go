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

import "github.com/DE-labtory/koa/crpyto"

// The implementation below is implemented for the abi spec of ethereum.
// https://solidity.readthedocs.io/en/develop/abi-spec.html

// Encode abi parameters
func Encode(params ...[]byte) []byte {
	return nil
}

// Encode abi parameters with function Selector
func EncodeWithSelector(selector []byte, params ...[]byte) []byte {
	return nil
}

// Get function selector(4bytes) from string of function signature
func Selector(functionSignature string) []byte {
	return crpyto.Keccak256([]byte(functionSignature))[:4]
}
