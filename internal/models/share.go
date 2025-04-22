package models

import "regexp"

type Share struct {
	Ticker    string
	HistPrice map[string]float32
}

type StatsShare struct {
	Count         int
	CountBuy      int
	AvgPriceBuy   float64
	LastPrice     float64
	SumPriceBuy   float64
	SumPriceTotal float64
	Div           float64
	SumDiv        float64
	DivPerc       float64
}

type StatsTOM struct {
	Count         int
	CountBuy      int
	AvgPriceBuy   float64
	LastPrice     float64
	SumPriceBuy   float64
	SumPriceTotal float64
}

func IsShare(ticker string) bool {
	prefixPattern := `[A-Z0-9]+`
	re := regexp.MustCompile(prefixPattern)
	return re.MatchString(ticker) && !IsBond(ticker) && !IsCurrency(ticker)
}
