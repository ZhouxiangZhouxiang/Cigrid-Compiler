package ir_translator 

import "cigrid/token"
import "cigrid/ast"
import "cigrid/ir"
import "strconv"

type IrFunction struct {
	functionName string // 记录了当前函数名
	irList       []ir.IntermediateRepresentation
	tempRegister int // how many temp registers have been used
	maxRegister  int // the maximum one
	variableMap  map[string]int // 记录是第几个变量
	addressMap   map[string]int // 记录相应变量的地址
	condition    int
}

func newIrFunc(fn string) *IrFunction {
	return &IrFunction{
		functionName: fn,
		irList: []ir.IntermediateRepresentation{},
		tempRegister: 0,
		maxRegister: 0,
		variableMap: make(map[string]int),
		addressMap: make(map[string]int),
	}
}

func (i *IrFunction) ReadName() string {
	return i.functionName
}

func (i *IrFunction) ReadIrList() []ir.IntermediateRepresentation {
	return i.irList 
}

func (i *IrFunction) ReadMaxRegister() int {
	return i.maxRegister
}

func (i *IrFunction) ReadVariableMap() map[string]int {
	return i.variableMap
}

func (i *IrFunction) ReadAddressMap() map[string]int {
	return i.addressMap
}

type IrTranslator struct {
	tree           *ast.ProgramLiteral // input
	irFunctionList []*IrFunction
	string_list    []string // record the string data
}

func New(tree *ast.ProgramLiteral) *IrTranslator {
	return &IrTranslator{
		tree: tree, 
		irFunctionList: []*IrFunction{}, 
	}
}

func (t *IrTranslator) ReadIrFunctionList() []*IrFunction {
	return t.irFunctionList
}

func (t *IrTranslator) ReadStringList() []string {
	return t.string_list
}

func (t *IrTranslator) translateFunction(fl *ast.FunctionLiteral) {
	integer_arguments := []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
	irFuncTemp := newIrFunc(fl.Name.String())
	t.irFunctionList = append(t.irFunctionList, irFuncTemp)
	for k, v := range fl.Param {
		varName := v.IdentifierLiteral.String()
		t.irFunctionList[len(t.irFunctionList) - 1].variableMap[varName] = 1
		varNameNew := varName + strconv.Itoa(
			t.irFunctionList[len(t.irFunctionList) - 1].variableMap[varName])
		t.irFunctionList[len(t.irFunctionList) - 1].addressMap[varNameNew] = 
			len(t.irFunctionList[len(t.irFunctionList) - 1].addressMap)
		ir_temp := ir.CalcInst{
			Operation: ir.MOV,
			Operand1: varNameNew, // 左值
			Operand2: integer_arguments[k],
		}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir_temp)
	}
	t.translateStatementBlock(fl.Body)
	for _, v := range fl.Param {
		varName := v.IdentifierLiteral.String()
		t.irFunctionList[len(t.irFunctionList) - 1].variableMap[varName] = 0
	}
}

func (t *IrTranslator) translateStatementBlock(bs *ast.BlockStatement) {
	tempVariable := []string{}
	for _, v := range bs.Statements {
		temp := t.translateStatement(v)
		if temp != "" {
			tempVariable = append(tempVariable, temp)
		}
	}
	for _, v := range tempVariable {
		t.irFunctionList[len(t.irFunctionList) - 1].variableMap[v]--
	}
}

