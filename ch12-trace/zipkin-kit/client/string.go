package client

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/longjoy/micro-go-book/ch12-trace/zipkin-kit/pb"
	endpts "github.com/longjoy/micro-go-book/ch12-trace/zipkin-kit/string-service/endpoint"
	"github.com/longjoy/micro-go-book/ch12-trace/zipkin-kit/string-service/service"
	"google.golang.org/grpc"
)

func StringDiff(conn *grpc.ClientConn, clientTracer kitgrpc.ClientOption) service.Service {

	var ep = grpctransport.NewClient(conn,
		"pb.StringService",
		"Diff",
		EncodeGRPCStringRequest,
		DecodeGRPCStringResponse,
		pb.StringResponse{},
		clientTracer,
	).Endpoint()

	StringEp := endpts.StringEndpoints{
		StringEndpoint: ep,
	}
	return StringEp
}
