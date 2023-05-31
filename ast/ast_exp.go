package ast

import (
	"bytes"
	"strings"

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

type Boolean struct {
	Token gtoken.Token
	Value bool
}

func (s *Boolean) expressionNode()      {}
func (s *Boolean) TokenLiteral() string { return s.Token.Literal }
func (s *Boolean) String() string       { return s.Token.Literal }

type IfExpression struct {
	Token       gtoken.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (s *IfExpression) expressionNode()      {}
func (s *IfExpression) TokenLiteral() string { return s.Token.Literal }
func (s *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(s.Condition.String())
	out.WriteString(" ")
	out.WriteString(s.Consequence.String())
	if s.Alternative != nil {
		out.WriteString("else")
		out.WriteString(s.Alternative.String())
	}
	return out.String()
}

type CallExpression struct {
	Token     gtoken.Token
	Function  Expression
	Arguments []Expression
}

func (s *CallExpression) expressionNode()      {}
func (s *CallExpression) TokenLiteral() string { return s.Token.Literal }
func (s *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, p := range s.Arguments {
		args = append(args, p.String())
	}
	out.WriteString(s.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

type FunctionLiteral struct {
	Token      gtoken.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (s *FunctionLiteral) expressionNode()      {}
func (s *FunctionLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range s.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(s.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(s.Body.String())
	return out.String()
}

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
