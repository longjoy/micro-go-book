package string_service

import (
	"context"
	"errors"
	"fmt"
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
	Concat(ctx context.Context, a, b string) (string, error)

	// a,b pkg string value
	Diff(ctx context.Context, a, b string) (string, error)
}

//ArithmeticService implement Service interface
type StringService struct {
}

func (s StringService) Concat(ctx context.Context, a, b string) (string, error) {
	// test for length overflow
	if len(a)+len(b) > StrMaxSize {
		return "", ErrMaxSize
	}
	fmt.Printf("StringService Concat return %s", a+b)
	return a + b, nil
}

func (s StringService) Diff(ctx context.Context, a, b string) (string, error) {
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

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
