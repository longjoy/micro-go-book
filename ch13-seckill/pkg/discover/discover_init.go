package discover

import (
	"errors"
	"fmt"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/bootstrap"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/common"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/loadbalance"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"os"
)

var ConsulService DiscoveryClient
var LoadBalance loadbalance.LoadBalance
var Logger *log.Logger
var NoInstanceExistedErr  = errors.New("no available client")

func init() {
	// 1.实例化一个 Consul 客户端，此处实例化了原生态实现版本
	ConsulService = New(bootstrap.DiscoverConfig.Host, bootstrap.DiscoverConfig.Port)
	LoadBalance = new(loadbalance.RandomLoadBalance)
	Logger = log.New(os.Stderr, "", log.LstdFlags)

}

//
func CheckHealth(writer http.ResponseWriter, reader *http.Request) {
	Logger.Println("Health check!")
	_, err := fmt.Fprintln(writer, "Server is OK!")
	if err != nil {
		Logger.Println(err)
	}
}

func DiscoveryService(serviceName string) (*common.ServiceInstance, error) {
	instances := ConsulService.DiscoverServices(serviceName, Logger)

	if len(instances) < 1 {
		Logger.Printf("no available client for %s.", serviceName)
		return nil, NoInstanceExistedErr
	}
	return LoadBalance.SelectService(instances)

}

func Register() {
	//// 实例失败，停止服务
	if ConsulService == nil {
		panic(0)
	}

	//判空 instanceId,通过 go.uuid 获取一个服务实例ID
	instanceId := bootstrap.DiscoverConfig.InstanceId

	if instanceId == "" {
		instanceId = bootstrap.DiscoverConfig.ServiceName + uuid.NewV4().String()
	}

	if !ConsulService.Register(instanceId, bootstrap.HttpConfig.Host, "/health",
		bootstrap.HttpConfig.Port, bootstrap.DiscoverConfig.ServiceName,
		bootstrap.DiscoverConfig.Weight,
		map[string]string{
			"rpcPort": bootstrap.RpcConfig.Port,
		}, nil, Logger) {
		Logger.Printf("register service %s failed.", bootstrap.DiscoverConfig.ServiceName)
		// 注册失败，服务启动失败
		panic(0)
	}
	Logger.Printf(bootstrap.DiscoverConfig.ServiceName+"-service for service %s success.", bootstrap.DiscoverConfig.ServiceName)

}

func Deregister() {
	//// 实例失败，停止服务
	if ConsulService == nil {
		panic(0)
	}
	//判空 instanceId,通过 go.uuid 获取一个服务实例ID
	instanceId := bootstrap.DiscoverConfig.InstanceId

	if instanceId == "" {
		instanceId = bootstrap.DiscoverConfig.ServiceName + "-" + uuid.NewV4().String()
	}
	if !ConsulService.DeRegister(instanceId, Logger) {
		Logger.Printf("deregister for service %s failed.", bootstrap.DiscoverConfig.ServiceName)
		panic(0)
	}
}
