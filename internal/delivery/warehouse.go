package delivery

import (
	"ecommerce/internal/app/usecase/warehouse"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type WarehouseHandler interface {
	SetStatusWarehouse(ctx *fiber.Ctx) (err error)
	CreateStock(ctx *fiber.Ctx) (err error)
	UpdateStock(ctx *fiber.Ctx) (err error)
	TransferStock(ctx *fiber.Ctx) (err error)
	CreateWarehouse(ctx *fiber.Ctx) (err error)
}

type warehouseHandler struct {
	service warehouse.Service
}

func NewWarehouseHandler(service warehouse.Service) WarehouseHandler {
	return &warehouseHandler{
		service: service,
	}
}

func (h *warehouseHandler) SetStatusWarehouse(c *fiber.Ctx) (err error) {
	ctx, cancel := helper.CreateContextWithTimeout()
	defer cancel()
	ctx = helper.SetValueToContext(ctx, c)

	request := entity.SetStatusWarehouseRequest{}
	if err = c.BodyParser(&request); err != nil {
		err = helper.Error(http.StatusBadRequest, err.Error(), err)
		return helper.ResponseError(ctx, err)
	}

	if err = h.service.SetStatusWarehouse(ctx, request); err != nil {
		return helper.ResponseError(ctx, err)
	}

	return helper.ResponseOK(c, constant.MsgWarehouseStatusUpdated, nil)
}

func (h *warehouseHandler) CreateStock(c *fiber.Ctx) (err error) {
	ctx, cancel := helper.CreateContextWithTimeout()
	defer cancel()
	ctx = helper.SetValueToContext(ctx, c)

	request := entity.CreateStockRequest{}
	if err = c.BodyParser(&request); err != nil {
		err = helper.Error(http.StatusBadRequest, err.Error(), err)
		return helper.ResponseError(ctx, err)
	}

	if err = h.service.CreateStock(ctx, request); err != nil {
		return helper.ResponseError(ctx, err)
	}

	return helper.ResponseOK(c, constant.MsgStockCreated, nil)
}

func (h *warehouseHandler) UpdateStock(c *fiber.Ctx) (err error) {
	ctx, cancel := helper.CreateContextWithTimeout()
	defer cancel()
	ctx = helper.SetValueToContext(ctx, c)

	id := c.Params("id")
	if id == "" {
		return helper.Error(http.StatusBadRequest, constant.MsgInvalidRequest, errors.New(constant.MsgInvalidRequest))
	}

	request := entity.UpdateStockRequest{}
	if err = c.BodyParser(&request); err != nil {
		err = helper.Error(http.StatusBadRequest, err.Error(), err)
		return helper.ResponseError(ctx, err)
	}

	if err = h.service.UpdateStock(ctx, id, request); err != nil {
		return helper.ResponseError(ctx, err)
	}

	return helper.ResponseOK(c, constant.MsgStockUpdated, nil)
}

func (h *warehouseHandler) TransferStock(c *fiber.Ctx) (err error) {
	ctx, cancel := helper.CreateContextWithTimeout()
	defer cancel()
	ctx = helper.SetValueToContext(ctx, c)

	request := entity.StockTransferRequest{}
	if err = c.BodyParser(&request); err != nil {
		err = helper.Error(http.StatusBadRequest, err.Error(), err)
		return helper.ResponseError(ctx, err)
	}

	if err = h.service.TransferStock(ctx, request); err != nil {
		return helper.ResponseError(ctx, err)
	}

	return helper.ResponseOK(c, constant.MsgStockTransferred, nil)
}

func (h *warehouseHandler) CreateWarehouse(c *fiber.Ctx) (err error) {
	ctx, cancel := helper.CreateContextWithTimeout()
	defer cancel()
	ctx = helper.SetValueToContext(ctx, c)

	request := entity.CreateWarehouseRequest{}
	if err = c.BodyParser(&request); err != nil {
		err = helper.Error(http.StatusBadRequest, err.Error(), err)
		return helper.ResponseError(ctx, err)
	}

	if err = h.service.CreateWarehouse(ctx, request); err != nil {
		return helper.ResponseError(ctx, err)
	}

	return helper.ResponseCreatedOK(c, "Warehouse created successfully", nil)
}
