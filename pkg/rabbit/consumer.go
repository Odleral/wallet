package rabbit

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Consumer struct {
	log       *zap.Logger
	queueName string
	conn      *Connection
	channel   *amqp091.Channel
	handler   func(ctx context.Context, msg amqp091.Delivery)
}

func NewConsumer(log *zap.Logger,
	conn *Connection,
	queueName string,
	handler func(ctx context.Context, msg amqp091.Delivery)) *Consumer {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return &Consumer{
		log:       log,
		queueName: queueName,
		conn:      conn,
		channel:   ch,
		handler:   handler,
	}
}

func (c *Consumer) Consume() error {
	// defer recovery if panic
	defer func() {
		if r := recover(); r != nil {
			c.log.Error("recovered from panic", zap.Any("panic", r))
		}
	}()

	messages, err := c.channel.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range messages {
		c.handler(context.Background(), msg)
	}

	return nil
}

func (c *Consumer) QueueName() string {
	return c.queueName
}
