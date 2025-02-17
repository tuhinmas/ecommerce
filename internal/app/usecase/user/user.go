package user

import (
	"context"
	"ecommerce/cmd/middleware"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"errors"
	"net/http"
	"time"
)

func (s *service) Signup(ctx context.Context, request *entity.SignupRequest) (err error) {

	if err := s.validator.Validate(request); err != nil {
		err = helper.Error(http.StatusBadRequest, err.Error(), err)
		return err
	}

	warehouse, err := s.repository.GetWarehouseById(ctx, request.WarehouseId)
	if err != nil {
		return err
	}

	if !warehouse {
		return helper.Error(http.StatusBadRequest, constant.MsgWarehouseNotFound, nil)
	}

	if request.Password != "" {
		request.Password = helper.EncryptPassword(request.Password)
	}

	check, err := s.repository.GetUserByPhone(ctx, request.Phone)
	if err != nil {
		statusCode, _, _ := helper.TrimMesssage(err)
		if statusCode != http.StatusNotFound {
			return err
		}
	}

	if check.Phone != "" {
		err = helper.Error(http.StatusConflict, constant.ErrorPhoneAlreadyExists, errors.New(constant.ErrorPhoneAlreadyExists))
		return
	}

	err = s.repository.Signup(ctx, *request)
	if err != nil {
		return
	}

	return nil
}

func (s *service) Login(ctx context.Context, request *entity.LoginRequest) (*entity.LoginResponse, error) {

	if err := s.validator.Validate(request); err != nil {
		err = helper.Error(http.StatusBadRequest, err.Error(), err)
		return nil, err
	}

	user, err := s.repository.GetUserByPhone(ctx, request.Phone)
	if err != nil {
		statusCode, _, _ := helper.TrimMesssage(err)
		if statusCode == http.StatusNotFound {
			err = helper.Error(http.StatusUnauthorized, constant.ErrorEmailNotFound, err)
			return nil, err
		}
		return nil, err
	}

	if user.Password == "" {
		err = helper.Error(http.StatusUnauthorized, constant.ErrorEmailNotFound, nil)
		return nil, err
	}

	if user.Password != helper.EncryptPassword(request.Password) {
		err = helper.Error(http.StatusUnauthorized, constant.ErrorPasswordWrong, nil)
		return nil, err
	}

	expirationTime := time.Now().Add(time.Hour * time.Duration(24))
	token, err := middleware.GenerateToken(user.Id, user.WarehouseId)
	if err != nil {
		return nil, err
	}

	resp := entity.LoginResponse{
		Token:     token,
		ExpiredAt: expirationTime.Format(time.RFC3339),
	}

	return &resp, nil
}
