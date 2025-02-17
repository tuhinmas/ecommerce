package user

import (
	"context"
	"net/http"
	"testing"

	mock_repository "ecommerce/internal/app/repository/user/mocks"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/identifier"
	"ecommerce/pkg/validator"

	"ecommerce/pkg/helper"

	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock_repository.NewMockUserRepository(ctrl)
	validator := validator.NewValidator(validatorv10.New())
	identifier := identifier.NewIdentifier()
	svc := NewUserService(mockRepository, validator, identifier)

	ctx := context.Background()

	t.Run("Validation Error", func(t *testing.T) {
		request := &entity.LoginRequest{} // Missing required fields

		resp, err := svc.Login(ctx, request)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
	})

	t.Run("Email Not Found", func(t *testing.T) {
		request := &entity.LoginRequest{
			Phone:    "081234567890",
			Password: "password123",
		}

		// Mock repository call to simulate user not found
		mockRepository.EXPECT().GetUserByPhone(ctx, request.Phone).Return(entity.GetUserDetailResponse{}, helper.Error(http.StatusNotFound, "User not found", nil))

		resp, err := svc.Login(ctx, request)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), constant.ErrorEmailNotFound)
	})

	t.Run("Incorrect Password", func(t *testing.T) {
		request := &entity.LoginRequest{
			Phone:    "081234567890",
			Password: "wrongpassword",
		}

		user := entity.GetUserDetailResponse{
			Id:       "1",
			Phone:    request.Phone,
			Password: helper.EncryptPassword("correctpassword"),
		}

		// Mock repository call to return the user
		mockRepository.EXPECT().GetUserByPhone(ctx, request.Phone).Return(user, nil)

		resp, err := svc.Login(ctx, request)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), constant.ErrorPasswordWrong)
	})

	t.Run("Successful Login", func(t *testing.T) {
		request := &entity.LoginRequest{
			Phone:    "081234567890",
			Password: "correctpassword",
		}

		user := entity.GetUserDetailResponse{
			Id:       "1",
			Phone:    request.Phone,
			Password: helper.EncryptPassword(request.Password),
			Gender:   "male",
		}

		mockRepository.EXPECT().GetUserByPhone(ctx, request.Phone).Return(user, nil)

		resp, err := svc.Login(ctx, request)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.ExpiredAt)
	})
}

func TestSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockRepository := mock_repository.NewMockUserRepository(ctrl)
	validator := validator.NewValidator(validatorv10.New())
	svc := NewUserService(mockRepository, validator, nil)

	t.Run("Validation Error", func(t *testing.T) {
		request := &entity.SignupRequest{}
		err := svc.Signup(ctx, request)
		assert.NotNil(t, err)
	})

	t.Run("Phone Already Exists", func(t *testing.T) {
		request := &entity.SignupRequest{
			WarehouseId: "1",
			Name:        "John Doe",
			Phone:       "081234567890",
			Password:    "password123",
			Gender:      "male",
		}

		mockRepository.EXPECT().GetWarehouseById(ctx, "1").Return(true, nil)
		mockRepository.EXPECT().GetUserByPhone(ctx, request.Phone).Return(entity.GetUserDetailResponse{
			Id:       "1",
			Phone:    request.Phone,
			Password: helper.EncryptPassword("password123"),
			Gender:   "male",
		}, nil)

		err := svc.Signup(ctx, request)
		assert.NotNil(t, err)
	})

	t.Run("Repository Error on GetUserByPhone", func(t *testing.T) {
		request := &entity.SignupRequest{
			WarehouseId: "1",
			Name:        "John Doe",
			Phone:       "081234567890",
			Password:    "password123",
			Gender:      "male",
		}

		mockRepository.EXPECT().GetWarehouseById(ctx, "1").Return(true, nil)
		mockRepository.EXPECT().GetUserByPhone(ctx, "081234567890").Return(entity.GetUserDetailResponse{}, helper.Error(http.StatusInternalServerError, "database error", nil))

		err := svc.Signup(ctx, request)
		assert.NotNil(t, err)
	})

	t.Run("Successful Signup", func(t *testing.T) {
		request := &entity.SignupRequest{
			WarehouseId: "1",
			Name:        "John Doe",
			Phone:       "081234567890",
			Password:    "password123",
			Gender:      "male",
		}

		mockRepository.EXPECT().GetWarehouseById(ctx, "1").Return(true, nil)
		mockRepository.EXPECT().GetUserByPhone(ctx, request.Phone).Return(entity.GetUserDetailResponse{}, helper.Error(http.StatusNotFound, "User not found", nil))
		mockRepository.EXPECT().Signup(ctx, gomock.Any()).Return(nil)

		err := svc.Signup(ctx, request)
		assert.Nil(t, err)
	})
}
