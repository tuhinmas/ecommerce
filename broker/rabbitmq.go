package broker

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitmqConfig struct {
	Host                    string
	Username                string
	Password                string
	Port                    int
	StockReversalExchange   string
	StockReversalQueue      string
	StockReversalRoutingKey string
	PendingPaymentQueue     string
}

type RabbitMQ interface {
	Connect() (err error)
	Close()
	Reconnect() error
	GetConfig() rabbitMQ
}

type rabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Err     chan error
	config  RabbitmqConfig
}

func NewConnection(config RabbitmqConfig) RabbitMQ {
	return &rabbitMQ{
		config: config,
		Err:    make(chan error),
	}
}

func (c *rabbitMQ) GetConfig() rabbitMQ {
	return *c
}

func (c *rabbitMQ) Connect() (err error) {
	connPattern := "amqp://%v:%v@%v:%v"

	clientUrl := fmt.Sprintf(connPattern,
		c.config.Username,
		c.config.Password,
		c.config.Host,
		c.config.Port,
	)

	c.Conn, err = amqp.Dial(clientUrl)
	if err != nil {
		if err = c.Retry(); err != nil {
			err = fmt.Errorf("failed to connect to rabbitmq: %v", err)
			return
		}
	}

	c.Channel, err = c.Conn.Channel()
	if err != nil {
		err = fmt.Errorf("failed to create channel to rabbitmq: %v", err)
		return
	}

	if err = c.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		err = fmt.Errorf("failed to setup queue to rabbitmq: %v", err)
		return
	}

	return
}

func (c *rabbitMQ) Retry() (err error) {
	fmt.Println("Retrying to connect to rabbitmq")
	connPattern := "amqp://%v:%v@%v:%v"

	clientUrl := fmt.Sprintf(connPattern,
		c.config.Username,
		c.config.Password,
		c.config.Host,
		c.config.Port,
	)

	conn, err := amqp.Dial(clientUrl)
	if err != nil {
		err = fmt.Errorf("failed to connect to rabbitmq: %v", err)
		return
	}

	c.Conn = conn

	return
}

func (c *rabbitMQ) Close() {
	c.Conn.Close()
}

func (c *rabbitMQ) Reconnect() error {
	fmt.Println("Retrying to connect to rabbitmq")
	if err := c.Connect(); err != nil {
		return err
	}
	return nil
}
