#### 安装 protobuf
rpc 部分涉及到 protobuf

```
protoc --version
go get github.com/golang/protobuf
go install github.com/golang/protobuf/protoc-gen-go/
```

```
protoc string.proto --go_out=plugins=grpc:.
```