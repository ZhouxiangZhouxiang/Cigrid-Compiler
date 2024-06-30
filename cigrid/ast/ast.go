package ast 

import "bytes"
//import "strings"
import "cigrid/token"
import "strconv"


type Node interface {
	TokenLiteral() string 
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type VarDef struct {
	Token     token.Token // type
	Name      *Identifier
	Dimension int // -1 int *, 0 int, x int[x]
	Value     []Expression
}
func (d *VarDef) statementNode()        {}
func (d *VarDef) TokenLiteral() string  { return d.Token.Literal }
func (d *VarDef) String() string {
	var out bytes.Buffer 
	out.WriteString(d.Token.Literal + "[")
	out.WriteString(strconv.Itoa(d.Dimension))
	out.WriteString("]")
	out.WriteString(d.Name.String())
	out.WriteString(" = [ ")
	for _, v := range d.Value {
		out.WriteString(v.String() + " ")
	}
	out.WriteString("]\n")
	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer 
	out.WriteString("BEGIN \n")
	for _, v := range bs.Statements {
		out.WriteString(v.String())
	}
	out.WriteString("END \n")
	return out.String()
}

type Expression interface {
	Node 
	expressionNode()
}

type Identifier struct {
	Token token.Token 
	Value string 
}
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token token.Token 
	Value int64
}
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token 
	Operator string 
	Right    Expression
}
func (pe *PrefixExpression) expressionNode()       {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token 	 token.Token 
	Left 	 Expression 
	Operator string 
	Right 	 Expression 
}
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer 
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}





