package ir 

import "bytes"
import "strconv"

type Op string 

// Op
const (
	ADD = "add"
	SUB = "sub"
	MOV = "mov"
	XOR = "xor"
	MUL = "imul"
	DIV = "idiv"
	NEG = "neg"
	PUSH = "push"
	POP = "pop"
	LEA = "lea"
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

type OneInst struct {
	Operation Op
	Operand1 interface{}
}
func (oi OneInst) IrString() string {
	var out bytes.Buffer
	out.WriteString(string(oi.Operation))
	out.WriteString(" ")
	if temp, ok := oi.Operand1.(int); ok {
		out.WriteString("temp" + strconv.Itoa(temp))
	} else if temp, ok := oi.Operand1.(string); ok {
		out.WriteString(temp)
	}
	return out.String()
}

type Label string 
func (l Label) IrString() string { return string(l) + ":" }

type Ret string
func (r Ret) IrString() string { return "ret" }

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

type CallInst struct {
	FuntionName string 
}
func (ci CallInst) IrString() string {
	var out bytes.Buffer 
	out.WriteString("call ")
	out.WriteString(ci.FuntionName)
	return out.String()
}