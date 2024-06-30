package parser 

import "cigrid/token"
import "cigrid/ast"
import "strconv"

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // + -
	PRODUCT     // * /
	PREFIX      // -X !X *X &X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int {
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

func lookupPrecedence(tokType token.TokenType) int {
	if p, ok := precedences[tokType]; ok {
		return p
	}
	return LOWEST
}

type Parser struct {
	tokList   []token.Token
	position  int 
	curToken  token.Token 
	peekToken token.Token
} 


func New(tokList []token.Token) *Parser {
	p := &Parser{tokList: tokList, position: 0, curToken: tokList[0]}
	if len(tokList) == 1 {
		p.peekToken = tokList[0]
	} else {
		p.peekToken = tokList[1]
	}
	return p
}

func (p *Parser) nextToken() {
	p.position = p.position + 1
	eofToken := token.Token{Type: token.EOF, Literal: ""}
	if p.position >= len(p.tokList) {
		p.curToken = eofToken
		p.peekToken = eofToken
	} else if p.position + 1 == len(p.tokList) {
		p.peekToken = eofToken
	} else {
		p.curToken = p.tokList[p.position]
		p.peekToken = p.tokList[p.position + 1]
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	if p.curToken.Type == token.IDENT {
		return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	} else if p.curToken.Type == token.INT {
		val, _ := strconv.Atoi(p.curToken.Literal)
		return &ast.IntegerLiteral{Token: p.curToken, Value: int64(val)}
	} else if p.curToken.Type == token.LPAREN {
		p.nextToken()
		expression := p.parseExpression(LOWEST)
		p.nextToken()
		return expression
	} else {
		// !, - , &, *
		expression := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
		p.nextToken()
		expression.Right = p.parseExpression(PREFIX)
		return expression
	}
	return &ast.IntegerLiteral{}
}


func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token: p.curToken, Operator: p.curToken.Literal, Left: left}
	precedence := precedences[p.curToken.Type]
	p.nextToken()
	right := p.parseExpression(precedence)
	expression.Right = right 
	return expression
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	leftExp := p.parsePrefixExpression()
	for precedence < lookupPrecedence(p.peekToken.Type) {
		p.nextToken()
		leftExp = p.parseInfixExpression(leftExp)
	} 
	return leftExp
}

func (p *Parser) parseDefineStatement() ast.Statement {
	statement := &ast.VarDef{Token: p.curToken, Dimension: 0}
	if p.peekToken.Type == token.ASTERISK {
		statement.Dimension = -1
		p.nextToken()
	}
	p.nextToken()
	statement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if p.peekToken.Type == token.LBRACKET {
		p.nextToken()
		p.nextToken()
		val, _ := strconv.Atoi(p.curToken.Literal)
		statement.Dimension = val
		p.nextToken()
	}
	p.nextToken()

	return statement

}

func (p *Parser) parseStatement() ast.Statement {
	var statement ast.Statement
	if p.curToken.Type == token.TSTRING || p.curToken.Type == token.TINT {
		statement = p.parseDefineStatement()
	}
	return statement
}

func (p *Parser) ParseProgram() ast.Statement {
	//leftExp := p.parseExpression(LOWEST)
	statement := p.parseStatement()
	return statement
}
