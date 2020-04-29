package main

import (
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-core/setup"
)

func main() {

	setup.InitZk()
	setup.InitRedis()
	setup.RunService()

}
