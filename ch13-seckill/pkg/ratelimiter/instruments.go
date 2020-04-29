package main

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
	"time"
)

var ErrLimitExceed = errors.New("Rate limit exceed!")

// NewTokenBucketLimitterWithBuildIn 使用x/time/rate创建限流中间件
func NewTokenBucketLimitterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

func DynamicLimitter(interval int, burst int) endpoint.Middleware {
	bucket := rate.NewLimiter(rate.Every(time.Second*time.Duration(interval)), burst)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bucket.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
