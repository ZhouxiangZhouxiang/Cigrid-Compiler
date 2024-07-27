package ir 

import "bytes"
import "strconv"

type Op string 

// Op
const (
	ADD = "add"
	SUB = "sub"
	MOV = "mov"
)

type IntermediateRepresentation interface {
	IrString() string
}

type Operand interface {
	OperandString() string 
}

type CalcInst struct {
	Operation Op
	Operand1 int
	Operand2 interface{}
}
func (ci CalcInst) IrString() string {
	var out bytes.Buffer 
	out.WriteString(string(ci.Operation))
	out.WriteString(" ")
	out.WriteString("temp" + strconv.Itoa(ci.Operand1))
	out.WriteString(" ")
	if temp, ok := ci.Operand2.(int); ok {
		out.WriteString("temp" + strconv.Itoa(temp))
	} else if temp, ok := ci.Operand2.(string); ok {
		out.WriteString(temp)
	}
	return out.String()
}

