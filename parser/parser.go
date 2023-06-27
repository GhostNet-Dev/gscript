package parser

import (
	"fmt"

	"github.com/GhostNet-Dev/glambda/ast"
	"github.com/GhostNet-Dev/glambda/gtoken"
	"github.com/GhostNet-Dev/glambda/lexer"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  gtoken.Token
	peekToken gtoken.Token

	prefixParseFns map[gtoken.TokenType]prefixParseFn
	infixParseFns  map[gtoken.TokenType]infixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.initExpression()

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
	case gtoken.TYPE:
		return p.parseTypeStatement()
	case gtoken.LET:
		return p.parseLetStatement()
	case gtoken.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.NextToken()

	for !p.curTokenIs(gtoken.RBRACE) && !p.curTokenIs(gtoken.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.NextToken()
	}
	return block
}

func (p *Parser) parseObjectBlockStatement() *ast.ObjectBlockStatement {
	block := &ast.ObjectBlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.NextToken()

	for !p.curTokenIs(gtoken.RBRACE) && !p.curTokenIs(gtoken.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.NextToken()
	}
	return block
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.NextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	for !p.curTokenIs(gtoken.SEMICOLON) && !p.curTokenIs(gtoken.EOF) {
		p.NextToken()
	}
	return stmt
}

func (p *Parser) parseTypeStatement() *ast.TypeStatement {
	stmt := &ast.TypeStatement{Token: p.curToken}
	if !p.expectPeek(gtoken.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(gtoken.STRUCT) {
		return nil
	}
	stmt.Type = &ast.IdentifierType{Token: p.curToken, Value: p.curToken.Literal}

	if stmt.Type.Value == "struct" && p.expectPeek(gtoken.LBRACE) {
		stmt.Body = p.parseObjectBlockStatement()
	}
	return stmt
}

func (p *Parser) parseStructStatement() *ast.StructStatement {
	stmt := &ast.StructStatement{Token: p.curToken}
	if !p.expectPeek(gtoken.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(gtoken.LBRACE) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()
	return stmt
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
	p.NextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if fl, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fl.Name = stmt.Name.Value
	}

	for !p.curTokenIs(gtoken.SEMICOLON) && !p.curTokenIs(gtoken.EOF) {
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
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t gtoken.TokenType) {
	msg := fmt.Sprintf("%d: expected next token to be %s, got %s instead",
		p.peekToken.Line, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
