package main

import (
	"context"
	"flag"
	"github.com/go-kit/kit/log"
	service "github.com/longjoy/micro-go-book/ch7-rpc/go-kit/string-service"
	"github.com/longjoy/micro-go-book/ch7-rpc/pb"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var svc service.Service
	svc = service.StringService{}

	// add logging middleware
	svc = service.LoggingMiddleware(logger)(svc)

	endpoint := service.MakeStringEndpoint(svc)

	//创建健康检查的Endpoint
	healthEndpoint := service.MakeHealthCheckEndpoint(svc)

	//把算术运算Endpoint和健康检查Endpoint封装至StringEndpoints
	endpts := service.StringEndpoints{
		StringEndpoint:      endpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	handler := service.NewStringServer(ctx, endpts)

	ls, _ := net.Listen("tcp", "127.0.0.1:8080")
	gRPCServer := grpc.NewServer()
	pb.RegisterStringServiceServer(gRPCServer, handler)
	gRPCServer.Serve(ls)

}
