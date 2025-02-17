package delivery

import (
	"ecommerce/internal/app/usecase/user"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	Signup(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
}

type userHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) UserHandler {
	return &userHandler{service}
}

func (h *userHandler) Signup(c *fiber.Ctx) error {
	ctx, cancel := helper.CreateContextWithTimeout()
	defer cancel()
	ctx = helper.SetValueToContext(ctx, c)

	request := new(entity.SignupRequest)
	if err := c.BodyParser(request); err != nil {
		err = helper.Error(http.StatusBadRequest, constant.MsgInvalidRequest, err)
		return helper.ResponseError(ctx, err)
	}

	err := h.service.Signup(ctx, request)
	if err != nil {
		return helper.ResponseError(ctx, err)
	}
	return helper.ResponseCreatedOK(c, "success", nil)
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	ctx, cancel := helper.CreateContextWithTimeout()
	defer cancel()
	ctx = helper.SetValueToContext(ctx, c)

	request := new(entity.LoginRequest)
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
