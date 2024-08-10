package ir_translator 

import "cigrid/token"
import "cigrid/ast"
import "cigrid/ir"
import "strconv"

type irTranslator struct {
	tree         *ast.ProgramLiteral
	irList       []ir.IntermediateRepresentation
	tempRegister int
	variableMap  map[string]int
	addressMap   map[string]int
	rsp          int
}

func New(tree *ast.ProgramLiteral) *irTranslator {
	return &irTranslator{
		tree: tree, 
		irList: []ir.IntermediateRepresentation{}, 
		tempRegister: 0,
		variableMap: make(map[string]int),
		addressMap: make(map[string]int),
		rsp: 0,
	}
}

func (t *irTranslator) translateStatementBlock(bs *ast.BlockStatement) {
	tempVariable := []string{}
	for _, v := range bs.Statements {
		temp := t.translateStatement(v)
		tempVariable = append(tempVariable, temp)
	}
	for _, v := range tempVariable {
		t.variableMap[v]--
	}
}

func (t *irTranslator) translateStatement(statement ast.Statement) string {
	if stmt, ok := statement.(*ast.VarAssign); ok {
		t.tempRegister = 0
		var op1_temp string 
		if id, ok := stmt.Left.(*ast.Identifier); ok {
			op1_temp = id.Value.Literal + strconv.Itoa(t.variableMap[id.Value.Literal])
		}
		ir_temp := ir.CalcInst{
			Operation: ir.MOV,
			Operand1: op1_temp,
			Operand2: t.translateExpression(stmt.Right),
		}
		t.irList = append(t.irList, ir_temp)
	} else if stmt, ok := statement.(*ast.VarDef); ok {
		t.tempRegister = 0
		varName := stmt.Name.Value.Literal
		if _, ok := t.variableMap[varName]; ok {
			t.variableMap[varName]++
		} else {
			t.variableMap[varName] = 1
		}
		varNameNew := varName + strconv.Itoa(t.variableMap[varName])
		ir_temp := ir.CalcInst{
			Operation: ir.MOV,
			Operand1: varNameNew,
			Operand2: t.translateExpression(stmt.Value),
		}
		t.irList = append(t.irList, ir_temp)
		t.addressMap[varNameNew] = t.rsp
		t.rsp++
		return varName
	} else if stmt, ok := statement.(*ast.IfStatement); ok {
		t.translateCondition(stmt.Condition, "r", "if", "else")
	}
	return ""
}

