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

package translate

import (
	"github.com/DE-labtory/koa/ast"
	"github.com/DE-labtory/koa/opcode"
	"github.com/pkg/errors"
)

type Compiler struct {
	// binary is the byte code which is compiled.
	binary []byte

	debug bool
}

func NewCompiler(debug bool) *Compiler {
	return &Compiler{
		binary: make([]byte,0),
		debug:  debug,
	}
}

// Compile() compiles statements in program.
// Statements would be compiled to byte code.
// TODO: implement w/ test cases :-)
func (c *Compiler) Compile(program ast.Program) error {
	for _, s := range program.Statements {
		if err := c.compileNode(s); err != nil {
			return err
		}
	}

	return nil
}

// compileNode() compiles a node in statement.
// This function will be executed recursively.
// TODO: implement w/ test cases :-)
func (c *Compiler) compileNode(node ast.Node) error {
	// Nodes are many kinds.
	switch node := node.(type) {
	case *ast.Identifier:
		return nil

	case *ast.AssignStatement:
		return nil

	case *ast.StringLiteral:
		return nil

	case *ast.IntegerLiteral:
		return nil

	case *ast.BooleanLiteral:
		return nil

	case *ast.PrefixExpression:
		return nil

	default:
		return errors.New("compileNode() error - "+node.String()+" could not compiled")
	}

	return nil
}

// emit() generates a byte code with operator and operands.
// Then, adds the byte code to binary in compiler
// and returns the position of this instruction.
// TODO: implement w/ test cases :-)
func (c *Compiler) emit(operator opcode.Type, operand ...[]byte) int {
	return 0
}
