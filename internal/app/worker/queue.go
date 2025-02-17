package worker

import (
	"context"
	"ecommerce/broker"
	transactionRepository "ecommerce/internal/app/repository/transaction"
)

type QueueService interface {
	PublishData(ctx context.Context, topic string, msg interface{}) (err error)
	ConsumeData(ctx context.Context, topic string) (err error)
}

type queueService struct {
	rabbitmq              broker.RabbitMQ
	cfg                   broker.RabbitmqConfig
	transactionRepository transactionRepository.TransactionRepository
}

func NewQueueService(
	rabbitmq broker.RabbitMQ,
	cfg broker.RabbitmqConfig,
	transactionRepository transactionRepository.TransactionRepository,
) QueueService {
	return &queueService{
		rabbitmq:              rabbitmq,
		cfg:                   cfg,
		transactionRepository: transactionRepository,
	}
}
