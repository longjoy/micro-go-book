package model

import (
	"github.com/gohouse/gorose/v2"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/mysql"
	"log"
)

const (
	ActivityStatusNormal  = 0
	ActivityStatusDisable = 1
	ActivityStatusExpire  = 2
)

type Activity struct {
	ActivityId   int    `json:"activity_id"`   //活动Id
	ActivityName string `json:"activity_name"` //活动名称
	ProductId    int    `json:"product_id"`    //商品Id
	StartTime    int64  `json:"start_time"`    //开始时间
	EndTime      int64  `json:"end_time"`      //结束时间
	Total        int    `json:"total"`         //商品总数
	Status       int    `json:"status"`        //状态

	StartTimeStr string  `json:"start_time_str"`
	EndTimeStr   string  `json:"end_time_str"`
	StatusStr    string  `json:"status_str"`
	Speed        int     `json:"speed"`
	BuyLimit     int     `json:"buy_limit"`
	BuyRate      float64 `json:"buy_rate"`
}

type SecProductInfoConf struct {
	ProductId         int     `json:"product_id"`           //商品Id
	StartTime         int64   `json:"start_time"`           //开始时间
	EndTime           int64   `json:"end_time"`             //结束时间
	Status            int     `json:"status"`               //状态
	Total             int     `json:"total"`                //商品总数
	Left              int     `json:"left"`                 //剩余商品数
	OnePersonBuyLimit int     `json:"one_person_buy_limit"` //一个人购买限制
	BuyRate           float64 `json:"buy_rate"`             //买中几率
	SoldMaxLimit      int     `json:"sold_max_limit"`       //每秒最多能卖多少个
}

type ActivityModel struct {
}

func NewActivityModel() *ActivityModel {
	return &ActivityModel{}
}

func (p *ActivityModel) getTableName() string {
	return "activity"
}

func (p *ActivityModel) GetActivityList() ([]gorose.Data, error) {
	conn := mysql.DB()
	list, err := conn.Table(p.getTableName()).Order("activity_id desc").Get()
	if err != nil {
		log.Printf("Error : %v", err)
		return nil, err
	}
	return list, nil
}

func (p *ActivityModel) CreateActivity(activity *Activity) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Data(
		map[string]interface{}{
			"activity_name": activity.ActivityName,
			"product_id":    activity.ProductId,
			"start_time":    activity.StartTime,
			"end_time":      activity.EndTime,
			"total":         activity.Total,
			"sec_speed":     activity.Speed,
			"buy_limit":     activity.BuyLimit,
			"buy_rate":      activity.BuyRate,
		},
	).Insert()
	if err != nil {
		return err
	}
	return nil
}
