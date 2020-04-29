package main

import (
	"flag"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/longjoy/micro-go-book/common/discover"
	"github.com/longjoy/micro-go-book/common/loadbalance"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// 创建环境变量
	var (
		consulHost = flag.String("consul.host", "127.0.0.1", "consul server ip address")
		consulPort = flag.Int("consul.port", 8500, "consul server port")
	)
	flag.Parse()

	//创建日志组件
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}



	consulClient, err := discover.NewKitDiscoverClient(*consulHost, *consulPort)

	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	//创建反向代理
	proxy := NewHystrixHandler(consulClient, new(loadbalance.RandomLoadBalance), log.New(os.Stderr, "", log.LstdFlags))

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	//开始监听
	go func() {
		logger.Log("transport", "HTTP", "addr", "9090")
		errc <- http.ListenAndServe(":9090", proxy)
	}()

	// 开始运行，等待结束
	logger.Log("exit", <-errc)
}
