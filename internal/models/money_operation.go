package models

import "time"

type OperationType int

const (
	Unknown       OperationType = iota // 0
	Replenishment                      // 1
	BuyOrSell                          // 2
	Commision                          // 3
	DFP                                // 4
	DVP                                // 5
	Coupon                             // 6
	BondExpire                         // 7
	Tax                                // 8
	TaxYear                            // 9
	TaxDiv                             // 10
	Dividends                          // 11
	Withdraw                           // 12
	Repo                           	   // 13
)

type MoneyOperation struct {
	Time          time.Time
	OperationType OperationType
	AmountIn      float64
	AmountOut     float64
	Comment       string
}
