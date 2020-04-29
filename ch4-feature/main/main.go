package main

import (
	"fmt"
	"github.com/longjoy/micro-go-book/ch4-feature/compute"
)

func main()  {

	params := &compute.IntParams{
		P1:1,
		P2:2,
	}
	fmt.Println(params.Add())

}