func (t *IrTranslator) translateStatement(statement ast.Statement) string {
	if stmt, ok := statement.(*ast.VarAssign); ok {
		if id, ok := stmt.Left.(*ast.Identifier); ok {
			// a = 1;
			// new statement, temp register reset
			t.irFunctionList[len(t.irFunctionList) - 1].tempRegister = 0
			index_temp := t.irFunctionList[len(t.irFunctionList) - 1].
				variableMap[id.Value.Literal]
			// e.g. x2
			op1_temp := id.Value.Literal + strconv.Itoa(index_temp)
			ir_temp := ir.CalcInst{
				Operation: ir.MOV,
				Operand1: op1_temp, // 左值
				Operand2: t.translateExpression(stmt.Right),
			}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ir_temp)
		} else if id, ok := stmt.Left.(*ast.PrefixExpression); ok {
			if id.Operator.Type == token.ASTERISK {
				// *x = 1;
				// new statement, temp register reset
				t.irFunctionList[len(t.irFunctionList) - 1].tempRegister = 0
				vv, _ := id.Right.(*ast.Identifier)
				index_temp := t.irFunctionList[len(t.irFunctionList) - 1].
					variableMap[vv.Value.Literal]
				variable := vv.Value.Literal + strconv.Itoa(index_temp)
				ir_temp := ir.CalcInst{
					Operation: ir.MOV,
					Operand1: "r8",
					Operand2: variable,
				}
				t.irFunctionList[len(t.irFunctionList) - 1].irList = 
					append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
					ir_temp)
				ir_temp = ir.CalcInst{
					Operation: ir.MOV,
					Operand1: "[r8]",
					Operand2: t.translateExpression(stmt.Right),
				}
				t.irFunctionList[len(t.irFunctionList) - 1].irList = 
					append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
					ir_temp)
			}
		}
		
	} else if stmt, ok := statement.(*ast.VarDef); ok {
		// int a = 1;
		// new statement, temp register reset
		t.irFunctionList[len(t.irFunctionList) - 1].tempRegister = 0
		varName := stmt.Name.Value.Literal
		// 判断该变量是否曾经出现过，如果出现则掩盖
		// int x = 1;
		// if (...) {
		//	  int x = 0;
		// }
		if _, ok := t.irFunctionList[len(t.irFunctionList) - 1].variableMap[varName]; ok {
			t.irFunctionList[len(t.irFunctionList) - 1].variableMap[varName]++
		} else {
			t.irFunctionList[len(t.irFunctionList) - 1].variableMap[varName] = 1
		}
		// varName: x
		// varName: x2
		varNameNew := varName + strconv.Itoa(
			t.irFunctionList[len(t.irFunctionList) - 1].variableMap[varName])
		ir_temp := ir.CalcInst{
			Operation: ir.MOV,
			Operand1: varNameNew, // 左值
			Operand2: t.translateExpression(stmt.Value),
		}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir_temp)
		// 查看该变量是否出现过
		// if (...) {
		// 	int x = 1;
		// }
		// else {
		// 	int x = 2;
		// }
		// 则if else中的x可以存在一个地址
		_, ok := t.irFunctionList[len(t.irFunctionList) - 1].addressMap[varNameNew]
		if !ok {
			// 如果变量未曾出现，需要另外分配
			t.irFunctionList[len(t.irFunctionList) - 1].addressMap[varNameNew] = 
				len(t.irFunctionList[len(t.irFunctionList) - 1].addressMap)
		}
		return varName
	} else if stmt, ok := statement.(*ast.IfStatement); ok {
		// new statement, temp register reset
		t.irFunctionList[len(t.irFunctionList) - 1].tempRegister = 0
		// condition block (cmp, jC, jmp)
		// if: 
		// ...
		// jmp end
		// else: 
		// ...
		// jmp end
		// end:
		// ...
		
		// translate condition
		condition_temp := "label" + 
			strconv.Itoa(t.irFunctionList[len(t.irFunctionList) - 1].condition)
		t.translateCondition(stmt.Condition, condition_temp, 
			condition_temp + "_if", condition_temp + "_else")
		t.irFunctionList[len(t.irFunctionList) - 1].condition++
		// translate if statement block
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir.Label(condition_temp + "_if"))
		t.translateStatementBlock(stmt.Consequence)
		ji1 := ir.JumpInst{JC: ir.MP, Addr: condition_temp + "_end"}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ji1)
		// translate else statement block
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir.Label(condition_temp + "_else"))
		if stmt.Alternative != nil {
			t.translateStatementBlock(stmt.Alternative)
		}
		ji2 := ir.JumpInst{JC: ir.MP, Addr: condition_temp + "_end"}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ji2)
		// add end flag
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir.Label(condition_temp + "_end"))
	} else if stmt, ok := statement.(*ast.WhileStatement); ok {
		// new statement, temp register reset
		t.irFunctionList[len(t.irFunctionList) - 1].tempRegister = 0
		// condition:
		// condition block (cmp, jC, jmp)
		// while: 
		// ...
		// jmp condition
		// end:
		
		condition_temp := "label" + 
			strconv.Itoa(t.irFunctionList[len(t.irFunctionList) - 1].condition)
		t.irFunctionList[len(t.irFunctionList) - 1].condition++
		// before judging the condition
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir.Label(condition_temp + "_condition"))
		// translate condition
		t.translateCondition(stmt.Condition, condition_temp, 
			condition_temp + "_while", condition_temp + "_end")
		// translate statement block
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir.Label(condition_temp + "_while"))
		t.translateStatementBlock(stmt.Consequence)
		ji1 := ir.JumpInst{JC: ir.MP, Addr: condition_temp + "_condition"}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ji1)
		// add end flag
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir.Label(condition_temp + "_end"))
	} else if stmt, ok := statement.(*ast.ReturnStatement); ok {
		// new statement, temp register reset
		t.irFunctionList[len(t.irFunctionList) - 1].tempRegister = 0
		if stmt.ReturnValue == nil {
			ir_temp := ir.CalcInst{
				Operation: ir.XOR,
				Operand1: "rax",
				Operand2: "rax",
			}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ir_temp)
		} else {
			reg1 := t.translateExpression(stmt.ReturnValue)
			ir_temp := ir.CalcInst{
				Operation: ir.MOV,
				Operand1: "rax",
				Operand2: reg1,
			}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ir_temp)
		}
		ir_temp := ir.Ret("")
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir_temp)
	} else if stmt, ok := statement.(*ast.CallStatement); ok {
		t.translateExpression(stmt.Value)
	}
	return ""
}

