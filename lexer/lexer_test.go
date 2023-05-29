package lexer

import (
	"bytes"
	"testing"

	"github.com/GhostNet-Dev/glambda/gtoken"
)

type ExpectedData struct {
	expectedType    gtoken.TokenType
	expectedLiteral string
}

func CodeTestSet2() (string, []ExpectedData) {
	var b bytes.Buffer
	test, expect := CodeTestSet()
	expectLen := len(expect)
	b.WriteString(test)
	b.WriteString(`!-/*5;
		5 < 10 > 5;
		if (5 < 10) {
			return true;
		} else {
			return false;
		}
		10 == 10;
		10 != 9;`)

	newExp := []ExpectedData{
		{gtoken.BANG, "!"},
		{gtoken.MINUS, "-"},
		{gtoken.SLASH, "/"},
		{gtoken.ASTERISK, "*"},
		{gtoken.INT, "5"},
		{gtoken.SEMICOLON, ";"},
		{gtoken.INT, "5"},
		{gtoken.LT, "<"},
		{gtoken.INT, "10"},
		{gtoken.RT, ">"},
		{gtoken.INT, "5"},
		{gtoken.SEMICOLON, ";"},
		{gtoken.IF, "if"},
		{gtoken.LPAREN, "("},
		{gtoken.INT, "5"},
		{gtoken.LT, "<"},
		{gtoken.INT, "10"},
		{gtoken.RPAREN, ")"},
		{gtoken.LBRACE, "{"},
		{gtoken.RETURN, "return"},
		{gtoken.TRUE, "true"},
		{gtoken.SEMICOLON, ";"},
		{gtoken.RBRACE, "}"},
		{gtoken.ELSE, "else"},
		{gtoken.LBRACE, "{"},
		{gtoken.RETURN, "return"},
		{gtoken.FALSE, "false"},
		{gtoken.SEMICOLON, ";"},
		{gtoken.RBRACE, "}"},
		{gtoken.INT, "10"},
		{gtoken.EQ, "=="},
		{gtoken.INT, "10"},
		{gtoken.SEMICOLON, ";"},
		{gtoken.INT, "10"},
		{gtoken.NOT_EQ, "!="},
		{gtoken.INT, "9"},
		{gtoken.SEMICOLON, ";"},
	}
	expect = append(expect[:expectLen-1], newExp...)
	return b.String(), expect
}

func CodeTestSet() (string, []ExpectedData) {
	return `let five = 5;
	let test = 10;
	let add = fn(x, y) {
		x+y;
	};
	let result = add(five, test);`, []ExpectedData{
			{gtoken.LET, "let"}, {gtoken.IDENT, "five"}, {gtoken.ASSIGN, "="}, {gtoken.INT, "5"}, {gtoken.SEMICOLON, ";"},
			{gtoken.LET, "let"}, {gtoken.IDENT, "test"}, {gtoken.ASSIGN, "="}, {gtoken.INT, "10"}, {gtoken.SEMICOLON, ";"},
			{gtoken.LET, "let"}, {gtoken.IDENT, "add"}, {gtoken.ASSIGN, "="}, {gtoken.FUNCTION, "fn"},
			{gtoken.LPAREN, "("}, {gtoken.IDENT, "x"}, {gtoken.COMMA, ","}, {gtoken.IDENT, "y"}, {gtoken.RPAREN, ")"},
			{gtoken.LBRACE, "{"}, {gtoken.IDENT, "x"}, {gtoken.PLUS, "+"}, {gtoken.IDENT, "y"}, {gtoken.SEMICOLON, ";"}, {gtoken.RBRACE, "}"}, {gtoken.SEMICOLON, ";"},
			{gtoken.LET, "let"}, {gtoken.IDENT, "result"}, {gtoken.ASSIGN, "="}, {gtoken.IDENT, "add"},
			{gtoken.LPAREN, "("}, {gtoken.IDENT, "five"}, {gtoken.COMMA, ","}, {gtoken.IDENT, "test"}, {gtoken.RPAREN, ")"}, {gtoken.SEMICOLON, ";"},
			{gtoken.EOF, ""},
		}
}

func SimpleTestSet() (string, []ExpectedData) {
	return `=+(){},;`,
		[]ExpectedData{
			{gtoken.ASSIGN, "="},
			{gtoken.PLUS, "+"},
			{gtoken.LPAREN, "("},
			{gtoken.RPAREN, ")"},
			{gtoken.LBRACE, "{"},
			{gtoken.RBRACE, "}"},
			{gtoken.COMMA, ","},
			{gtoken.SEMICOLON, ";"},
			{gtoken.EOF, ""},
		}
}

func TestNextToken(t *testing.T) {
	input, tests := CodeTestSet2()

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextTokenMake()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong, expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
