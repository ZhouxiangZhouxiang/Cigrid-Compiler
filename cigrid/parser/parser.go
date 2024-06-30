package parser 

import "cigrid/token"
import "cigrid/ast"
import "strconv"
//import "fmt"

func lookupPrecedence(tokType token.TokenType) int {
	if tokType == token.ASTERISK || tokType == token.SLASH { // * /
		return 8
	} else if tokType == token.PLUS || tokType == token.MINUS { // 
		return 7
	} else if tokType == token.LT || tokType == token.GT || 
			  tokType == token.L_EQ || tokType == token.G_EQ {
		return 6		
	} else if tokType == token.EQ || tokType == token.NOT_EQ {
		return 5
	} else if tokType == token.AND {
		return 4
	} else if tokType == token.OR {
		return 3
	}
	return 0
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
		p.curToken = p.tokList[p.position]
		p.peekToken = eofToken
	} else {
		p.curToken = p.tokList[p.position]
		p.peekToken = p.tokList[p.position + 1]
	}
}

func (p *Parser) parseIndexExpression() ast.Expression {
	ie := &ast.IndexExpression{}
	ie.Name = &ast.Identifier{Value: p.curToken}
	p.nextToken()
	p.nextToken()
	ie.Index = p.parseExpression(0)
	p.nextToken()
	return ie
}

func (p *Parser) parseCallExpression() ast.Expression {
	ce := &ast.CallExpression{}
	ce.Name = &ast.Identifier{Value: p.curToken}
	p.nextToken()
	list := []ast.Expression{}
	if p.peekToken.Type == token.RPAREN {
		ce.Params = list 
		p.nextToken()
		return ce 
	}
	p.nextToken()
	list = append(list, p.parseExpression(0))
	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(0))
	}
	ce.Params = list
	p.nextToken()
	return ce
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	if p.curToken.Type == token.IDENT {
		if p.peekToken.Type == token.LBRACKET {
			// a[0]，indexExpression
			return p.parseIndexExpression()
		} else if p.peekToken.Type == token.LPAREN {
			// printf("hello")，
			return p.parseCallExpression()
		}
		// i
		return &ast.Identifier{Value: p.curToken}
	} else if p.curToken.Type == token.INT {
		return &ast.IntegerLiteral{Value: p.curToken}
	} else if p.curToken.Type == token.STRING {
		return &ast.StringLiteral{Value: p.curToken}
	} else if p.curToken.Type == token.LPAREN {
		p.nextToken()
		expression := p.parseExpression(0)
		p.nextToken()
		return expression
	} else if p.curToken.Type == token.LBRACE {
		// {1, 2, 3, 4}
		// 空数组不支持
		list := []ast.Expression{}
		p.nextToken()
		list = append(list, p.parseExpression(0))
		for p.peekToken.Type == token.COMMA {
			p.nextToken()
			p.nextToken()
			list = append(list, p.parseExpression(0))
		}
		p.nextToken()
		return &ast.ArrayLiteral{Elements: list}
	} else {
		// !, - , &, *
		expression := &ast.PrefixExpression{Operator: p.curToken}
		p.nextToken()
		expression.Right = p.parsePrefixExpression()
		return expression
	}
	return &ast.IntegerLiteral{}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{Operator: p.curToken, Left: left}
	precedence := lookupPrecedence(p.curToken.Type)
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



func (p *Parser) parseVarDefStatement() ast.Statement {
	statement := &ast.VarDef{}
	statement.VarType = p.parseType()
	p.nextToken()
	statement.Name = &ast.Identifier{Value: p.curToken}
	p.nextToken()
	p.nextToken()
	statement.Value = p.parseExpression(0)
	p.nextToken()
	return statement
}

func (p *Parser) parseVarAssignStatement() ast.Statement {
	statement := &ast.VarAssign{}
	statement.Left = p.parseExpression(0)
	p.nextToken()
	p.nextToken()
	statement.Right = p.parseExpression(0)
	p.nextToken()
	return statement
}

