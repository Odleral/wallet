package domain

import "time"

type Wallet struct {
	ID         string
	Owner      string
	Balance    float64
	Currency   string
	Authorised bool
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}
