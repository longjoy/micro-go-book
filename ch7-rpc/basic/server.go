package main

import (
	"github.com/longjoy/micro-go-book/ch7-rpc/basic/string-service"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	stringService := new(service.StringService)
	registerError := rpc.Register(stringService)
	if registerError != nil {
		log.Fatal("Register error: ", registerError)
	}
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "127.0.0.1:1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
