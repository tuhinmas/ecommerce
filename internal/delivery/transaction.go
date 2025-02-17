package delivery

import (
	"ecommerce/internal/app/usecase/transaction"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler interface {
	Order(ctx *fiber.Ctx) (err error)
}

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) TransactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) Order(c *fiber.Ctx) error {
	ctx, cancel := helper.CreateContextWithTimeout()
	defer cancel()
	ctx = helper.SetValueToContext(ctx, c)

	request := new(entity.OrderRequest)
	if err := c.BodyParser(request); err != nil {
		err = helper.Error(http.StatusBadRequest, constant.MsgInvalidRequest, err)
		return helper.ResponseError(ctx, err)
	}

	err := h.service.PlaceOrder(ctx, request)
	if err != nil {
		fmt.Println(err)
		return helper.ResponseError(ctx, err)
	}
	return helper.ResponseCreatedOK(c, "success", nil)
}
