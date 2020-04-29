package main

import "fmt"

func main()  {

	var classMates1 [3]string

	classMates1[0] = "小明"
	classMates1[1] = "小红"
	classMates1[2] = "小李"
	fmt.Println(classMates1)
	fmt.Println("The No.1 student is " + classMates1[0])

	classMates2  := [...]string{"小明", "小红", "小李"}
	fmt.Println(classMates2)


	classMates3 := new([3]string)
	classMates3[0] = "小明"
	classMates3[1] = "小红"
	classMates3[2] = "小李"
	fmt.Println(*classMates3)


}