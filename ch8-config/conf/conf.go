package conf

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/streadway/amqp"
	"log"
	"net/http"
	"strings"
)

const (
	kAppName       = "APP_NAME"
	kConfigServer  = "CONFIG_SERVER"
	kConfigLabel   = "CONFIG_LABEL"
	kConfigProfile = "CONFIG_PROFILE"
	kConfigType    = "CONFIG_TYPE"
	kAmqpURI       = "AmqpURI"
)

var (
	Resume ResumeConfig
)

type ResumeConfig struct {
	Name string
	Age  int
	Sex  string
}

func init() {
	viper.AutomaticEnv()
	initDefault()
	go StartListener(viper.GetString(kAppName), viper.GetString(kAmqpURI), "springCloudBus")

	if err := loadRemoteConfig(); err != nil {
		log.Fatal("Fail to load config", err)
	}

	if err := sub("resume", &Resume); err != nil {
		log.Fatal("Fail to parse config", err)
	}
}

func initDefault() {
	viper.SetDefault(kAppName, "client-demo")
	viper.SetDefault(kConfigServer, "http://localhost:8888")
	viper.SetDefault(kConfigLabel, "master")
	viper.SetDefault(kConfigProfile, "dev")
	viper.SetDefault(kConfigType, "yaml")
	viper.SetDefault(kAmqpURI, "amqp://admin:admin@114.67.98.210:5672")

}
func handleRefreshEvent(body []byte, consumerTag string) {
	updateToken := &UpdateToken{}
	err := json.Unmarshal(body, updateToken)
	if err != nil {
		log.Printf("Problem parsing UpdateToken: %v", err.Error())
	} else {
		log.Println(consumerTag, updateToken.DestinationService)
		if strings.Contains(updateToken.DestinationService, consumerTag) {
			log.Println("Reloading Viper config from Spring Cloud Config server")
			loadRemoteConfig()
			log.Println(viper.GetString("resume.name"))
		}
	}
}

func loadRemoteConfig() (err error) {
	confAddr := fmt.Sprintf("%v/%v/%v-%v.%v",
		viper.Get(kConfigServer), viper.Get(kConfigLabel),
		viper.Get(kAppName), viper.Get(kConfigProfile),
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
	log.Println("Load config from: ", confAddr)
	return
}

func sub(key string, value interface{}) error {
	log.Printf("配置文件的前缀为：%v", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}
