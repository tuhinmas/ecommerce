package shop

import (
	"context"
	"ecommerce/pkg/validator"
	"testing"

	mock_repository "ecommerce/internal/app/repository/shop/mocks"
	"ecommerce/internal/entity"

	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestCreateShop(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockShopRepository := mock_repository.NewMockShopRepository(ctl)
	validator := validator.NewValidator(validatorv10.New())
	shopService := NewShopService(mockShopRepository, validator)
	ctx := context.Background()

	t.Run("success create shop", func(t *testing.T) {
		mockShopRepository.EXPECT().CreateShop(gomock.Any(), gomock.Any()).Return(nil)
		err := shopService.CreateShop(ctx, entity.CreateShopRequest{
			Name: "Shop 1",
		})
		assert.Nil(t, err)
	})
}
