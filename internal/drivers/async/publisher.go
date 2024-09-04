package async

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"wallet/internal/config"
	"wallet/internal/domain"
	"wallet/pkg/rabbit"
)

type Publisher struct {
	pub *rabbit.Pub
	cfg config.Publisher
	rk  string
}

type transactionMessage struct {
	TransactionID string `json:"transaction_id"`
	CorrelationID string `json:"correlation_id"`
}

func New(log *zap.Logger, cfg config.Publisher, conn *rabbit.Connection) *Publisher {
	pub, err := conn.NewPublisher(log, cfg.Exchange, cfg.RoutingKey, cfg.QueueName)
	if err != nil {
		log.Error("conn.NewPublisher failed", zap.Error(err))
	}

	return &Publisher{
		pub: pub,
		cfg: cfg,
		rk:  cfg.RoutingKey,
	}
}

func (p *Publisher) Publish(t domain.Transaction) error {
	m := transactionMessage{
		TransactionID: t.ID,
		CorrelationID: t.CorrelationID,
	}

	raw, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return p.pub.Publish(p.rk, amqp091.Publishing{
		ContentType: "application/json",
		Body:        raw,
	})
}
