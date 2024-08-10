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
	Operand1 interface{}
	Operand2 interface{}
}
func (ci CalcInst) IrString() string {
	var out bytes.Buffer 
	out.WriteString(string(ci.Operation))
	out.WriteString(" ")
	if temp, ok := ci.Operand1.(int); ok {
		out.WriteString("temp" + strconv.Itoa(temp))
	} else if temp, ok := ci.Operand1.(string); ok {
		out.WriteString(temp)
	}
	out.WriteString(" ")
	if temp, ok := ci.Operand2.(int); ok {
		out.WriteString("temp" + strconv.Itoa(temp))
	} else if temp, ok := ci.Operand2.(string); ok {
		out.WriteString(temp)
	}
	return out.String()
}

type Label string 
func (l Label) IrString() string { return string(l) + ":" }

type JumpType string 
const (
	MP = "mp"
	E = "e"
	NE = "ne"
	G = "g"
	L = "l"
	GE = "ge"
	LE = "le"
)
type JumpInst struct {
	JC   JumpType
	Addr string 
}
func (ji JumpInst) IrString() string {
	var out bytes.Buffer 
	out.WriteString("j" + string(ji.JC))
	out.WriteString(" ")
	out.WriteString(ji.Addr)
	return out.String()
}

type CmpInst struct {
	Left int 
	Right int 
}
func (ci CmpInst) IrString() string {
	var out bytes.Buffer 
	out.WriteString("cmp ")
	out.WriteString("temp" + strconv.Itoa(ci.Left))
	out.WriteString(" ")
	out.WriteString("temp" + strconv.Itoa(ci.Right))
	return out.String()
}