func (t *irTranslator) translateCondition(expression ast.Expression,
										  curNode string,
										  trueNode string, 
										  falseNode string) {
	if exp, ok := expression.(*ast.InfixExpression); ok {
		label := ir.Label(curNode)
		t.irList = append(t.irList, label)
		infix_temp := exp.Operator.Type
		switch infix_temp {
		case token.LT: 
			reg1 := t.translateExpression(exp.Left)
			reg2 := t.translateExpression(exp.Right)
			ci := ir.CmpInst{Left: reg1, Right: reg2}
			t.irList = append(t.irList, ci)
			ji1 := ir.JumpInst{JC: ir.L, Addr: trueNode}
			t.irList = append(t.irList, ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irList = append(t.irList, ji2)
		case token.GT:
			reg1 := t.translateExpression(exp.Left)
			reg2 := t.translateExpression(exp.Right)
			ci := ir.CmpInst{Left: reg1, Right: reg2}
			t.irList = append(t.irList, ci)
			ji1 := ir.JumpInst{JC: ir.G, Addr: trueNode}
			t.irList = append(t.irList, ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irList = append(t.irList, ji2)
		case token.L_EQ:
			reg1 := t.translateExpression(exp.Left)
			reg2 := t.translateExpression(exp.Right)
			ci := ir.CmpInst{Left: reg1, Right: reg2}
			t.irList = append(t.irList, ci)
			ji1 := ir.JumpInst{JC: ir.LE, Addr: trueNode}
			t.irList = append(t.irList, ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irList = append(t.irList, ji2)
		case token.G_EQ:
			reg1 := t.translateExpression(exp.Left)
			reg2 := t.translateExpression(exp.Right)
			ci := ir.CmpInst{Left: reg1, Right: reg2}
			t.irList = append(t.irList, ci)
			ji1 := ir.JumpInst{JC: ir.GE, Addr: trueNode}
			t.irList = append(t.irList, ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irList = append(t.irList, ji2)
		case token.EQ:
			reg1 := t.translateExpression(exp.Left)
			reg2 := t.translateExpression(exp.Right)
			ci := ir.CmpInst{Left: reg1, Right: reg2}
			t.irList = append(t.irList, ci)
			ji1 := ir.JumpInst{JC: ir.E, Addr: trueNode}
			t.irList = append(t.irList, ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irList = append(t.irList, ji2)
		case token.NOT_EQ:
			reg1 := t.translateExpression(exp.Left)
			reg2 := t.translateExpression(exp.Right)
			ci := ir.CmpInst{Left: reg1, Right: reg2}
			t.irList = append(t.irList, ci)
			ji1 := ir.JumpInst{JC: ir.NE, Addr: trueNode}
			t.irList = append(t.irList, ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irList = append(t.irList, ji2)
		case token.AND: 
			leftNode := curNode + "0"
			rightNode := curNode + "1"
			t.translateCondition(exp.Left, leftNode, rightNode, falseNode)
			t.translateCondition(exp.Right, rightNode, trueNode, falseNode)
		case token.OR:
			leftNode := curNode + "0"
			rightNode := curNode + "1"
			t.translateCondition(exp.Left, leftNode, trueNode, rightNode)
			t.translateCondition(exp.Right, rightNode, trueNode, falseNode)
		default:
			reg1 := t.translateExpression(exp)
			ir_temp := ir.CalcInst{
				Operation: ir.MOV, 
				Operand1: t.tempRegister, 
				Operand2: "0",
			}
			t.irList = append(t.irList, ir_temp)
			t.tempRegister++ 
			ci := ir.CmpInst{Left: reg1, Right: t.tempRegister - 1}
			t.irList = append(t.irList, ci)
			ji1 := ir.JumpInst{JC: ir.NE, Addr: trueNode}
			t.irList = append(t.irList, ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irList = append(t.irList, ji2)
		}
	} else if exp, ok := expression.(*ast.PrefixExpression); ok {
		if exp.Operator.Type == token.BANG {
			t.translateCondition(exp.Right, curNode, falseNode, trueNode)
		} else {
			reg1 := t.translateExpression(exp.Right)
			ir_temp := ir.CalcInst{
				Operation: ir.MOV, 
				Operand1: t.tempRegister, 
				Operand2: "0",
			}
			t.irList = append(t.irList, ir_temp)
			t.tempRegister++ 
			ci := ir.CmpInst{Left: reg1, Right: t.tempRegister - 1}
			t.irList = append(t.irList, ci)
			ji1 := ir.JumpInst{JC: ir.NE, Addr: trueNode}
			t.irList = append(t.irList, ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irList = append(t.irList, ji2)
		}
	} else {
		reg1 := t.translateExpression(expression)
		ir_temp := ir.CalcInst{
			Operation: ir.MOV, 
			Operand1: t.tempRegister, 
			Operand2: "0",
		}
		t.irList = append(t.irList, ir_temp)
		t.tempRegister++ 
		ci := ir.CmpInst{Left: reg1, Right: t.tempRegister - 1}
		t.irList = append(t.irList, ci)
		ji1 := ir.JumpInst{JC: ir.NE, Addr: trueNode}
		t.irList = append(t.irList, ji1)
		ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
		t.irList = append(t.irList, ji2)
	}
}

func (t *irTranslator) translateExpression(expression ast.Expression) int {
	if exp, ok := expression.(*ast.IntegerLiteral); ok {
		ir_temp := ir.CalcInst{
			Operation: ir.MOV, 
			Operand1: t.tempRegister, 
			Operand2: exp.Value.Literal,
		}
		t.irList = append(t.irList, ir_temp)
		t.tempRegister++ 
		return t.tempRegister - 1
	} else if exp, ok := expression.(*ast.Identifier); ok {
		ir_temp := ir.CalcInst{
			Operation: ir.MOV,
			Operand1: t.tempRegister,
			Operand2: exp.Value.Literal + strconv.Itoa(t.variableMap[exp.Value.Literal]),
		}
		t.irList = append(t.irList, ir_temp)
		t.tempRegister++
		return t.tempRegister - 1
	} else if exp, ok := expression.(*ast.InfixExpression); ok {
		var infix_temp ir.Op
		switch exp.Operator.Type {
		case token.PLUS: 
			infix_temp = ir.ADD
		case token.MINUS:
			infix_temp = ir.SUB
		default:
		}
		o1 := t.translateExpression(exp.Left)
		o2 := t.translateExpression(exp.Right)
		ir_temp := ir.CalcInst{
			Operation: infix_temp, 
			Operand1: o1,
			Operand2: o2,
		}
		t.irList = append(t.irList, ir_temp)
		return o1
	}
	return 0
}

func (t *irTranslator) Translate() []ir.IntermediateRepresentation {
	function := t.tree.GlobalList[0]
	functionff, _ := function.(*ast.FunctionLiteral)
	statement := functionff.Body
	t.translateStatementBlock(statement)
	
	
	return t.irList
}