package plugins

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/longjoy/micro-go-book/ch13-seckill/user-service/service"
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

func (mw loggingMiddleware) Check(ctx context.Context, a, b string) (ret int64, err error) {

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Check",
			"username", a,
			"pwd", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err = mw.Service.Check(ctx, a, b)
	return ret, err
}

func (mw loggingMiddleware) HealthCheck() (result bool) {
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
