package worker

import (
	"context"
	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/logger"
	"encoding/json"
	"errors"
	"strings"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (q queueService) ConsumeData(ctx context.Context, queueName string) (err error) {
	if q.rabbitmq == nil {
		return errors.New("consumer failed to start")
	}
	cfg := q.rabbitmq.GetConfig()
	notify := cfg.Conn.NotifyClose(make(chan *amqp.Error)) // Notify channel for connection errors

	ch, err := cfg.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	//logger.LogInfo(constant.CONSUMER, fmt.Sprintf(constant.START_LISTENING_TOPIC_FROM_BROKER, topic))

	for {
		select {
		case notifyErr := <-notify:
			if notifyErr != nil {
				// Log the error
				//logger.LogError(constant.CONSUMER, constant.ErrConsumeQueueToBroker, notifyErr.Error(), "", "", "")

				// Try to reconnect
				for i := 0; i < 3; i++ {
					err = q.rabbitmq.Reconnect()
					if err == nil {
						break
					}
				}
			}
		case msg := <-msgs:

			switch msg.RoutingKey {
			case q.cfg.StockReversalRoutingKey:
				logger.LogInfo(constant.CONSUMER, "consume check stock reversal")
				var payloadQueue entity.PayloadOrderQueue
				if err = json.Unmarshal(msg.Body, &payloadQueue); err != nil {
					logger.LogError(constant.CONSUMER, err.Error(), payloadQueue, "", "")
					msg.Ack(false)
					continue
				}

				order, err := q.transactionRepository.GetOrderById(ctx, payloadQueue.OrderId)
				if err != nil {
					logger.LogError(constant.CONSUMER, err.Error(), payloadQueue, "", "")
					msg.Ack(false)
					continue
				}

				if strings.ToLower(order.Status) == "pending" {
					tx, err := q.transactionRepository.BeginTx(ctx)
					if err != nil {
						logger.LogError(constant.CONSUMER, err.Error(), payloadQueue, "", "")
						msg.Nack(false, true)
						continue
					}

					if err := q.transactionRepository.UpdateOrderStatus(ctx, tx, payloadQueue.OrderId, "cancelled"); err != nil {
						logger.LogError(constant.CONSUMER, err.Error(), payloadQueue, "", "")
						msg.Nack(false, true)
						continue
					}

					var wg sync.WaitGroup
					errCh := make(chan error, len(payloadQueue.Sku)) // Buffered channel for errors

					for _, sku := range payloadQueue.Sku {
						wg.Add(1)
						go func(sku entity.SkuRequest) {
							defer wg.Done()

							if err := q.transactionRepository.ReverseStock(ctx, tx, payloadQueue.WarehouseId, sku); err != nil {
								errCh <- err
							}
						}(sku)
					}

					wg.Wait()
					close(errCh)

					// Check if any error occurred
					for err := range errCh {
						if err != nil {
							logger.LogError(constant.CONSUMER, err.Error(), payloadQueue, "", "")
							tx.Rollback()
							msg.Nack(false, true)
							continue
						}
					}

					if err := q.transactionRepository.CommitTx(ctx, tx); err != nil {
						logger.LogError(constant.CONSUMER, err.Error(), payloadQueue, "", "")
						msg.Nack(false, true)
						continue
					}
				}

				if ackErr := msg.Ack(false); ackErr != nil {
					logger.LogError(constant.CONSUMER, ackErr.Error(), payloadQueue, "", "")
					msg.Nack(false, true)
				}
				logger.LogInfoWithData(constant.CONSUMER, payloadQueue, "", "", constant.MsgAcknowledgeMessage)
			}
		}
	}
}
