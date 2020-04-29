package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var (
		consulHost  = flag.String("consul.host", "114.67.98.210", "consul ip address")
		consulPort  = flag.String("consul.port", "8500", "consul port")
		serviceHost = flag.String("service.host", "localhost", "service ip address")
		servicePort = flag.String("service.port", "8080", "service port")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var svc Service
	svc = StringService{}

	// add logging middleware
	svc = LoggingMiddleware(logger)(svc)

	endpoint := MakeStringEndpoint(svc)

	//创建健康检查的Endpoint
	healthEndpoint := MakeHealthCheckEndpoint(svc)

	//把算术运算Endpoint和健康检查Endpoint封装至StringEndpoints
	endpts := StringEndpoints{
		StringEndpoint:      endpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	//创建http.Handler
	r := MakeHttpHandler(ctx, endpts, logger)

	//创建注册对象
	//TODO replace with pkg consul
	registar := Register(*consulHost, *consulPort, *serviceHost, *servicePort, logger)

	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		//启动前执行注册
		registar.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务退出取消注册
	registar.Deregister()
	fmt.Println(error)
}
