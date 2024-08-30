package domain

type Product struct {
	ID                                string
	Name                              string
	Description                       string
	MaxTransactionAmount              float64
	MinTransactionAmount              float64
	MaxDailyTransactionAmount         float64
	MaxMonthlyTransactionAmount       float64
	AuthorisedTransactionMinAmount    float64
	NotAuthorisedTransactionMaxAmount float64
}
