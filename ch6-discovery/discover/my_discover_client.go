package discover

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// 服务实例结构体
type InstanceInfo struct {
	ID                string                     `json:"ID"`                // 服务实例ID
	Service           string                     `json:"Service,omitempty"` // 服务发现时返回的服务名
	Name              string                     `json:"Name"`              // 服务名
	Tags              []string                   `json:"Tags,omitempty"`    // 标签，可用于进行服务过滤
	Address           string                     `json:"Address"`           // 服务实例HOST
	Port              int                        `json:"Port"`              // 服务实例端口
	Meta              map[string]string          `json:"Meta,omitempty"`    // 元数据
	EnableTagOverride bool                       `json:"EnableTagOverride"` // 是否允许标签覆盖
	Check             `json:"Check,omitempty"`   // 健康检查相关配置
	Weights           `json:"Weights,omitempty"` // 权重
}

type Check struct {
	DeregisterCriticalServiceAfter string   `json:"DeregisterCriticalServiceAfter"` // 多久之后注销服务
	Args                           []string `json:"Args,omitempty"`                 // 请求参数
	HTTP                           string   `json:"HTTP"`                           // 健康检查地址
	Interval                       string   `json:"Interval,omitempty"`             // Consul 主动检查间隔
	TTL                            string   `json:"TTL,omitempty"`                  // 服务实例主动维持心跳间隔，与Interval只存其一
}

type Weights struct {
	Passing int `json:"Passing"`
	Warning int `json:"Warning"`
}

type MyDiscoverClient struct {
	Host string // Consul 的 Host
	Port int    // Consul 的 端口
}

func NewMyDiscoverClient(consulHost string, consulPort int) (DiscoveryClient, error) {
	return &MyDiscoverClient{
		Host: consulHost,
		Port: consulPort,
	}, nil
}


func (consulClient *MyDiscoverClient) Register(serviceName, instanceId, healthCheckUrl string,instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool {


	// 1.封装服务实例的元数据
	instanceInfo := &InstanceInfo{
		ID:                instanceId,
		Name:              serviceName,
		Address:           instanceHost,
		Port:              instancePort,
		Meta:              meta,
		EnableTagOverride: false,
		Check: Check{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval:                       "15s",
		},
		Weights: Weights{
			Passing: 10,
			Warning: 1,
		},
	}

	byteData, _ := json.Marshal(instanceInfo)

	// 2. 向 Consul 发送服务注册的请求
	req, err := http.NewRequest("PUT",
		"http://"+consulClient.Host+":"+strconv.Itoa(consulClient.Port)+"/v1/agent/service/register",
		bytes.NewReader(byteData))

	if err == nil {
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		client := http.Client{}
		resp, err := client.Do(req)

		// 3. 检查注册结果
		if err != nil {
			log.Println("Register Service Error!")
		} else {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				log.Println("Register Service Success!")
				return true
			} else {
				log.Println("Register Service Error!")
			}
		}
	}
	return false
}

func (consulClient *MyDiscoverClient) DeRegister(instanceId string, logger *log.Logger) bool {
	// 1.发送注销请求
	req, err := http.NewRequest("PUT",
		"http://"+consulClient.Host+":"+strconv.Itoa(consulClient.Port)+"/v1/agent/service/deregister/"+instanceId, nil)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Deregister Service Error!")
	} else {
		resp.Body.Close()
		if resp.StatusCode == 200 {
			log.Println("Deregister Service Success!")
			return true
		} else {
			log.Println("Deregister Service Error!")
		}
	}
	return false
}


func (consulClient *MyDiscoverClient) DiscoverServices(serviceName string, logger *log.Logger) []interface{} {
	// 1. 从 Consul 中获取服务实例列表
	req, err := http.NewRequest("GET",
		"http://"+consulClient.Host+":"+strconv.Itoa(consulClient.Port)+"/v1/health/service/"+serviceName, nil)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Discover Service Error!")
	} else if resp.StatusCode == 200 {

		var serviceList []struct {
			Service InstanceInfo `json:"Service"`
		}
		err = json.NewDecoder(resp.Body).Decode(&serviceList)
		resp.Body.Close()
		if err == nil {
			instances := make([]interface{}, len(serviceList))
			for i := 0; i < len(instances); i++ {
				instances[i] = serviceList[i].Service
			}
			return instances
		}
	}
	return nil
}
