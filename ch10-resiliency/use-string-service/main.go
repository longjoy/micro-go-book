package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/longjoy/micro-go-book/ch10-resiliency/use-string-service/config"
	"github.com/longjoy/micro-go-book/ch10-resiliency/use-string-service/endpoint"
	"github.com/longjoy/micro-go-book/ch10-resiliency/use-string-service/service"
	"github.com/longjoy/micro-go-book/ch10-resiliency/use-string-service/transport"
	"github.com/longjoy/micro-go-book/common/discover"
	"github.com/longjoy/micro-go-book/common/loadbalance"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	var (
		servicePort = flag.Int("service.port", 10086, "service port")
		serviceHost = flag.String("service.host", "127.0.0.1", "service host")
		consulPort = flag.Int("consul.port", 8500, "consul port")
		consulHost = flag.String("consul.host", "127.0.0.1", "consul host")
		serviceName = flag.String("service.name", "use-string", "service name")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)
	var discoveryClient discover.DiscoveryClient
	discoveryClient, err := discover.NewKitDiscoverClient(*consulHost, *consulPort)

	if err != nil{
		config.Logger.Println("Get Consul Client failed")
		os.Exit(-1)

	}
	var svc service.Service
	svc = service.NewUseStringService(discoveryClient, &loadbalance.RandomLoadBalance{} )
	useStringEndpoint := endpoint.MakeUseStringEndpoint(svc)
	useStringEndpoint = circuitbreaker.Hystrix(service.StringServiceCommandName)(useStringEndpoint)

	//创建健康检查的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	//把算术运算Endpoint和健康检查Endpoint封装至StringEndpoints
	endpts := endpoint.UseStringEndpoints{
		UseStringEndpoint:      useStringEndpoint,
		HealthCheckEndpoint: 	healthEndpoint,
	}

	//创建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, config.KitLogger)

	instanceId := *serviceName + "-" + uuid.NewV4().String()

	//http server
	go func() {

		config.Logger.Println("Http Server start at port:" + strconv.Itoa(*servicePort))
		//启动前执行注册
		if !discoveryClient.Register(*serviceName, instanceId, "/health", *serviceHost,  *servicePort, nil, config.Logger){
			config.Logger.Printf("use-string-service for service %s failed.", serviceName)
			// 注册失败，服务启动失败
			os.Exit(-1)
		}
		handler := r
		errChan <- http.ListenAndServe(":"  + strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()


	error := <-errChan
	//服务退出取消注册
	discoveryClient.DeRegister(instanceId, config.Logger)
	config.Logger.Println(error)
}
