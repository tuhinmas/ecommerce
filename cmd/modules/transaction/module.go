package transaction

import (
	"context"
	"ecommerce/cmd/modules"
	"ecommerce/cmd/routes"
	repository "ecommerce/internal/app/repository/transaction"
	"ecommerce/internal/app/usecase/transaction"
	"ecommerce/internal/app/worker"
	handler "ecommerce/internal/delivery"
)

func InitModule(container modules.Container) {
	repository := repository.NewTransactionRepository(container.Db)
	queueService := worker.NewQueueService(container.Rmq, container.WorkerConfig, repository)

	transactionService := transaction.NewTransactionService(repository, container.Validator, container.Identifier, queueService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	routes.TransactionRouter(container.App, transactionHandler)

	go queueService.ConsumeData(context.Background(), container.WorkerConfig.StockReversalQueue)
}
