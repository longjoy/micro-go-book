package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func printInput(ch chan string)  {


	// 使用 for 循环从 channel 中读取数据
	for val := range ch{
		// 读取到结束符号
		if val == "EOF"{
			break
		}
		fmt.Printf("Input is %s\n", val)
	}


}

func consume(ch chan int)  {

	// 线程休息 100s 再从 channel 读取数据
	time.Sleep(time.Second * 100)
	<- ch

}



func main()  {

	// 创建一个无缓冲的 channel
	ch := make(chan string)
	go printInput(ch)

	// 从命令行读取输入
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		val := scanner.Text()
		ch <- val
		if val == "EOF"{
			fmt.Println("End the game!")
			break
		}
	}
	// 程序最后关闭 ch
	defer close(ch)


	//// 创建一个长度为 2 的 channel
	//ch := make(chan int, 2)
	//go consume(ch)
	//
	//ch <- 0
	//ch <- 1
	//// 发送数据不被阻塞
	//fmt.Println("I am free!")
	//ch <- 2
	//fmt.Println("I can not go there within 100s!")
	//
	//time.Sleep(time.Second)


	//ch1 := make(chan int)
	//ch2 := make(chan int)
	//
	//go send(ch1, 0)
	//go send(ch2, 10)
	//
	//// 主 goroutine 休眠 1s，保证调度成功
	//time.Sleep(time.Second)
	//
	//for {
	//	select {
	//	case val := <- ch1: // 从 ch1 读取数据
	//		fmt.Printf("get value %d from ch1\n", val)
	//	case val := <- ch2 : // 从 ch2 读取数据
	//		fmt.Printf("get value %d from ch2\n", val)
	//	case <-time.After(2 * time.Second): // 超时设置
	//		fmt.Println("Time out")
	//		return
	//	}
	//}







}

func send(ch chan int, begin int )  {

	for i :=begin ; i< begin + 10 ;i++{
		ch <- i

	}

}
