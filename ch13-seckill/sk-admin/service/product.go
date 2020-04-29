package service

import (
	"github.com/gohouse/gorose/v2"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-admin/model"
	"log"
)

type ProductService interface {
	CreateProduct(product *model.Product) error
	GetProductList() ([]gorose.Data, error)
}

type ProductServiceMiddleware func(ProductService) ProductService

type ProductServiceImpl struct {
}

func (p ProductServiceImpl) CreateProduct(product *model.Product) error {
	productEntity := model.NewProductModel()
	err := productEntity.CreateProduct(product)
	if err != nil {
		log.Printf("ProductEntity.CreateProduct, err : %v", err)
		return err
	}
	return nil
}

func (p ProductServiceImpl) GetProductList() ([]gorose.Data, error) {
	productEntity := model.NewProductModel()
	productList, err := productEntity.GetProductList()
	if err != nil {
		log.Printf("ProductEntity.CreateProduct, err : %v", err)
		return nil, err
	}
	return productList, nil
}
