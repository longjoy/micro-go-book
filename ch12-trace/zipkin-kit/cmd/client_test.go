package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/longjoy/micro-go-book/ch12-trace/zipkin-kit/client"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	var (
		grpcAddr    = flag.String("addr", ":9008", "gRPC address")
		serviceHost = flag.String("service.host", "localhost", "service ip address")
		servicePort = flag.String("service.port", "8009", "service port")
		zipkinURL   = flag.String("zipkin.url", "http://114.67.98.210:9411/api/v2/spans", "Zipkin server url")
	)
	flag.Parse()
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
			hostPort      = *serviceHost + ":" + *servicePort
			serviceName   = "test-service"
			useNoopTracer = (*zipkinURL == "")
			reporter      = zipkinhttp.NewReporter(*zipkinURL)
		)
		defer reporter.Close()
		zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
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
	tr := zipkinTracer
	parentSpan := tr.StartSpan("test")
	defer parentSpan.Flush()

	ctx := zipkin.NewContext(context.Background(), parentSpan)

	clientTracer := kitzipkin.GRPCClientTrace(tr)
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		fmt.Println("gRPC dial err:", err)
	}
	defer conn.Close()

	svr := client.StringDiff(conn, clientTracer)
	result, err := svr.Diff(ctx, "Add", "ppsdd")
	if err != nil {
		fmt.Println("Diff error", err.Error())

	}

	fmt.Println("result =", result)
}
