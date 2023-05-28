package compiler

type Lexer struct {
	input            string
	position         int
	nextReadPosition int
	ch               byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.nextReadPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.nextReadPosition]
	}
	l.position = l.nextReadPosition
	l.nextReadPosition++
}

func (l *Lexer) NextToken() Token {
	var tok Token

	switch l.ch {
	case '=':
		tok = NewToken(ASSIGN, l.ch)
	case ';':
		tok = NewToken(SEMICOLON, l.ch)
	case '(':
		tok = NewToken(LPAREN, l.ch)
	case ')':
		tok = NewToken(RPAREN, l.ch)
	case ',':
		tok = NewToken(COMMA, l.ch)
	case '+':
		tok = NewToken(PLUS, l.ch)
	case '{':
		tok = NewToken(LBRACE, l.ch)
	case '}':
		tok = NewToken(RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	}

	l.readChar()
	return tok
}
