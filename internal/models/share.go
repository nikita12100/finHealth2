package models

import "regexp"

type Share struct {
	Ticker    string
	HistPrice map[string]float32
}

type StatsShare struct {
	Weight        float64
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

type StatsBond struct {
	Count               int
	CountBuy            int
	AvgPriceBuy         float64
	LastPrice           float64
	CouponValue         float64
	SumPriceBuy         float64
	SumPriceTotal       float64
	CouponPeriodPerYear int
	Coup2025            float64
}

func IsShare(ticker string) bool {
	prefixPattern := `[A-Z0-9]+`
	re := regexp.MustCompile(prefixPattern)
	return re.MatchString(ticker) && !IsBond(ticker) && !IsCurrency(ticker)
}
