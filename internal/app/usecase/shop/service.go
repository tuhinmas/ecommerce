package shop

import (
	"context"
	"ecommerce/internal/app/repository/shop"
	"ecommerce/internal/entity"
	"ecommerce/pkg/helper"
	"net/http"

	"ecommerce/pkg/constant"
	"ecommerce/pkg/validator"
)

type ShopService interface {
	CreateShop(ctx context.Context, shop entity.CreateShopRequest) (err error)
}

type shopService struct {
	shopRepository shop.ShopRepository
	validator      validator.Validator
}

func NewShopService(shopRepository shop.ShopRepository, validator validator.Validator) ShopService {
	return &shopService{
		shopRepository: shopRepository,
		validator:      validator,
	}
}

func (s *shopService) CreateShop(ctx context.Context, shop entity.CreateShopRequest) (err error) {
	if err = s.validator.Validate(shop); err != nil {
		err = helper.Error(http.StatusBadRequest, constant.MsgInvalidRequest, err)
		return
	}

	return s.shopRepository.CreateShop(ctx, shop)
}
