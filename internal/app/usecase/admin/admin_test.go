package admin

import (
	"context"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"ecommerce/pkg/identifier"
	"ecommerce/pkg/validator"
	"net/http"
	"testing"

	mock_repository "ecommerce/internal/app/repository/admin/mocks"

	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock_repository.NewMockAdminRepository(ctrl)
	validator := validator.NewValidator(validatorv10.New())
	identifier := identifier.NewIdentifier()
	svc := NewAdminService(mockRepository, validator, identifier)

	ctx := context.Background()

	t.Run("Validation Error", func(t *testing.T) {
		request := &entity.LoginAdminRequest{} // Missing required fields

		resp, err := svc.Login(ctx, request)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
	})

	t.Run("Username Not Found", func(t *testing.T) {
		request := &entity.LoginAdminRequest{
			Username: "admin",
			Password: "password123",
		}

		// Mock repository call to simulate user not found
		mockRepository.EXPECT().GetAdminByUsername(ctx, request.Username).Return(entity.GetAdminDetailResponse{}, helper.Error(http.StatusNotFound, "User not found", nil))

		resp, err := svc.Login(ctx, request)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), constant.ErrorEmailNotFound)
	})

	t.Run("Incorrect Password", func(t *testing.T) {
		request := &entity.LoginAdminRequest{
			Username: "admin",
			Password: "wrongpassword",
		}

		user := &entity.GetAdminDetailResponse{
			Id:       "1",
			Username: request.Username,
			Password: helper.EncryptPassword("correctpassword"),
		}

		// Mock repository call to return the user
		mockRepository.EXPECT().GetAdminByUsername(ctx, request.Username).Return(*user, nil)

		resp, err := svc.Login(ctx, request)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), constant.ErrorPasswordWrong)
	})

	t.Run("Successful Login", func(t *testing.T) {
		request := &entity.LoginAdminRequest{
			Username: "admin",
			Password: "correctpassword",
		}

		user := &entity.GetAdminDetailResponse{
			Id:       "1",
			Username: request.Username,
			Password: helper.EncryptPassword(request.Password),
		}

		mockRepository.EXPECT().GetAdminByUsername(ctx, request.Username).Return(*user, nil)

		resp, err := svc.Login(ctx, request)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.ExpiredAt)
	})
}
