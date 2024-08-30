package domain

type Transaction struct {
	ID       string
	Amount   float64
	Currency string
	From     string
	To       string
}