func (p *Parser) parseReturnStatement() ast.Statement {
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
		rs := &ast.ReturnStatement{ReturnValue: nil}
		return rs
	}
	p.nextToken()
	rs := &ast.ReturnStatement{ReturnValue: p.parseExpression(0)}
	p.nextToken()
	return rs
}

func (p *Parser) parseIfStatement() ast.Statement {
	ifstat := &ast.IfStatement{}
	p.nextToken()
	ifstat.Condition = p.parseExpression(0)
	p.nextToken()
	ifstat.Consequence = p.parseBlockStatement()
	if p.peekToken.Type == token.ELSE {
		p.nextToken()
		p.nextToken()
		ifstat.Alternative = p.parseBlockStatement()
	}
	return ifstat
}

func (p *Parser) parseWhileStatement() ast.Statement {
	whilestat := &ast.WhileStatement{}
	p.nextToken()
	whilestat.Condition = p.parseExpression(0)
	p.nextToken()
	whilestat.Consequence = p.parseBlockStatement()
	return whilestat
}

func (p *Parser) parseStatement() ast.Statement {
	var statement ast.Statement
	if p.curToken.Type == token.TSTRING || p.curToken.Type == token.TINT {
		statement = p.parseVarDefStatement()
	} else if p.curToken.Type == token.RETURN {
		statement = p.parseReturnStatement()
	} else if p.curToken.Type == token.IF {
		statement = p.parseIfStatement()
	} else if p.curToken.Type == token.WHILE {
		statement = p.parseWhileStatement()
	} else if p.curToken.Type == token.IDENT && p.peekToken.Type == token.LPAREN {
		// printf("hello\n");
		temp, _ := (p.parseCallExpression()).(*ast.CallExpression)
		statement = &ast.CallStatement{Value: temp}
		p.nextToken()
	} else if p.curToken.Type == token.IDENT {
		// x = 1;
		statement = p.parseVarAssignStatement()
	} else if p.curToken.Type == token.ASTERISK {
		// *x = 1;
		statement = p.parseVarAssignStatement()
	}
	return statement
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken, Statements: []ast.Statement{}} 
	p.nextToken()
	for p.curToken.Type != token.RBRACE {
		statement := p.parseStatement()
		block.Statements = append(block.Statements, statement)
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFunctionLiteral() ast.Global {
	fl := &ast.FunctionLiteral{}
	fl.ReturnType = p.parseType()
	p.nextToken()
	fl.Name = &ast.Identifier{Value: p.curToken}
	p.nextToken()
	params := []*ast.TypeIdentifierPair{}
	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		fl.Param = params
		p.nextToken()
		fl.Body = p.parseBlockStatement()
		return fl
	} else {
		p.nextToken()
		param := &ast.TypeIdentifierPair{}
		param.TypeLiteral = p.parseType()
		p.nextToken()
		param.IdentifierLiteral = &ast.Identifier{Value: p.curToken}
		params = append(params, param)
		for p.peekToken.Type == token.COMMA {
			p.nextToken()
			p.nextToken()
			param = &ast.TypeIdentifierPair{}
			param.TypeLiteral = p.parseType()
			p.nextToken()
			param.IdentifierLiteral = &ast.Identifier{Value: p.curToken}
			params = append(params, param)
		}
		p.nextToken()
		fl.Param = params
		p.nextToken()
		fl.Body = p.parseBlockStatement()
		return fl
	}
	return nil
}

func (p *Parser) ParseProgram() *ast.ProgramLiteral {
	program := &ast.ProgramLiteral{}
	global := []ast.Global{} 
	for p.curToken.Type != token.EOF {
		global = append(global, p.parseFunctionLiteral())
		p.nextToken()
	}
	program.GlobalList = global
	return program
}

func (p *Parser) parseType() *ast.Type {
	varType := &ast.Type{Dtype: p.curToken}
	if p.peekToken.Type == token.ASTERISK {
		p.nextToken()
		varType.Dimension = -1
	} else if p.peekToken.Type == token.LBRACKET {
		p.nextToken()
		p.nextToken()
		varType.Dimension, _ = strconv.Atoi(p.curToken.Literal)
		p.nextToken()
	} else {
		varType.Dimension = 0
	}
	return varType
}
