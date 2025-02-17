package warehouse

import (
	"context"
	"testing"

	mock_repository "ecommerce/internal/app/repository/warehouse/mocks"
	"ecommerce/internal/entity"
	"ecommerce/pkg/validator"

	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateWarehouse(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWarehouseRepository := mock_repository.NewMockWarehouseRepository(ctl)
	validator := validator.NewValidator(validatorv10.New())
	warehouseService := NewWarehouseService(mockWarehouseRepository, validator)
	ctx := context.Background()

	t.Run("success create warehouse", func(t *testing.T) {
		mockWarehouseRepository.EXPECT().IsExistShopId(gomock.Any(), gomock.Any()).Return(true, nil)
		mockWarehouseRepository.EXPECT().CreateWarehouse(gomock.Any(), gomock.Any()).Return(nil)
		err := warehouseService.CreateWarehouse(ctx, entity.CreateWarehouseRequest{
			Location: "Jl. Imam Bonjol",
			Address:  "Jl. Imam Bonjol",
			ShopId:   "1",
		})
		assert.Nil(t, err)
	})
}

func TestSetStatusWarehouse(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWarehouseRepository := mock_repository.NewMockWarehouseRepository(ctl)
	validator := validator.NewValidator(validatorv10.New())
	warehouseService := NewWarehouseService(mockWarehouseRepository, validator)
	ctx := context.Background()

	t.Run("success set status warehouse", func(t *testing.T) {
		mockWarehouseRepository.EXPECT().SetStatusWarehouse(gomock.Any(), gomock.Any()).Return(nil)
		err := warehouseService.SetStatusWarehouse(ctx, entity.SetStatusWarehouseRequest{
			WarehouseId: "1",
		})
		assert.Nil(t, err)
	})
}

func TestCreateStock(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWarehouseRepository := mock_repository.NewMockWarehouseRepository(ctl)
	validator := validator.NewValidator(validatorv10.New())
	warehouseService := NewWarehouseService(mockWarehouseRepository, validator)
	ctx := context.Background()

	t.Run("success create stock", func(t *testing.T) {
		mockWarehouseRepository.EXPECT().IsExistStockByWarehouseIdAndSkuId(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
		mockWarehouseRepository.EXPECT().CreateStock(gomock.Any(), gomock.Any()).Return(nil)
		err := warehouseService.CreateStock(ctx, entity.CreateStockRequest{
			WarehouseId: "1",
			SkuId:       "1",
			Stock:       10,
		})
		assert.Nil(t, err)
	})
}

func TestUpdateStock(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWarehouseRepository := mock_repository.NewMockWarehouseRepository(ctl)
	validator := validator.NewValidator(validatorv10.New())
	warehouseService := NewWarehouseService(mockWarehouseRepository, validator)
	ctx := context.Background()

	t.Run("success update stock", func(t *testing.T) {
		mockWarehouseRepository.EXPECT().GetStockById(gomock.Any(), gomock.Any()).Return(true, nil)
		mockWarehouseRepository.EXPECT().UpdateStock(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		err := warehouseService.UpdateStock(ctx, "1", entity.UpdateStockRequest{
			Stock: 10,
		})
		assert.Nil(t, err)
	})
}

func TestTransferStock(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWarehouseRepository := mock_repository.NewMockWarehouseRepository(ctl)
	validator := validator.NewValidator(validatorv10.New())
	warehouseService := NewWarehouseService(mockWarehouseRepository, validator)
	ctx := context.Background()

	t.Run("success transfer stock", func(t *testing.T) {
		mockWarehouseRepository.EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
		mockWarehouseRepository.EXPECT().GetStockByWarehouseIdAndSkuId(gomock.Any(), gomock.Any(), gomock.Any()).Return(10, nil)
		mockWarehouseRepository.EXPECT().IsExistStockByWarehouseIdAndSkuId(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
		mockWarehouseRepository.EXPECT().IncreaseStock(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockWarehouseRepository.EXPECT().DecreaseStock(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockWarehouseRepository.EXPECT().CreateStockTransfer(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockWarehouseRepository.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return(nil)
		err := warehouseService.TransferStock(ctx, entity.StockTransferRequest{
			From:     "1",
			To:       "2",
			SkuId:    "1",
			Quantity: 10,
		})
		assert.Nil(t, err)
	})
}
