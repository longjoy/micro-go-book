package srv_redis

import (
	"encoding/json"
	"fmt"
	conf "github.com/longjoy/micro-go-book/ch13-seckill/pkg/config"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-core/config"
	"log"
	"time"
)

func RunProcess() {
	for i := 0; i < conf.SecKill.CoreReadRedisGoroutineNum; i++ {
		go HandleReader()
	}

	for i := 0; i < conf.SecKill.CoreWriteRedisGoroutineNum; i++ {
		go HandleWrite()
	}

	for i := 0; i < conf.SecKill.CoreHandleGoroutineNum; i++ {
		go HandleUser()
	}

	log.Printf("all process goroutine started")
	return
}

func HandleReader() {
	log.Printf("read goroutine running %v", conf.Redis.Proxy2layerQueueName)
	for {
		conn := conf.Redis.RedisConn
		for {
			//从Redis队列中取出数据
			data, err := conn.BRPop(time.Second, conf.Redis.Proxy2layerQueueName).Result()
			if err != nil {
				continue
			}
			log.Printf("brpop from proxy to layer queue, data : %s\n", data)

			//转换数据结构
			var req config.SecRequest
			err = json.Unmarshal([]byte(data[1]), &req)
			if err != nil {
				log.Printf("unmarshal to secrequest failed, err : %v", err)
				continue
			}

			//判断是否超时
			nowTime := time.Now().Unix()
			//int64(config.SecLayerCtx.SecLayerConf.MaxRequestWaitTimeout)
			fmt.Println(nowTime, " ", req.SecTime, " ", 100)
			if nowTime-req.SecTime >= int64(conf.SecKill.MaxRequestWaitTimeout) {
				log.Printf("req[%v] is expire", req)
				continue
			}

			//设置超时时间
			timer := time.NewTicker(time.Millisecond * time.Duration(conf.SecKill.CoreWaitResultTimeout))
			select {
			case config.SecLayerCtx.Read2HandleChan <- &req:
			case <-timer.C:
				log.Printf("send to handle chan timeout, req : %v", req)
				break
			}
		}
	}
}

func HandleWrite() {
	log.Println("handle write running")

	for res := range config.SecLayerCtx.Handle2WriteChan {
		fmt.Println("===", res)
		err := sendToRedis(res)
		if err != nil {
			log.Printf("send to redis, err : %v, res : %v", err, res)
			continue
		}
	}
}

//将数据推入到Redis队列
func sendToRedis(res *config.SecResult) (err error) {
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("marshal failed, err : %v", err)
		return
	}

	fmt.Printf("推入队列前~~ %v", conf.Redis.Layer2proxyQueueName)
	conn := conf.Redis.RedisConn
	err = conn.LPush(conf.Redis.Layer2proxyQueueName, string(data)).Err()
	fmt.Println("推入队列后~~")
	if err != nil {
		log.Printf("rpush layer to proxy redis queue failed, err : %v", err)
		return
	}
	log.Printf("lpush layer to proxy success. data[%v]", string(data))

	return
}
