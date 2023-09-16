package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/GhostNet-Dev/gscript/gtoken"
)

type Expression interface {
	Node
	expressionNode()
}

type TypeIdentifier struct {
	Token    gtoken.Token
	Value    string
	Variable *Identifier
}

func (i *TypeIdentifier) expressionNode()      {}
func (i *TypeIdentifier) TokenLiteral() string { return i.Token.Literal }
func (i *TypeIdentifier) String() string       { return i.Value }

type IdentifierType struct {
	Token gtoken.Token
	Value string
}

func (i *IdentifierType) expressionNode()      {}
func (i *IdentifierType) TokenLiteral() string { return i.Token.Literal }
func (i *IdentifierType) String() string       { return i.Value }

type Identifier struct {
	Token gtoken.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type StringLiteral struct {
	Token gtoken.Token
	Value string
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string       { return s.TokenLiteral() }

type IntegerLiteral struct {
	Token gtoken.Token
	Value int64
}

func (s *IntegerLiteral) expressionNode()      {}
func (s *IntegerLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *IntegerLiteral) String() string       { return s.TokenLiteral() }

type Null struct {
	Token gtoken.Token
	Value string
}

func (s *Null) expressionNode()      {}
func (s *Null) TokenLiteral() string { return s.Token.Literal }
func (s *Null) String() string       { return s.Token.Literal }

type Boolean struct {
	Token gtoken.Token
	Value bool
}

func (s *Boolean) expressionNode()      {}
func (s *Boolean) TokenLiteral() string { return s.Token.Literal }
func (s *Boolean) String() string       { return s.Token.Literal }

type IndexExpression struct {
	Token gtoken.Token
	Left  Expression
	Index Expression
}

func (s *IndexExpression) expressionNode()      {}
func (s *IndexExpression) TokenLiteral() string { return s.Token.Literal }
func (s *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(s.Left.String())
	out.WriteString("[")
	out.WriteString(s.Index.String())
	out.WriteString("])")
	return out.String()
}

type ArrayLiteral struct {
	Token    gtoken.Token
	Elements []Expression
}

func (s *ArrayLiteral) expressionNode()      {}
func (s *ArrayLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range s.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type HashLiteral struct {
	Token gtoken.Token
	Pairs map[Expression]Expression
}

func (s *HashLiteral) expressionNode()      {}
func (s *HashLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *HashLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for key, value := range s.Pairs {
		elements = append(elements, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")
	return out.String()
}

type ForExpression struct {
	Token       gtoken.Token
	Init        Statement
	Increment   Expression
	Condition   Expression
	Consequence *BlockStatement
}

func (s *ForExpression) expressionNode()      {}
func (s *ForExpression) TokenLiteral() string { return s.Token.Literal }
func (s *ForExpression) String() string {
	var out bytes.Buffer
	out.WriteString("for")
	out.WriteString(s.Init.String())
	out.WriteString(s.Condition.String())
	out.WriteString(s.Increment.String())
	out.WriteString(" ")
	out.WriteString(s.Consequence.String())
	return out.String()
}

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
	Name       string
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
	if s.Name != "" {
		out.WriteString(fmt.Sprintf("<%s>", s.Name))
	}
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
