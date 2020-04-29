package service

// Service Define a service interface
type Service interface {
	// HealthCheck check service health status
	HealthCheck() bool
}

//UserService implement Service interface
type SkAdminService struct {
}

// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (s SkAdminService) HealthCheck() bool {
	return true
}

type ServiceMiddleware func(Service) Service
