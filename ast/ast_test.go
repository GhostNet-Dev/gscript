package ast

import (
	"testing"

	"github.com/GhostNet-Dev/gscript/gtoken"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: gtoken.Token{Type: gtoken.LET, Literal: "let"},
				Name: &Identifier{
					Token: gtoken.Token{Type: gtoken.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: gtoken.Token{Type: gtoken.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
