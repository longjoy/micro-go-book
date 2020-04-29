package main

import (
	"errors"
	"strings"
)

// Service constants
const (
	StrMaxSize = 1024
)

// Service errors
var (
	ErrMaxSize = errors.New("maximum size of 1024 bytes exceeded")

	ErrStrValue = errors.New("maximum size of 1024 bytes exceeded")
)

// Service Define a service interface
type Service interface {
	// Concat a and b
	Concat(a, b string) (string, error)

	// a,b pkg string value
	Diff(a, b string) (string, error)

	// HealthCheck check service health status
	HealthCheck() bool
}

//ArithmeticService implement Service interface
type StringService struct {
}

func (s StringService) Concat(a, b string) (string, error) {
	// test for length overflow
	if len(a)+len(b) > StrMaxSize {
		return "", ErrMaxSize
	}
	return a + b, nil
}

func (s StringService) Diff(a, b string) (string, error) {
	if len(a) < 1 || len(b) < 1 {
		return "", nil
	}
	res := ""
	if len(a) >= len(b) {
		for _, char := range b {
			if strings.Contains(a, string(char)) {
				res = res + string(char)
			}
		}
	} else {
		for _, char := range a {
			if strings.Contains(b, string(char)) {
				res = res + string(char)
			}
		}
	}
	return res, nil
}

// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true。
func (s StringService) HealthCheck() bool {
	return true
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
