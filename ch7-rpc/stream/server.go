package main

import (
	"flag"
	pb "github.com/longjoy/micro-go-book/ch7-rpc/stream-pb"
	"github.com/longjoy/micro-go-book/ch7-rpc/stream/string-service"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
	"net"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	stringService := new(string_service.StringService)
	pb.RegisterStringServiceServer(grpcServer, stringService)
	grpcServer.Serve(lis)
}
