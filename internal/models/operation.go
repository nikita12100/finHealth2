package models

type Operation struct {
	IsBuy    bool    `json:"is_buy"`
	Ticker   string  `json:"ticker"`
	Price    float64 `json:"price"`
	Count    int     `json:"count"`
	Currency string  `json:"currency"`
	Fee      float32 `json:"fee"`
}
