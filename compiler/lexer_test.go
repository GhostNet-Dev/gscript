package compiler

import "testing"

type ExpectedData struct {
	expectedType    TokenType
	expectedLiteral string
}

func CodeTestSet() (string, []ExpectedData) {
	return `let five = 5;
	let test = 10;
	let add = fn(x, t) {
		x+y;
	}
	let result = add(five, test);`, []ExpectedData {
		{LET, "let"}, {IDENT, "five"}, {ASSIGN, "="}, {INT, "5"}, {SEMICOLON, ";"},
		{LET, "let"}, {IDENT, "test"}, {ASSIGN, "="}, {INT, "10"}, {SEMICOLON, ";"},
		{LET, "let"}, {IDENT, "add"}, {ASSIGN, "="}, {FUNCTION, "fn"}, 
		{LPAREN, "("}, {IDENT, "x"}, {COMMA, ","}, {IDENT, "y"}, {RPAREN, ")"},
		{LBRACE, "{"}, {IDENT, "x"}, {PLUS, "+"}, {IDENT, "y"}, {RBRACE, "}"}, {SEMICOLON, ";"},
		{LET, "let"}, {IDENT, "result"}, {ASSIGN, "="}, {IDENT, "add"},
		{LPAREN, "("}, {IDENT, "five"}, {COMMA, ","}, {IDENT, "test"}, {RPAREN, ")"}, {SEMICOLON, ";"},
		{EOF, ""},
	}
}


func SimpleTestSet() (string, []ExpectedData) {
	return `=+(){},;`, 
	[]ExpectedData{
		{ASSIGN, "="},
		{PLUS, "+"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RBRACE, "}"},
		{COMMA, ","},
		{SEMICOLON, ";"},
		{EOF, ""},
	}
}

func TestNextToken(t *testing.T) {
	input, tests := CodeTestSet()

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong, expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
	}
}
