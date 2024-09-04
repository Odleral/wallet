package domain

import "time"

type Product struct {
	ID                             string
	Name                           string
	Description                    string
	MaxTransactionAmount           float64
	MinTransactionAmount           float64
	AuthorisedMaxTransactionAmount float64
	CreatedAt                      time.Time
	UpdatedAt                      time.Time
}
