package asm 

import "cigrid/ir"
import "cigrid/ir_translator"
import "strconv"

func address(addressMap map[string]int, reg interface{}) (string, bool) {
	if temp, ok := reg.(int); ok {
		// int type, refers to a temporary register
		return "qword [rsp + " + 
			strconv.Itoa((temp + len(addressMap)) * 8) + "]", true
	} else if temp, ok := reg.(string); ok {
		// string type
		if value, ok := addressMap[temp]; ok {
			// if variable name
			return "qword [rsp + " + strconv.Itoa(value * 8) + "]", true
		} else if temp[0] == 91 && temp[len(temp) - 1] == 93 {
			// like [r10]
			return temp, true
		} else {
			// if existing register, like rax, rdi
			return temp, false
		}
	}
	return "", false
}

func generateSingleAsm(i *ir_translator.IrFunction) []string {
	result := []string{}
	functionName := i.ReadName()
	result = append(result, functionName + ": ")
	addressMap := i.ReadAddressMap()
	stack_depth := (i.ReadMaxRegister() + len(addressMap)) * 8
	result = append(result, "sub rsp, " + strconv.Itoa(stack_depth))
	callee_register := []string{"rbp", "rbx", "r12", "r13", "r14", "r15"}
	for _, v := range(callee_register) {
		result = append(result, "push " + v)
	}
	for _, v := range(i.ReadIrList()) {
		if value, ok := v.(ir.Label); ok {
			result = append(result, string(value) + ": ")
		} else if value, ok := v.(ir.CalcInst); ok {
			if value.Operation == ir.ADD || value.Operation == ir.SUB ||
			   value.Operation == ir.MOV || value.Operation == ir.XOR {
				temp := string(value.Operation)
				r1, o1 := address(addressMap, value.Operand1)
				r2, o2 := address(addressMap, value.Operand2)
				if o1 && o2 {
					// Binary instructions (e.g., add) cannot use two memory operands.
					mov_temp := "mov r10, " + r2
					temp += " " + r1 + ", r10"
					result = append(result, mov_temp)
					result = append(result, temp)
				} else {
					temp += " " + r1 + ", " + r2
					result = append(result, temp)
				}
			} else if value.Operation == ir.MUL || value.Operation == ir.DIV {
				temp := string(value.Operation)
				r1, _ := address(addressMap, value.Operand1)
				r2, _ := address(addressMap, value.Operand2)
				result = append(result, "mov rax, " + r1)
				temp += " " + r2 
				result = append(result, temp)
				result = append(result, "mov " + r1 + ", rax")
			} else if value.Operation == ir.LEA {
				temp := "lea r10, "
				r1, _ := address(addressMap, value.Operand1)
				r2, _ := address(addressMap, value.Operand2)
				result = append(result, temp + r2[6:])
				result = append(result, "mov " + r1 + ", r10")
			}
		} else if _, ok := v.(ir.Ret); ok {
			for i := len(callee_register) - 1; i >= 0; i-- {
				result = append(result, "pop " + callee_register[i])
			}
			result = append(result, "add rsp, " + strconv.Itoa(stack_depth))
			result = append(result, "ret")
		} else if value, ok := v.(ir.OneInst); ok {
			if value.Operation == ir.NEG || value.Operation == ir.PUSH || 
			   value.Operation == ir.POP {
				temp := string(value.Operation)
				r1, _ := address(addressMap, value.Operand1)
				temp += " " + r1
				result = append(result, temp)
			}
		} else if value, ok := v.(ir.CmpInst); ok {
			temp := "cmp"
			r1, o1 := address(addressMap, value.Left)
			r2, o2 := address(addressMap, value.Right)
			if o1 && o2 {
				// Binary instructions (e.g., add) cannot use two memory operands.
				mov_temp := "mov r10, " + r2
				temp += " " + r1 + ", r10"
				result = append(result, mov_temp)
				result = append(result, temp)
			} else {
				temp += " " + r1 + ", " + r2
				result = append(result, temp)
			}
		} else if value, ok := v.(ir.JumpInst); ok {
			result = append(result, value.IrString())
		} else if value, ok := v.(ir.CallInst); ok {
			result = append(result, value.IrString())
		}
	}
	return result
}

func GenerateAsm(t *ir_translator.IrTranslator) []string {
	list := t.ReadIrFunctionList()
	result := []string{}
	result = append(result, "global main")
	result = append(result, "extern printf")
	result = append(result, "section .data")
	for k, v := range(t.ReadStringList()) {
		pre := "str" + strconv.Itoa(k + 1) + ": db "
		// 判断是否存在换行符
		// 假设只在字符串末尾出现换行符
		if v[len(v) - 2] == 92 && v[len(v) - 1] == 110 {
			v = "\"" + v[:len(v) - 2] + "\"" + ", 10, 0"
		} else {
			v = "\"" + v[:len(v) - 2] + "\"" + ", 0"
		}
		result = append(result, pre + v)
	}
	result = append(result, "section .text")
	for _, v := range(list) {
		result = append(result, generateSingleAsm(v)...)
	}
	return result
}

