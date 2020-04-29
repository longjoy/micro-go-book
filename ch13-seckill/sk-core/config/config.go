package config

import (
	"github.com/go-kit/kit/log"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/bootstrap"
	_ "github.com/longjoy/micro-go-book/ch13-seckill/pkg/bootstrap"
	conf "github.com/longjoy/micro-go-book/ch13-seckill/pkg/config"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-core/service/srv_product"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-core/service/srv_user"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	_ "github.com/openzipkin/zipkin-go/reporter/recorder"
	"github.com/spf13/viper"
	"os"
	"sync"
)

const (
	kConfigType = "CONFIG_TYPE"
)

var ZipkinTracer *zipkin.Tracer
var Logger log.Logger

func init() {
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)
	Logger = log.With(Logger, "caller", log.DefaultCaller)
	viper.AutomaticEnv()
	initDefault()

	if err := conf.LoadRemoteConfig(); err != nil {
		Logger.Log("Fail to load remote config", err)
	}

	if err := conf.Sub("mysql", &conf.MysqlConfig); err != nil {
		Logger.Log("Fail to parse mysql", err)
	}
	if err := conf.Sub("trace", &conf.TraceConfig); err != nil {
		Logger.Log("Fail to parse trace", err)
	}

	if err := conf.Sub("redis", &conf.Redis); err != nil {
		Logger.Log("Fail to parse trace", err)
	}

	if err := conf.Sub("service", &conf.SecKill); err != nil {
		Logger.Log("Fail to parse trace", err)
	}

	zipkinUrl := "http://" + conf.TraceConfig.Host + ":" + conf.TraceConfig.Port + conf.TraceConfig.Url
	Logger.Log("zipkin url", zipkinUrl)
	initTracer(zipkinUrl)
}

func initDefault() {
	viper.SetDefault(kConfigType, "yaml")
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

var SecLayerCtx = &SecLayerContext{
	Read2HandleChan:  make(chan *SecRequest, 1024),
	Handle2WriteChan: make(chan *SecResult, 1024),
	HistoryMap:       make(map[int]*srv_user.UserBuyHistory, 1024),
	ProductCountMgr:  srv_product.NewProductCountMgr(),
}
var CoreCtx = &SkAppCtx{}

type SecResult struct {
	ProductId int    `json:"product_id"` //商品ID
	UserId    int    `json:"user_id"`    //用户ID
	Token     string `json:"token"`      //Token
	TokenTime int64  `json:"token_time"` //Token生成时间
	Code      int    `json:"code"`       //状态码
}

type SecRequest struct {
	ProductId     int             `json:"product_id"` //商品ID
	Source        string          `json:"source"`
	AuthCode      string          `json:"auth_code"`
	SecTime       int64           `json:"sec_time"`
	Nance         string          `json:"nance"`
	UserId        int             `json:"user_id"`
	UserAuthSign  string          `json:"user_auth_sign"` //用户授权签名
	ClientAddr    string          `json:"client_addr"`
	ClientRefence string          `json:"client_refence"`
	CloseNotify   <-chan bool     `json:"-"`
	ResultChan    chan *SecResult `json:"-"`
}

type SkAppCtx struct {
	SecReqChan       chan *SecRequest
	SecReqChanSize   int
	RWSecProductLock sync.RWMutex

	UserConnMap     map[string]chan *SecResult
	UserConnMapLock sync.Mutex
}

const (
	ProductStatusNormal       = 0 //商品状态正常
	ProductStatusSaleOut      = 1 //商品售罄
	ProductStatusForceSaleOut = 2 //商品强制售罄
)

type SecLayerContext struct {
	RWSecProductLock sync.RWMutex

	WaitGroup sync.WaitGroup

	Read2HandleChan  chan *SecRequest
	Handle2WriteChan chan *SecResult

	HistoryMap     map[int]*srv_user.UserBuyHistory
	HistoryMapLock sync.Mutex

	ProductCountMgr *srv_product.ProductCountMgr //商品计数
}
