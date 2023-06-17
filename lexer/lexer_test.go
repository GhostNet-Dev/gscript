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

func TestCodeTestSet(t *testing.T) {
	testLexing(t, `
	package tester
	import test
	class test {}
	for () {
		break;
	}
	`, []ExpectedData{
		{gtoken.PACKAGE, "package"}, {gtoken.IDENT, "tester"},
		{gtoken.IMPORT, "import"}, {gtoken.IDENT, "test"},
		{gtoken.CLASS, "class"}, {gtoken.IDENT, "test"}, {gtoken.LBRACE, "{"}, {gtoken.RBRACE, "}"},
		{gtoken.FOR, "for"}, {gtoken.LPAREN, "("}, {gtoken.RPAREN, ")"}, {gtoken.LBRACE, "{"},
		{gtoken.BREAK, "break"}, {gtoken.SEMICOLON, ";"},
		{gtoken.RBRACE, "}"},
		{gtoken.EOF, ""},
	})
}

func TestCodeSet2(t *testing.T) {
	var b bytes.Buffer
	test, expect := testCodeSet()
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
	testLexing(t, b.String(), expect)
}

func testCodeSet() (string, []ExpectedData) {
	return `let five = 5;
	let test = 10;
	let add = fn(x, y) {
		x+y;
	};
	let result = add(five, test);
	"foobar"
	"foo bar"
	[1, 2];
	{"foo": "bar"};
	`, []ExpectedData{
			{gtoken.LET, "let"}, {gtoken.IDENT, "five"}, {gtoken.ASSIGN, "="}, {gtoken.INT, "5"}, {gtoken.SEMICOLON, ";"},
			{gtoken.LET, "let"}, {gtoken.IDENT, "test"}, {gtoken.ASSIGN, "="}, {gtoken.INT, "10"}, {gtoken.SEMICOLON, ";"},
			{gtoken.LET, "let"}, {gtoken.IDENT, "add"}, {gtoken.ASSIGN, "="}, {gtoken.FUNCTION, "fn"},
			{gtoken.LPAREN, "("}, {gtoken.IDENT, "x"}, {gtoken.COMMA, ","}, {gtoken.IDENT, "y"}, {gtoken.RPAREN, ")"},
			{gtoken.LBRACE, "{"}, {gtoken.IDENT, "x"}, {gtoken.PLUS, "+"}, {gtoken.IDENT, "y"}, {gtoken.SEMICOLON, ";"}, {gtoken.RBRACE, "}"}, {gtoken.SEMICOLON, ";"},
			{gtoken.LET, "let"}, {gtoken.IDENT, "result"}, {gtoken.ASSIGN, "="}, {gtoken.IDENT, "add"},
			{gtoken.LPAREN, "("}, {gtoken.IDENT, "five"}, {gtoken.COMMA, ","}, {gtoken.IDENT, "test"}, {gtoken.RPAREN, ")"}, {gtoken.SEMICOLON, ";"},
			{gtoken.STRING, "foobar"}, {gtoken.STRING, "foo bar"},
			{gtoken.LBRACKET, "["}, {gtoken.INT, "1"}, {gtoken.COMMA, ","}, {gtoken.INT, "2"}, {gtoken.RBRACKET, "]"}, {gtoken.SEMICOLON, ";"},
			{gtoken.LBRACE, "{"}, {gtoken.STRING, "foo"}, {gtoken.COLON, ":"}, {gtoken.STRING, "bar"}, {gtoken.RBRACE, "}"}, {gtoken.SEMICOLON, ";"},
			{gtoken.EOF, ""},
		}
}

func TestCodeSet(t *testing.T) {
	input, tests := testCodeSet()
	testLexing(t, input, tests)
}

func TestSimpleSet(t *testing.T) {
	testLexing(t, `=+(){},;`,
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
		})
}

func testLexing(t *testing.T, input string, tests []ExpectedData) {
	t.Helper()
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
