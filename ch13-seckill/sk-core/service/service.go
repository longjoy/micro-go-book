package service

// Service Define a service interface
type Service interface {

	// Divide calculate a/b
	SecKill() (int, error)
}

//ArithmeticService implement Service interface
type SecKillService struct {
}

// Add implement Add method
func (s SecKillService) SecKill(a, b int) int {
	return a + b
}