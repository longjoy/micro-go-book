package srv_product

import (
	"sync"
)

//商品数量管理
type ProductCountMgr struct {
	productCount map[int]int
	lock         sync.RWMutex
}

func NewProductCountMgr() *ProductCountMgr {
	productMgr := &ProductCountMgr{
		productCount: make(map[int]int, 128),
	}
	return productMgr
}

//商品数量
func (p *ProductCountMgr) Count(productId int) (count int) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	count = p.productCount[productId]
	return
}

//添加商品
func (p *ProductCountMgr) Add(productId, count int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	cur, ok := p.productCount[productId]
	if !ok {
		cur = count
	} else {
		cur += count
	}
	p.productCount[productId] = cur
}
