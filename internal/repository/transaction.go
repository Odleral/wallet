package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"wallet/internal/domain"
	"wallet/internal/errs"
)

const (
	_qCreateTransaction = `INSERT INTO transaction (amount, currency, from, to, status) 
							VALUES ($1, $2, $3, $4, $5) RETURNING id`

	_qTransferUpdateBalanceSum = `UPDATE wallet SET balance = balance + $1 WHERE id = $2;`
	_qTransferUpdateTranStatus = `UPDATE transaction SET status = $1 WHERE id = $2;`
)

type TransactionRepo struct {
	DB *sqlx.DB
}

func NewTransactionRepo(db *sqlx.DB) *TransactionRepo {
	return &TransactionRepo{
		DB: db,
	}
}

func (t *TransactionRepo) Create(ctx context.Context, transaction domain.Transaction) (string, error) {
	tx, err := t.DB.Beginx()
	if err != nil {
		return "", errs.Wrap(err)
	}

	result, err := tx.QueryContext(ctx, _qCreateTransaction,
		transaction.Amount,
		transaction.Currency,
		transaction.CorrelationID,
		transaction.WalletID,
		transaction.Status)
	if err != nil {
		return "", errs.Wrap(err)
	}

	if err := tx.Commit(); err != nil {
		return "", errs.Wrap(err)
	}

	var id string
	for result.Next() {
		if err := result.Scan(&id); err != nil {
			return "", errs.Wrap(err)
		}
	}

	return id, nil
}

func (t *TransactionRepo) Update(ctx context.Context, transaction domain.Transaction) error {
	tx, err := t.DB.Beginx()
	if err != nil {
		return errs.Wrap(err)
	}

	_, err = tx.ExecContext(ctx, _qTransferUpdateTranStatus, transaction.Status, transaction.ID)
	if err != nil {
		err = tx.Rollback() //nolint:errcheck

		return errs.Wrap(err)
	}

	return nil
}

func (t *TransactionRepo) Transfer(ctx context.Context, transaction domain.Transaction) error {
	tx, err := t.DB.Beginx()
	if err != nil {
		return errs.Wrap(err)
	}

	_, err = tx.ExecContext(ctx, _qTransferUpdateBalanceSum, transaction.Amount, transaction.WalletID)
	if err != nil {
		tx.Rollback() //nolint:errcheck
		return errs.Wrap(err)
	}

	_, err = tx.ExecContext(ctx, _qTransferUpdateTranStatus, domain.Success, transaction.ID)
	if err != nil {
		tx.Rollback() //nolint:errcheck
		return errs.Wrap(err)
	}

	if err := tx.Commit(); err != nil {
		return errs.Wrap(err)
	}

	return nil
}
