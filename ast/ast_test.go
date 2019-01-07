package ast

import "testing"

func TestAssignStatement_String(t *testing.T) {
	tests := []struct {
		input  AssignStatement
		expect string
	}{
		{
			input: AssignStatement{
				Type:     DataStructure{Type: DataStructureType(0), Val: "int"},
				Variable: Identifier{Value: "a"},
				Value:    &IntegerLiteral{Value: 1},
			},
			expect: "int a = 1",
		},
		{
			input: AssignStatement{
				Type:     DataStructure{Type: DataStructureType(0), Val: "bool"},
				Variable: Identifier{Value: "asdf"},
				Value:    &StringLiteral{Value: "hello, world"},
			},
			// type mismatch is not considered here
			expect: "bool asdf = \"hello, world\"",
		},
		{
			input: AssignStatement{
				Type:     DataStructure{Type: DataStructureType(0), Val: "string"},
				Variable: Identifier{Value: "ff"},
				Value:    &BooleanLiteral{Value: true},
			},
			// type mismatch is not considered here
			expect: "string ff = true",
		},
	}

	for i, tt := range tests {
		result := tt.input.String()
		if result != tt.expect {
			t.Errorf("test[%d] - String() wrong result. expected=\"%s\", got=\"%s\"", i, tt.expect, result)
		}
	}
}
