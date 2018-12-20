package parse

import "testing"

func TestToken_String(t *testing.T) {
	token := Token{}
	if token.String() != "" {
		t.Fatalf("token String() wrong, expected=%q, got=%q", "", token.String())
	}
}
