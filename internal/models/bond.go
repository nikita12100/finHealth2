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

func IsBond(ticker string) bool {
	prefixPattern := `^[A-Z]{2}[0-9]{2}[A-Z0-9]+`
	re := regexp.MustCompile(prefixPattern)
	return re.MatchString(ticker)
}
