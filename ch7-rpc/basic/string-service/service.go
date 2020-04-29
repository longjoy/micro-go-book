package service

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

type StringRequest struct {
	A string
	B string
}

type Service interface {
	// Concat a and b
	Concat(req StringRequest, ret *string) error

	// a,b pkg string value
	Diff(req StringRequest, ret *string) error
}

//ArithmeticService implement Service interface
type StringService struct {
}

func (s StringService) Concat(req StringRequest, ret *string) error {
	// test for length overflow
	if len(req.A)+len(req.B) > StrMaxSize {
		*ret = ""
		return ErrMaxSize
	}
	*ret = req.A + req.B
	return nil
}

func (s StringService) Diff(req StringRequest, ret *string) error {
	if len(req.A) < 1 || len(req.B) < 1 {
		*ret = ""
		return nil
	}
	res := ""
	if len(req.A) >= len(req.B) {
		for _, char := range req.B {
			if strings.Contains(req.A, string(char)) {
				res = res + string(char)
			}
		}
	} else {
		for _, char := range req.A {
			if strings.Contains(req.B, string(char)) {
				res = res + string(char)
			}
		}
	}
	*ret = res
	return nil
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
