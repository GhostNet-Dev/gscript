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
	INDEX       // array[index]
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
	gtoken.LPAREN:   CALL,
	gtoken.LBRACKET: INDEX,
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
	p.registerPrefix(gtoken.TRUE, p.parseBoolean)
	p.registerPrefix(gtoken.FALSE, p.parseBoolean)
	p.registerPrefix(gtoken.LPAREN, p.parseGroupedExpresion)
	p.registerPrefix(gtoken.IF, p.parseIfExpresion)
	p.registerPrefix(gtoken.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(gtoken.STRING, p.parseStringLiteral)
	p.registerPrefix(gtoken.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(gtoken.LBRACE, p.parseHashLiteral)

	p.infixParseFns = make(map[gtoken.TokenType]infixParseFn)
	p.registerInfix(gtoken.PLUS, p.parseInfixExpression)
	p.registerInfix(gtoken.MINUS, p.parseInfixExpression)
	p.registerInfix(gtoken.SLASH, p.parseInfixExpression)
	p.registerInfix(gtoken.ASTERISK, p.parseInfixExpression)
	p.registerInfix(gtoken.EQ, p.parseInfixExpression)
	p.registerInfix(gtoken.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(gtoken.LT, p.parseInfixExpression)
	p.registerInfix(gtoken.RT, p.parseInfixExpression)
	p.registerInfix(gtoken.LPAREN, p.parseCallExpression)
	p.registerInfix(gtoken.LBRACKET, p.parseIndexExpression)
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(gtoken.RBRACE) {
		p.NextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek((gtoken.COLON)) {
			return nil
		}
		p.NextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value
		if !p.peekTokenIs(gtoken.RBRACE) && !p.expectPeek(gtoken.COMMA) {
			return nil
		}
	}
	if !p.expectPeek(gtoken.RBRACE) {
		return nil
	}
	return hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.NextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(gtoken.RBRACKET) {
		return nil
	}
	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(gtoken.RBRACKET)
	return array
}

func (p *Parser) parseExpressionList(end gtoken.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.NextToken()
		return list
	}

	p.NextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(gtoken.COMMA) {
		p.NextToken()
		p.NextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(gtoken.RPAREN)
	return exp
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(gtoken.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(gtoken.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatemnt()
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(gtoken.RPAREN) {
		p.NextToken()
		return identifiers
	}
	p.NextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)
	for p.peekTokenIs(gtoken.COMMA) {
		p.NextToken()
		p.NextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !p.expectPeek(gtoken.RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseIfExpresion() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(gtoken.LPAREN) {
		return nil
	}
	p.NextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(gtoken.RPAREN) {
		return nil
	}
	if !p.expectPeek(gtoken.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatemnt()
	if p.peekTokenIs(gtoken.ELSE) {
		p.NextToken()
		if !p.expectPeek(gtoken.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatemnt()
	}
	return expression
}

func (p *Parser) parseGroupedExpresion() ast.Expression {
	p.NextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(gtoken.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(gtoken.TRUE)}
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
