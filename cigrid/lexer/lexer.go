package lexer

import "cigrid/token"

type Lexer struct {
	input    string
	position int 
	ch       byte
	peekCh   byte 
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func New(input string) *Lexer{
	l := &Lexer{input: input, position: 0}
	if len(input) == 0 {
		l.ch = 0
		l.peekCh = 0
	} else if len(input) == 1 {
		l.ch = input[0]
		l.peekCh = 0
	} else {
		l.ch = input[0]
		l.peekCh = input[1]
	}
	return l
}

func (l *Lexer) readChar() {
	l.position += 1
	if l.position >= len(l.input) {
		l.ch = 0
		l.peekCh = 0
	} else if l.position + 1 == len(l.input) {
		l.ch = l.input[l.position]
		l.peekCh = 0
	} else {
		l.ch = l.input[l.position]
		l.peekCh = l.input[l.position + 1]
	}
}

func (l *Lexer) nextToken() token.Token {
	l.skipWhiteSpace()
	var tok token.Token 
	switch l.ch {
	case '=': 
		if l.peekCh == '=' {
			tok.Type = token.EQ
			tok.Literal = "=="
			l.readChar()
		} else {
			tok.Type = token.ASSIGN
			tok.Literal = "="
		}
	case '!':
		if l.peekCh == '=' {
			tok.Type = token.NOT_EQ
			tok.Literal = "!="
			l.readChar()
		} else {
			tok.Type = token.BANG
			tok.Literal = "!"
		}
	case '+':
		tok.Type = token.PLUS
		tok.Literal = "+"
	case '-':
		tok.Type = token.MINUS
		tok.Literal = "-"
	case '*':
		tok.Type = token.ASTERISK
		tok.Literal = "*"
	case '/':
		tok.Type = token.SLASH
		tok.Literal = "/"
	case '&':
		if l.peekCh == '&' {
			tok.Type = token.AND
			tok.Literal = "&&"
			l.readChar()
		} else {
			tok.Type = token.ET 
			tok.Literal = "&"
		}
	case '|':
		if l.peekCh == '|' {
			tok.Type = token.OR 
			tok.Literal = "||"
			l.readChar()
		}
	case '<':
		if l.peekCh == '=' {
			tok.Type = token.L_EQ
			tok.Literal = "<="
			l.readChar()
		} else {
			tok.Type = token.LT
			tok.Literal = "<"
		}
	case '>':
		if l.peekCh == '=' {
			tok.Type = token.G_EQ 
			tok.Literal = ">="
			l.readChar()
		} else {
			tok.Type = token.GT
			tok.Literal = ">"
		}
	case '(':
		tok.Type = token.LPAREN
		tok.Literal = "("
	case ')':
		tok.Type = token.RPAREN
		tok.Literal = ")"
	case '[':
		tok.Type = token.LBRACKET
		tok.Literal = "["
	case ']':
		tok.Type = token.RBRACKET
		tok.Literal = "]"
	case '{':
		tok.Type = token.LBRACE
		tok.Literal = "{"
	case '}':
		tok.Type = token.RBRACE
		tok.Literal = "}"
	case ',':
		tok.Type = token.COMMA
		tok.Literal = ","
	case ';':
		tok.Type = token.SEMICOLON
		tok.Literal = ";"
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdent()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok 
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT 
			return tok
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) Scan() []token.Token {
	result := []token.Token{}
	tok := l.nextToken()
	result = append(result, tok)
	for tok.Type != token.EOF {
		tok = l.nextToken()
		result = append(result, tok)
	} 
	return result
}

func (l *Lexer)readIdent() string {
	position := l.position 
	for (isLetter(l.ch)) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer)readNumber() string {
	position := l.position 
	for (isDigit(l.ch)) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer)readString() string {
	position := l.position + 1
	l.readChar()
	for l.ch != '"' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer)skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}