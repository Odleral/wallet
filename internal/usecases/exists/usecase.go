package exists

import (
	"context"
	"go.uber.org/zap"
	"wallet/internal/domain"
)

const (
	_namespace = "internal.usecases.exists"

	_cacheWalletPrefix = "wallet:"
)

type UseCase struct {
	log    *zap.Logger
	cache  Cache
	wallet WalletRepo
}

type (
	Cache interface {
		Get(ctx context.Context, id string, dest any) error
	}

	WalletRepo interface {
		GetByID(ctx context.Context, id string) (domain.Wallet, error)
	}
)

func New(l *zap.Logger, c Cache, w WalletRepo) *UseCase {
	return &UseCase{
		log:    l,
		cache:  c,
		wallet: w,
	}
}

func (uc *UseCase) Execute(ctx context.Context, id string) (domain.Wallet, error) {
	l := uc.log.Named(_namespace)

	l.Info("checking wallet exists", zap.String("wallet_id", id))

	var wallet domain.Wallet

	if err := uc.cache.Get(ctx, _cacheWalletPrefix+id, &wallet); err == nil {
		l.Info("wallet exists in cache")

		return wallet, nil
	}

	wallet, err := uc.wallet.GetByID(ctx, id)
	if err != nil {
		l.Error("wallet.GetByID failed", zap.Error(err))

		return domain.Wallet{}, err
	}

	return wallet, nil
}
