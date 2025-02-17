package admin

import (
	"context"
	"ecommerce/cmd/middleware"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"net/http"
	"time"
)

func (s *service) Login(ctx context.Context, request *entity.LoginAdminRequest) (*entity.LoginResponse, error) {

	if err := s.validator.Validate(request); err != nil {
		err = helper.Error(http.StatusBadRequest, err.Error(), err)
		return nil, err
	}

	user, err := s.repository.GetAdminByUsername(ctx, request.Username)
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
	token, err := middleware.GenerateTokenAdmin(user.Id)
	if err != nil {
		return nil, err
	}

	resp := entity.LoginResponse{
		Token:     token,
		ExpiredAt: expirationTime.Format(time.RFC3339),
	}

	return &resp, nil
}
