package plugins

import (
	"github.com/go-kit/kit/log"
	"github.com/gohouse/gorose/v2"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-admin/model"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-admin/service"
	"time"
)

// loggingMiddleware Make a new type
// that contains Service interface and logger instance
type skAdminLoggingMiddleware struct {
	service.Service
	logger log.Logger
}

type activityLoggingMiddleware struct {
	service.ActivityService
	logger log.Logger
}

type productLoggingMiddleware struct {
	service.ProductService
	logger log.Logger
}

// LoggingMiddleware make logging middleware
func SkAdminLoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return skAdminLoggingMiddleware{next, logger}
	}
}

func ActivityLoggingMiddleware(logger log.Logger) service.ActivityServiceMiddleware {
	return func(next service.ActivityService) service.ActivityService {
		return activityLoggingMiddleware{next, logger}
	}
}

func ProductLoggingMiddleware(logger log.Logger) service.ProductServiceMiddleware {
	return func(next service.ProductService) service.ProductService {
		return productLoggingMiddleware{next, logger}
	}
}

func (mw productLoggingMiddleware) CreateProduct(product *model.Product) (err error) {

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Check",
			"product", product,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.ProductService.CreateProduct(product)
	return err
}

func (mw productLoggingMiddleware) GetProductList() ([]gorose.Data, error) {

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Check",
			"took", time.Since(begin),
		)
	}(time.Now())

	data, err := mw.ProductService.GetProductList()
	return data, err
}

func (mw activityLoggingMiddleware) GetActivityList() ([]gorose.Data, error) {

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Check",
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err := mw.ActivityService.GetActivityList()
	return ret, err
}

func (mw activityLoggingMiddleware) CreateActivity(activity *model.Activity) error {

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Check",
			"activity", activity,
			"took", time.Since(begin),
		)
	}(time.Now())

	err := mw.ActivityService.CreateActivity(activity)
	return err
}

func (mw skAdminLoggingMiddleware) HealthCheck() (result bool) {
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
