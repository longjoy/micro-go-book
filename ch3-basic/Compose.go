package main

import "fmt"

// 游泳特性
type Swimming struct {
}

func (swim *Swimming) swim()  {
	fmt.Println("swimming is my ability")
}

// 飞行特性
type Flying struct {
}

func (fly *Flying) fly()  {
	fmt.Println("flying is my ability")
}

// 野鸭，具备飞行和游泳特性
type WildDuck struct {
	Swimming
	Flying
}

// 家鸭，具备游泳特性
type DomesticDuck struct {
	Swimming
}

func main()  {

	// 声明一只野鸭，可以飞，也可以游泳
	wild := WildDuck{}
	wild.fly()
	wild.swim()

	// 声明一只家鸭，只会游泳
	domestic := DomesticDuck{}
	domestic.swim()

	
}