package models

type Operation struct {
	IsBuy bool
	Ticker string
	Price float64
	Count int
	Currency string
	Fee float32
}