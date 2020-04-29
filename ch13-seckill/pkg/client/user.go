package client

import (
	"context"
	"github.com/longjoy/micro-go-book/ch13-seckill/pb"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/discover"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/loadbalance"
	"github.com/opentracing/opentracing-go"
)

type UserClient interface {
	CheckUser(ctx context.Context, tracer opentracing.Tracer, request *pb.UserRequest) (*pb.UserResponse, error)
}

type UserClientImpl struct {
	/**
	* 可以配置负载均衡策略，重试、等机制。也可以配置invokeAfter和invokerBefore
	 */
	manager     ClientManager
	serviceName string
	loadBalance loadbalance.LoadBalance
	tracer      opentracing.Tracer
}

func (impl *UserClientImpl) CheckUser(ctx context.Context, tracer opentracing.Tracer, request *pb.UserRequest) (*pb.UserResponse, error) {
	response := new(pb.UserResponse)
	if err := impl.manager.DecoratorInvoke("/pb.UserService/Check", "user_check", tracer, ctx, request, response); err == nil {
		return response, nil
	} else {
		return nil, err
	}
}

func NewUserClient(serviceName string, lb loadbalance.LoadBalance, tracer opentracing.Tracer) (UserClient, error) {
	if serviceName == "" {
		serviceName = "user"
	}
	if lb == nil {
		lb = defaultLoadBalance
	}

	return &UserClientImpl{
		manager: &DefaultClientManager{
			serviceName: serviceName,
			loadBalance: lb,
			discoveryClient:discover.ConsulService,
			logger:discover.Logger,
		},
		serviceName: serviceName,
		loadBalance: lb,
		tracer:      tracer,
	}, nil

}
