package ast

import (
	"bytes"

	"github.com/GhostNet-Dev/gscript/gtoken"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type TypeStatement struct {
	Token gtoken.Token
	Name  *Identifier
	Type  *IdentifierType
	Body  *ObjectBlockStatement
}

func (s *TypeStatement) statementNode()       {}
func (s *TypeStatement) TokenLiteral() string { return s.Token.Literal }
func (s *TypeStatement) String() string {
	var out bytes.Buffer
	out.WriteString("type " + s.Name.String() + " ")
	out.WriteString(s.TokenLiteral())
	out.WriteString(" {\n")
	if s.Body != nil {
		out.WriteString(s.Body.String())
	}
	out.WriteString("}")
	return out.String()
}

type StructStatement struct {
	Token gtoken.Token
	Name  *Identifier
	Body  *BlockStatement
}

func (s *StructStatement) statementNode()       {}
func (s *StructStatement) TokenLiteral() string { return s.Token.Literal }
func (s *StructStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Name.String())
	out.WriteString(" {\n")
	if s.Body != nil {
		out.WriteString(s.Body.String())
	}
	out.WriteString("}")
	return out.String()
}

type LetStatement struct {
	Token gtoken.Token
	Name  *Identifier
	Value Expression
}

func (s *LetStatement) statementNode()       {}
func (s *LetStatement) TokenLiteral() string { return s.Token.Literal }
func (s *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Name.String())
	out.WriteString(" = ")
	if s.Value != nil {
		out.WriteString(s.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token       gtoken.Token
	ReturnValue Expression
}

func (s *ReturnStatement) statementNode()       {}
func (s *ReturnStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	if s.ReturnValue != nil {
		out.WriteString(s.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ObjectBlockStatement struct {
	Token      gtoken.Token
	Statements []Statement
}

func (s *ObjectBlockStatement) statementNode()       {}
func (s *ObjectBlockStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ObjectBlockStatement) String() string {
	var out bytes.Buffer
	for _, statement := range s.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token      gtoken.Token
	Statements []Statement
}

func (s *BlockStatement) statementNode()       {}
func (s *BlockStatement) TokenLiteral() string { return s.Token.Literal }
func (s *BlockStatement) String() string {
	var out bytes.Buffer
	for _, statement := range s.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}

type ExpressionStatement struct {
	Token      gtoken.Token
	Expression Expression
}

func (s *ExpressionStatement) statementNode()       {}
func (s *ExpressionStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ExpressionStatement) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}
	return ""
}
