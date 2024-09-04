package rabbit

import (
	"github.com/rabbitmq/amqp091-go"
	"wallet/internal/config"
)

type Connection struct {
	conn   *amqp091.Connection
	chanel *amqp091.Channel
}

func NewConnection(cfg config.Config) *Connection {
	conn, err := amqp091.Dial(cfg.AMQP)
	if err != nil {
		panic(err)
	}

	chanel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return &Connection{
		chanel: chanel,
		conn:   conn,
	}
}

func (c *Connection) Channel() (*amqp091.Channel, error) {
	return c.conn.Channel()
}

func (c *Connection) Conn() *amqp091.Connection {
	return c.conn
}

func (c *Connection) Close() error {
	if err := c.chanel.Close(); err != nil {
		return err
	}

	if err := c.conn.Close(); err != nil {
		return err
	}

	return nil
}
