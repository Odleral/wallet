package transaction

import (
	"context"
	tt "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"wallet/internal/domain"
	"wallet/internal/errs"
)

const (
	_namespace = "internal.usecases.transaction"
)

type UseCase struct {
	log    *zap.Logger
	cache  Cache
	wallet WalletRepo
	tracer Tracer
	pub    TranPublisher
}

type (
	Cache interface {
		GetByID(ctx context.Context, id string) (domain.Wallet, error)
		Exists(ctx context.Context, id string) (bool, error)
		Update(ctx context.Context, w domain.Wallet) error
	}

	WalletRepo interface {
		GetByID(ctx context.Context, id string) (domain.Wallet, error)
		Update(ctx context.Context, w domain.Wallet) error
	}

	TranRepo interface {
		Create(ctx context.Context, t domain.Transaction) error
		Update(ctx context.Context, t domain.Transaction) error
	}

	TranPublisher interface {
		Publish(ctx context.Context, t domain.Transaction) error
	}

	Tracer interface {
		StartSpan(ctx context.Context, namespace string) (context.Context, tt.Span)
		Error(span tt.Span, method string, err error)
	}
)

func New(l *zap.Logger, c Cache, w WalletRepo, t Tracer, pub TranPublisher) *UseCase {
	return &UseCase{
		log:    l,
		cache:  c,
		wallet: w,
		tracer: t,
		pub:    pub,
	}
}

func (uc *UseCase) Execute(ctx context.Context, t domain.Transaction) error {
	ctx, span := uc.tracer.StartSpan(ctx, _namespace)
	defer span.End()

	l := uc.log.With(zap.String("wallet_id", t.From)).Named(_namespace)

	l.Info("transaction started")

	wallet, err := uc.cache.GetByID(ctx, t.From)
	if err != nil {
		l.Error("failed to get wallet", zap.Error(err))
		uc.tracer.Error(span, "GetByID", err)

		return err
	}

	if wallet.Balance < t.Amount {
		err = errs.ErrInsufficientBalance
		l.Error("insufficient balance", zap.Error(err))
		uc.tracer.Error(span, "insufficient balance", err)

		return err
	}

	incomeWalletExists, err := uc.cache.Exists(ctx, t.To)
	if err != nil {
		l.Error("failed to check wallet", zap.Error(err))
		uc.tracer.Error(span, "Exists", err)

		return err
	}

	if !incomeWalletExists {
		err = errs.ErrNotFound
		l.Error("wallet not found", zap.Error(err))
		uc.tracer.Error(span, "wallet not found", err)

		return err
	}

	if err := uc.pub.Publish(ctx, t); err != nil {
		l.Error("failed to publish transaction", zap.Error(err))
		uc.tracer.Error(span, "Publish", err)

		return err
	}

	return nil
}

func (uc *UseCase) TransactionUpdate(ctx context.Context, t domain.Transaction, s domain.StatusCode) error {
	ctx, span := uc.tracer.StartSpan(ctx, _namespace)
	defer span.End()

	l := uc.log.With(zap.String("wallet_id", t.From)).Named(_namespace)

	l.Info("transaction update", zap.String("status", string(s)), zap.String("transaction_id", t.ID))

	return nil
}
