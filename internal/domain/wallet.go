package domain

type Wallet struct {
	ID         string
	Balance    float64
	Currency   string
	Authorised bool
}
