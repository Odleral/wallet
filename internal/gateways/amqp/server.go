package amqp

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"sync"
	"wallet/internal/config"
	"wallet/internal/domain"
	"wallet/internal/worker/transaction"
	"wallet/pkg/rabbit"
)

type AMQP struct {
	wg        *sync.WaitGroup
	log       *zap.Logger
	Consumers []*rabbit.Consumer

	transactionWorker *transaction.Worker
}

func NewAMQP(log *zap.Logger, cfg config.Config, conn *rabbit.Connection,
	tw *transaction.Worker) *AMQP {
	wg := &sync.WaitGroup{}

	amqp := &AMQP{
		wg:                wg,
		log:               log,
		transactionWorker: tw,
	}

	consumers := []*rabbit.Consumer{
		rabbit.NewConsumer(log, conn, cfg.TransactionSub.QueueName, amqp.transactionHandler),
	}

	return &AMQP{
		wg:        wg,
		Consumers: consumers,
	}
}

func (amqp *AMQP) Run() error {
	for _, consumer := range amqp.Consumers {
		go func(consumer *rabbit.Consumer) {
			defer amqp.wg.Done()

			if err := consumer.Consume(); err != nil {
				amqp.log.Error("failed to consume", zap.Error(err), zap.String("queue", consumer.QueueName()))
			}
		}(consumer)
	}

	return nil
}

type transactionMessage struct {
	TransactionID string `json:"transaction_id"`
	CorrelationID string `json:"correlation_id"`
}

func (amqp *AMQP) transactionHandler(ctx context.Context, msg amqp091.Delivery) {
	log := amqp.log.Named("transactionHandler")

	var m transactionMessage
	if err := json.Unmarshal(msg.Body, &m); err != nil {
		log.Error("failed to unmarshal message", zap.Error(err))

		_ = msg.Nack(false, false)

		return
	}

	log.Info("received message", zap.String("transaction_id", m.TransactionID),
		zap.String("correlation_id", m.CorrelationID))

	if m.TransactionID == "" || m.CorrelationID == "" {
		log.Error("invalid message")
		_ = msg.Nack(false, false)

		return
	}

	// transactionWorker
	if err := amqp.transactionWorker.Execute(ctx, domain.Transaction{
		ID:            m.TransactionID,
		CorrelationID: m.CorrelationID,
	}); err != nil {
		log.Error("failed to execute transaction", zap.Error(err))

		_ = msg.Nack(false, false)

		return
	}

	_ = msg.Ack(false)
}
