package product

import (
	"context"
	"database/sql"
	"ecommerce/internal/entity"
	"ecommerce/pkg/identifier"
	"ecommerce/pkg/validator"
	"testing"

	mock_repository "ecommerce/internal/app/repository/product/mocks"

	validatorv10 "github.com/go-playground/validator/v10"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock_repository.NewMockProductRepository(ctrl)
	validator := validator.NewValidator(validatorv10.New())
	identifier := identifier.NewIdentifier()
	svc := NewProductService(mockRepository, validator, identifier)

	ctx := context.Background()

	t.Run("success case", func(t *testing.T) {
		request := entity.CreateProductRequest{
			ProductId: "1",
			Name:      "Test Product",
			ShopId:    "1",
			Sku: []entity.CreateSkuRequest{
				{
					Variant: "Red",
					Price:   100000,
					Uom:     "pcs",
					Image:   "image.jpg",
				},
			},
		}
		tx := &sql.Tx{} // Mock transaction object
		mockRepository.EXPECT().BeginTx(ctx).Return(tx, nil)
		mockRepository.EXPECT().CreateProduct(ctx, tx, request).Return("1", nil)
		mockRepository.EXPECT().CreateMultipleSku(ctx, tx, request).Return(nil)
		mockRepository.EXPECT().CommitTx(ctx, tx).Return(nil)

		err := svc.CreateProduct(ctx, request)
		assert.NoError(t, err)
	})
}

func TestProductList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock_repository.NewMockProductRepository(ctrl)
	validator := validator.NewValidator(validatorv10.New())
	identifier := identifier.NewIdentifier()
	svc := NewProductService(mockRepository, validator, identifier)

	ctx := context.Background()

	t.Run("success case", func(t *testing.T) {
		request := entity.QueryRequest{
			Page:   1,
			Limit:  10,
			SortBy: "created_at",
		}
		mockRepository.EXPECT().GetProduct(ctx, request).Return([]*entity.GetProductListResponse{}, nil)
		mockRepository.EXPECT().GetSku(ctx, gomock.Any()).Return([]*entity.ProductDetailResponse{
			{
				Id:      "1",
				Image:   "image.jpg",
				Uom:     "pcs",
				Price:   100000,
				Stock:   100,
				Variant: "Red",
			},
		}, nil)
		mockRepository.EXPECT().GetTotalProduct(ctx).Return(10, nil)

		resp, pagination, err := svc.ProductList(ctx, request)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, pagination)
	})
}
