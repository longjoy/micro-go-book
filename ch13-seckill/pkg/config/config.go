package conf

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/bootstrap"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/discover"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	_ "github.com/openzipkin/zipkin-go/reporter/recorder"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strconv"
)

const (
	kConfigType = "CONFIG_TYPE"
)

var ZipkinTracer *zipkin.Tracer
var Logger log.Logger

func initDefault() {
	viper.SetDefault(kConfigType, "yaml")

}

func init() {
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)
	Logger = log.With(Logger, "caller", log.DefaultCaller)
	viper.AutomaticEnv()
	initDefault()

	if err := LoadRemoteConfig(); err != nil {
		Logger.Log("Fail to load remote config", err)
	}

	//if err := Sub("mysql", &MysqlConfig); err != nil {
	//	Logger.Log("Fail to parse mysql", err)
	//}
	if err := Sub("trace", &TraceConfig); err != nil {
		Logger.Log("Fail to parse trace", err)
	}
	zipkinUrl := "http://" + TraceConfig.Host + ":" + TraceConfig.Port + TraceConfig.Url
	Logger.Log("zipkin url", zipkinUrl)
	initTracer(zipkinUrl)
}

func initTracer(zipkinURL string) {
	var (
		err           error
		useNoopTracer = zipkinURL == ""
		reporter      = zipkinhttp.NewReporter(zipkinURL)
	)
	//defer reporter.Close()
	zEP, _ := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName, bootstrap.HttpConfig.Port)
	ZipkinTracer, err = zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer),
	)
	if err != nil {
		Logger.Log("err", err)
		os.Exit(1)
	}
	if !useNoopTracer {
		Logger.Log("tracer", "Zipkin", "type", "Native", "URL", zipkinURL)
	}
}
func LoadRemoteConfig() (err error) {
	serviceInstance, err := discover.DiscoveryService(bootstrap.ConfigServerConfig.Id)
	if err != nil {
		return
	}
	configServer := "http://" + serviceInstance.Host + ":" + strconv.Itoa(serviceInstance.Port)
	confAddr := fmt.Sprintf("%v/%v/%v-%v.%v",
		configServer, bootstrap.ConfigServerConfig.Label,
		bootstrap.DiscoverConfig.ServiceName, bootstrap.ConfigServerConfig.Profile,
		viper.Get(kConfigType))
	resp, err := http.Get(confAddr)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	viper.SetConfigType(viper.GetString(kConfigType))
	if err = viper.ReadConfig(resp.Body); err != nil {
		return
	}
	Logger.Log("Load config from: ", confAddr)
	return
}

func Sub(key string, value interface{}) error {
	Logger.Log("配置文件的前缀为：", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}
