package srv_limit

import (
	"fmt"
	conf "github.com/longjoy/micro-go-book/ch13-seckill/pkg/config"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-app/model"
	"log"
	"sync"
)

//限制管理
type SecLimitMgr struct {
	UserLimitMap map[int]*Limit
	IpLimitMap   map[string]*Limit
	lock         sync.Mutex
}

var SecLimitMgrVars = &SecLimitMgr{
	UserLimitMap: make(map[int]*Limit),
	IpLimitMap:   make(map[string]*Limit),
}

//防作弊
func AntiSpam(req *model.SecRequest) (err error) {
	//判断用户Id是否在黑名单
	_, ok := conf.SecKill.IDBlackMap[req.UserId]
	if ok {
		err = fmt.Errorf("invalid request")
		log.Printf("user[%v] is block by id black", req.UserId)
		return
	}

	//判断客户端IP是否在黑名单
	_, ok = conf.SecKill.IPBlackMap[req.ClientAddr]
	if ok {
		err = fmt.Errorf("invalid request")
		log.Printf("userId[%v] ip[%v] is block by ip black", req.UserId, req.ClientAddr)
	}

	var secIdCount, minIdCount, secIpCount, minIpCount int
	//加锁
	SecLimitMgrVars.lock.Lock()
	{
		//用户Id频率控制
		limit, ok := SecLimitMgrVars.UserLimitMap[req.UserId]
		if !ok {
			limit = &Limit{
				secLimit: &SecLimit{},
				minLimit: &MinLimit{},
			}
			SecLimitMgrVars.UserLimitMap[req.ClientAddr] = limit
		}

		secIdCount = limit.secLimit.Count(req.AccessTime) //获取该秒内该用户访问次数
		minIdCount = limit.secLimit.Count(req.AccessTime) //获取该分钟内该用户访问次数

		//客户端Ip频率控制
		limit, ok = SecLimitMgrVars.IpLimitMap[req.ClientAddr]
		if !ok {
			limit = &Limit{
				secLimit: &SecLimit{},
				minLimit: &MinLimit{},
			}
			SecLimitMgrVars.IpLimitMap[req.ClientAddr] = limit
		}

		secIpCount = limit.secLimit.Count(req.AccessTime) //获取该秒内该IP访问次数
		minIpCount = limit.minLimit.Count(req.AccessTime) //获取该分钟内该IP访问次数
	}
	//释放锁
	SecLimitMgrVars.lock.Unlock()

	//判断该用户一秒内访问次数是否大于配置的最大访问次数
	if secIdCount > conf.SecKill.AccessLimitConf.UserSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	//判断该用户一分钟内访问次数是否大于配置的最大访问次数
	if minIdCount > conf.SecKill.AccessLimitConf.UserMinAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	//判断该IP一秒内访问次数是否大于配置的最大访问次数
	if secIpCount > conf.SecKill.AccessLimitConf.IPSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	//判断该IP一分钟内访问次数是否大于配置的最大访问次数
	if minIpCount > conf.SecKill.AccessLimitConf.IPMinAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	return
}
