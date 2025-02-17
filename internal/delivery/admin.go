package delivery

import (
	"ecommerce/internal/app/usecase/admin"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler interface {
	Login(ctx *fiber.Ctx) (err error)
}

type adminHandler struct {
	service admin.Service
}

func NewAdminHandler(service admin.Service) AdminHandler {
	return &adminHandler{service}
}

func (h *adminHandler) Login(c *fiber.Ctx) (err error) {
	ctx, cancel := helper.CreateContextWithTimeout()
	defer cancel()
	ctx = helper.SetValueToContext(ctx, c)

	request := new(entity.LoginAdminRequest)
	if err := c.BodyParser(request); err != nil {
		err = helper.Error(http.StatusBadRequest, constant.MsgInvalidRequest, err)
		return helper.ResponseError(ctx, err)
	}

	resp, err := h.service.Login(ctx, request)
	if err != nil {
		return helper.ResponseError(ctx, err)
	}
	return helper.ResponseCreatedOK(c, "success", resp)
}
