package main

import "fmt"

// Cat接口
type Cat interface {
	// 抓老鼠
	CatchMouse()
}

// Dog接口
type Dog interface {
	// 吠叫
	Bark()
}

type CatDog struct {
	Name string
}

//  实现Cat接口
func (catDog *CatDog) CatchMouse()  {
	fmt.Printf("%v caught the mouse and ate it!\n", catDog.Name)
}

// Dog接口
func (catDog *CatDog) Bark()  {
	fmt.Printf("%v barked loudly!\n", catDog.Name)
}

func main()  {
	catDog := &CatDog{
		"Lucy",
	}

	// 声明一个Cat接口，并将catDog指针类型赋值给cat
	var cat Cat
	cat = catDog
	cat.CatchMouse()

	// 声明一个Dog接口，并将catDog指针类型赋值给dog
	var dog Dog
	dog = catDog
	dog.Bark()
}