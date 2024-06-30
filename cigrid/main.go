package main

import "cigrid/token"
import "cigrid/lexer"
import "cigrid/parser"
import "fmt"

func printTokenList(tokList []token.Token) {
	for id, tok := range tokList {
		fmt.Printf("%-4v $ %-10v $ %-15v\n", id, tok.Type, tok.Literal)
	}
}

func main() {
	input := `
{
int i = a + b; int c = d;
}`
	l := lexer.New(input)
	tokList := l.Scan()
	// printTokenList(tokList)
	p := parser.New(tokList)
	leftExp := p.ParseProgram()
	fmt.Println(leftExp.String())
}