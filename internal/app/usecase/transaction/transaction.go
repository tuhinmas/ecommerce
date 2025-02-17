package transaction

import (
	"context"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"
	"net/http"
	"os"
	"sync"
)

func (s *service) PlaceOrder(ctx context.Context, request *entity.OrderRequest) (err error) {

	if err := s.validator.Validate(request); err != nil {
		err = helper.Error(http.StatusBadRequest, err.Error(), err)
		return err
	}

	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return err
	}

	skus, err := s.repository.GetMultipleSku(ctx, request.Sku)
	if err != nil {
		return err
	}

	// Map sku to skuMap
	skuMap := make(map[string]entity.WarehouseStock)
	for _, sku := range skus {
		skuMap[sku.SkuId] = sku
	}

	// Validate sku and calculate amount
	amount := 0
	for i, sku := range request.Sku {
		value, ok := skuMap[sku.Id]
		// If sku not found, return error
		if !ok {
			err = helper.Error(http.StatusBadRequest, constant.MsgSkuNotFound, nil)
			return err
		}

		// If stock is not enough, return error
		if value.Stock < sku.Quantity {
			err = helper.Error(http.StatusBadRequest, constant.MsgStockNotEnough, nil)
			return err
		}

		request.Sku[i].Price = value.Price
		amount += value.Price * sku.Quantity
	}

	request.Amount = amount
	id, err := s.repository.CreateOrder(ctx, tx, *request)
	if err != nil {
		return err
	}

	request.OrderId = id
	err = s.repository.CreateOrderItem(ctx, tx, *request)
	if err != nil {
		tx.Rollback()
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(skus)) // Buffered channel for errors

	// update sku stock by concurrent
	for _, sku := range request.Sku {
		wg.Add(1)
		go func(sku entity.SkuRequest) {
			defer wg.Done()

			if err := s.repository.UpdateStock(ctx, tx, sku); err != nil {
				errCh <- err
			}
		}(sku)
	}

	wg.Wait()
	close(errCh)

	// Check if any error occurred
	for err := range errCh {
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := s.repository.CommitTx(ctx, tx); err != nil {
		return err
	}

	valueCtx := helper.GetValueContext(ctx)
	payloadQueue := entity.PayloadOrderQueue{
		OrderId:     request.OrderId,
		WarehouseId: valueCtx.WarehouseId,
		Sku:         request.Sku,
	}

	go s.worker.PublishData(ctx, os.Getenv("WORKER_PENDING_PAYMENT_QUEUE"), payloadQueue)

	return nil
}
