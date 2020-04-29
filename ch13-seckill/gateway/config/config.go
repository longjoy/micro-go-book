package config

import (
	"github.com/go-kit/kit/log"
	conf "github.com/longjoy/micro-go-book/ch13-seckill/pkg/config"
	"github.com/spf13/viper"
	"os"
)

const (
	kConfigType = "CONFIG_TYPE"
)

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
	if err := conf.Sub("auth", &AuthPermitConfig); err != nil {
		Logger.Log("Fail to parse config", err)
	}
}
func initDefault() {
	viper.SetDefault(kConfigType, "yaml")
}
