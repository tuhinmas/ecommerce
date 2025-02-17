package product

import (
	"context"
	"ecommerce/internal/entity"
	"ecommerce/pkg/helper"
	"net/http"
)

func (s *service) ProductList(ctx context.Context, filter entity.QueryRequest) (resp []entity.GetProductListResponse, pagination *helper.Pagination, err error) {
	product, err := s.repository.GetProduct(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	productIds := make([]string, len(product))
	for i, p := range product {
		productIds[i] = p.Id
	}

	skus, err := s.repository.GetSku(ctx, productIds)
	if err != nil {
		return nil, nil, err
	}

	productMap := make(map[string]entity.GetProductListResponse)

	for _, p := range skus {

		product, ok := productMap[p.Id]
		if !ok {
			product = entity.GetProductListResponse{
				Id:   p.Id,
				Name: p.Name,
				Sku:  []entity.Sku{},
			}
		}

		product.Sku = append(product.Sku, entity.Sku{
			Id:      p.SkuId,
			Image:   p.Image,
			Uom:     p.Uom,
			Price:   p.Price,
			Stock:   p.Stock,
			Variant: p.Variant,
		})

		productMap[p.Id] = product

	}

	for _, p := range productMap {
		resp = append(resp, p)
	}

	totalDT, err := s.repository.GetTotalProduct(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Calculate pagination
	pagination, err = helper.CalculatePagination(ctx, filter.Limit, totalDT)
	if err != nil {
		return nil, nil, err
	}

	pagination.Page = filter.Page

	return
}

func (s *service) CreateProduct(ctx context.Context, request entity.CreateProductRequest) (err error) {
	if err := s.validator.Validate(request); err != nil {
		err = helper.Error(http.StatusBadRequest, err.Error(), err)
		return err
	}

	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return err
	}

	productId, err := s.repository.CreateProduct(ctx, tx, request)
	if err != nil {
		return err
	}

	request.ProductId = productId
	err = s.repository.CreateMultipleSku(ctx, tx, request)
	if err != nil {
		s.repository.RollbackTx(ctx, tx)
		return err
	}

	return s.repository.CommitTx(ctx, tx)

}
