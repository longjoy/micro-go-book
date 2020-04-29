package plugins

import (
	"github.com/go-kit/kit/log"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-app/model"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-app/service"
	"time"
)

// loggingMiddleware Make a new type
// that contains Service interface and logger instance
type skAppLoggingMiddleware struct {
	service.Service
	logger log.Logger
}

// LoggingMiddleware make logging middleware
func SkAppLoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return skAppLoggingMiddleware{next, logger}
	}
}

func (mw skAppLoggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "HealthChcek",
			"result", result,
			"took", time.Since(begin),
		)
	}(time.Now())

	result = mw.Service.HealthCheck()
	return
}

func (mw skAppLoggingMiddleware) SecInfo(productId int) map[string]interface{} {

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Check",
			"took", time.Since(begin),
		)
	}(time.Now())

	ret := mw.Service.SecInfo(productId)
	return ret
}

func (mw skAppLoggingMiddleware) SecInfoList() ([]map[string]interface{}, int, error) {

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Check",
			"took", time.Since(begin),
		)
	}(time.Now())

	data, num, error := mw.Service.SecInfoList()
	return data, num, error
}

func (mw skAppLoggingMiddleware) SecKill(req *model.SecRequest) (map[string]interface{}, int, error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Check",
			"took", time.Since(begin),
		)
	}(time.Now())

	result, num, error := mw.Service.SecKill(req)
	return result, num, error
}
