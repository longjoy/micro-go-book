package main

import (
	"flag"
	"fmt"
)

func main()  {


	//参数依次是命令行参数的名称，默认值，提示
	surname := flag.String("surname", "王", "您的姓")
	//除了返回结果，还可以直接传入变量地址获取参数值
	var personalName string
	flag.StringVar(&personalName, "personalName", "小二", "您的名")
	id := flag.Int("id", 0, "您的ID")
	//解析命令行参数
	flag.Parse()
	fmt.Printf("I am %v %v, and my id is %v\n", *surname, personalName, *id)


}