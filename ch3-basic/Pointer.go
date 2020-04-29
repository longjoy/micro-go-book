package main

import "fmt"

func main()  {


	str := "Golang is Good!"
	strPrt := &str

	fmt.Printf("str type is %T, value is %v, address is %p\n", str, str, &str)
	fmt.Printf("strPtr type is %T, and value is %v\n", strPrt, strPrt)

	newStr := *strPrt	//获取指针对应变量的值
	fmt.Printf("newStr type is %T, value is %v, and address is %p\n", newStr, newStr, &newStr)

	*strPrt = "Java is Good too!"	//通过指针对变量进行赋值
	fmt.Printf("newStr type is %T, value is %v, and address is %p\n", newStr, newStr, &newStr)
	fmt.Printf("str type is %T, value is %v, address is %p\n", str, str, &str)






}