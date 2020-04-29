package service

import (
	"fmt"
	conf "github.com/longjoy/micro-go-book/ch13-seckill/pkg/config"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-app/config"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-app/model"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-app/service/srv_err"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-app/service/srv_limit"
	"log"
	"math/rand"
	"time"
)

// Service Define a service interface
type Service interface {
	// HealthCheck check service health status
	HealthCheck() bool
	SecInfo(productId int) (date map[string]interface{})
	SecKill(req *model.SecRequest) (map[string]interface{}, int, error)
	SecInfoList() ([]map[string]interface{}, int, error)
}

//UserService implement Service interface
type SkAppService struct {
}

// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (s SkAppService) HealthCheck() bool {
	return true
}

type ServiceMiddleware func(Service) Service

func (s SkAppService) SecInfo(productId int) (date map[string]interface{}) {
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.RUnlock()

	v, ok := conf.SecKill.SecProductInfoMap[productId]
	if !ok {
		return nil
	}

	data := make(map[string]interface{})
	data["product_id"] = productId
	data["start_time"] = v.StartTime
	data["end_time"] = v.EndTime
	data["status"] = v.Status

	return data
}

func (s SkAppService) SecKill(req *model.SecRequest) (map[string]interface{}, int, error) {
	//对Map加锁处理
	//config.SkAppContext.RWSecProductLock.RLock()
	//defer config.SkAppContext.RWSecProductLock.RUnlock()
	var code int
	err := srv_limit.AntiSpam(req)
	if err != nil {
		code = srv_err.ErrUserServiceBusy
		log.Printf("userId antiSpam [%d] failed, req[%v]", req.UserId, err)
		return nil, code, err
	}

	data, code, err := SecInfoById(req.ProductId)
	if err != nil {
		log.Printf("userId[%d] secInfoById Id failed, req[%v]", req.UserId, req)
		return nil, code, err
	}

	userKey := fmt.Sprintf("%d_%d", req.UserId, req.ProductId)
	ResultChan := make(chan *model.SecResult, 1)
	config.SkAppContext.UserConnMapLock.Lock()
	config.SkAppContext.UserConnMap[userKey] = ResultChan
	config.SkAppContext.UserConnMapLock.Unlock()

	//将请求送入通道并推入到redis队列当中
	config.SkAppContext.SecReqChan <- req

	ticker := time.NewTicker(time.Millisecond * time.Duration(conf.SecKill.AppWaitResultTimeout))

	defer func() {
		ticker.Stop()
		config.SkAppContext.UserConnMapLock.Lock()
		delete(config.SkAppContext.UserConnMap, userKey)
		config.SkAppContext.UserConnMapLock.Unlock()
	}()

	select {
	case <-ticker.C:
		code = srv_err.ErrProcessTimeout
		err = fmt.Errorf("request timeout")
		return nil, code, err
	case <-req.CloseNotify:
		code = srv_err.ErrClientClosed
		err = fmt.Errorf("client already closed")
		return nil, code, err
	case result := <-ResultChan:
		code = result.Code
		if code != 1002 {
			return data, code, srv_err.GetErrMsg(code)
		}
		log.Printf("secKill success")
		data["product_id"] = result.ProductId
		data["token"] = result.Token
		data["user_id"] = result.UserId
		return data, code, nil
	}
}

func NewSecRequest() *model.SecRequest {
	secRequest := &model.SecRequest{
		ResultChan: make(chan *model.SecResult, 1),
	}
	return secRequest
}

func (s SkAppService) SecInfoList() ([]map[string]interface{}, int, error) {
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.RUnlock()

	var data []map[string]interface{}
	for _, v := range conf.SecKill.SecProductInfoMap {
		item, _, err := SecInfoById(v.ProductId)
		if err != nil {
			log.Printf("get sec info, err : %v", err)
			continue
		}
		data = append(data, item)
	}
	return data, 0, nil
}

func SecInfoById(productId int) (map[string]interface{}, int, error) {
	//对Map加锁处理
	//config.SkAppContext.RWSecProductLock.RLock()
	//defer config.SkAppContext.RWSecProductLock.RUnlock()

	var code int
	v, ok := conf.SecKill.SecProductInfoMap[productId]

	if !ok {
		return nil, srv_err.ErrNotFoundProductId, fmt.Errorf("not found product_id:%d", productId)
	}
	start := false      //秒杀活动是否开始
	end := false        //秒杀活动是否结束
	status := "success" //状态
	var err error
	nowTime := time.Now().Unix()
	//秒杀活动没有开始
	if nowTime-v.StartTime < 0 {
		start = false
		end = false
		status = "second kill not start"
		code = srv_err.ErrActiveNotStart
		err = fmt.Errorf(status)
	}

	//秒杀活动已经开始
	if nowTime-v.StartTime > 0 {
		start = true
	}

	//秒杀活动已经结束
	if nowTime-v.EndTime > 0 {
		start = false
		end = true
		status = "second kill is already end"
		code = srv_err.ErrActiveAlreadyEnd
		err = fmt.Errorf(status)

	}

	//商品已经被停止或售磬
	if v.Status == config.ProductStatusForceSaleOut || v.Status == config.ProductStatusSaleOut {
		start = false
		end = false
		status = "product is sale out"
		code = srv_err.ErrActiveSaleOut
		err = fmt.Errorf(status)

	}

	curRate := rand.Float64()
	/**
	 * 放大于购买比率的1.5倍的请求进入core层
	 */
	if curRate > v.BuyRate*1.5 {
		start = false
		end = false
		status = "retry"
		code = srv_err.ErrRetry
		err = fmt.Errorf(status)
	}

	//组装数据
	data := map[string]interface{}{
		"product_id": productId,
		"start":      start,
		"end":        end,
		"status":     status,
	}
	return data, code, err
}
