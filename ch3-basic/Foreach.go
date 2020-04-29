package main

import (
	"fmt"
)

func main()  {

	// 数组的遍历
	nums := [...]int{1,2,3,4,5,6,7,8}

	for k, v:= range nums{
		// k为下标，v为对应的值
		fmt.Println(k, v, " ")
	}

	fmt.Println()

	// 切片的遍历
	slis := []int{1,2,3,4,5,6,7,8}
	for k, v:= range slis{
		// k为下标，v为对应的值
		fmt.Println(k, v, " ")
	}

	fmt.Println()


	// 字典的遍历
	tmpMap := map[int]string{
		0 : "小明",
		1 : "小红",
		2 : "小张",
	}

	for k, v:= range tmpMap{
		// k为键值，v为对应值
		fmt.Println(k, v, " ")
	}

}
