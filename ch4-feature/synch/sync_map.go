package main

import (
	"fmt"
	"strconv"
	"sync"
)
var syncMap sync.Map
var waitGroup sync.WaitGroup

func main()  {

	routineSize := 5
	// 让主线程等待数据添加完毕
	waitGroup.Add(routineSize)
	// 并发添加数据
	for i := 0 ; i < routineSize; i++{
		go addNumber(i * 10)
	}

	// 开始等待
	waitGroup.Wait()
	var size int
	// 统计数量
	syncMap.Range(func(key, value interface{}) bool {
		size++
		//fmt.Println("key-value pair is", key, value, " ")
		return true
	})
	fmt.Println("syncMap current size is " + strconv.Itoa(size))
	// 获取键为 0 的值
	value, ok := syncMap.Load(0); if ok{
		fmt.Println("key 0 has value", value, " ")
	}

}

func addNumber(begin int)  {
	// 往 syncMap 中放入数据
	for i := begin ; i < begin + 3 ; i++{
		syncMap.Store(i, i)
	}
	// 通知数据已添加完毕
	waitGroup.Done()
}
