package plugins

import (
	"github.com/go-kit/kit/log"
	"github.com/longjoy/micro-go-book/ch10-resiliency/use-string-service/service"
	"time"
)

// loggingMiddleware Make a new type
// that contains Service interface and logger instance
type loggingMiddleware struct {
	service.Service
	logger log.Logger
}

// LoggingMiddleware make logging middleware
func LoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return loggingMiddleware{next, logger}
	}
}


func (mw loggingMiddleware) UseStringService(operationType, a, b string) (ret string, err error) {

	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "UseStringService",
			"a", a,
			"b", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err = mw.Service.UseStringService(operationType, a, b)
	return ret, err
}

func (mw loggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "HealthChcek",
			"result", result,
			"took", time.Since(begin),
		)
	}(time.Now())
	result = mw.Service.HealthCheck()
	return
}


