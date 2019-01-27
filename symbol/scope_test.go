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

package symbol

import "testing"

func testNewEnclosedScope(t *testing.T) {
	outer := NewScope()
	s := NewEnclosedScope(outer)

	if s.outer != outer {
		t.Fatalf("testNewEnclosedScope() failed. outer must be set")
	}

	if len(s.store) > 0 {
		t.Fatalf("testNewEnclosedScope() failed. store's size must be 0")
	}
}

func testNewScope(t *testing.T) {
	s := NewScope()
	if s.outer != nil {
		t.Fatalf("testNewScope() failed. outer must be nil")
	}

	if len(s.store) > 0 {
		t.Fatalf("testNewScope() failed. store's size must be 0")
	}
}

// TODO implement me w/ test cases :-)
func testScopeGetter(t *testing.T) {

}

// TODO implement me w/ test cases :-)
func testScopeSetter(t *testing.T) {

}

// TODO implement me w/ test cases :-)
func testScopeString(t *testing.T) {

}
