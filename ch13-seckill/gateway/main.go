package main

import (
	"flag"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/log"
	_ "github.com/longjoy/micro-go-book/ch13-seckill/gateway/config"
	"github.com/longjoy/micro-go-book/ch13-seckill/gateway/route"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/bootstrap"
	register "github.com/longjoy/micro-go-book/ch13-seckill/pkg/discover"
	"github.com/openzipkin/zipkin-go"
	zipkinhttpsvr "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// 创建环境变量
	var (
		zipkinURL = flag.String("zipkin.url", "http://114.67.98.210:9411/api/v2/spans", "Zipkin server url")
	)
	flag.Parse()

	//创建日志组件
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var zipkinTracer *zipkin.Tracer
	{
		var (
			err           error
			useNoopTracer = *zipkinURL == ""
			reporter      = zipkinhttp.NewReporter(*zipkinURL)
		)
		defer reporter.Close()
		zEP, _ := zipkin.NewEndpoint(bootstrap.HttpConfig.Host, bootstrap.HttpConfig.Port)
		zipkinTracer, err = zipkin.NewTracer(
			reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer),
		)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		if !useNoopTracer {
			logger.Log("tracer", "Zipkin", "type", "Native", "URL", *zipkinURL)
		}
	}
	register.Register()

	tags := map[string]string{
		"component": "gateway_server",
	}

	hystrixRouter := route.Routes(zipkinTracer, "Circuit Breaker:Service unavailable", logger)

	handler := zipkinhttpsvr.NewServerMiddleware(
		zipkinTracer,
		zipkinhttpsvr.SpanName(bootstrap.DiscoverConfig.ServiceName),
		zipkinhttpsvr.TagResponseSize(true),
		zipkinhttpsvr.ServerTags(tags),
	)(hystrixRouter)

	errc := make(chan error)

	//启用hystrix实时监控，监听端口为9010
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go func() {
		errc <- http.ListenAndServe(net.JoinHostPort("", "9010"), hystrixStreamHandler)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	//开始监听
	go func() {
		logger.Log("transport", "HTTP", "addr", "9090")
		register.Register()
		errc <- http.ListenAndServe(":9090", handler)
	}()

	// 开始运行，等待结束
	error := <-errc
	//服务退出取消注册
	register.Deregister()
	logger.Log("exit", error)
}
