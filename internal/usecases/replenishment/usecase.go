package replenishment

import (
	"context"
	tt "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"wallet/internal/domain"
	"wallet/pkg/tracer"
)

const (
	_namespace = "internal.usecases.transaction"
)

type UseCase struct {
	log    *zap.Logger
	tracer Tracer

	cache Cache

	tran   TranRepo
	wallet WalletRepo

	pub TranPublisher
}

type (
	Cache interface {
		Exists(ctx context.Context, id string) (bool, error)
	}

	WalletRepo interface {
		GetByID(ctx context.Context, id string) (domain.Wallet, error)
	}

	TranRepo interface {
		Create(ctx context.Context, t domain.Transaction) (string, error)
	}

	TranPublisher interface {
		Publish(t domain.Transaction) error
	}

	Tracer interface {
		StartSpan(ctx context.Context, namespace string) (context.Context, tt.Span)
	}
)

func New(l *zap.Logger, c Cache, w WalletRepo,
	tran TranRepo, t Tracer, pub TranPublisher) *UseCase {
	return &UseCase{
		log:    l,
		cache:  c,
		wallet: w,
		tracer: t,
		tran:   tran,
		pub:    pub,
	}
}

func (uc *UseCase) Execute(ctx context.Context, t domain.Transaction) error {
	ctx, span := uc.tracer.StartSpan(ctx, _namespace)
	defer span.End()

	l := uc.log.Named(_namespace).With(zap.String("CorrelationID", t.CorrelationID))

	var wallet domain.Wallet

	l.Info("transaction start")

	walletExists, err := uc.cache.Exists(ctx, t.WalletID)
	if err != nil {
		l.Error("uc.cache.Exists", zap.Error(err))
		tracer.Error(span, "uc.cache.Exists", err)

		return err
	}

	if !walletExists {
		wallet, err = uc.wallet.GetByID(ctx, t.WalletID)
		if err != nil {
			l.Error("uc.cache.GetByID", zap.Error(err))
			tracer.Error(span, "uc.cache.GetByID", err)

			return err
		}
	}

	l.Info("wallet exists", zap.String("WalletID", wallet.ID))

	id, err := uc.tran.Create(ctx, t)
	if err != nil {
		l.Error("uc.tran.Create", zap.Error(err))
		tracer.Error(span, "uc.tran.Create", err)

		return err
	}

	t.ID = id

	if err = uc.pub.Publish(t); err != nil {
		l.Error("uc.pub.Publish", zap.Error(err))
		tracer.Error(span, "uc.pub.Publish", err)

		return err
	}

	return nil
}
