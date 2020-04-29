package setup

import (
	"fmt"
	register "github.com/longjoy/micro-go-book/ch13-seckill/pkg/discover"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-core/service/srv_redis"
	"os"
	"os/signal"
	"syscall"
)

func RunService() {
	//启动处理线程
	srv_redis.RunProcess()
	errChan := make(chan error)
	//http server
	go func() {
		//启动前执行注册
		register.Register()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务退出取消注册
	register.Deregister()
	fmt.Println(error)

}
