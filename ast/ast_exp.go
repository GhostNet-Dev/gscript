package ast

import (
	"bytes"

	"github.com/GhostNet-Dev/glambda/gtoken"
)

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token gtoken.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token gtoken.Token
	Value int64
}

func (s *IntegerLiteral) expressionNode()      {}
func (s *IntegerLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *IntegerLiteral) String() string       { return s.TokenLiteral() }

type PrefixExpression struct {
	Token    gtoken.Token // ex. !, -
	Operator string
	Right    Expression
}

func (s *PrefixExpression) expressionNode()      {}
func (s *PrefixExpression) TokenLiteral() string { return s.Token.Literal }
func (s *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(s.Operator)
	out.WriteString(s.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    gtoken.Token // ex. +, -, *...
	Left     Expression
	Operator string
	Right    Expression
}

func (s *InfixExpression) expressionNode()      {}
func (s *InfixExpression) TokenLiteral() string { return s.Token.Literal }
func (s *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(s.Left.String())
	out.WriteString(" " + s.Operator + " ")
	out.WriteString(s.Right.String())
	out.WriteString(")")
	return out.String()
}
