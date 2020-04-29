package srv_user

import "sync"

//用户购买历史记录
type UserBuyHistory struct {
	History map[int]int
	Lock    sync.RWMutex
}

func (p *UserBuyHistory) GetProductBuyCount(productId int) int {
	p.Lock.RLock()
	defer p.Lock.RUnlock()

	count, _ := p.History[productId]
	return count
}

func (p *UserBuyHistory) Add(productId, count int) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	cur, ok := p.History[productId]
	if !ok {
		cur = count
	} else {
		cur += count
	}

	p.History[productId] = cur
}
