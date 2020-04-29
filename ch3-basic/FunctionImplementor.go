package main

import "fmt"

type Printer interface {
	Print(p interface{})
}


// 函数定义为类型
type FuncCaller func(p interface{})

// 实现Printer的Print方法
func (funcCaller FuncCaller) Print(p interface{}) {
	// 调用funcCaller函数本体
	funcCaller(p)
}

func main()  {
	var printer Printer
	// 将匿名函数强转为FuncCaller赋值给printer
	printer = FuncCaller(func(p interface{}) {
		fmt.Println(p)
	})
	printer.Print("Golang is Good!")
}
