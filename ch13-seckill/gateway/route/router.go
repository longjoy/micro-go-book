package route

import (
	"context"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/log"
	"github.com/longjoy/micro-go-book/ch13-seckill/gateway/config"
	"github.com/longjoy/micro-go-book/ch13-seckill/pb"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/client"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/discover"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/loadbalance"
	"github.com/openzipkin/zipkin-go"
	zipkinhttpsvr "github.com/openzipkin/zipkin-go/middleware/http"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

// HystrixRouter hystrix路由
type HystrixRouter struct {
	svcMap      *sync.Map      //服务实例，存储已经通过hystrix监控服务列表
	logger      log.Logger     //日志工具
	fallbackMsg string         //回调消息
	tracer      *zipkin.Tracer //服务追踪对象
	loadbalance loadbalance.LoadBalance
}

func Routes(zipkinTracer *zipkin.Tracer, fbMsg string, logger log.Logger) http.Handler {
	return HystrixRouter{
		svcMap:      &sync.Map{},
		logger:      logger,
		fallbackMsg: fbMsg,
		tracer:      zipkinTracer,
		loadbalance: &loadbalance.RandomLoadBalance{},
	}
}

func preFilter(r *http.Request) bool {
	//查询原始请求路径，如：/string-service/calculate/10/5
	reqPath := r.URL.Path
	if reqPath == "" {
		return false
	}

	res := config.Match(reqPath)

	if res {
		return true
	}
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		return false
	}
	oauthClient, _ := client.NewOAuthClient("oauth", nil, nil)


	resp, remoteErr := oauthClient.CheckToken(context.Background(), nil, &pb.CheckTokenRequest{
		Token: authToken,
	})


	if remoteErr != nil || resp == nil {
		return false
	} else {
		return true
	}
}

func postFilter() {
	// for custom filter
}
func (router HystrixRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//查询原始请求路径，如：/string-service/calculate/10/5
	reqPath := r.URL.Path
	router.logger.Log("reqPath: ", reqPath)

	// 健康检查直接返回
	if reqPath == "/health" {
		w.WriteHeader(200)
		return
	}


	var err error
	if reqPath == "" || !preFilter(r) {
		err = errors.New("illegal request!")
		w.WriteHeader(403)
		w.Write([]byte(err.Error()))
		return
	}




	//按照分隔符'/'对路径进行分解，获取服务名称serviceName
	pathArray := strings.Split(reqPath, "/")
	serviceName := pathArray[1]

	//检查是否已经加入监控
	if _, ok := router.svcMap.Load(serviceName); !ok {
		//把serviceName作为命令对象，设置参数
		hystrix.ConfigureCommand(serviceName, hystrix.CommandConfig{Timeout: 1000})
		router.svcMap.Store(serviceName, serviceName)
	}

	//执行命令
	err = hystrix.Do(serviceName, func() (err error) {

		//调用consul api查询serviceNam
		serviceInstance, err := discover.DiscoveryService(serviceName)
		if err != nil {
			return err
		}

		director := func(req *http.Request) {
			//重新组织请求路径，去掉服务名称部分
			destPath := strings.Join(pathArray[2:], "/")

			//随机选择一个服务实例
			router.logger.Log("service id", serviceInstance.Host, serviceInstance.Port)

			//设置代理服务地址信息
			req.URL.Scheme = "http"
			req.URL.Host = fmt.Sprintf("%s:%d", serviceInstance.Host, serviceInstance.Port)
			req.URL.Path = "/" + destPath
		}

		var proxyError error = nil
		// 为反向代理增加追踪逻辑，使用如下RoundTrip代替默认Transport
		roundTrip, _ := zipkinhttpsvr.NewTransport(router.tracer, zipkinhttpsvr.TransportTrace(true))

		//反向代理失败时错误处理
		errorHandler := func(ew http.ResponseWriter, er *http.Request, err error) {
			proxyError = err
		}

		proxy := &httputil.ReverseProxy{
			Director:     director,
			Transport:    roundTrip,
			ErrorHandler: errorHandler,
		}
		proxy.ServeHTTP(w, r)

		return proxyError

	}, func(err error) error {
		//run执行失败，返回fallback信息
		router.logger.Log("fallback error description", err.Error())

		return errors.New(router.fallbackMsg)
	})

	// Do方法执行失败，响应错误信息
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}
