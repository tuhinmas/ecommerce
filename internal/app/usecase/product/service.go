package product

import (
	"context"
	repository "ecommerce/internal/app/repository/product"
	"ecommerce/internal/entity"
	"ecommerce/pkg/helper"
	"ecommerce/pkg/identifier"
	"ecommerce/pkg/validator"
)

type Service interface {
	ProductList(ctx context.Context, filter entity.QueryRequest) (resp []entity.GetProductListResponse, pagination *helper.Pagination, err error)
	CreateProduct(ctx context.Context, request entity.CreateProductRequest) (err error)
}

type service struct {
	repository repository.ProductRepository
	validator  validator.Validator
	identifier identifier.Identifier
}

func NewProductService(
	repository repository.ProductRepository,
	validator validator.Validator,
	identifier identifier.Identifier,
) Service {
	return &service{
		repository: repository,
		validator:  validator,
		identifier: identifier,
	}
}
