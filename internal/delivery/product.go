package delivery

import (
	"ecommerce/internal/app/usecase/product"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler interface {
	ProductList(ctx *fiber.Ctx) (err error)
	CreateProduct(ctx *fiber.Ctx) (err error)
}

type productHandler struct {
	service product.Service
}

func NewProductHandler(service product.Service) ProductHandler {
	return &productHandler{service}
}

func (h *productHandler) ProductList(c *fiber.Ctx) (err error) {
	ctx, cancel := helper.CreateContextWithTimeout()
	defer cancel()
	ctx = helper.SetValueToContext(ctx, c)

	page, err := strconv.Atoi(strings.TrimSpace(c.Query(constant.PAGE)))
	if err != nil || page == 0 {
		page = constant.DefaultPage
	}

	limit, err := strconv.Atoi(strings.TrimSpace(c.Query(constant.LIMIT)))
	if err != nil || limit == 0 {
		limit = constant.DefaultLimitPerPage
	}

	filter := entity.QueryRequest{
		Page:  page,
		Limit: limit,
	}

	resp, pagination, err := h.service.ProductList(ctx, filter)
	if err != nil {
		return helper.ResponseError(ctx, err)
	}

	return helper.ResponseOkWithPagination(c, "success", resp, pagination)
}

func (h *productHandler) CreateProduct(c *fiber.Ctx) (err error) {
	ctx, cancel := helper.CreateContextWithTimeout()
	defer cancel()
	ctx = helper.SetValueToContext(ctx, c)

	request := entity.CreateProductRequest{}
	if err := c.BodyParser(&request); err != nil {
		err = helper.Error(http.StatusBadRequest, err.Error(), err)
		return helper.ResponseError(ctx, err)
	}

	err = h.service.CreateProduct(ctx, request)
	if err != nil {
		return helper.ResponseError(ctx, err)
	}

	return helper.ResponseCreatedOK(c, "success", nil)
}
