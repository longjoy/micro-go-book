package main

import (
	"container/list"
	"fmt"
)

func main()  {


	tmpList := list.New()

	for i:= 1 ; i <= 10 ; i++ {
		tmpList.PushBack(i)
	}

	first := tmpList.PushFront(0)
	tmpList.Remove(first)

	for l := tmpList.Front(); l != nil; l = l.Next(){
		fmt.Print(l.Value, " ")
	}

}
