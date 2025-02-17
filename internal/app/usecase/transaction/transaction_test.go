package transaction

import (
	"context"
	"database/sql"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/identifier"
	"ecommerce/pkg/validator"
	"sync"
	"testing"

	mock_repository "ecommerce/internal/app/repository/transaction/mocks"
	mock_worker "ecommerce/internal/app/worker/mocks"

	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPlaceOrder(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.HeaderContext, entity.ValueContext{
		UserId:      "1",
		WarehouseId: "1",
	})

	identifier := identifier.NewIdentifier()
	validator := validator.NewValidator(validatorv10.New())
	mockWorker := mock_worker.NewMockQueueService(ctl)

	mockRepository := mock_repository.NewMockTransactionRepository(ctl)
	feature := NewTransactionService(mockRepository, validator, identifier, mockWorker) // Adjust as needed

	t.Run("successful place order", func(t *testing.T) {
		request := &entity.OrderRequest{
			OrderId:       "1",
			PaymentMethod: "1",
			Address:       "Jl. Imam Bonjol",
			Amount:        10000,
			WarehouseId:   "1",
			Sku:           []entity.SkuRequest{{Id: "sku1", Quantity: 1, Price: 10000}},
		}

		tx := &sql.Tx{}
		var wg sync.WaitGroup
		wg.Add(1)

		// Set up expected calls and return values
		mockRepository.EXPECT().BeginTx(gomock.Any()).Return(tx, nil)
		mockRepository.EXPECT().GetMultipleSku(gomock.Any(), gomock.Any()).Return([]entity.WarehouseStock{{SkuId: "sku1", Stock: 2, Price: 10000}}, nil)
		mockRepository.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), *request).Return("1", nil)
		mockRepository.EXPECT().CreateOrderItem(gomock.Any(), gomock.Any(), *request).Return(nil)
		mockRepository.EXPECT().UpdateStock(gomock.Any(), gomock.Any(), request.Sku[0]).Return(nil)
		mockRepository.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return(nil)
		mockWorker.EXPECT().PublishData(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, queue string, payload interface{}) error {
			defer wg.Done() // Ensure we track when PublishData is actually called
			return nil
		})

		// Execute the test
		err := feature.PlaceOrder(ctx, request)
		assert.Nil(t, err)
		// Wait for goroutine to finish
		wg.Wait()
	})

}
