package string_service

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

// loggingMiddleware Make a new type
// that contains Service interface and logger instance
type loggingMiddleware struct {
	Service
	logger log.Logger
}

// LoggingMiddleware make logging middleware
func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{next, logger}
	}
}

func (mw loggingMiddleware) Concat(ctx context.Context, a, b string) (ret string, err error) {

	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Concat",
			"a", a,
			"b", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err = mw.Service.Concat(ctx, a, b)
	return ret, err
}

func (mw loggingMiddleware) Diff(ctx context.Context, a, b string) (ret string, err error) {

	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Diff",
			"a", a,
			"b", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err = mw.Service.Diff(ctx, a, b)
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
	result = true
	return
}
