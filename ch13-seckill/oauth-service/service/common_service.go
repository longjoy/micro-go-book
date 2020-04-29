package service



type Service interface {
	// HealthCheck check service health status
	HealthCheck() bool
}

type CommentService struct {

}

// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (s *CommentService) HealthCheck() bool {
	return true
}

func NewCommentService() *CommentService {
	return &CommentService{}
}

