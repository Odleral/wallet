package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
	"wallet/internal/domain"
)

type ProductRepo struct {
	DB *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) *ProductRepo {
	return &ProductRepo{
		DB: db,
	}
}

func (p *ProductRepo) Get(ctx context.Context) (domain.Product, error) {
	var dbPro dbProduct
	err := p.DB.GetContext(ctx, &dbPro, "SELECT * FROM product WHERE name = $1", "product")
	if err != nil {
		return domain.Product{}, err
	}

	return dbPro.toDomain(), nil
}

type dbProduct struct {
	ID                             string     `db:"id"`
	Name                           string     `db:"name"`
	Description                    string     `db:"description"`
	MaxTransactionAmount           float64    `db:"max_transaction_amount"`
	MinTransactionAmount           float64    `db:"min_transaction_amount"`
	AuthorisedMaxTransactionAmount float64    `db:"authorised_max_transaction_amount"`
	CreatedAt                      time.Time  `db:"created_at"`
	UpdatedAt                      *time.Time `db:"updated_at"`
}

func (dto *dbProduct) toDomain() domain.Product {
	var updatedAt time.Time

	if dto.UpdatedAt != nil {
		updatedAt = *dto.UpdatedAt
	}
	return domain.Product{
		ID:                             dto.ID,
		Name:                           dto.Name,
		Description:                    dto.Description,
		MaxTransactionAmount:           dto.MaxTransactionAmount,
		MinTransactionAmount:           dto.MinTransactionAmount,
		AuthorisedMaxTransactionAmount: dto.AuthorisedMaxTransactionAmount,
		CreatedAt:                      dto.CreatedAt,
		UpdatedAt:                      updatedAt,
	}
}
