package gtoken

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	EQ       = "=="
	NOT_EQ   = "!="

	LT = "<"
	RT = ">"

	COMMA     = ","
	DOT       = "."
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// keyword
	FUNCTION = "FUNCTION"
	LET      = "LET"
	CONST    = "CONST"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	STRING   = "STRING"
	FOR      = "FOR"
	BREAK    = "BREAK"
	CLASS    = "CLASS"
	IMPORT   = "IMPORT"
	PACKAGE  = "PACKAGE"
	TYPE     = "TYPE"
	STRUCT    = "STRUCT"
)

var keywords = map[string]TokenType{
	"fn":      FUNCTION,
	"let":     LET,
	"const":   CONST,
	"true":    TRUE,
	"false":   FALSE,
	"if":      IF,
	"else":    ELSE,
	"return":  RETURN,
	"for":     FOR,
	"break":   BREAK,
	"import":  IMPORT,
	"package": PACKAGE,
	"type":    TYPE,
	"struct":    STRUCT,
}

func NewToken(tokenType TokenType, ch byte, line int) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: line}
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
