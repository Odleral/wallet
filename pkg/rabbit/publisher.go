package rabbit

import (
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Pub struct {
	log        *zap.Logger
	conn       *amqp091.Connection
	channel    *amqp091.Channel
	exchange   string
	RoutingKey string
	QueueName  string
}

func (c *Connection) NewPublisher(log *zap.Logger,
	exchange, rk, q string) (*Pub, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}

	if err = ch.ExchangeDeclare(exchange,
		"direct",
		true,
		false,
		false,
		false, nil); err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(q, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	if err = ch.QueueBind(queue.Name, rk, exchange, false, nil); err != nil {
		return nil, err
	}

	return &Pub{
		log:      log,
		conn:     c.conn,
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
