package domain

import "time"

type StatusCode string

const (
	Success StatusCode = "success"
	Failed  StatusCode = "failed"
	Pending StatusCode = "pending"
)

type Transaction struct {
	ID            string
	Amount        float64
	Currency      string
	CorrelationID string
	WalletID      string
	Status        StatusCode
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}
