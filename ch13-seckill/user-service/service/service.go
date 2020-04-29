package service

import (
	"context"
	"github.com/longjoy/micro-go-book/ch13-seckill/user-service/model"
	"log"
)

// Service Define a service interface
type Service interface {
	Check(ctx context.Context, username, password string) (int64, error)

	// HealthCheck check service health status
	HealthCheck() bool
}

//UserService implement Service interface
type UserService struct {
}

// Add implement check method
func (s UserService) Check(ctx context.Context, username string, password string) (int64, error) {
	userEntity := model.NewUserModel()
	res, err := userEntity.CheckUser(username, password)
	if err != nil {
		log.Printf("UserEntity.CreateUser, err : %v", err)
		return 0, err
	}
	return res.UserId, nil
}

// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (s UserService) HealthCheck() bool {
	return true
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
