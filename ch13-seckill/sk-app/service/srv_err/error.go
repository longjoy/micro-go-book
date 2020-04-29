package srv_err

import (
	"errors"
)

const (
	ErrInvalidRequest      = 1101
	ErrNotFoundProductId   = 1102
	ErrUserCheckAuthFailed = 1103
	ErrUserServiceBusy     = 1104
	ErrActiveNotStart      = 1105
	ErrActiveAlreadyEnd    = 1106
	ErrActiveSaleOut       = 1107
	ErrProcessTimeout      = 1108
	ErrClientClosed        = 1109
)

const (
	ErrServiceBusy     = 1001
	ErrSecKillSucc     = 1002
	ErrNotFoundProduct = 1003
	ErrSoldout         = 1004
	ErrRetry           = 1005
	ErrAlreadyBuy      = 1006
)

var errMsg = map[int]string{
	ErrServiceBusy:     "服务器错误",
	ErrSecKillSucc:     "抢购成功",
	ErrNotFoundProduct: "没有该商品",
	ErrSoldout:         "商品售罄",
	ErrRetry:           "请重试",
	ErrAlreadyBuy:      "已经抢购",
}

func GetErrMsg(code int) error {
	return errors.New(errMsg[code])
}
