package main

import "fmt"


type Person struct {
	Name string	// 姓名
	Birth string	// 生日
	ID int64	// 身份证号
}

// 指针类型，修改个人信息
func (person *Person) changeName(name string)  {
	person.Name = name
}

// 非指针类型，打印个人信息
func (person Person) printMess()  {
	fmt.Printf("My name is %v, and my birthday is %v, and my id is %v\n",
		person.Name, person.Birth, person.ID)

	// 尝试修改个人信息，但是对原接收器并没有影响
	// person.ID = 3

}

func main()  {
		p1 := Person{
			Name:"王小二",
			Birth: "1990-12-23",
			ID:1,
		}

		p1.printMess()
		p1.changeName("王老二")
		p1.printMess()

}


//func main()  {
//	// 声明实例化
//	var p1 Person
//	p1.Name =  "王小二"
//	p1.Birth = "1990-12-11"
//
//
//	// new函数实例化
//	p2 := new(Person)
//	p2.Name = "王二小"
//	p2.Birth = "1990-12-22"
//
//
//	// 取址实例化
//	p3 := &Person{}
//	p3.Name = "王三小"
//	p3.Birth = "1990-12-23"
//
//	// 初始化
//	p4 := Person{
//		Name:"王小四",
//		Birth: "1990-12-23",
//	}
//
//	// 初始化
//	p5 := &Person{
//		"王五",
//		"1990-12-23",
//		5,
//	}
//
//
//
//
//
//}
