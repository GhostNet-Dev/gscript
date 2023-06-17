package lexer

import "github.com/GhostNet-Dev/glambda/gtoken"

type Lexer struct {
	input            string
	position         int
	nextReadPosition int
	ch               byte
	line             int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1,
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

func (l *Lexer) NextTokenMake() gtoken.Token {
	var tok gtoken.Token
	l.skipWhiteSpace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = gtoken.Token{Type: gtoken.EQ, Literal: literal}
		} else {
			tok = gtoken.NewToken(gtoken.ASSIGN, l.ch, l.line)
		}
	case '+':
		tok = gtoken.NewToken(gtoken.PLUS, l.ch, l.line)
	case '-':
		tok = gtoken.NewToken(gtoken.MINUS, l.ch, l.line)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = gtoken.Token{Type: gtoken.NOT_EQ, Literal: literal}
		} else {
			tok = gtoken.NewToken(gtoken.BANG, l.ch, l.line)
		}
	case '/':
		tok = gtoken.NewToken(gtoken.SLASH, l.ch, l.line)
	case '*':
		tok = gtoken.NewToken(gtoken.ASTERISK, l.ch, l.line)
	case '<':
		tok = gtoken.NewToken(gtoken.LT, l.ch, l.line)
	case '>':
		tok = gtoken.NewToken(gtoken.RT, l.ch, l.line)
	case ';':
		tok = gtoken.NewToken(gtoken.SEMICOLON, l.ch, l.line)
	case ':':
		tok = gtoken.NewToken(gtoken.COLON, l.ch, l.line)
	case '(':
		tok = gtoken.NewToken(gtoken.LPAREN, l.ch, l.line)
	case ')':
		tok = gtoken.NewToken(gtoken.RPAREN, l.ch, l.line)
	case ',':
		tok = gtoken.NewToken(gtoken.COMMA, l.ch, l.line)
	case '{':
		tok = gtoken.NewToken(gtoken.LBRACE, l.ch, l.line)
	case '}':
		tok = gtoken.NewToken(gtoken.RBRACE, l.ch, l.line)
	case '[':
		tok = gtoken.NewToken(gtoken.LBRACKET, l.ch, l.line)
	case ']':
		tok = gtoken.NewToken(gtoken.RBRACKET, l.ch, l.line)
	case '"':
		tok.Type = gtoken.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = gtoken.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = gtoken.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = gtoken.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = gtoken.NewToken(gtoken.ILLEGAL, l.ch, l.line)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.nextReadPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextReadPosition]
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
