package ir_translator 

import "cigrid/ast"
import "cigrid/ir"

type irTranslator struct {
	tree        *ast.ProgramLiteral
	irList      []ir.IntermediateRepresentation
	tempRegister int
}

func New(tree *ast.ProgramLiteral) *irTranslator {
	return &irTranslator{
		tree: tree, 
		irList: []ir.IntermediateRepresentation{}, 
		tempRegister: 0,
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
			Operand2: exp.Value.Literal,
		}
		t.irList = append(t.irList, ir_temp)
		t.tempRegister++
		return t.tempRegister - 1
	} else if exp, ok := expression.(*ast.InfixExpression); ok {
		o1 := t.translateExpression(exp.Left)
		o2 := t.translateExpression(exp.Right)
		ir_temp := ir.CalcInst{
			Operation: ir.ADD, 
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
	statement := functionff.Body.Statements[0]
	if tt, ok := statement.(*ast.VarAssign); ok {
		t.translateExpression(tt.Right)
	}
	
	return t.irList
}