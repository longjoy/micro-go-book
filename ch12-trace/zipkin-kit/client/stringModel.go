package client

import (
	"context"
	"github.com/longjoy/micro-go-book/ch12-trace/zipkin-kit/pb"
	"github.com/longjoy/micro-go-book/ch12-trace/zipkin-kit/string-service/endpoint"
	"github.com/pkg/errors"
)

func EncodeGRPCStringRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.StringRequest)
	return &pb.StringRequest{
		RequestType: string(req.RequestType),
		A:           string(req.A),
		B:           string(req.B),
	}, nil
}

func DecodeGRPCStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.StringRequest)
	return endpoint.StringRequest{
		RequestType: string(req.RequestType),
		A:           string(req.A),
		B:           string(req.B),
	}, nil
}

func EncodeGRPCStringResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.StringResponse)

	if resp.Error != nil {
		return &pb.StringResponse{
			Result: string(resp.Result),
			Err:    resp.Error.Error(),
		}, nil
	}

	return &pb.StringResponse{
		Result: string(resp.Result),
		Err:    "",
	}, nil
}

func DecodeGRPCStringResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.StringResponse)
	return endpoint.StringResponse{
		Result: string(resp.Result),
		Error:  errors.New(resp.Err),
	}, nil
}
