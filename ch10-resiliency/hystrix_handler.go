package main

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
	"github.com/longjoy/micro-go-book/common/discover"
	"github.com/longjoy/micro-go-book/common/loadbalance"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

var (
	ErrNoInstances = errors.New("query service instance error")

)

type HystrixHandler struct {

	// 记录hystrix是否已配置
	hystrixs      map[string]bool
	hystrixsMutex *sync.Mutex

	discoveryClient discover.DiscoveryClient
	loadbalance loadbalance.LoadBalance
	logger       *log.Logger
}

func NewHystrixHandler(discoveryClient discover.DiscoveryClient, loadbalance loadbalance.LoadBalance, logger *log.Logger) *HystrixHandler {

	return &HystrixHandler{
		discoveryClient: discoveryClient,
		logger:        	logger,
		hystrixs:      	make(map[string]bool),
		loadbalance:	loadbalance,
		hystrixsMutex: 	&sync.Mutex{},
	}

}

func (hystrixHandler *HystrixHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	reqPath := req.URL.Path
	if reqPath == "" {
		return
	}
	//按照分隔符'/'对路径进行分解，获取服务名称serviceName
	pathArray := strings.Split(reqPath, "/")
	serviceName := pathArray[1]

	if serviceName == "" {
		// 路径不存在
		rw.WriteHeader(404)
		return
	}

	if _, ok := hystrixHandler.hystrixs[serviceName]; !ok {
		hystrixHandler.hystrixsMutex.Lock()
		if _, ok := hystrixHandler.hystrixs[serviceName]; !ok {
			//把serviceName作为 hystrix 命令命名
			hystrix.ConfigureCommand(serviceName, hystrix.CommandConfig{
				// 进行 hystrix 命令自定义
			})
			hystrixHandler.hystrixs[serviceName] = true
		}
		hystrixHandler.hystrixsMutex.Unlock()
	}

	err := hystrix.Do(serviceName, func() error {

		//调用consul api查询serviceName的服务实例列表
		instances := hystrixHandler.discoveryClient.DiscoverServices(serviceName, hystrixHandler.logger)
		instanceList := make([]*api.AgentService, len(instances))
		for i := 0; i < len(instances); i++ {
			instanceList[i] = instances[i].(*api.AgentService)
		}
		// 使用负载均衡算法选取实例
		selectInstance, err := hystrixHandler.loadbalance.SelectService(instanceList)

		if err != nil{
			return ErrNoInstances
		}

		//创建Director
		director := func(req *http.Request) {

			//重新组织请求路径，去掉服务名称部分
			destPath := strings.Join(pathArray[2:], "/")

			hystrixHandler.logger.Println("service id ", selectInstance.ID)

			//设置代理服务地址信息
			req.URL.Scheme = "http"
			req.URL.Host = fmt.Sprintf("%s:%d", selectInstance.Address, selectInstance.Port)
			req.URL.Path = "/" + destPath
		}

		var proxyError error

		// 返回代理异常，用于记录 hystrix.Do 执行失败
		errorHandler := func(ew http.ResponseWriter, er *http.Request, err error) {
			proxyError = err
		}

		proxy := &httputil.ReverseProxy{
			Director:     director,
			ErrorHandler: errorHandler,
		}

		proxy.ServeHTTP(rw, req)

		// 将执行异常反馈 hystrix
		return proxyError

	}, func(e error) error {
		hystrixHandler.logger.Println("proxy error ", e)
		return errors.New("fallback excute")
	})

	// hystrix.Do 返回执行异常
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(err.Error()))
	}

}