func (t *IrTranslator) translateCondition(expression ast.Expression,
										  curNode string,
										  trueNode string, 
										  falseNode string) {
	t.irFunctionList[len(t.irFunctionList) - 1].irList = 
		append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
		ir.Label(curNode))
	if exp, ok := expression.(*ast.InfixExpression); ok {
		// label := ir.Label(curNode)
		// t.irFunctionList[len(t.irFunctionList) - 1].irList = 
		// 	append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
		// 	label)
		infix_temp := exp.Operator.Type
		switch infix_temp {
		case token.LT: 
			reg1 := t.translateExpression(exp.Left)
			reg2 := t.translateExpression(exp.Right)
			ci := ir.CmpInst{Left: reg1, Right: reg2}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ci)
			ji1 := ir.JumpInst{JC: ir.L, Addr: trueNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji2)
		case token.GT:
			reg1 := t.translateExpression(exp.Left)
			reg2 := t.translateExpression(exp.Right)
			ci := ir.CmpInst{Left: reg1, Right: reg2}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ci)
			ji1 := ir.JumpInst{JC: ir.G, Addr: trueNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji2)
		case token.L_EQ:
			reg1 := t.translateExpression(exp.Left)
			reg2 := t.translateExpression(exp.Right)
			ci := ir.CmpInst{Left: reg1, Right: reg2}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ci)
			ji1 := ir.JumpInst{JC: ir.LE, Addr: trueNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji2)
		case token.G_EQ:
			reg1 := t.translateExpression(exp.Left)
			reg2 := t.translateExpression(exp.Right)
			ci := ir.CmpInst{Left: reg1, Right: reg2}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ci)
			ji1 := ir.JumpInst{JC: ir.GE, Addr: trueNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji2)
		case token.EQ:
			reg1 := t.translateExpression(exp.Left)
			reg2 := t.translateExpression(exp.Right)
			ci := ir.CmpInst{Left: reg1, Right: reg2}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ci)
			ji1 := ir.JumpInst{JC: ir.E, Addr: trueNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji2)
		case token.NOT_EQ:
			reg1 := t.translateExpression(exp.Left)
			reg2 := t.translateExpression(exp.Right)
			ci := ir.CmpInst{Left: reg1, Right: reg2}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ci)
			ji1 := ir.JumpInst{JC: ir.NE, Addr: trueNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji2)
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
				Operand1: t.irFunctionList[len(t.irFunctionList) - 1].tempRegister, 
				Operand2: "0",
			}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ir_temp)
			t.irFunctionList[len(t.irFunctionList) - 1].tempRegister++ 
			// 更新maxRegister
			if (t.irFunctionList[len(t.irFunctionList) - 1].tempRegister > 
				t.irFunctionList[len(t.irFunctionList) - 1].maxRegister) {
				t.irFunctionList[len(t.irFunctionList) - 1].maxRegister = 
				t.irFunctionList[len(t.irFunctionList) - 1].tempRegister
			}
			ci := ir.CmpInst{
				Left: reg1, 
				Right: t.irFunctionList[len(t.irFunctionList) - 1].tempRegister - 1,
			}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ci)
			ji1 := ir.JumpInst{JC: ir.NE, Addr: trueNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji2)
		}
	} else if exp, ok := expression.(*ast.PrefixExpression); ok {
		if exp.Operator.Type == token.BANG {
			t.translateCondition(exp.Right, curNode, falseNode, trueNode)
		} else {
			reg1 := t.translateExpression(exp.Right)
			ir_temp := ir.CalcInst{
				Operation: ir.MOV, 
				Operand1: t.irFunctionList[len(t.irFunctionList) - 1].tempRegister, 
				Operand2: "0",
			}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ir_temp)
			t.irFunctionList[len(t.irFunctionList) - 1].tempRegister++
			// 更新maxRegister
			if (t.irFunctionList[len(t.irFunctionList) - 1].tempRegister > 
				t.irFunctionList[len(t.irFunctionList) - 1].maxRegister) {
				t.irFunctionList[len(t.irFunctionList) - 1].maxRegister = 
				t.irFunctionList[len(t.irFunctionList) - 1].tempRegister
			} 
			ci := ir.CmpInst{
				Left: reg1, 
				Right: t.irFunctionList[len(t.irFunctionList) - 1].tempRegister - 1,
			}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ci)
			ji1 := ir.JumpInst{JC: ir.NE, Addr: trueNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji1)
			ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ji2)
		}
	} else {
		// like x, 1
		reg1 := t.translateExpression(expression)
		ir_temp := ir.CalcInst{
			Operation: ir.MOV, 
			Operand1: t.irFunctionList[len(t.irFunctionList) - 1].tempRegister, 
			Operand2: "0",
		}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir_temp)
		t.irFunctionList[len(t.irFunctionList) - 1].tempRegister++
		// 更新maxRegister
		if (t.irFunctionList[len(t.irFunctionList) - 1].tempRegister > 
			t.irFunctionList[len(t.irFunctionList) - 1].maxRegister) {
			t.irFunctionList[len(t.irFunctionList) - 1].maxRegister = 
			t.irFunctionList[len(t.irFunctionList) - 1].tempRegister
		} 
		ci := ir.CmpInst{
			Left: reg1, 
			Right: t.irFunctionList[len(t.irFunctionList) - 1].tempRegister - 1,
		}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ci)
		ji1 := ir.JumpInst{JC: ir.NE, Addr: trueNode}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ji1)
		ji2 := ir.JumpInst{JC: ir.MP, Addr: falseNode}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ji2)
	}
}

