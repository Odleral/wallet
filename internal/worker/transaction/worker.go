package transaction

import (
	"context"
	"errors"
	tt "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
	"wallet/internal/domain"
	"wallet/internal/errs"
	"wallet/pkg/tracer"
)

const (
	_namespace = "internal.usecases.transaction"

	_cacheKey          = "product"
	_cacheWalletPrefix = "wallet:"
)

type Worker struct {
	log    *zap.Logger
	tracer Tracer

	cache Cache

	tran    TranRepo
	wallet  WalletRepo
	product ProductRepo
}

type (
	Cache interface {
		Exists(ctx context.Context, id string) (bool, error)
		Set(ctx context.Context, key string, val any, ttl time.Duration) error
		Get(ctx context.Context, key string, dest any) error
	}

	WalletRepo interface {
		GetByID(ctx context.Context, id string) (domain.Wallet, error)
		Update(ctx context.Context, w domain.Wallet) error
	}

	TranRepo interface {
		Update(ctx context.Context, t domain.Transaction) error
		Transfer(ctx context.Context, t domain.Transaction) error
	}

	TranPublisher interface {
		Publish(ctx context.Context, t domain.Transaction) error
	}

	ProductRepo interface {
		Get(ctx context.Context) (domain.Product, error)
	}

	Tracer interface {
		StartSpan(ctx context.Context, namespace string) (context.Context, tt.Span)
	}
)

func New(log *zap.Logger, tracer Tracer, cache Cache, tran TranRepo, wallet WalletRepo, product ProductRepo) *Worker {
	return &Worker{
		log:    log,
		tracer: tracer,

		cache: cache,

		tran:    tran,
		wallet:  wallet,
		product: product,
	}
}

// Execute performs a transaction
func (w *Worker) Execute(ctx context.Context, t domain.Transaction) error {
	ctx, span := w.tracer.StartSpan(ctx, _namespace)
	defer span.End()

	l := w.log.With(zap.String("wallet_id", t.WalletID)).Named(_namespace)

	var wallet domain.Wallet

	err := w.cache.Get(ctx, _cacheWalletPrefix+t.WalletID, &wallet)
	if err != nil {
		tracer.Error(span, "w.cache.GetByID", err)
		l.Error("w.cache.GetByID", zap.Error(err))

		w.TransactionUpdate(ctx, t, domain.Failed)

		return err
	}

	limit, err := w.GetLimit(ctx, wallet)
	if err != nil {
		l.Error("uc.GetLimit", zap.Error(err))
		tracer.Error(span, "uc.GetLimit", err)

		w.TransactionUpdate(ctx, t, domain.Failed)

		return err
	}

	if limit < wallet.Balance+t.Amount {
		l.Error("balance over limit", zap.Error(errs.ErrOverLimit))

		w.TransactionUpdate(ctx, t, domain.Failed)

		return err
	}

	if err = w.tran.Transfer(ctx, t); err != nil {
		l.Error("uc.tran.Transfer", zap.Error(err))
		tracer.Error(span, "uc.tran.Transfer", err)

		w.TransactionUpdate(ctx, t, domain.Failed)

		return err
	}

	w.TransactionUpdate(ctx, t, domain.Success)

	wallet, err = w.wallet.GetByID(ctx, t.WalletID)
	if err != nil {
		l.Error("uc.wallet.GetByID", zap.Error(err))
		tracer.Error(span, "uc.wallet.GetByID", err)

		return err
	}

	if err = w.cache.Set(ctx, _cacheWalletPrefix+t.WalletID, wallet, 0); err != nil {
		l.Error("uc.cache.Set", zap.Error(err))
		tracer.Error(span, "uc.cache.Set", err)

		return err
	}

	return nil
}

func (w *Worker) GetLimit(ctx context.Context, wallet domain.Wallet) (float64, error) {
	ctx, span := w.tracer.StartSpan(ctx, _namespace+".GetLimit")
	defer span.End()

	l := w.log.Named(_namespace + ".GetLimit")

	var product domain.Product

	err := w.cache.Get(ctx, _cacheKey, &product)
	if err == nil {
		l.Info("limit by product returned from cache")

		if wallet.Authorised {
			return product.AuthorisedMaxTransactionAmount, nil
		} else {
			return product.MaxTransactionAmount, nil
		}
	}

	if !errors.Is(err, errs.ErrNotFound) {
		l.Error("uc.cache.Get", zap.Error(err))
		tracer.Error(span, "uc.pub.Publish", err)

		return 0, err
	}

	l.Warn("product not fount in cache")

	product, err = w.product.Get(ctx)
	if err != nil {
		l.Error("uc.product.Get", zap.Error(err))
		tracer.Error(span, "uc.product.Get", err)

		return 0, err
	}

	l.Info("limit by product returned from DB")

	if wallet.Authorised {
		return product.AuthorisedMaxTransactionAmount, nil
	} else {
		return product.MaxTransactionAmount, nil
	}
}

func (w *Worker) TransactionUpdate(ctx context.Context, t domain.Transaction, s domain.StatusCode) {
	ctx, span := w.tracer.StartSpan(ctx, _namespace+".TransactionUpdate")
	defer span.End()

	l := w.log.Named(_namespace + ".TransactionUpdate")

	t.Status = s

	err := w.tran.Update(ctx, t)
	if err != nil {
		l.Error("uc.tran.Update", zap.Error(err))
		tracer.Error(span, "uc.tran.Update", err)
	}
}
