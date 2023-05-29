package parser

import (
	"fmt"

	"github.com/GhostNet-Dev/glambda/ast"
	"github.com/GhostNet-Dev/glambda/gtoken"
	"github.com/GhostNet-Dev/glambda/lexer"
)

type Parser struct {
	l *lexer.Lexer

	errors    []string
	curToken  gtoken.Token
	peekToken gtoken.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.NextToken()
	p.NextToken()

	return p
}

func (p *Parser) NextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextTokenMake()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != gtoken.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case gtoken.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(gtoken.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(gtoken.ASSIGN) {
		return nil
	}

	for !p.curTokenIs(gtoken.SEMICOLON) {
		p.NextToken()
	}
	return stmt
}

func (p *Parser) curTokenIs(t gtoken.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t gtoken.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t gtoken.TokenType) bool {
	if p.peekTokenIs(t) {
		p.NextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t gtoken.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
