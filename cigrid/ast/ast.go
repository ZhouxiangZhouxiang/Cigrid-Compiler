package ast 

import "bytes"
import "strings"
import "cigrid/token"
import "strconv"

var ident = 0

func identFunc() string {
	var out bytes.Buffer
	for i := 0; i < ident; i++ {
		out.WriteString("    ")
	} 
	return out.String()
}

type Node interface {
	String() string
}

type Global interface {
	Node 
	GlobalNode()
}

type FunctionLiteral struct {
	ReturnType *Type
	Name       *Identifier 
	Param      []*TypeIdentifierPair
	Body       *BlockStatement
}
func (fl *FunctionLiteral) GlobalNode() {}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer 
	out.WriteString(fl.ReturnType.String())
	out.WriteString(fl.Name.String())
	params := []string{}
	for _, v := range fl.Param {
		params = append(params, v.String())
	}
	out.WriteString("(" + strings.Join(params, ", ") +")")
	out.WriteString("\n")
	out.WriteString(fl.Body.String())
	return out.String()
}

type Statement interface {
	Node
	statementNode()
}

type VarDef struct {
	VarType *Type   
	Name    *Identifier
	Value   Expression
}
func (d *VarDef) statementNode() {}
func (d *VarDef) String() string {
	var out bytes.Buffer 
	out.WriteString(d.VarType.String())
	out.WriteString(d.Name.String())
	out.WriteString(" = ")
	out.WriteString(d.Value.String())
	return out.String()
}

type VarAssign struct {
	Left Expression 
	Right Expression
}
func (va *VarAssign) statementNode() {}
func (va *VarAssign) String() string {
	var out bytes.Buffer 
	out.WriteString(va.Left.String())
	out.WriteString(" = ")
	out.WriteString(va.Right.String())
	return out.String()
}

type ReturnStatement struct {
	ReturnValue Expression 
}
func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer 
	out.WriteString("return ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	return out.String()
}

type IfStatement struct {
	Condition   Expression 
	Consequence *BlockStatement
	Alternative *BlockStatement
}
func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string {
	var out bytes.Buffer 
	out.WriteString("if (")
	out.WriteString(is.Condition.String())
	out.WriteString(")\n")
	out.WriteString(is.Consequence.String())
	if is.Alternative != nil {
		out.WriteString(identFunc())
		out.WriteString("else\n")
		out.WriteString(is.Alternative.String())
	}
	return out.String()
}

type WhileStatement struct {
	Condition   Expression 
	Consequence *BlockStatement 
}
func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) String() string {
	var out bytes.Buffer 
	out.WriteString("while (")
	out.WriteString(ws.Condition.String())
	out.WriteString(")\n")
	out.WriteString(ws.Consequence.String())
	return out.String()
}

type CallStatement struct {
	Value *CallExpression
}
func (cs *CallStatement) statementNode() {}
func (cs *CallStatement) String() string { return cs.Value.String() }

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}
func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer 
	out.WriteString(identFunc())
	ident++
	out.WriteString("BEGIN \n")
	for _, v := range bs.Statements {
		out.WriteString(identFunc())
		out.WriteString(v.String())
		out.WriteString("\n")
	}
	ident--
	out.WriteString(identFunc())
	out.WriteString("END")
	return out.String()
}

type Expression interface {
	Node 
	expressionNode()
}

type Identifier struct {
	Value token.Token
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value.Literal }

type IntegerLiteral struct {
	Value token.Token
}
func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string  { return il.Value.Literal }

type StringLiteral struct {
	Value token.Token
}
func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string { return "\"" + sl.Value.Literal + "\"" }

type ArrayLiteral struct {
	Elements []Expression
}
func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) String() string  {
	var out bytes.Buffer 
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type PrefixExpression struct {
	Operator token.Token
	Right    Expression
}
func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator.Literal)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Left 	 Expression 
	Operator token.Token
	Right 	 Expression 
}
func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer 
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator.Literal + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

type IndexExpression struct {
	Name *Identifier
	Index Expression
}
func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) String() string {
	var out bytes.Buffer 
	out.WriteString("(")
	out.WriteString(ie.Name.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	out.WriteString(")")
	return out.String()
}

type CallExpression struct {
	Name *Identifier
	Params []Expression 
}
func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	var out bytes.Buffer 
	out.WriteString(ce.Name.String())
	out.WriteString("(")
	params := []string{}
	for _, v := range ce.Params {
		params = append(params, v.String())
	} 
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	return out.String()
}


type Type struct {
	Dtype     token.Token 
	Dimension int
}
func (t *Type) String() string {
	var out bytes.Buffer
	out.WriteString(t.Dtype.Literal + " ")
	out.WriteString("[" + strconv.Itoa(t.Dimension) + "] ")
	return out.String()
}

type TypeIdentifierPair struct {
	TypeLiteral       *Type 
	IdentifierLiteral *Identifier
}
func (tip *TypeIdentifierPair) String() string {
	var out bytes.Buffer 
	out.WriteString(tip.TypeLiteral.String())
	out.WriteString(tip.IdentifierLiteral.String())
	return out.String()
}

type ProgramLiteral struct {
	GlobalList []Global
}
func (pl *ProgramLiteral) String() string {
	var out bytes.Buffer 
	for _, v := range pl.GlobalList {
		out.WriteString(v.String())
		out.WriteString("\n")
	}
	return out.String()
}