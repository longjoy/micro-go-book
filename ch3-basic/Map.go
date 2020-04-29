package main

import "fmt"

func main() {

	classMates1 := make(map[int]string)

	// 添加映射关系
	classMates1[0] = "小明"
	classMates1[1] = "小红"
	classMates1[2] = "小张"

	fmt.Printf("id %v is %v\n", 1, classMates1[1])

	// 在声明时初始化数据
	classMates2 := map[int]string{
		0 : "小明",
		1 : "小红",
		2 : "小张",
	}

	fmt.Printf("id %v is %v\n", 3, classMates2[3])

}