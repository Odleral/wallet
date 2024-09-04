package rabbit

import (
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Pub struct {
	log      *zap.Logger
	conn     *amqp091.Connection
	channel  *amqp091.Channel
	exchange string
}

func NewPublisher(log *zap.Logger, conn *amqp091.Connection, exchange string) (*Pub, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := ch.ExchangeDeclare(exchange,
		"direct",
		true,
		false,
		false,
		false, nil); err != nil {
		return nil, err
	}

	return &Pub{
		log:      log,
		conn:     conn,
		channel:  ch,
		exchange: exchange,
	}, nil
}

func (p *Pub) Publish(routingKey string, msg amqp091.Publishing) error {
	if err := p.channel.Publish(p.exchange, routingKey, false, false, msg); err != nil {
		p.log.Error("p.channel.Publish failed", zap.Error(err))

		return err
	}

	return nil
}