func (t *IrTranslator) translateExpression(expression ast.Expression) int {
	if exp, ok := expression.(*ast.IntegerLiteral); ok {
		// int字面量，1
		ir_temp := ir.CalcInst{
			Operation: ir.MOV, 
			Operand1: t.irFunctionList[len(t.irFunctionList) - 1].tempRegister, 
			Operand2: exp.Value.Literal,
		}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir_temp)
		t.irFunctionList[len(t.irFunctionList) - 1].tempRegister++
		// 更新maxRegister
		if (t.irFunctionList[len(t.irFunctionList) - 1].tempRegister > 
			t.irFunctionList[len(t.irFunctionList) - 1].maxRegister) {
			t.irFunctionList[len(t.irFunctionList) - 1].maxRegister = 
			t.irFunctionList[len(t.irFunctionList) - 1].tempRegister
		}	 
		return t.irFunctionList[len(t.irFunctionList) - 1].tempRegister - 1
	} else if exp, ok := expression.(*ast.StringLiteral); ok {
		// string字面量
		t.string_list = append(t.string_list, exp.Value.Literal)
		ir_temp := ir.CalcInst{
			Operation: ir.MOV, 
			Operand1: t.irFunctionList[len(t.irFunctionList) - 1].tempRegister, 
			Operand2: "str" + strconv.Itoa(len(t.string_list)),
		}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir_temp)
		t.irFunctionList[len(t.irFunctionList) - 1].tempRegister++
		// 更新maxRegister
		if (t.irFunctionList[len(t.irFunctionList) - 1].tempRegister > 
			t.irFunctionList[len(t.irFunctionList) - 1].maxRegister) {
			t.irFunctionList[len(t.irFunctionList) - 1].maxRegister = 
			t.irFunctionList[len(t.irFunctionList) - 1].tempRegister
		}	 
		return t.irFunctionList[len(t.irFunctionList) - 1].tempRegister - 1
	} else if exp, ok := expression.(*ast.Identifier); ok {
		// x
		ir_temp := ir.CalcInst{
			Operation: ir.MOV,
			Operand1: t.irFunctionList[len(t.irFunctionList) - 1].tempRegister,
			Operand2: exp.Value.Literal + 
				strconv.Itoa(t.irFunctionList[len(t.irFunctionList) - 1].
				variableMap[exp.Value.Literal]),
		}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, ir_temp)
		t.irFunctionList[len(t.irFunctionList) - 1].tempRegister++
		// 更新maxRegister
		if (t.irFunctionList[len(t.irFunctionList) - 1].tempRegister > 
			t.irFunctionList[len(t.irFunctionList) - 1].maxRegister) {
			t.irFunctionList[len(t.irFunctionList) - 1].maxRegister = 
			t.irFunctionList[len(t.irFunctionList) - 1].tempRegister
		}
		return t.irFunctionList[len(t.irFunctionList) - 1].tempRegister - 1
	} else if exp, ok := expression.(*ast.InfixExpression); ok {
		// 1 + 2
		var infix_temp ir.Op
		switch exp.Operator.Type {
		case token.PLUS: 
			infix_temp = ir.ADD
		case token.MINUS:
			infix_temp = ir.SUB
		case token.ASTERISK:
			infix_temp = ir.MUL 
		case token.SLASH: 
			infix_temp = ir.DIV
		default:
		}
		o1 := t.translateExpression(exp.Left)
		o2 := t.translateExpression(exp.Right)
		ir_temp := ir.CalcInst{
			Operation: infix_temp, 
			Operand1: o1,
			Operand2: o2,
		}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir_temp)
		return o1
	} else if exp, ok := expression.(*ast.PrefixExpression); ok {
		if exp.Operator.Type == token.MINUS {
			// -1
			o1 := t.translateExpression(exp.Right)
			ir_temp := ir.OneInst{
				Operation: ir.NEG, 
				Operand1: o1,
			}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ir_temp)
			return o1
		} else if exp.Operator.Type == token.ET {
			// &x or &x[0]
			if er, ok := exp.Right.(*ast.Identifier); ok {
				// &x
				ir_temp := ir.CalcInst{
					Operation: ir.LEA,
					Operand1: t.irFunctionList[len(t.irFunctionList) - 1].tempRegister,
					Operand2: er.Value.Literal + 
						strconv.Itoa(t.irFunctionList[len(t.irFunctionList) - 1].
						variableMap[er.Value.Literal]),
				}
				t.irFunctionList[len(t.irFunctionList) - 1].irList = 
					append(t.irFunctionList[len(t.irFunctionList) - 1].irList, ir_temp)
				t.irFunctionList[len(t.irFunctionList) - 1].tempRegister++
				// 更新maxRegister
				if (t.irFunctionList[len(t.irFunctionList) - 1].tempRegister > 
					t.irFunctionList[len(t.irFunctionList) - 1].maxRegister) {
					t.irFunctionList[len(t.irFunctionList) - 1].maxRegister = 
					t.irFunctionList[len(t.irFunctionList) - 1].tempRegister
				}
				return t.irFunctionList[len(t.irFunctionList) - 1].tempRegister - 1
			}
		} else if exp.Operator.Type == token.ASTERISK {
			// *x
			if er, ok := exp.Right.(*ast.Identifier); ok {
				// &x
				// 首先将地址mov到r9
				ir_temp := ir.CalcInst{
					Operation: ir.MOV,
					Operand1: "r9",
					Operand2: er.Value.Literal + 
						strconv.Itoa(t.irFunctionList[len(t.irFunctionList) - 1].
						variableMap[er.Value.Literal]),
				}
				t.irFunctionList[len(t.irFunctionList) - 1].irList = 
					append(t.irFunctionList[len(t.irFunctionList) - 1].irList, ir_temp)
				// 将[r9]移动到temp register
				ir_temp = ir.CalcInst{
					Operation: ir.MOV,
					Operand1: t.irFunctionList[len(t.irFunctionList) - 1].tempRegister,
					Operand2: "[r9]",
				}
				t.irFunctionList[len(t.irFunctionList) - 1].irList = 
					append(t.irFunctionList[len(t.irFunctionList) - 1].irList, ir_temp)
				t.irFunctionList[len(t.irFunctionList) - 1].tempRegister++
				// 更新maxRegister
				if (t.irFunctionList[len(t.irFunctionList) - 1].tempRegister > 
					t.irFunctionList[len(t.irFunctionList) - 1].maxRegister) {
					t.irFunctionList[len(t.irFunctionList) - 1].maxRegister = 
					t.irFunctionList[len(t.irFunctionList) - 1].tempRegister
				}
				return t.irFunctionList[len(t.irFunctionList) - 1].tempRegister - 1
			}
		}
	} else if exp, ok := expression.(*ast.CallExpression); ok {
		integer_arguments := []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
		reg_list := []int{}
		// prepare for the input arguments
		// 最多有6个arguments
		for _, v := range(exp.Params) {
			reg := t.translateExpression(v)
			reg_list = append(reg_list, reg)
		}
		for k, v := range(reg_list) {
			ir_temp := ir.CalcInst{
				Operation: ir.MOV, 
				Operand1: integer_arguments[k], 
				Operand2: v,
			}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				ir_temp)
		}
		// reset xor before call func
		ir_temp := ir.CalcInst{
			Operation: ir.XOR,
			Operand1: "rax",
			Operand2: "rax",
		}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir_temp)
		// caller saved register
		// push
		for _, v := range(integer_arguments) {
			temp := ir.OneInst{
				Operation: ir.PUSH, 
				Operand1: v,
			}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				temp)
		}
		// call function
		call_temp := ir.CallInst{
			FuntionName: exp.Name.String(),
		}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			call_temp)
		// caller saved register
		// pop
		for i := len(integer_arguments) - 1; i >= 0; i-- {
			temp := ir.OneInst{
				Operation: ir.POP, 
				Operand1: integer_arguments[i],
			}
			t.irFunctionList[len(t.irFunctionList) - 1].irList = 
				append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
				temp)
		}
		// move return value from rax to a temp register
		ir_temp = ir.CalcInst{
			Operation: ir.MOV, 
			Operand1: t.irFunctionList[len(t.irFunctionList) - 1].tempRegister, 
			Operand2: "rax",
		}
		t.irFunctionList[len(t.irFunctionList) - 1].irList = 
			append(t.irFunctionList[len(t.irFunctionList) - 1].irList, 
			ir_temp)
		t.irFunctionList[len(t.irFunctionList) - 1].tempRegister++
		// 更新maxRegister
		if (t.irFunctionList[len(t.irFunctionList) - 1].tempRegister > 
			t.irFunctionList[len(t.irFunctionList) - 1].maxRegister) {
			t.irFunctionList[len(t.irFunctionList) - 1].maxRegister = 
			t.irFunctionList[len(t.irFunctionList) - 1].tempRegister
		}	 
		return t.irFunctionList[len(t.irFunctionList) - 1].tempRegister - 1
	}
	return 0
}

func (t *IrTranslator) Translate() {
	for _, value := range t.tree.GlobalList {
		if v, ok := value.(*ast.FunctionLiteral); ok {
			t.translateFunction(v)
		}
	}
}
