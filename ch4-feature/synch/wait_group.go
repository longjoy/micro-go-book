package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main()  {

	var waitGroup sync.WaitGroup


	waitGroup.Add(5)

	for i := 0 ; i < 5 ; i++{
		go func(i int) {
			fmt.Println("work " + strconv.Itoa(i) + " is done at " + time.Now().String())
			// 等待 1 s 后减少等待数
			time.Sleep(time.Second)
			waitGroup.Done()
		}(i)
	}

	waitGroup.Wait()

	fmt.Println("all works are done at " + time.Now().String())

}
