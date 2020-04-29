package client

import (
	"context"
	"github.com/longjoy/micro-go-book/ch13-seckill/pb"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/discover"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/loadbalance"
	"github.com/opentracing/opentracing-go"
)

type OAuthClient interface {
	CheckToken(ctx context.Context, tracer opentracing.Tracer, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)
}

type OAuthClientImpl struct {
	manager     ClientManager
	serviceName string
	loadBalance loadbalance.LoadBalance
	tracer      opentracing.Tracer
}

func (impl *OAuthClientImpl) CheckToken(ctx context.Context, tracer opentracing.Tracer, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	response := new(pb.CheckTokenResponse)
	if err := impl.manager.DecoratorInvoke("/pb.OAuthService/CheckToken", "token_check", tracer, ctx, request, response); err == nil {
		return response, nil
	} else {
		return nil, err
	}
}
func NewOAuthClient(serviceName string, lb loadbalance.LoadBalance, tracer opentracing.Tracer) (OAuthClient, error) {
	if serviceName == "" {
		serviceName = "oauth"
	}
	if lb == nil {
		lb = defaultLoadBalance
	}

	return &OAuthClientImpl{
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
