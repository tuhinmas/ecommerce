package warehouse

import (
	"context"
	"ecommerce/internal/app/repository/warehouse"
	"ecommerce/internal/entity"
	"ecommerce/pkg/validator"
)

type Service interface {
	SetStatusWarehouse(ctx context.Context, request entity.SetStatusWarehouseRequest) (err error)
	CreateStock(ctx context.Context, request entity.CreateStockRequest) (err error)
	UpdateStock(ctx context.Context, id string, request entity.UpdateStockRequest) (err error)
	TransferStock(ctx context.Context, request entity.StockTransferRequest) (err error)
	CreateWarehouse(ctx context.Context, request entity.CreateWarehouseRequest) (err error)
}

type service struct {
	warehouseRepository warehouse.WarehouseRepository
	validator           validator.Validator
}

func NewWarehouseService(
	warehouseRepository warehouse.WarehouseRepository,
	validator validator.Validator,
) Service {
	return &service{
		warehouseRepository: warehouseRepository,
		validator:           validator,
	}
}
