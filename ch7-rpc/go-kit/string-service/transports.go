package string_service

import (
	"context"
	"errors"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/longjoy/micro-go-book/ch7-rpc/pb"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

type grpcServer struct {
	concat grpc.Handler
	diff   grpc.Handler
}

func (s *grpcServer) Concat(ctx context.Context, r *pb.StringRequest) (*pb.StringResponse, error) {
	_, resp, err := s.concat.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.StringResponse), nil
}

func (s *grpcServer) Diff(ctx context.Context, r *pb.StringRequest) (*pb.StringResponse, error) {
	_, resp, err := s.diff.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.StringResponse), nil
}

func NewStringServer(ctx context.Context, endpoints StringEndpoints) pb.StringServiceServer {
	return &grpcServer{
		concat: grpc.NewServer(
			endpoints.StringEndpoint,
			DecodeConcatStringRequest,
			EncodeStringResponse,
		),
		diff: grpc.NewServer(
			endpoints.StringEndpoint,
			DecodeDiffStringRequest,
			EncodeStringResponse,
		),
	}
}

func DecodeConcatStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.StringRequest)
	return StringRequest{
		RequestType: "Concat",
		A:           string(req.A),
		B:           string(req.B),
	}, nil
}

func DecodeDiffStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.StringRequest)
	return StringRequest{
		RequestType: "Diff",
		A:           string(req.A),
		B:           string(req.B),
	}, nil
}

func EncodeStringResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(StringResponse)

	if resp.Error != nil {
		return &pb.StringResponse{
			Ret: resp.Result,
			Err: resp.Error.Error(),
		}, nil
	}

	return &pb.StringResponse{
		Ret: resp.Result,
		Err: "",
	}, nil
}
