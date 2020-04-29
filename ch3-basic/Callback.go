package main

import "fmt"

func proc(input string, processor func(str string))  {
	processor(input)
}

func main()  {

	proc("王小二", func(str string) {
		for _, v := range str{
			fmt.Printf("%c\n", v)
		}
	})
}