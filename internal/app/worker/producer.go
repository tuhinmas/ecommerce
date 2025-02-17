package worker

import (
	"context"
	"ecommerce/pkg/logger"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"ecommerce/pkg/constant"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (q queueService) PublishData(ctx context.Context, topic string, msg interface{}) (err error) {
	fmt.Println("Publishing data...")
	cfg := q.rabbitmq.GetConfig()

	select {
	case err := <-cfg.Err:
		if err != nil {
			// Check and handle the error returned by Reconnect
			if reconnectErr := q.rabbitmq.Reconnect(); reconnectErr != nil {
				return reconnectErr
			}
		}
	default:
	}

	ch, err := cfg.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(msg)
	if err != nil {
		return
	}

	// Declare Dead Letter Exchange (DLX) for stock reversal
	err = ch.ExchangeDeclare(q.cfg.StockReversalExchange, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Declare the stock reversal queue (final processing after delay)
	_, err = ch.QueueDeclare(q.cfg.StockReversalQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Bind DLX to stock reversal queue
	err = ch.QueueBind(q.cfg.StockReversalQueue, q.cfg.StockReversalRoutingKey, q.cfg.StockReversalExchange, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Declare Delay Queue (Pending Payment Queue) with TTL N minutes
	_, err = ch.QueueDeclare(q.cfg.PendingPaymentQueue, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    q.cfg.StockReversalExchange,
		"x-dead-letter-routing-key": q.cfg.StockReversalRoutingKey,
		"x-message-ttl":             300000, // TTL 5 minutes (in milliseconds)
	})
	if err != nil {
		log.Fatal(err)
	}

	// Publish a message to the delay queue
	err = ch.Publish(
		"",
		q.cfg.PendingPaymentQueue,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,    // keeping message if broker restart
			ContentType:  "application/json", // XXX: We will revisit this in future episodes
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	logger.LogInfoWithData(constant.PRODUCER, msg, "", "", constant.MsgMessageSent)

	log.Println("Message successfully sent to pending_payment_queue with TTL 5 minutes!")

	return

}
