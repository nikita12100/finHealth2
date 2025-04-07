package models

import "time"

type Operation struct {
	NoTransaction string    `json:"id_transaction"`
	IsBuy         bool      `json:"is_buy"`
	Ticker        string    `json:"ticker"`
	Price         float64   `json:"price"`
	Count         int       `json:"count"`
	Currency      string    `json:"currency"`
	Fee           float32   `json:"fee"`
	Date          time.Time `json:"date"`
}
