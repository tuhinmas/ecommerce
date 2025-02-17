package transaction

import (
	"context"
	repository "ecommerce/internal/app/repository/transaction"
	worker "ecommerce/internal/app/worker"
	"ecommerce/internal/entity"
	"ecommerce/pkg/identifier"
	"ecommerce/pkg/validator"
)

type Service interface {
	PlaceOrder(ctx context.Context, request *entity.OrderRequest) (err error)
}

type service struct {
	repository repository.TransactionRepository
	validator  validator.Validator
	identifier identifier.Identifier
	worker     worker.QueueService
}

func NewTransactionService(
	repository repository.TransactionRepository,
	validator validator.Validator,
	identifier identifier.Identifier,
	worker worker.QueueService,
) Service {
	return &service{
		repository: repository,
		validator:  validator,
		identifier: identifier,
		worker:     worker,
	}
}
