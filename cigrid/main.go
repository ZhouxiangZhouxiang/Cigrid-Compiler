package main

import "cigrid/token"
import "cigrid/lexer"
import "cigrid/parser"
import "cigrid/ir"
import "cigrid/ir_translator"
import "fmt"

func printTokenList(tokList []token.Token) {
	for id, tok := range tokList {
		fmt.Printf("%-4v $ %-10v $ %-15v\n", id, tok.Type, tok.Literal)
	}
}

func printIrList(irList []ir.IntermediateRepresentation) {
	for _, v := range irList {
		fmt.Println(v.IrString())
	}
}

func main() {
	input := `
int print(string a, int b) {
	if (!a && b) {
		i = 2;
	}
}
`
	l := lexer.New(input)
	tokList := l.Scan()
	//printTokenList(tokList)
	p := parser.New(tokList)
	leftExp := p.ParseProgram()
	fmt.Println(leftExp.String())
	t := ir_translator.New(leftExp)
	ll := t.Translate()
	printIrList(ll)
}