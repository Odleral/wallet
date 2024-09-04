package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
	"wallet/internal/domain"
	"wallet/internal/errs"
)

const (
	_qGetWalletByID = "SELECT * FROM wallet WHERE id = $1"
	_qUpdateWallet  = "UPDATE wallet SET balance = $1 WHERE id = $2"
)

type WalletRepo struct {
	DB *sqlx.DB
}

func NewWalletRepo(db *sqlx.DB) *WalletRepo {
	return &WalletRepo{
		DB: db,
	}
}

func (w *WalletRepo) GetByID(ctx context.Context, id string) (domain.Wallet, error) {
	var dtoWallet dbWallet

	tx, err := w.DB.Beginx()
	if err != nil {
		return domain.Wallet{}, errs.Wrap(err)
	}

	err = tx.GetContext(ctx, &dtoWallet, _qGetWalletByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Wallet{}, errs.ErrNotFound
		}

		return domain.Wallet{}, err
	}

	if err := tx.Commit(); err != nil {
		return domain.Wallet{}, errs.Wrap(err)
	}

	return dtoWallet.toModel(), nil
}

func (w *WalletRepo) Update(ctx context.Context, wallet domain.Wallet) error {
	tx, err := w.DB.Beginx()
	if err != nil {
		return errs.Wrap(err)
	}

	_, err = tx.ExecContext(ctx, _qUpdateWallet, wallet.Balance, wallet.ID)
	if err != nil {
		return errs.Wrap(err)
	}

	if err := tx.Commit(); err != nil {
		return errs.Wrap(err)
	}

	return nil
}

type dbWallet struct {
	ID         string     `db:"id"`
	Balance    float64    `db:"balance"`
	Owner      string     `db:"owner"`
	Currency   string     `db:"currency"`
	Authorised bool       `db:"authorised"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

func (d dbWallet) toModel() domain.Wallet {
	var updatedAt *time.Time

	if d.UpdatedAt != nil {
		updatedAt = d.UpdatedAt
	}

	return domain.Wallet{
		ID:         d.ID,
		Balance:    d.Balance,
		Owner:      d.Owner,
		Currency:   d.Currency,
		Authorised: d.Authorised,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  updatedAt,
	}
}
