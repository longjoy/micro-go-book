package loadbalance

import (
	"errors"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/common"
	"math/rand"
)

// 负载均衡器
type LoadBalance interface {
	SelectService(service []*common.ServiceInstance) (*common.ServiceInstance, error)
}

type RandomLoadBalance struct {
}

// 随机负载均衡
func (loadBalance *RandomLoadBalance) SelectService(services []*common.ServiceInstance) (*common.ServiceInstance, error) {

	if services == nil || len(services) == 0 {
		return nil, errors.New("service instances are not exist")
	}

	return services[rand.Intn(len(services))], nil
}

type WeightRoundRobinLoadBalance struct {
}

// 权重平滑负载均衡
func (loadBalance *WeightRoundRobinLoadBalance) SelectService(services []*common.ServiceInstance) (best *common.ServiceInstance, err error) {

	if services == nil || len(services) == 0 {
		return nil, errors.New("service instances are not exist")
	}

	total := 0
	for i := 0; i < len(services); i++ {
		w := services[i]
		if w == nil {
			continue
		}

		w.CurWeight += w.Weight

		total += w.Weight
		if best == nil || w.CurWeight > best.CurWeight {
			best = w
		}
	}

	if best == nil {
		return nil, nil
	}

	best.CurWeight -= total
	return best, nil
}
