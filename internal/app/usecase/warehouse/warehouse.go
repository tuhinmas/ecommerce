package warehouse

import (
	"context"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"errors"
	"net/http"
)

func (s *service) SetStatusWarehouse(ctx context.Context, request entity.SetStatusWarehouseRequest) (err error) {
	return s.warehouseRepository.SetStatusWarehouse(ctx, request)
}

func (s *service) CreateStock(ctx context.Context, request entity.CreateStockRequest) (err error) {

	if err := s.validator.Validate(request); err != nil {
		err = helper.Error(http.StatusBadRequest, err.Error(), err)
		return err
	}

	isExist, err := s.warehouseRepository.IsExistStockByWarehouseIdAndSkuId(ctx, request.WarehouseId, request.SkuId)
	if err != nil {
		return err
	}
	if isExist {
		return helper.Error(http.StatusBadRequest, constant.MsgStockExist, errors.New(constant.MsgStockExist))
	}
	return s.warehouseRepository.CreateStock(ctx, request)
}

func (s *service) UpdateStock(ctx context.Context, id string, request entity.UpdateStockRequest) (err error) {
	isExist, err := s.warehouseRepository.GetStockById(ctx, id)
	if err != nil {
		return err
	}
	if !isExist {
		return helper.Error(http.StatusBadRequest, constant.MsgStockNotFound, errors.New(constant.MsgStockNotFound))
	}
	return s.warehouseRepository.UpdateStock(ctx, id, request)
}

func (s *service) TransferStock(ctx context.Context, request entity.StockTransferRequest) (err error) {
	tx, err := s.warehouseRepository.BeginTx(ctx)
	if err != nil {
		return err
	}

	stock, err := s.warehouseRepository.GetStockByWarehouseIdAndSkuId(ctx, request.From, request.SkuId)
	if err != nil {
		return err
	}

	if stock == 0 {
		return helper.Error(http.StatusBadRequest, constant.MsgStockNotFound, errors.New(constant.MsgStockNotFound))
	}

	if stock < request.Quantity {
		return helper.Error(http.StatusBadRequest, constant.MsgStockNotEnough, errors.New(constant.MsgStockNotEnough))
	}

	exist, err := s.warehouseRepository.IsExistStockByWarehouseIdAndSkuId(ctx, request.To, request.SkuId)
	if err != nil {
		return err
	}
	if !exist {
		return helper.Error(http.StatusBadRequest, constant.MsgStockNotFound, errors.New(constant.MsgStockNotFound))
	}

	err = s.warehouseRepository.DecreaseStock(ctx, tx, request)
	if err != nil {
		return err
	}

	err = s.warehouseRepository.IncreaseStock(ctx, tx, request)
	if err != nil {
		s.warehouseRepository.RollbackTx(ctx, tx)
		return err
	}

	err = s.warehouseRepository.CreateStockTransfer(ctx, tx, request)
	if err != nil {
		s.warehouseRepository.RollbackTx(ctx, tx)
		return err
	}

	return s.warehouseRepository.CommitTx(ctx, tx)
}

func (s *service) CreateWarehouse(ctx context.Context, request entity.CreateWarehouseRequest) (err error) {
	if err := s.validator.Validate(request); err != nil {
		err = helper.Error(http.StatusBadRequest, err.Error(), err)
		return err
	}

	isExist, err := s.warehouseRepository.IsExistShopId(ctx, request.ShopId)
	if err != nil {
		return err
	}

	if !isExist {
		return helper.Error(http.StatusBadRequest, "Shop not found", errors.New("shop not found"))
	}

	return s.warehouseRepository.CreateWarehouse(ctx, request)
}
