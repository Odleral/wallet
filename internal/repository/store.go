package repository

import "github.com/jmoiron/sqlx"

type Store struct {
	DB *sqlx.DB

	wallet *WalletRepo
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		DB: db,
	}
}
