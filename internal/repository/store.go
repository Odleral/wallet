package repository

import "github.com/jmoiron/sqlx"

type Store struct {
	DB *sqlx.DB

	wallet      *WalletRepo
	transaction *TransactionRepo
	product     *ProductRepo
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		DB: db,
	}
}

func (s *Store) WalletRepo() *WalletRepo {
	if s.wallet != nil {
		return s.wallet
	}

	s.wallet = &WalletRepo{
		DB: s.DB,
	}

	return s.wallet
}

func (s *Store) TransactionRepo() *TransactionRepo {
	if s.transaction != nil {
		return s.transaction
	}

	s.transaction = &TransactionRepo{
		DB: s.DB,
	}

	return s.transaction
}

func (s *Store) ProductRepo() *ProductRepo {
	if s.product != nil {
		return s.product
	}

	s.product = &ProductRepo{
		DB: s.DB,
	}

	return s.product
}

func (s *Store) Close() error {
	if err := s.DB.Close(); err != nil {
		return err
	}

	return nil
}
