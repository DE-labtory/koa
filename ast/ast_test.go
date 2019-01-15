package ast

import "testing"

func TestAssignStatement_String(t *testing.T) {
	tests := []struct {
		input    AssignStatement
		expected string
	}{
		{
			input: AssignStatement{
				Type:     DataStructure{Type: DataStructureType(0), Val: "int"},
				Variable: Identifier{Value: "a"},
				Value:    &IntegerLiteral{Value: 1},
			},
			expected: "int a = 1",
		},
		{
			input: AssignStatement{
				Type:     DataStructure{Type: DataStructureType(0), Val: "bool"},
				Variable: Identifier{Value: "asdf"},
				Value:    &StringLiteral{Value: "hello, world"},
			},
			// type mismatch is not considered here
			expected: "bool asdf = \"hello, world\"",
		},
		{
			input: AssignStatement{
				Type:     DataStructure{Type: DataStructureType(0), Val: "string"},
				Variable: Identifier{Value: "ff"},
				Value:    &BooleanLiteral{Value: true},
			},
			// type mismatch is not considered here
			expected: "string ff = true",
		},
	}

	for _, tt := range tests {
		result := tt.input.String()
		testString(t, result, tt.expected)
	}
}

func TestIdentifier_String(t *testing.T) {
	tests := []struct {
		input    Identifier
		expected string
	}{
		{
			input:    Identifier{Value: "a"},
			expected: "a",
		},
		{
			input:    Identifier{Value: "zcf"},
			expected: "zcf",
		},
		{
			input:    Identifier{Value: "zcf asdf"},
			expected: "zcf asdf",
		},
		{
			input:    Identifier{Value: "123"},
			expected: "123",
		},
		{
			input:    Identifier{Value: ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		result := tt.input.String()
		testString(t, result, tt.expected)
	}
}

func TestStringLiteral_String(t *testing.T) {
	tests := []struct{
		input StringLiteral
		expected string
	}{
		{
			StringLiteral{"hello"},
			"\"hello\"",
		},
		{
			StringLiteral{"hello, world"},
			"\"hello, world\"",
		},
		{
			StringLiteral{"123"},
			"\"123\"",
		},
		{
			StringLiteral{"123, hello"},
			"\"123, hello\"",
		},
		{
			StringLiteral{""},
			"\"\"",
		},
	}

	for _, tt := range tests {
		result := tt.input.String()
		testString(t, result, tt.expected)
	}
}

func TestIntegerLiteral_String(t *testing.T) {
	tests := []struct{
		input IntegerLiteral
		expected string
	}{
		{
			IntegerLiteral{Value: 123},
			"123",
		},
		{
			IntegerLiteral{Value: -1},
			"-1",
		},
	}

	for _, tt := range tests {
		result := tt.input.String()
		testString(t, result, tt.expected)
	}
}

func TestBooleanLiteral_String(t *testing.T) {
	tests := []struct{
		input BooleanLiteral
		expected string
	}{
		{
			BooleanLiteral{true},
			"true",
		},
		{
			BooleanLiteral{false},
			"false",
		},
	}

	for _, tt := range tests {
		result := tt.input.String()
		testString(t, result, tt.expected)
	}
}

func testString(t *testing.T, got, expected string) {
	t.Helper()
	if got != expected {
		t.Errorf("String() wrong result. expected=\"%s\", got=\"%s\"", expected, got)
	}
}
