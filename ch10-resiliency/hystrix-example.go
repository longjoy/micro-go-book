package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"strconv"
	"time"
)

func main() {

	hystrix.ConfigureCommand("test_command", hystrix.CommandConfig{
		// 设置参数
		Timeout: hystrix.DefaultTimeout,
	})

	err := hystrix.Do("test_command", func() error {
		// 远程调用&或者其他需要保护的方法
		return nil
	}, func(err error) error {
		// 失败回滚方法
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	resultChan := make(chan interface{}, 1)
	errChan := hystrix.Go("test_command", func() error {
		// 远程调用&或者其他需要保护的方法
		resultChan <- "success"
		return nil
	}, func(e error) error {
		// 失败回滚方法
		return nil
	})

	select {
	case err := <-errChan:
		// 执行失败
		fmt.Println(err)
	case result := <-resultChan:
		// 执行成功
		fmt.Println(result)
	case <-time.After(2 * time.Second): // 超时设置
		fmt.Println("Time out")
		return
	}

	circuit, _, _ := hystrix.GetCircuit("test_command")
	fmt.Println("command test_command's circuit open is " + strconv.FormatBool(circuit.IsOpen()))

}
