package parser

import (
	"fmt"
	"strconv"

	"github.com/GhostNet-Dev/glambda/ast"
	"github.com/GhostNet-Dev/glambda/gtoken"
)

const (
	_ int = iota
	LOWEST
	EQUALS      //==
	LESSGREATER // > OR <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X OR !X
	CALL        // myFunction(X)
)

var precedences = map[gtoken.TokenType]int{
	gtoken.EQ:       EQUALS,
	gtoken.NOT_EQ:   EQUALS,
	gtoken.LT:       LESSGREATER,
	gtoken.RT:       LESSGREATER,
	gtoken.PLUS:     SUM,
	gtoken.MINUS:    SUM,
	gtoken.SLASH:    PRODUCT,
	gtoken.ASTERISK: PRODUCT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) initExpression() {
	p.prefixParseFns = make(map[gtoken.TokenType]prefixParseFn)
	p.registerPrefix(gtoken.IDENT, p.parseIdentifier)
	p.registerPrefix(gtoken.INT, p.parseIntergerLiteral)
	p.registerPrefix(gtoken.BANG, p.parsePrefixExpression)
	p.registerPrefix(gtoken.MINUS, p.parsePrefixExpression)
	p.infixParseFns = make(map[gtoken.TokenType]infixParseFn)
	p.registerInfix(gtoken.PLUS, p.parseInfixExpression)
	p.registerInfix(gtoken.MINUS, p.parseInfixExpression)
	p.registerInfix(gtoken.SLASH, p.parseInfixExpression)
	p.registerInfix(gtoken.ASTERISK, p.parseInfixExpression)
	p.registerInfix(gtoken.EQ, p.parseInfixExpression)
	p.registerInfix(gtoken.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(gtoken.LT, p.parseInfixExpression)
	p.registerInfix(gtoken.RT, p.parseInfixExpression)
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(gtoken.SEMICOLON) {
		p.NextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(gtoken.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.NextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.NextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.NextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntergerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t gtoken.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType gtoken.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType gtoken.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
