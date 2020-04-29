package transport

import (
	"context"
	"github.com/longjoy/micro-go-book/ch13-seckill/pb"
	"github.com/longjoy/micro-go-book/ch13-seckill/user-service/endpoint"
)

func EncodeGRPCUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(pb.UserRequest)
	return &pb.UserRequest{
		Username: string(req.Username),
		Password: string(req.Password),
	}, nil
}

func DecodeGRPCUserRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UserRequest)
	return endpoint.UserRequest{
		Username: string(req.Username),
		Password: string(req.Password),
	}, nil
}

func EncodeGRPCUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.UserResponse)

	if resp.Error != nil {
		return &pb.UserResponse{
			Result: bool(resp.Result),
			Err:    "error",
		}, nil
	}

	return &pb.UserResponse{
		Result: bool(resp.Result),
		UserId:resp.UserId,
		Err:    "",
	}, nil
}

func DecodeGRPCUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.UserResponse)
	return pb.UserResponse{
		Result: bool(resp.Result),
		Err:    resp.Err,
	}, nil
}
