package client

import (
	"context"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/bootstrap"
	conf "github.com/longjoy/micro-go-book/ch13-seckill/pkg/config"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/discover"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/loadbalance"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"
)

var (
	ErrRPCService = errors.New("no rpc service")
)

var defaultLoadBalance loadbalance.LoadBalance = &loadbalance.RandomLoadBalance{}

type ClientManager interface {
	DecoratorInvoke(path string, hystrixName string, tracer opentracing.Tracer,
		ctx context.Context, inputVal interface{}, outVal interface{}) (err error)
}

type DefaultClientManager struct {
	serviceName     string
	logger          *log.Logger
	discoveryClient discover.DiscoveryClient
	loadBalance     loadbalance.LoadBalance
	after           []InvokerAfterFunc
	before          []InvokerBeforeFunc
}

type InvokerAfterFunc func() (err error)

type InvokerBeforeFunc func() (err error)

func (manager *DefaultClientManager) DecoratorInvoke(path string, hystrixName string,
	tracer opentracing.Tracer, ctx context.Context, inputVal interface{}, outVal interface{}) (err error) {

	for _, fn := range manager.before {
		if err = fn(); err != nil {
			return err
		}
	}

	if err = hystrix.Do(hystrixName, func() error {

		instances := manager.discoveryClient.DiscoverServices(manager.serviceName, manager.logger)
		if instance, err := manager.loadBalance.SelectService(instances); err == nil {
			if instance.GrpcPort > 0 {
				if conn, err := grpc.Dial(instance.Host+":"+strconv.Itoa(instance.GrpcPort), grpc.WithInsecure(),
					grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(genTracer(tracer), otgrpc.LogPayloads())), grpc.WithTimeout(1*time.Second)); err == nil {
					if err = conn.Invoke(ctx, path, inputVal, outVal); err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				return ErrRPCService
			}
		} else {
			return err
		}
		return nil
	}, func(e error) error {
		return e
	}); err != nil {
		return err
	} else {
		for _, fn := range manager.after {
			if err = fn(); err != nil {
				return err
			}
		}
		return nil
	}
}

func genTracer(tracer opentracing.Tracer) opentracing.Tracer {
	if tracer != nil {
		return tracer
	}
	zipkinUrl := "http://" + conf.TraceConfig.Host + ":" + conf.TraceConfig.Port + conf.TraceConfig.Url
	zipkinRecorder := bootstrap.HttpConfig.Host + ":" + bootstrap.HttpConfig.Port
	collector, err := zipkin.NewHTTPCollector(zipkinUrl)
	if err != nil {
		log.Fatalf("zipkin.NewHTTPCollector err: %v", err)
	}

	recorder := zipkin.NewRecorder(collector, false, zipkinRecorder, bootstrap.DiscoverConfig.ServiceName)

	res, err := zipkin.NewTracer(
		recorder, zipkin.ClientServerSameSpan(true),
	)
	if err != nil {
		log.Fatalf("zipkin.NewTracer err: %v", err)
	}
	return res

}
