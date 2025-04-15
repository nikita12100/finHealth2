package models

import "regexp"

type Bond struct {
	Ticker string
}

type StockBondInfo struct {
	CouponValue  float64 `json:"coupon_value"`
	LastPrice    float64 `json:"last_price"`
	FaceValue    float64 `json:"face_value"`
	CouponPeriod float64 `json:"coupon_period"`
	FaceUnit     string  `json:"face_unit"`
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

func IsBond(ticker string) bool {
	prefixPattern := `^[A-Z]{2}[0-9]{2}[A-Z0-9]+`
	re := regexp.MustCompile(prefixPattern)
	return re.MatchString(ticker)
}

func IsCurrency(ticker string) bool {
	return len(ticker) > 5 && ticker[len(ticker)-4:] == "_TOM"
}
