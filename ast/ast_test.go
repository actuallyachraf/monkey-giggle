package ast

import (
	"testing"

	"github.com/actuallyachraf/monkey-giggle/token"
)

func TestAST(t *testing.T) {

	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.NewLiteral(token.LET, "let"),
				Name: &Identifier{
					Token: token.NewLiteral(token.IDENT, "myVar"),
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.NewLiteral(token.IDENT, "anotherVar"),
					Value: "anotherVar",
				},
			},
		},
	}
	expected := "let myVar = anotherVar;"
	if program.String() != expected {
		t.Fatalf("program.String() failed got %s expected %s", program.String(), expected)
	}
}
