package string_service

import (
	"context"
	"errors"
	"github.com/longjoy/micro-go-book/ch7-rpc/stream-pb"
	"io"
	"log"
	"strings"
)

const (
	StrMaxSize = 1024
)

// Service errors
var (
	ErrMaxSize = errors.New("maximum size of 1024 bytes exceeded")

	ErrStrValue = errors.New("maximum size of 1024 bytes exceeded")
)

type StringService struct{}

func (s *StringService) LotsOfServerStream(req *stream_pb.StringRequest, qs stream_pb.StringService_LotsOfServerStreamServer) error {
	response := stream_pb.StringResponse{Ret: req.A + req.B}
	for i := 0; i < 10; i++ {
		qs.Send(&response)
	}
	return nil
}

func (s *StringService) LotsOfClientStream(qs stream_pb.StringService_LotsOfClientStreamServer) error {
	var params []string
	for {
		in, err := qs.Recv()
		if err == io.EOF {
			qs.SendAndClose(&stream_pb.StringResponse{Ret: strings.Join(params, "")})
			return nil
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
			return err
		}
		params = append(params, in.A, in.B)
	}
}
func (s *StringService) LotsOfServerAndClientStream(qs stream_pb.StringService_LotsOfServerAndClientStreamServer) error {
	for {
		in, err := qs.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("failed to recv %v", err)
			return err
		}
		qs.Send(&stream_pb.StringResponse{Ret: in.A + in.B})
	}
	return nil
}

func (s *StringService) Concat(ctx context.Context, req *stream_pb.StringRequest) (*stream_pb.StringResponse, error) {
	if len(req.A)+len(req.B) > StrMaxSize {
		response := stream_pb.StringResponse{Ret: ""}
		return &response, nil
	}
	response := stream_pb.StringResponse{Ret: req.A + req.B}
	return &response, nil
}
