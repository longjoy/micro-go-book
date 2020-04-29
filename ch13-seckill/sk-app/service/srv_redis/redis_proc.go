package srv_redis

import (
	"encoding/json"
	"fmt"
	conf "github.com/longjoy/micro-go-book/ch13-seckill/pkg/config"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-app/config"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-app/model"
	"log"
	"time"
)

//写数据到Redis
func WriteHandle() {
	for {
		fmt.Println("wirter data to redis.")
		req := <-config.SkAppContext.SecReqChan
		fmt.Println("accessTime : ", req.AccessTime)
		conn := conf.Redis.RedisConn

		data, err := json.Marshal(req)
		if err != nil {
			log.Printf("json.Marshal req failed. Error : %v, req : %v", err, req)
			continue
		}

		err = conn.LPush(conf.Redis.Proxy2layerQueueName, string(data)).Err()
		if err != nil {
			log.Printf("lpush req failed. Error : %v, req : %v", err, req)
			continue
		}
		log.Printf("lpush req success. req : %v", string(data))
	}
}

//从redis读取数据
func ReadHandle() {
	for {
		conn := conf.Redis.RedisConn
		//阻塞弹出
		data, err := conn.BRPop(time.Second, conf.Redis.Layer2proxyQueueName).Result()
		if err != nil {
			continue
		}

		var result *model.SecResult
		err = json.Unmarshal([]byte(data[1]), &result)
		if err != nil {
			log.Printf("json.Unmarshal failed. Error : %v", err)
			continue
		}

		userKey := fmt.Sprintf("%d_%d", result.UserId, result.ProductId)
		fmt.Println("userKey : ", userKey)
		config.SkAppContext.UserConnMapLock.Lock()
		resultChan, ok := config.SkAppContext.UserConnMap[userKey]
		config.SkAppContext.UserConnMapLock.Unlock()
		if !ok {
			log.Printf("user not found : %v", userKey)
			continue
		}
		log.Printf("request result send to chan")

		resultChan <- result
		log.Printf("request result send to chan succeee, userKey : %v", userKey)
	}
}
