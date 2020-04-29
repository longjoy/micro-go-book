package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/longjoy/micro-go-book/ch10-resiliency/use-string-service/service"
)

// CalculateEndpoint define endpoint
type UseStringEndpoints struct {
	UseStringEndpoint      endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}


// StringRequest define request struct
type UseStringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

// StringResponse define response struct
type UseStringResponse struct {
	Result string `json:"result"`
	Error  string  `json:"error"`
}

//// MakeStringEndpoint make endpoint
//func MakeUseStringEndpoint(svc service.Service) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		req := request.(UseStringRequest)
//
//		var (
//			res, a, b, opErrorString string
//			opError   error
//		)
//
//		a = req.A
//		b = req.B
//
//		res, opError = svc.UseStringService(req.RequestType, a, b)
//
//		if opError != nil{
//			opErrorString = opError.Error()
//		}
//
//		return UseStringResponse{Result: res, Error: opErrorString}, nil
//	}
//}

func MakeUseStringEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UseStringRequest)

		var (
			res, a, b string
			opError   error
		)

		a = req.A
		b = req.B

		res, opError = svc.UseStringService(req.RequestType, a, b)

		return UseStringResponse{Result: res}, opError
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
