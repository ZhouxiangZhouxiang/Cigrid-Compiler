package main

import "cigrid/token"
import "cigrid/lexer"
import "cigrid/parser"
import "cigrid/ir_translator"
import "strconv"
import "cigrid/asm"
import "fmt"
import "bytes"
import "os"

func printTokenList(tokList []token.Token) {
	for id, tok := range tokList {
		fmt.Printf("%-4v $ %-10v $ %-15v\n", id, tok.Type, tok.Literal)
	}
}

func printIrList(t *ir_translator.IrTranslator) {
	stringList := t.ReadStringList()
	for k, v := range stringList {
		fmt.Println("str" + strconv.Itoa(k + 1), v)
	}
	irFunctionList := t.ReadIrFunctionList()
	for _, v1 := range irFunctionList {
		fmt.Println("-----------------------------------------------")
		fmt.Println("<<" + v1.ReadName() + ">>")
		for _, v2 := range v1.ReadIrList() {
			fmt.Println(v2.IrString())
		}
		for k, v := range v1.ReadAddressMap() {
			fmt.Println(k + " " + strconv.Itoa(v))
		}
		fmt.Println(v1.ReadMaxRegister())
		fmt.Println("-----------------------------------------------")
	}
}

func printAsm(a []string) {
	var out bytes.Buffer 
	for _, v := range(a) {
		out.WriteString(v)
		out.WriteString("\n")
	}
	fmt.Println(out.String())
	file, err := os.OpenFile("hello.asm", 
		os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println("error")
		return
	} 
	file.WriteString(out.String())
}

func main() {
	input := `
void swap(int * x, int * y) {
	*x = *x + *y;
	return;
}
int main() {
	int x = 2;
	int y = 30;
	swap(&x, &y);
	printf("%d %d\n", x, y);
	return 0;
}
`
	l := lexer.New(input)
	tokList := l.Scan()
	printTokenList(tokList) // 打印token list
	p := parser.New(tokList)
	tree := p.ParseProgram()
	fmt.Println(tree.String()) // 打印ast
	t := ir_translator.New(tree)
	t.Translate()
	printIrList(t) // 打印IR
	asm_list := asm.GenerateAsm(t) 
	printAsm(asm_list) // 打印x86-64
}

