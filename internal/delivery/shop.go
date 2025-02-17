package delivery

import (
	"ecommerce/internal/app/usecase/shop"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ShopHandler interface {
	CreateShop(ctx *fiber.Ctx) error
}

type shopHandler struct {
	shopService shop.ShopService
}

func NewShopHandler(shopService shop.ShopService) ShopHandler {
	return &shopHandler{
		shopService: shopService,
	}
}

func (h *shopHandler) CreateShop(c *fiber.Ctx) error {
	ctx, cancel := helper.CreateContextWithTimeout()
	defer cancel()
	ctx = helper.SetValueToContext(ctx, c)

	shop := entity.CreateShopRequest{}
	if err := c.BodyParser(&shop); err != nil {
		err = helper.Error(http.StatusBadRequest, constant.MsgInvalidRequest, err)
		return helper.ResponseError(ctx, err)
	}

	if err := h.shopService.CreateShop(ctx, shop); err != nil {
		return helper.ResponseError(ctx, err)
	}

	return helper.ResponseCreatedOK(c, "Shop created successfully", nil)
}
