package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/pborman/uuid"
	"os"
	"strconv"
)

func Register(consulHost, consulPort, svcHost, svcPort string, logger log.Logger) (registar sd.Registrar) {

	// 创建Consul客户端连接
	var client consul.Client
	{
		consulCfg := api.DefaultConfig()
		consulCfg.Address = consulHost + ":" + consulPort
		consulClient, err := api.NewClient(consulCfg)
		if err != nil {
			logger.Log("create consul client error:", err)
			os.Exit(1)
		}

		client = consul.NewClient(consulClient)
	}

	// 设置Consul对服务健康检查的参数
	check := api.AgentServiceCheck{
		HTTP:     "http://" + svcHost + ":" + svcPort + "/health",
		Interval: "10s",
		Timeout:  "1s",
		Notes:    "Consul check service health status.",
	}

	port, _ := strconv.Atoi(svcPort)

	//设置微服务想Consul的注册信息
	reg := api.AgentServiceRegistration{
		ID:      "string-service" + uuid.New(),
		Name:    "string-service",
		Address: svcHost,
		Port:    port,
		Tags:    []string{"string-service"},
		Check:   &check,
	}

	// 执行注册
	registar = consul.NewRegistrar(client, &reg, logger)
	return
}
