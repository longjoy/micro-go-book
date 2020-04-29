package main

import "fmt"

type aliasInt = int // 定义一个类型别名
type myInt int // 定义一个新的类型

func main()  {

	var alias aliasInt
	fmt.Printf("alias value is %v, type is %T\n", alias, alias)

	var myint myInt
	fmt.Printf("myint value is %v, type is %T\n", myint, myint)

	name := "小红"
	switch name {
	case "小明":
		fmt.Println("扫地")
	case "小红":
		fmt.Println("擦黑板")
	case "小刚":
		fmt.Println("倒垃圾")
	default:
		fmt.Println("没人干活")
	}

	score := 90
	switch  {
	case score < 100 && score >= 90:
		fmt.Println("优秀")
	case score < 90 && score >= 80:
		fmt.Println("良好")
	case score < 80 && score >= 60:
		fmt.Println("及格")
	case score < 60 :
		fmt.Println("不及格")
	default:
		fmt.Println("分数错误")

	}





}

