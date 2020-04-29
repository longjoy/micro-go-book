package model

import (
	"github.com/gohouse/gorose/v2"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/mysql"
	"log"
)

type Product struct {
	ProductId   int    `json:"product_id"`   //商品Id
	ProductName string `json:"product_name"` //商品名称
	Total       int    `json:"total"`        //商品数量
	Status      int    `json:"status"`       //商品状态
}

type ProductModel struct {
}

func NewProductModel() *ProductModel {
	return &ProductModel{}
}

func (p *ProductModel) getTableName() string {
	return "product"
}

func (p *ProductModel) GetProductList() ([]gorose.Data, error) {
	conn := mysql.DB()
	list, err := conn.Table(p.getTableName()).Get()
	if err != nil {
		log.Printf("Error : %v", err)
		return nil, err
	}
	return list, nil
}

func (p *ProductModel) CreateProduct(product *Product) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"product_name": product.ProductName,
		"total":        product.Total,
		"status":       product.Status,
	}).Insert()
	if err != nil {
		log.Printf("Error : %v", err)
		return err
	}
	return nil
}
