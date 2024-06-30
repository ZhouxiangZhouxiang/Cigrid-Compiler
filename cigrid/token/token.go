package token

type TokenType string 

type Token struct {
	Type    TokenType 
	Literal string
}

// TokenType
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	TVOID   = "TVOID"
	TSTRING = "TSTRING"
	TINT    = "TINT"
	IF     = "IF"
	ELSE   = "ELSE"
	WHILE  = "WHILE"
	RETURN = "RETURN"

	ASSIGN = "="
	BANG = "!"
	PLUS = "+"
	MINUS = "-"
	ASTERISK = "*"
	SLASH = "/"
	ET = "&"
	AND = "&&"
	OR = "||"

	LT = "<"
	GT = ">"
	L_EQ = "<="
	G_EQ = ">="
	EQ = "=="
	NOT_EQ = "!="

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACKET = "["
	RBRACKET = "]"

	COMMA     = ","
	SEMICOLON = ";"
)

var keywords = map[string]TokenType {
	"void": TVOID,
	"string": TSTRING,
	"int": TINT,
	"if": IF,
	"else": ELSE,
	"while": WHILE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tokType, ok := keywords[ident]; ok {
		return tokType 
	}
	return IDENT
}

