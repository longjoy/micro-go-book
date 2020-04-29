package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-app/model"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-app/service"
)

// CalculateEndpoint define endpoint
type SkAppEndpoints struct {
	SecKillEndpoint        endpoint.Endpoint
	HeathCheckEndpoint     endpoint.Endpoint
	GetSecInfoEndpoint     endpoint.Endpoint
	GetSecInfoListEndpoint endpoint.Endpoint
	TestEndpoint           endpoint.Endpoint
}

func (ue SkAppEndpoints) HealthCheck() bool {
	return false
}

type SecInfoRequest struct {
	productId int `json:"id"`
}

type Response struct {
	Result map[string]interface{} `json:"result"`
	Error  error                  `json:"error"`
	Code   int                    `json:"code"`
}

type SecInfoListResponse struct {
	Result []map[string]interface{} `json:"result"`
	Num    int                      `json:"num"`
	Error  error                    `json:"error"`
}

//  make endpoint
func MakeSecInfoEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SecInfoRequest)

		ret := svc.SecInfo(req.productId)
		return Response{Result: ret, Error: nil}, nil
	}
}

func MakeSecInfoListEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ret, num, error := svc.SecInfoList()
		return SecInfoListResponse{ret, num, error}, nil
	}
}

//  make endpoint
func MakeSecKillEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.SecRequest)
		ret, code, calError := svc.SecKill(&req)
		return Response{Result: ret, Code: code, Error: calError}, nil
	}
}

func MakeTestEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Response{Result: nil, Code: 1, Error: nil}, nil
	}
}

// HealthRequest 健康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{status}, nil
	}
}
