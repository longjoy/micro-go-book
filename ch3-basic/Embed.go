package main

import "fmt"

type Wheel struct {
	shape string
}

type Car struct {
	Wheel
	Name string
}

func main()  {

	car := &Car{
		Wheel{
			"圆形的",
		},
		"福特",
	}
	fmt.Println(car.Name, car.shape, " ")
}
